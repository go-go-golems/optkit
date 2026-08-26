package space

import (
	"fmt"
	"regexp"
	"strings"
)

type SectionID string

var sectionPattern = regexp.MustCompile(`^[a-z][a-z0-9_.-]*$`)

func (id SectionID) Validate() error {
	if !sectionPattern.MatchString(string(id)) {
		return fmt.Errorf("invalid section ID %q", id)
	}
	return nil
}

type Section struct {
	ID        SectionID            `json:"id"`
	Label     string               `json:"label"`
	Short     string               `json:"short"`
	Long      string               `json:"long"`
	Variables []VariableDescriptor `json:"variables"`
}

func (s Section) Validate() error {
	if err := s.ID.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(s.Label) == "" {
		return fmt.Errorf("section %s requires a label", s.ID)
	}
	if strings.TrimSpace(s.Short) == "" {
		return fmt.Errorf("section %s requires short documentation", s.ID)
	}
	if strings.TrimSpace(s.Long) == "" {
		return fmt.Errorf("section %s requires long documentation", s.ID)
	}
	if len(s.Variables) == 0 {
		return fmt.Errorf("section %s requires at least one variable", s.ID)
	}
	for index, variable := range s.Variables {
		if err := variable.Validate(); err != nil {
			return fmt.Errorf("section %s variable %d: %w", s.ID, index, err)
		}
	}
	return nil
}

func cloneSection(section Section) Section {
	copySection := section
	copySection.Variables = make([]VariableDescriptor, len(section.Variables))
	for index, variable := range section.Variables {
		copyVariable := variable
		copyVariable.Default = append([]byte(nil), variable.Default...)
		copyVariable.Probes = append([]string(nil), variable.Probes...)
		copyVariable.Value = cloneValueSpec(variable.Value)
		copySection.Variables[index] = copyVariable
	}
	return copySection
}

func cloneValueSpec(spec ValueSpec) ValueSpec {
	copySpec := spec
	if spec.IntegerRange != nil {
		value := *spec.IntegerRange
		copySpec.IntegerRange = &value
	}
	if spec.FloatRange != nil {
		value := *spec.FloatRange
		copySpec.FloatRange = &value
	}
	copySpec.Choices = make([]ChoiceSpec, len(spec.Choices))
	for index, choice := range spec.Choices {
		copySpec.Choices[index] = ChoiceSpec{Value: append([]byte(nil), choice.Value...), Label: choice.Label}
	}
	if spec.String != nil {
		value := *spec.String
		copySpec.String = &value
	}
	if spec.Artifact != nil {
		value := *spec.Artifact
		copySpec.Artifact = &value
	}
	return copySpec
}
