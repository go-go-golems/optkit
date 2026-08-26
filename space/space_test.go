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
	value := 2
	variable, err := NewVariable(
		VariableMetadata{ID: "math.multiplier", Key: "multiplier", Label: "Multiplier", Short: "Integer multiplier.", Long: "Integer multiplier used by the test configuration.", BindingVersion: "test.multiplier/v1"},
		Lens[testConfig, int]{
			Get: func(c testConfig) int { return c.Multiplier },
			Put: func(c testConfig, value int) (testConfig, error) { c.Multiplier = value; return c, nil },
		},
		IntRange(1, 10), codec, &value,
	)
	if err != nil {
		panic(err)
	}
	return variable
}

func modeVariable() Variable[testConfig, string] {
	codec := NewJSONCodec[string]("schema:test.mode/v1")
	domain := Choices(map[string]string{"none": "None", "small": "Small"})
	value := "none"
	variable, err := NewVariable(
		VariableMetadata{ID: "noise.mode", Key: "mode", Label: "Noise", Short: "Noise choice.", Long: "Noise choice used by the test configuration.", BindingVersion: "test.noise/v1"},
		Lens[testConfig, string]{
			Get: func(c testConfig) string { return c.Mode },
			Put: func(c testConfig, value string) (testConfig, error) { c.Mode = value; return c, nil },
		},
		domain, codec, &value,
	)
	if err != nil {
		panic(err)
	}
	return variable
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

func TestVerifySnapshotValueRejectsRecordValueDrift(t *testing.T) {
	ctx := context.Background()
	store := memory.New()
	codec := NewJSONCodec[testConfig]("schema:test.config/v1")
	base, err := MaterializeSnapshot(ctx, store, "system:test/v1", codec, testConfig{Multiplier: 2, Mode: "none"})
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifySnapshotValue(base, codec); err != nil {
		t.Fatalf("verify materialized snapshot: %v", err)
	}

	drifted := base
	drifted.Value.Multiplier = 3
	if err := VerifySnapshotValue(drifted, codec); err == nil {
		t.Fatal("snapshot value drift unexpectedly verified")
	}

	forged := base
	forged.ID = "snapshot:forged"
	if err := VerifySnapshotValue(forged, codec); err == nil {
		t.Fatal("forged snapshot identity unexpectedly verified")
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
