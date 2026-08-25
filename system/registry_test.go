package system

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/go-go-golems/optkit/artifact"
	artifactmemory "github.com/go-go-golems/optkit/artifact/memory"
	"github.com/go-go-golems/optkit/episode"
	"github.com/go-go-golems/optkit/record"
	"github.com/go-go-golems/optkit/space"
)

var (
	fixtureSystem       = record.SystemID("system:fixture/v1")
	fixtureConfigSchema = record.SchemaID("schema:fixture.config/v1")
	fixtureCaseSchema   = record.SchemaID("schema:fixture.case/v1")
	fixtureSnapshot     = record.SnapshotID("snapshot:fixture")
)

type fixtureFactory struct {
	prepared Prepared
	err      error
}

func (fixtureFactory) SystemID() record.SystemID     { return fixtureSystem }
func (fixtureFactory) ConfigSchema() record.SchemaID { return fixtureConfigSchema }
func (fixtureFactory) CaseSchema() record.SchemaID   { return fixtureCaseSchema }
func (f fixtureFactory) Prepare(context.Context, artifact.Store, space.SnapshotRecord) (Prepared, error) {
	return f.prepared, f.err
}

type fixturePrepared struct {
	system   record.SystemID
	snapshot record.SnapshotID
	caseType record.SchemaID
}

func (p fixturePrepared) SystemID() record.SystemID     { return p.system }
func (p fixturePrepared) SnapshotID() record.SnapshotID { return p.snapshot }
func (p fixturePrepared) CaseSchema() record.SchemaID   { return p.caseType }
func (fixturePrepared) Run(context.Context, artifact.Ref, int64, episode.Sink) (episode.RunResult, error) {
	return episode.RunResult{Status: episode.StatusCompleted}, nil
}

func TestRegistryPreparesExactSnapshotIdentity(t *testing.T) {
	registry := NewRegistry()
	prepared := fixturePrepared{system: fixtureSystem, snapshot: fixtureSnapshot, caseType: fixtureCaseSchema}
	requireNoError(t, registry.Register(fixtureFactory{prepared: prepared}))
	requireEqual(t, []record.SystemID{fixtureSystem}, registry.Systems())
	got, err := registry.Prepare(context.Background(), artifactmemory.New(), space.SnapshotRecord{
		ID: fixtureSnapshot, System: fixtureSystem, Schema: fixtureConfigSchema,
	})
	requireNoError(t, err)
	if got.SystemID() != fixtureSystem || got.SnapshotID() != fixtureSnapshot || got.CaseSchema() != fixtureCaseSchema {
		t.Fatalf("prepared identity = %s/%s/%s", got.SystemID(), got.SnapshotID(), got.CaseSchema())
	}
}

func TestRegistryRejectsDuplicateUnknownAndMismatchedPreparation(t *testing.T) {
	registry := NewRegistry()
	prepared := fixturePrepared{system: fixtureSystem, snapshot: fixtureSnapshot, caseType: fixtureCaseSchema}
	requireNoError(t, registry.Register(fixtureFactory{prepared: prepared}))
	if err := registry.Register(fixtureFactory{prepared: prepared}); err == nil {
		t.Fatal("duplicate registration succeeded")
	}
	store := artifactmemory.New()
	_, err := registry.Prepare(context.Background(), store, space.SnapshotRecord{ID: fixtureSnapshot, System: "system:unknown/v1", Schema: fixtureConfigSchema})
	requireErrorContains(t, err, "not registered")
	_, err = registry.Prepare(context.Background(), store, space.SnapshotRecord{ID: fixtureSnapshot, System: fixtureSystem, Schema: "schema:wrong/v1"})
	requireErrorContains(t, err, "does not match")

	bad := NewRegistry()
	requireNoError(t, bad.Register(fixtureFactory{prepared: fixturePrepared{system: fixtureSystem, snapshot: "snapshot:wrong", caseType: fixtureCaseSchema}}))
	_, err = bad.Prepare(context.Background(), store, space.SnapshotRecord{ID: fixtureSnapshot, System: fixtureSystem, Schema: fixtureConfigSchema})
	requireErrorContains(t, err, "identity does not match")
}

func TestRegistryPropagatesCancellationAndPreparationFailure(t *testing.T) {
	registry := NewRegistry()
	requireNoError(t, registry.Register(fixtureFactory{err: errors.New("provider unavailable")}))
	store := artifactmemory.New()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := registry.Prepare(ctx, store, space.SnapshotRecord{ID: fixtureSnapshot, System: fixtureSystem, Schema: fixtureConfigSchema})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled prepare error = %v", err)
	}
	_, err = registry.Prepare(context.Background(), store, space.SnapshotRecord{ID: fixtureSnapshot, System: fixtureSystem, Schema: fixtureConfigSchema})
	requireErrorContains(t, err, "provider unavailable")
}

func requireNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func requireErrorContains(t testing.TB, err error, text string) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), text) {
		t.Fatalf("error = %v, want containing %q", err, text)
	}
}

func requireEqual[T comparable](t testing.TB, want, got []T) {
	t.Helper()
	if len(want) != len(got) {
		t.Fatalf("length = %d, want %d", len(got), len(want))
	}
	for index := range want {
		if want[index] != got[index] {
			t.Fatalf("item %d = %v, want %v", index, got[index], want[index])
		}
	}
}
