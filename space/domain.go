package space

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/record"
)

// Domain validates typed values and produces the lossless serialized legality
// specification used by catalogs. Descriptor values are encoded by the same
// codec used for assignments.
type Domain[V any] interface {
	Validate(V) error
	Descriptor(Codec[V]) (ValueSpec, error)
}

type IntRangeDomain struct {
	Minimum int
	Maximum int
}

func IntRange(minimum, maximum int) IntRangeDomain {
	return IntRangeDomain{Minimum: minimum, Maximum: maximum}
}

func (d IntRangeDomain) Validate(v int) error {
	if d.Minimum > d.Maximum {
		return fmt.Errorf("invalid integer domain [%d,%d]", d.Minimum, d.Maximum)
	}
	if v < d.Minimum || v > d.Maximum {
		return fmt.Errorf("value %d is outside [%d,%d]", v, d.Minimum, d.Maximum)
	}
	return nil
}

func (d IntRangeDomain) Descriptor(Codec[int]) (ValueSpec, error) {
	spec := ValueSpec{Kind: ValueInt, IntegerRange: &IntegerRangeSpec{Minimum: int64(d.Minimum), Maximum: int64(d.Maximum)}}
	return spec, spec.Validate()
}

type FloatRangeDomain struct {
	Minimum float64
	Maximum float64
}

func FloatRange(minimum, maximum float64) FloatRangeDomain {
	return FloatRangeDomain{Minimum: minimum, Maximum: maximum}
}

func (d FloatRangeDomain) Validate(v float64) error {
	if !finite(d.Minimum) || !finite(d.Maximum) || d.Minimum > d.Maximum {
		return fmt.Errorf("invalid finite float domain [%g,%g]", d.Minimum, d.Maximum)
	}
	if !finite(v) {
		return fmt.Errorf("float value must be finite")
	}
	if v < d.Minimum || v > d.Maximum {
		return fmt.Errorf("value %g is outside [%g,%g]", v, d.Minimum, d.Maximum)
	}
	return nil
}

func (d FloatRangeDomain) Descriptor(Codec[float64]) (ValueSpec, error) {
	spec := ValueSpec{Kind: ValueFloat, FloatRange: &FloatRangeSpec{Minimum: d.Minimum, Maximum: d.Maximum}}
	return spec, spec.Validate()
}

type BoolDomain struct{}

func Boolean() BoolDomain { return BoolDomain{} }

func (BoolDomain) Validate(bool) error { return nil }

func (BoolDomain) Descriptor(Codec[bool]) (ValueSpec, error) {
	spec := ValueSpec{Kind: ValueBool}
	return spec, spec.Validate()
}

type StringDomain struct {
	Pattern string
}

func StringMatching(pattern string) StringDomain { return StringDomain{Pattern: pattern} }

func (d StringDomain) Validate(value string) error {
	if d.Pattern == "" {
		return nil
	}
	pattern, err := regexp.Compile(d.Pattern)
	if err != nil {
		return fmt.Errorf("invalid string domain pattern: %w", err)
	}
	if !pattern.MatchString(value) {
		return fmt.Errorf("value %q does not match %q", value, d.Pattern)
	}
	return nil
}

func (d StringDomain) Descriptor(Codec[string]) (ValueSpec, error) {
	spec := ValueSpec{Kind: ValueString, String: &StringSpec{Pattern: d.Pattern}}
	return spec, spec.Validate()
}

type ChoiceDomain[V comparable] struct {
	values []V
	labels map[V]string
}

func Choices[V comparable](values map[V]string) ChoiceDomain[V] {
	copyLabels := make(map[V]string, len(values))
	keys := make([]V, 0, len(values))
	for value, label := range values {
		copyLabels[value] = label
		keys = append(keys, value)
	}
	return ChoiceDomain[V]{values: keys, labels: copyLabels}
}

func (d ChoiceDomain[V]) Validate(v V) error {
	if len(d.labels) == 0 {
		return fmt.Errorf("choice domain requires at least one value")
	}
	if _, ok := d.labels[v]; !ok {
		return fmt.Errorf("value %v is not in the choice domain", v)
	}
	return nil
}

func (d ChoiceDomain[V]) Descriptor(codec Codec[V]) (ValueSpec, error) {
	choices := make([]ChoiceSpec, 0, len(d.values))
	for _, value := range d.values {
		canonical, err := codec.EncodeCanonical(value)
		if err != nil {
			return ValueSpec{}, fmt.Errorf("encode choice value: %w", err)
		}
		choices = append(choices, ChoiceSpec{Value: append([]byte(nil), canonical...), Label: d.labels[value]})
	}
	sort.Slice(choices, func(i, j int) bool {
		return bytes.Compare(choices[i].Value, choices[j].Value) < 0
	})
	spec := ValueSpec{Kind: ValueChoice, Choices: choices}
	return spec, spec.Validate()
}

type ArtifactRefDomain struct {
	Schema record.SchemaID
}

func Artifacts(schema record.SchemaID) ArtifactRefDomain {
	return ArtifactRefDomain{Schema: schema}
}

func (d ArtifactRefDomain) Validate(ref artifact.Ref) error {
	if err := d.Schema.Validate(); err != nil {
		return fmt.Errorf("artifact domain schema: %w", err)
	}
	if err := ref.Validate(); err != nil {
		return fmt.Errorf("artifact ref: %w", err)
	}
	if ref.Schema == nil {
		return fmt.Errorf("artifact ref requires schema %s", d.Schema)
	}
	if *ref.Schema != d.Schema {
		return fmt.Errorf("artifact schema %s does not match required %s", *ref.Schema, d.Schema)
	}
	return nil
}

func (d ArtifactRefDomain) Descriptor(Codec[artifact.Ref]) (ValueSpec, error) {
	spec := ValueSpec{Kind: ValueArtifact, Artifact: &ArtifactSpec{Schema: d.Schema}}
	return spec, spec.Validate()
}

var _ Domain[int] = IntRangeDomain{}
var _ Domain[float64] = FloatRangeDomain{}
var _ Domain[bool] = BoolDomain{}
var _ Domain[string] = StringDomain{}
var _ Domain[artifact.Ref] = ArtifactRefDomain{}
