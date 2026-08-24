package space

import (
	"context"
	"testing"

	"github.com/go-go-golems/optkit/artifact/memory"
	"github.com/go-go-golems/optkit/record"
)

type testConfig struct {
	Multiplier int    `json:"multiplier"`
	Mode       string `json:"mode"`
}

func multiplierVariable() Variable[testConfig, int] {
	codec := NewJSONCodec[int]("schema:test.int/v1")
	return Variable[testConfig, int]{
		Descriptor: VariableDescriptor{ID: "math.multiplier", Name: "Multiplier", ValueSchema: codec.Schema(), Domain: IntRange(1, 10).Descriptor()},
		Lens: Lens[testConfig, int]{
			Get: func(c testConfig) int { return c.Multiplier },
			Put: func(c testConfig, value int) (testConfig, error) { c.Multiplier = value; return c, nil },
		},
		Domain: IntRange(1, 10), Codec: codec,
	}
}

func modeVariable() Variable[testConfig, string] {
	codec := NewJSONCodec[string]("schema:test.mode/v1")
	domain := Choices(map[string]string{"none": "None", "small": "Small"})
	return Variable[testConfig, string]{
		Descriptor: VariableDescriptor{ID: "noise.mode", Name: "Noise", ValueSchema: codec.Schema(), Domain: domain.Descriptor()},
		Lens: Lens[testConfig, string]{
			Get: func(c testConfig) string { return c.Mode },
			Put: func(c testConfig, value string) (testConfig, error) { c.Mode = value; return c, nil },
		},
		Domain: domain, Codec: codec,
	}
}

func TestLensLaws(t *testing.T) {
	if err := CheckLensLaws(multiplierVariable().Lens, testConfig{Multiplier: 2, Mode: "none"}, 3, 4); err != nil {
		t.Fatal(err)
	}
}

func TestPatchOrderIsCanonicalAndConflictsFail(t *testing.T) {
	ctx := context.Background()
	store := memory.New()
	codec := NewJSONCodec[testConfig]("schema:test.config/v1")
	base, err := MaterializeSnapshot(ctx, store, record.SystemID("system:test/v1"), codec, testConfig{Multiplier: 2, Mode: "none"})
	if err != nil {
		t.Fatal(err)
	}
	first := NewPatchBuilder(base, store, codec)
	if err := Set(first, multiplierVariable(), 3); err != nil {
		t.Fatal(err)
	}
	if err := Set(first, modeVariable(), "small"); err != nil {
		t.Fatal(err)
	}
	p1, child1, err := first.Build(ctx)
	if err != nil {
		t.Fatal(err)
	}

	second := NewPatchBuilder(base, store, codec)
	if err := Set(second, modeVariable(), "small"); err != nil {
		t.Fatal(err)
	}
	if err := Set(second, multiplierVariable(), 3); err != nil {
		t.Fatal(err)
	}
	p2, child2, err := second.Build(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if p1.ID != p2.ID || child1.ID != child2.ID {
		t.Fatalf("assignment order changed identity: %s/%s vs %s/%s", p1.ID, child1.ID, p2.ID, child2.ID)
	}

	conflict := NewPatchBuilder(base, store, codec)
	if err := Set(conflict, multiplierVariable(), 3); err != nil {
		t.Fatal(err)
	}
	if err := Set(conflict, multiplierVariable(), 4); err == nil {
		t.Fatal("conflicting assignment unexpectedly succeeded")
	}
}

func TestPatchRejectsOutOfDomainValue(t *testing.T) {
	ctx := context.Background()
	store := memory.New()
	codec := NewJSONCodec[testConfig]("schema:test.config/v1")
	base, err := MaterializeSnapshot(ctx, store, "system:test/v1", codec, testConfig{Multiplier: 2, Mode: "none"})
	if err != nil {
		t.Fatal(err)
	}
	builder := NewPatchBuilder(base, store, codec)
	if err := Set(builder, multiplierVariable(), 99); err != nil {
		t.Fatal(err)
	}
	if _, _, err := builder.Build(ctx); err == nil {
		t.Fatal("out-of-domain patch unexpectedly built")
	}
}
