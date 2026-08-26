package space

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/record"
)

func TestScalarDomainDescriptorsAndBounds(t *testing.T) {
	intCodec := NewJSONCodec[int]("schema:test.int/v1")
	intDomain := IntRange(1, 3)
	intSpec, err := intDomain.Descriptor(intCodec)
	if err != nil {
		t.Fatal(err)
	}
	if intSpec.Kind != ValueInt || intSpec.IntegerRange.Minimum != 1 || intSpec.IntegerRange.Maximum != 3 {
		t.Fatalf("unexpected integer spec: %+v", intSpec)
	}
	for _, value := range []int{1, 3} {
		if err := intDomain.Validate(value); err != nil {
			t.Fatalf("boundary %d rejected: %v", value, err)
		}
	}
	for _, value := range []int{0, 4} {
		if err := intDomain.Validate(value); err == nil {
			t.Fatalf("outside value %d accepted", value)
		}
	}

	floatCodec := NewJSONCodec[float64]("schema:test.float/v1")
	floatDomain := FloatRange(0.5, 2.5)
	if _, err := floatDomain.Descriptor(floatCodec); err != nil {
		t.Fatal(err)
	}
	for _, value := range []float64{math.NaN(), math.Inf(1), math.Inf(-1), 0.49, 2.51} {
		if err := floatDomain.Validate(value); err == nil {
			t.Fatalf("invalid float %v accepted", value)
		}
	}
	if _, err := FloatRange(math.NaN(), 1).Descriptor(floatCodec); err == nil {
		t.Fatal("non-finite float domain accepted")
	}

	boolSpec, err := Boolean().Descriptor(NewJSONCodec[bool]("schema:test.bool/v1"))
	if err != nil || boolSpec.Kind != ValueBool {
		t.Fatalf("unexpected bool spec/err: %+v %v", boolSpec, err)
	}
	stringDomain := StringMatching(`^[a-z]+$`)
	if err := stringDomain.Validate("lowercase"); err != nil {
		t.Fatal(err)
	}
	if err := stringDomain.Validate("not-lower-1"); err == nil {
		t.Fatal("non-matching string accepted")
	}
}

func TestChoiceDescriptorPreservesCanonicalMachineValues(t *testing.T) {
	domain := Choices(map[string]string{"small": "Small seeded noise", "none": "No noise"})
	spec, err := domain.Descriptor(NewJSONCodec[string]("schema:test.choice/v1"))
	if err != nil {
		t.Fatal(err)
	}
	if spec.Kind != ValueChoice || len(spec.Choices) != 2 {
		t.Fatalf("unexpected choice spec: %+v", spec)
	}
	got := map[string]string{}
	for _, choice := range spec.Choices {
		got[string(choice.Value)] = choice.Label
	}
	if got[`"none"`] != "No noise" || got[`"small"`] != "Small seeded noise" {
		t.Fatalf("machine values or labels lost: %#v", got)
	}
	first, _ := json.Marshal(spec)
	for range 10 {
		repeated, err := domain.Descriptor(NewJSONCodec[string]("schema:test.choice/v1"))
		if err != nil {
			t.Fatal(err)
		}
		raw, _ := json.Marshal(repeated)
		if string(raw) != string(first) {
			t.Fatalf("choice descriptor is nondeterministic: %s != %s", raw, first)
		}
	}
}

func TestArtifactDomainRequiresExactSchema(t *testing.T) {
	required := record.SchemaID("schema:test.prompt/v1")
	other := record.SchemaID("schema:test.other/v1")
	domain := Artifacts(required)
	spec, err := domain.Descriptor(NewJSONCodec[artifact.Ref]("schema:test.artifact-ref/v1"))
	if err != nil {
		t.Fatal(err)
	}
	if spec.Kind != ValueArtifact || spec.Artifact.Schema != required {
		t.Fatalf("unexpected artifact spec: %+v", spec)
	}
	valid := artifact.Ref{Digest: record.SumBytes([]byte("prompt")), MediaType: "text/plain", Schema: &required, Size: 6, Sensitivity: artifact.SensitivityInternal}
	if err := domain.Validate(valid); err != nil {
		t.Fatal(err)
	}
	wrong := valid
	wrong.Schema = &other
	if err := domain.Validate(wrong); err == nil {
		t.Fatal("wrong artifact schema accepted")
	}
	missing := valid
	missing.Schema = nil
	if err := domain.Validate(missing); err == nil {
		t.Fatal("missing artifact schema accepted")
	}
}

func TestNewVariableDerivesAndChecksValueMetadata(t *testing.T) {
	codec := NewJSONCodec[int]("schema:test.derived-int/v1")
	defaultValue := 2
	variable, err := NewVariable(
		VariableMetadata{ID: "test.count", Key: "count", Label: "Count", Short: "A test count.", Long: "A test count with complete documentation.", CostHint: CostLow, BindingVersion: "test.count/v1"},
		Lens[testConfig, int]{Get: func(c testConfig) int { return c.Multiplier }, Put: func(c testConfig, value int) (testConfig, error) { c.Multiplier = value; return c, nil }},
		IntRange(1, 3), codec, &defaultValue,
	)
	if err != nil {
		t.Fatal(err)
	}
	if string(variable.Descriptor.Default) != "2" {
		t.Fatalf("default is not canonical: %q", variable.Descriptor.Default)
	}
	variable.Descriptor.Value.IntegerRange.Maximum = 4
	if err := variable.Validate(); err == nil {
		t.Fatal("descriptor/domain drift accepted")
	}
}

func TestValueSpecRejectsMalformedDiscriminators(t *testing.T) {
	cases := []ValueSpec{
		{Kind: "future"},
		{Kind: ValueInt},
		{Kind: ValueBool, String: &StringSpec{}},
		{Kind: ValueChoice, Choices: []ChoiceSpec{{Value: json.RawMessage(` 1 `), Label: "one"}}},
		{Kind: ValueArtifact, Artifact: &ArtifactSpec{Schema: "bad"}},
	}
	for index, spec := range cases {
		if err := spec.Validate(); err == nil {
			t.Fatalf("malformed spec %d accepted: %+v", index, spec)
		}
	}
}
