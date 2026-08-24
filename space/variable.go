package space

import (
	"fmt"
	"regexp"

	"github.com/go-go-golems/optkit/record"
)

type VariableID string

var variablePattern = regexp.MustCompile(`^[a-z][a-z0-9_.-]*$`)

func (id VariableID) Validate() error {
	if !variablePattern.MatchString(string(id)) {
		return fmt.Errorf("invalid variable ID %q", id)
	}
	return nil
}

type VariableDescriptor struct {
	ID          VariableID       `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	ValueSchema record.SchemaID  `json:"value_schema"`
	Domain      DomainDescriptor `json:"domain"`
	Sensitive   bool             `json:"sensitive"`
	Probes      []string         `json:"probes,omitempty"`
}

type Variable[C, V any] struct {
	Descriptor VariableDescriptor
	Lens       Lens[C, V]
	Domain     Domain[V]
	Codec      Codec[V]
}

func (v Variable[C, V]) Validate() error {
	if err := v.Descriptor.ID.Validate(); err != nil {
		return err
	}
	if v.Descriptor.Name == "" {
		return fmt.Errorf("variable %s requires a name", v.Descriptor.ID)
	}
	if err := v.Lens.Validate(); err != nil {
		return fmt.Errorf("variable %s: %w", v.Descriptor.ID, err)
	}
	if v.Domain == nil {
		return fmt.Errorf("variable %s has nil domain", v.Descriptor.ID)
	}
	if v.Codec == nil {
		return fmt.Errorf("variable %s has nil codec", v.Descriptor.ID)
	}
	if v.Descriptor.ValueSchema != v.Codec.Schema() {
		return fmt.Errorf("variable %s descriptor schema %s differs from codec schema %s", v.Descriptor.ID, v.Descriptor.ValueSchema, v.Codec.Schema())
	}
	return nil
}
