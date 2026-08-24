package space

import (
	"fmt"
	"sort"
)

type DomainDescriptor struct {
	Kind    string   `json:"kind"`
	Minimum *int     `json:"minimum,omitempty"`
	Maximum *int     `json:"maximum,omitempty"`
	Choices []string `json:"choices,omitempty"`
}

type Domain[V any] interface {
	Validate(V) error
	Descriptor() DomainDescriptor
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

func (d IntRangeDomain) Descriptor() DomainDescriptor {
	min, max := d.Minimum, d.Maximum
	return DomainDescriptor{Kind: "integer_range", Minimum: &min, Maximum: &max}
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
	if _, ok := d.labels[v]; !ok {
		return fmt.Errorf("value %v is not in the choice domain", v)
	}
	return nil
}

func (d ChoiceDomain[V]) Descriptor() DomainDescriptor {
	labels := make([]string, 0, len(d.values))
	for _, value := range d.values {
		labels = append(labels, d.labels[value])
	}
	sort.Strings(labels)
	return DomainDescriptor{Kind: "choice", Choices: labels}
}
