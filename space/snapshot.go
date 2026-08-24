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
	if err := record.ValidateID("system", string(system)); err != nil {
		return Snapshot[C]{}, err
	}
	bytesValue, err := codec.EncodeCanonical(config)
	if err != nil {
		return Snapshot[C]{}, fmt.Errorf("encode snapshot config: %w", err)
	}
	configRef, err := store.Put(ctx, artifact.PutRequest{
		MediaType:   "application/vnd.optkit.record+json",
		Schema:      schemaPtr(codec.Schema()),
		Sensitivity: artifact.SensitivityInternal,
	}, bytes.NewReader(bytesValue))
	if err != nil {
		return Snapshot[C]{}, fmt.Errorf("store snapshot config: %w", err)
	}
	identity := struct {
		System record.SystemID `json:"system"`
		Schema record.SchemaID `json:"schema"`
		Config record.Digest   `json:"config_digest"`
	}{System: system, Schema: codec.Schema(), Config: configRef.Digest}
	digest, _, err := record.SemanticDigest("schema:optkit.snapshot-identity/v1", identity)
	if err != nil {
		return Snapshot[C]{}, err
	}
	rawID, err := record.ContentID("snapshot", digest)
	if err != nil {
		return Snapshot[C]{}, err
	}
	return Snapshot[C]{
		SnapshotRecord: SnapshotRecord{
			ID:     record.SnapshotID(rawID),
			System: system,
			Schema: codec.Schema(),
			Config: configRef,
		},
		Value: config,
	}, nil
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
	if materialized.ID != rec.ID || materialized.SnapshotRecord.Config.Digest != rec.Config.Digest {
		return Snapshot[C]{}, fmt.Errorf("snapshot identity verification failed: got %s/%s, want %s/%s", materialized.ID, materialized.SnapshotRecord.Config.Digest, rec.ID, rec.Config.Digest)
	}
	return Snapshot[C]{SnapshotRecord: rec, Value: config}, nil
}

func schemaPtr(id record.SchemaID) *record.SchemaID { return &id }
