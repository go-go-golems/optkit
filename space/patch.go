package space

import (
	"bytes"
	"context"
	"fmt"
	"sort"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/record"
)

type AssignmentRecord struct {
	Variable          VariableID    `json:"variable"`
	ExpectedOldDigest record.Digest `json:"expected_old_digest"`
	Value             artifact.Ref  `json:"value"`
}

type PatchRecord struct {
	ID          record.PatchID     `json:"id"`
	Base        record.SnapshotID  `json:"base"`
	Assignments []AssignmentRecord `json:"assignments"`
	Result      record.SnapshotID  `json:"result"`
}

type Patch[C any] struct {
	PatchRecord
	ResultConfig C `json:"-"`
}

type assignment[C any] interface {
	variableID() VariableID
	apply(context.Context, artifact.Store, C) (C, AssignmentRecord, error)
}

type typedAssignment[C, V any] struct {
	variable Variable[C, V]
	value    V
}

func (a typedAssignment[C, V]) variableID() VariableID { return a.variable.Descriptor.ID }

func (a typedAssignment[C, V]) apply(ctx context.Context, store artifact.Store, config C) (C, AssignmentRecord, error) {
	if err := a.variable.Validate(); err != nil {
		return config, AssignmentRecord{}, err
	}
	if err := a.variable.Domain.Validate(a.value); err != nil {
		return config, AssignmentRecord{}, fmt.Errorf("variable %s: %w", a.variable.Descriptor.ID, err)
	}
	oldValue := a.variable.Lens.Get(config)
	oldBytes, err := a.variable.Codec.EncodeCanonical(oldValue)
	if err != nil {
		return config, AssignmentRecord{}, fmt.Errorf("encode old value for %s: %w", a.variable.Descriptor.ID, err)
	}
	newBytes, err := a.variable.Codec.EncodeCanonical(a.value)
	if err != nil {
		return config, AssignmentRecord{}, fmt.Errorf("encode new value for %s: %w", a.variable.Descriptor.ID, err)
	}
	updated, err := a.variable.Lens.Put(config, a.value)
	if err != nil {
		return config, AssignmentRecord{}, fmt.Errorf("apply variable %s: %w", a.variable.Descriptor.ID, err)
	}
	ref, err := store.Put(ctx, artifact.PutRequest{
		MediaType:   "application/vnd.optkit.record+json",
		Schema:      schemaPtr(a.variable.Codec.Schema()),
		Sensitivity: sensitivityForVariable(a.variable.Descriptor),
	}, bytes.NewReader(newBytes))
	if err != nil {
		return config, AssignmentRecord{}, fmt.Errorf("store value for %s: %w", a.variable.Descriptor.ID, err)
	}
	return updated, AssignmentRecord{
		Variable:          a.variable.Descriptor.ID,
		ExpectedOldDigest: record.SumBytes(oldBytes),
		Value:             ref,
	}, nil
}

type PatchBuilder[C any] struct {
	base        Snapshot[C]
	store       artifact.Store
	configCodec Codec[C]
	items       map[VariableID]assignment[C]
}

func NewPatchBuilder[C any](base Snapshot[C], store artifact.Store, configCodec Codec[C]) *PatchBuilder[C] {
	return &PatchBuilder[C]{base: base, store: store, configCodec: configCodec, items: make(map[VariableID]assignment[C])}
}

func Set[C, V any](builder *PatchBuilder[C], variable Variable[C, V], value V) error {
	if builder == nil {
		return fmt.Errorf("patch builder is nil")
	}
	if err := variable.Validate(); err != nil {
		return err
	}
	id := variable.Descriptor.ID
	if _, exists := builder.items[id]; exists {
		return fmt.Errorf("conflicting assignment for variable %s", id)
	}
	builder.items[id] = typedAssignment[C, V]{variable: variable, value: value}
	return nil
}

func (b *PatchBuilder[C]) Build(ctx context.Context) (Patch[C], Snapshot[C], error) {
	if b == nil || b.store == nil || b.configCodec == nil {
		return Patch[C]{}, Snapshot[C]{}, fmt.Errorf("patch builder is incomplete")
	}
	if len(b.items) == 0 {
		return Patch[C]{}, Snapshot[C]{}, fmt.Errorf("patch requires at least one assignment")
	}
	ids := make([]VariableID, 0, len(b.items))
	for id := range b.items {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	config := b.base.Value
	records := make([]AssignmentRecord, 0, len(ids))
	for _, id := range ids {
		updated, rec, err := b.items[id].apply(ctx, b.store, config)
		if err != nil {
			return Patch[C]{}, Snapshot[C]{}, err
		}
		config = updated
		records = append(records, rec)
	}
	child, err := MaterializeSnapshot(ctx, b.store, b.base.System, b.configCodec, config)
	if err != nil {
		return Patch[C]{}, Snapshot[C]{}, err
	}
	identity := struct {
		Base        record.SnapshotID  `json:"base"`
		Assignments []AssignmentRecord `json:"assignments"`
		Result      record.SnapshotID  `json:"result"`
	}{Base: b.base.ID, Assignments: records, Result: child.ID}
	digest, _, err := record.SemanticDigest("schema:optkit.patch-identity/v1", identity)
	if err != nil {
		return Patch[C]{}, Snapshot[C]{}, err
	}
	rawID, err := record.ContentID("patch", digest)
	if err != nil {
		return Patch[C]{}, Snapshot[C]{}, err
	}
	patch := Patch[C]{
		PatchRecord: PatchRecord{
			ID:          record.PatchID(rawID),
			Base:        b.base.ID,
			Assignments: records,
			Result:      child.ID,
		},
		ResultConfig: config,
	}
	return patch, child, nil
}

func (p Patch[C]) VerifyBase(base Snapshot[C]) error {
	if p.Base != base.ID {
		return fmt.Errorf("patch base %s does not match snapshot %s", p.Base, base.ID)
	}
	return nil
}

func sensitivityForVariable(desc VariableDescriptor) artifact.Sensitivity {
	if desc.Sensitive {
		return artifact.SensitivityRestricted
	}
	return artifact.SensitivityInternal
}
