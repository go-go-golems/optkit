// Package system binds immutable Optkit snapshots to product-owned executable
// factories without placing provider handles or credentials in durable state.
package system

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/episode"
	"github.com/go-go-golems/optkit/record"
	"github.com/go-go-golems/optkit/space"
)

// Factory is the composition-root registration for one executable system.
// Implementations may retain process-local providers and credentials. Prepare
// must verify and decode the immutable snapshot before returning an executable.
type Factory interface {
	SystemID() record.SystemID
	ConfigSchema() record.SchemaID
	CaseSchema() record.SchemaID
	Prepare(context.Context, artifact.Store, space.SnapshotRecord) (Prepared, error)
}

// Prepared executes cases against one exact immutable snapshot. Case payloads
// remain artifact references so workers can lease compact, durable work items.
type Prepared interface {
	SystemID() record.SystemID
	SnapshotID() record.SnapshotID
	CaseSchema() record.SchemaID
	Run(context.Context, artifact.Ref, int64, episode.Sink) (episode.RunResult, error)
}

type Registry struct {
	mu        sync.RWMutex
	factories map[record.SystemID]Factory
}

func NewRegistry() *Registry {
	return &Registry{factories: make(map[record.SystemID]Factory)}
}

func (r *Registry) Register(factory Factory) error {
	if r == nil {
		return fmt.Errorf("system registry is nil")
	}
	if factory == nil {
		return fmt.Errorf("system factory is required")
	}
	if err := validateFactory(factory); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	id := factory.SystemID()
	if _, duplicate := r.factories[id]; duplicate {
		return fmt.Errorf("system %s is already registered", id)
	}
	r.factories[id] = factory
	return nil
}

func (r *Registry) Prepare(ctx context.Context, store artifact.Store, snapshot space.SnapshotRecord) (Prepared, error) {
	if r == nil {
		return nil, fmt.Errorf("system registry is nil")
	}
	if store == nil {
		return nil, fmt.Errorf("artifact store is required to prepare a system")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.RLock()
	factory, ok := r.factories[snapshot.System]
	r.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("system %s is not registered", snapshot.System)
	}
	if snapshot.Schema != factory.ConfigSchema() {
		return nil, fmt.Errorf("snapshot schema %s does not match system %s config schema %s", snapshot.Schema, snapshot.System, factory.ConfigSchema())
	}
	prepared, err := factory.Prepare(ctx, store, snapshot)
	if err != nil {
		return nil, fmt.Errorf("prepare system %s snapshot %s: %w", snapshot.System, snapshot.ID, err)
	}
	if prepared == nil {
		return nil, fmt.Errorf("system %s returned a nil prepared executable", snapshot.System)
	}
	if prepared.SystemID() != snapshot.System || prepared.SnapshotID() != snapshot.ID || prepared.CaseSchema() != factory.CaseSchema() {
		return nil, fmt.Errorf("prepared system identity does not match registered snapshot and case schema")
	}
	return prepared, nil
}

func (r *Registry) Systems() []record.SystemID {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := make([]record.SystemID, 0, len(r.factories))
	for id := range r.factories {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids
}

func validateFactory(factory Factory) error {
	if err := record.ValidateID("system", string(factory.SystemID())); err != nil {
		return err
	}
	if err := record.ValidateID("schema", string(factory.ConfigSchema())); err != nil {
		return fmt.Errorf("invalid system config schema: %w", err)
	}
	if err := record.ValidateID("schema", string(factory.CaseSchema())); err != nil {
		return fmt.Errorf("invalid system case schema: %w", err)
	}
	return nil
}
