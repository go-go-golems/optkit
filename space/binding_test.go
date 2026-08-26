package space

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/go-go-golems/optkit/artifact/memory"
	"github.com/go-go-golems/optkit/record"
)

func testRegistry(t *testing.T) *Registry[testConfig] {
	t.Helper()
	builder := NewRegistryBuilder[testConfig]()
	if err := builder.AddSection(SectionMetadata{ID: "math", Label: "Math", Short: "Arithmetic controls.", Long: "Controls deterministic arithmetic behavior."}); err != nil {
		t.Fatal(err)
	}
	if err := builder.AddSection(SectionMetadata{ID: "noise", Label: "Noise", Short: "Noise controls.", Long: "Controls deterministic seeded-noise behavior."}); err != nil {
		t.Fatal(err)
	}
	if err := Register(builder, "math", multiplierVariable()); err != nil {
		t.Fatal(err)
	}
	if err := Register(builder, "noise", modeVariable()); err != nil {
		t.Fatal(err)
	}
	registry, err := builder.Build()
	if err != nil {
		t.Fatal(err)
	}
	return registry
}

func TestBindingPureAndDurablePathsAgree(t *testing.T) {
	registry := testRegistry(t)
	binding, err := registry.Binding("math.multiplier")
	if err != nil {
		t.Fatal(err)
	}
	if binding.ID() != "math.multiplier" || binding.Descriptor().ID != binding.ID() {
		t.Fatalf("descriptor/binding mismatch: %+v", binding.Descriptor())
	}
	current, err := binding.ReadCanonical(testConfig{Multiplier: 2, Mode: "none"})
	if err != nil || string(current) != "2" {
		t.Fatalf("read canonical: %q %v", current, err)
	}
	normalized, err := binding.Normalize(json.RawMessage(" 3 "))
	if err != nil || string(normalized) != "3" {
		t.Fatalf("normalize: %q %v", normalized, err)
	}
	pure, err := binding.ApplyPure(testConfig{Multiplier: 2, Mode: "none"}, normalized)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	store := memory.New()
	codec := NewJSONCodec[testConfig]("schema:test.binding-config/v1")
	base, err := MaterializeSnapshot(ctx, store, record.SystemID("system:test-binding/v1"), codec, testConfig{Multiplier: 2, Mode: "none"})
	if err != nil {
		t.Fatal(err)
	}
	builder := NewPatchBuilder(base, store, codec)
	if err := binding.Assign(builder, normalized); err != nil {
		t.Fatal(err)
	}
	_, child, err := builder.Build(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if pure != child.Value {
		t.Fatalf("pure/durable mismatch: %+v != %+v", pure, child.Value)
	}
}

func TestBindingRejectsMalformedAndIllegalValues(t *testing.T) {
	binding, err := testRegistry(t).Binding("math.multiplier")
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		raw  json.RawMessage
		kind BindingErrorKind
	}{
		{raw: json.RawMessage(`"3"`), kind: BindingErrorDecode},
		{raw: json.RawMessage(`3 4`), kind: BindingErrorDecode},
		{raw: json.RawMessage(`99`), kind: BindingErrorDomain},
	} {
		_, err := binding.Normalize(test.raw)
		if err == nil {
			t.Fatalf("invalid value %q accepted", test.raw)
		}
		var bindingErr *BindingError
		if !errors.As(err, &bindingErr) {
			t.Fatalf("invalid value %q returned unclassified error %T: %v", test.raw, err, err)
		}
		if bindingErr.Variable != binding.ID() || bindingErr.Kind != test.kind {
			t.Fatalf("invalid value %q classification = %s/%s, want %s/%s", test.raw, bindingErr.Variable, bindingErr.Kind, binding.ID(), test.kind)
		}
	}
	if _, err := testRegistry(t).Binding("missing.variable"); err == nil {
		t.Fatal("unknown variable accepted")
	}
}

func TestRegistryConstructionRejectsDriftAndDuplicates(t *testing.T) {
	builder := NewRegistryBuilder[testConfig]()
	if err := builder.AddSection(SectionMetadata{ID: "math", Label: "Math", Short: "Arithmetic controls.", Long: "Controls deterministic arithmetic behavior."}); err != nil {
		t.Fatal(err)
	}
	if err := builder.AddSection(SectionMetadata{ID: "math", Label: "Again", Short: "Duplicate.", Long: "Duplicate section metadata."}); err == nil {
		t.Fatal("duplicate section accepted")
	}
	if err := Register(builder, "missing", multiplierVariable()); err == nil {
		t.Fatal("unknown section registration accepted")
	}
	if err := Register(builder, "math", multiplierVariable()); err != nil {
		t.Fatal(err)
	}
	if err := Register(builder, "math", multiplierVariable()); err == nil {
		t.Fatal("duplicate variable accepted")
	}

	other := modeVariable()
	other.Descriptor.ID = "math.other"
	other.Descriptor.Key = "multiplier"
	if err := Register(builder, "math", other); err == nil {
		t.Fatal("duplicate local key accepted")
	}

	drifted := multiplierVariable()
	drifted.Descriptor.Value.IntegerRange.Maximum = 11
	drifted.Descriptor.ID = "math.drifted"
	drifted.Descriptor.Key = "drifted"
	if err := Register(builder, "math", drifted); err == nil {
		t.Fatal("descriptor/domain drift accepted")
	}
}

func TestRegistryCatalogAndDescriptorsAreDetached(t *testing.T) {
	registry := testRegistry(t)
	catalog := registry.Catalog()
	sections := catalog.Sections()
	sections[0].Variables[0].Label = "mutated"
	descriptor, ok := registry.Catalog().Lookup("math.multiplier")
	if !ok || descriptor.Label == "mutated" {
		t.Fatal("catalog mutation escaped detached copy")
	}
	binding, _ := registry.Lookup("math.multiplier")
	first := binding.Descriptor()
	first.Label = "mutated"
	if binding.Descriptor().Label == "mutated" {
		t.Fatal("binding exposed mutable descriptor")
	}
}
