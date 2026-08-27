package space

import (
	"bytes"
	"context"
	"fmt"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/record"
)

type SnapshotRecord struct {
	ID     record.SnapshotID `json:"id"`
	System record.SystemID   `json:"system"`
	Schema record.SchemaID   `json:"schema"`
	Config artifact.Ref      `json:"config"`
}

type Snapshot[C any] struct {
	SnapshotRecord
	Value C `json:"-"`
}

func MaterializeSnapshot[C any](ctx context.Context, store artifact.Store, system record.SystemID, codec Codec[C], config C) (Snapshot[C], error) {
	derived, bytesValue, err := deriveSnapshotValue(system, codec, config)
	if err != nil {
		return Snapshot[C]{}, err
	}
	configRef, err := store.Put(ctx, artifact.PutRequest{
		MediaType:      derived.Config.MediaType,
		Schema:         schemaPtr(codec.Schema()),
		Sensitivity:    derived.Config.Sensitivity,
		ExpectedDigest: &derived.Config.Digest,
	}, bytes.NewReader(bytesValue))
	if err != nil {
		return Snapshot[C]{}, fmt.Errorf("store snapshot config: %w", err)
	}
	if configRef.Digest != derived.Config.Digest || configRef.Size != derived.Config.Size ||
		configRef.MediaType != derived.Config.MediaType || configRef.Sensitivity != derived.Config.Sensitivity ||
		configRef.Schema == nil || *configRef.Schema != codec.Schema() {
		return Snapshot[C]{}, fmt.Errorf("stored snapshot config ref differs from its pure derivation")
	}
	return derived, nil
}

// DeriveSnapshotValue computes exactly the snapshot record that
// MaterializeSnapshot would create without writing its canonical config bytes
// to a store. It is suitable for pure draft compilation; durable callers must
// materialize before publishing the record as stored history.
func DeriveSnapshotValue[C any](system record.SystemID, codec Codec[C], config C) (Snapshot[C], error) {
	derived, _, err := deriveSnapshotValue(system, codec, config)
	return derived, err
}

func deriveSnapshotValue[C any](system record.SystemID, codec Codec[C], config C) (Snapshot[C], []byte, error) {
	if err := record.ValidateID("system", string(system)); err != nil {
		return Snapshot[C]{}, nil, err
	}
	bytesValue, err := codec.EncodeCanonical(config)
	if err != nil {
		return Snapshot[C]{}, nil, fmt.Errorf("encode snapshot config: %w", err)
	}
	configRef := artifact.Ref{
		Digest: record.SumBytes(bytesValue), MediaType: "application/vnd.optkit.record+json",
		Schema: schemaPtr(codec.Schema()), Size: int64(len(bytesValue)), Sensitivity: artifact.SensitivityInternal,
	}
	rawID, err := snapshotContentID(system, codec.Schema(), configRef.Digest)
	if err != nil {
		return Snapshot[C]{}, nil, err
	}
	return Snapshot[C]{
		SnapshotRecord: SnapshotRecord{ID: rawID, System: system, Schema: codec.Schema(), Config: configRef},
		Value:          config,
	}, bytesValue, nil
}

// VerifySnapshotValue verifies a loaded snapshot record against its in-memory
// value without reading or writing an artifact store. Proposal compilation uses
// this to reject forged or stale parent/value pairs while remaining pure.
func VerifySnapshotValue[C any](snapshot Snapshot[C], codec Codec[C]) error {
	if err := record.ValidateID("system", string(snapshot.System)); err != nil {
		return fmt.Errorf("snapshot system: %w", err)
	}
	if err := record.ValidateID("snapshot", string(snapshot.ID)); err != nil {
		return fmt.Errorf("snapshot ID: %w", err)
	}
	if snapshot.Schema != codec.Schema() {
		return fmt.Errorf("snapshot schema %s does not match codec schema %s", snapshot.Schema, codec.Schema())
	}
	if err := snapshot.Config.Validate(); err != nil {
		return fmt.Errorf("snapshot config ref: %w", err)
	}
	if snapshot.Config.MediaType != "application/vnd.optkit.record+json" {
		return fmt.Errorf("snapshot config media type %q is unsupported", snapshot.Config.MediaType)
	}
	if snapshot.Config.Sensitivity != artifact.SensitivityInternal {
		return fmt.Errorf("snapshot config sensitivity %q is unsupported", snapshot.Config.Sensitivity)
	}
	if snapshot.Config.Schema == nil || *snapshot.Config.Schema != codec.Schema() {
		return fmt.Errorf("snapshot config schema does not match codec schema %s", codec.Schema())
	}
	bytesValue, err := codec.EncodeCanonical(snapshot.Value)
	if err != nil {
		return fmt.Errorf("encode snapshot config: %w", err)
	}
	if got := record.SumBytes(bytesValue); got != snapshot.Config.Digest {
		return fmt.Errorf("snapshot config digest mismatch: got %s, want %s", got, snapshot.Config.Digest)
	}
	if int64(len(bytesValue)) != snapshot.Config.Size {
		return fmt.Errorf("snapshot config size mismatch: got %d, want %d", len(bytesValue), snapshot.Config.Size)
	}
	expectedID, err := snapshotContentID(snapshot.System, snapshot.Schema, snapshot.Config.Digest)
	if err != nil {
		return err
	}
	if snapshot.ID != expectedID {
		return fmt.Errorf("snapshot identity mismatch: got %s, want %s", snapshot.ID, expectedID)
	}
	return nil
}

func snapshotContentID(system record.SystemID, schema record.SchemaID, config record.Digest) (record.SnapshotID, error) {
	identity := struct {
		System record.SystemID `json:"system"`
		Schema record.SchemaID `json:"schema"`
		Config record.Digest   `json:"config_digest"`
	}{System: system, Schema: schema, Config: config}
	digest, _, err := record.SemanticDigest("schema:optkit.snapshot-identity/v1", identity)
	if err != nil {
		return "", err
	}
	rawID, err := record.ContentID("snapshot", digest)
	if err != nil {
		return "", err
	}
	return record.SnapshotID(rawID), nil
}

func LoadSnapshot[C any](ctx context.Context, store artifact.Store, rec SnapshotRecord, codec Codec[C]) (Snapshot[C], error) {
	if rec.Schema != codec.Schema() {
		return Snapshot[C]{}, fmt.Errorf("snapshot schema %s does not match codec schema %s", rec.Schema, codec.Schema())
	}
	data, err := artifact.ReadAll(ctx, store, rec.Config)
	if err != nil {
		return Snapshot[C]{}, err
	}
	config, err := codec.Decode(data)
	if err != nil {
		return Snapshot[C]{}, err
	}
	materialized, err := MaterializeSnapshot(ctx, store, rec.System, codec, config)
	if err != nil {
		return Snapshot[C]{}, err
	}
	if materialized.ID != rec.ID || materialized.Config.Digest != rec.Config.Digest {
		return Snapshot[C]{}, fmt.Errorf("snapshot identity verification failed: got %s/%s, want %s/%s", materialized.ID, materialized.Config.Digest, rec.ID, rec.Config.Digest)
	}
	return Snapshot[C]{SnapshotRecord: rec, Value: config}, nil
}

func schemaPtr(id record.SchemaID) *record.SchemaID { return &id }
