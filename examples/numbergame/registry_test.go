package numbergame

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-go-golems/optkit/artifact/memory"
	"github.com/go-go-golems/optkit/record"
	"github.com/go-go-golems/optkit/space"
)

func TestRegistrySerializedMutationMatchesDirectTypedPatch(t *testing.T) {
	registry, err := Registry()
	if err != nil {
		t.Fatal(err)
	}
	sections := registry.Catalog().Sections()
	if len(sections) != 2 || sections[0].ID != "math" || sections[1].ID != "noise" {
		t.Fatalf("unexpected ordered catalog: %+v", sections)
	}
	binding, err := registry.Binding("math.multiplier")
	if err != nil {
		t.Fatal(err)
	}
	baseConfig := Config{Multiplier: 2, Noise: NoiseNone}
	preview, err := binding.ApplyPure(baseConfig, json.RawMessage(" 3 "))
	if err != nil {
		t.Fatal(err)
	}
	if preview != (Config{Multiplier: 3, Noise: NoiseNone}) {
		t.Fatalf("unexpected preview: %+v", preview)
	}

	ctx := context.Background()
	serializedStore := memory.New()
	serializedBase, err := space.MaterializeSnapshot(ctx, serializedStore, SystemID, ConfigCodec(), baseConfig)
	if err != nil {
		t.Fatal(err)
	}
	serializedBuilder := space.NewPatchBuilder(serializedBase, serializedStore, ConfigCodec())
	if err := binding.Assign(serializedBuilder, json.RawMessage("3")); err != nil {
		t.Fatal(err)
	}
	serializedPatch, serializedChild, err := serializedBuilder.Build(ctx)
	if err != nil {
		t.Fatal(err)
	}

	directStore := memory.New()
	directBase, err := space.MaterializeSnapshot(ctx, directStore, SystemID, ConfigCodec(), baseConfig)
	if err != nil {
		t.Fatal(err)
	}
	directBuilder := space.NewPatchBuilder(directBase, directStore, ConfigCodec())
	if err := space.Set(directBuilder, MultiplierVariable(), 3); err != nil {
		t.Fatal(err)
	}
	directPatch, directChild, err := directBuilder.Build(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if preview != serializedChild.Value || serializedChild.ID != directChild.ID || serializedPatch.ID != directPatch.ID {
		t.Fatalf("serialized/direct mismatch: preview=%+v serialized=%s/%s direct=%s/%s", preview, serializedPatch.ID, serializedChild.ID, directPatch.ID, directChild.ID)
	}
}

func TestRegistryCatalogIdentityIsDeterministic(t *testing.T) {
	first, err := Registry()
	if err != nil {
		t.Fatal(err)
	}
	second, err := Registry()
	if err != nil {
		t.Fatal(err)
	}
	if first.Catalog().ID != second.Catalog().ID || first.Catalog().SemanticID != second.Catalog().SemanticID {
		t.Fatal("numbergame catalog identities are nondeterministic")
	}
	if err := first.Catalog().SemanticID.Validate(); err != nil {
		t.Fatal(err)
	}
	if first.Catalog().SemanticID == record.SumBytes(nil) {
		t.Fatal("catalog identity unexpectedly equals empty digest")
	}
}
