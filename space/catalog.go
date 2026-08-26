package space

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-go-golems/optkit/record"
)

const (
	CatalogSchema         record.SchemaID = "schema:optkit.catalog/v1"
	CatalogSemanticSchema record.SchemaID = "schema:optkit.catalog-semantic/v1"
)

// Catalog is immutable after construction. Sections returns detached copies so
// shared registries cannot be mutated by readers.
type Catalog struct {
	Schema     record.SchemaID
	SemanticID record.Digest
	ID         record.Digest
	sections   []Section
}

func NewCatalog(sections []Section) (Catalog, error) {
	if len(sections) == 0 {
		return Catalog{}, fmt.Errorf("catalog requires at least one section")
	}
	copySections := make([]Section, len(sections))
	seenSections := make(map[SectionID]struct{}, len(sections))
	seenVariables := make(map[VariableID]struct{})
	for index, section := range sections {
		if err := section.Validate(); err != nil {
			return Catalog{}, fmt.Errorf("catalog section %d: %w", index, err)
		}
		if _, exists := seenSections[section.ID]; exists {
			return Catalog{}, fmt.Errorf("duplicate section %s", section.ID)
		}
		seenSections[section.ID] = struct{}{}
		for _, variable := range section.Variables {
			if _, exists := seenVariables[variable.ID]; exists {
				return Catalog{}, fmt.Errorf("duplicate variable %s", variable.ID)
			}
			seenVariables[variable.ID] = struct{}{}
		}
		copySections[index] = cloneSection(section)
	}

	semanticID, _, err := record.SemanticDigest(CatalogSemanticSchema, semanticCatalogProjection(copySections))
	if err != nil {
		return Catalog{}, fmt.Errorf("catalog semantic identity: %w", err)
	}
	fullID, _, err := record.SemanticDigest(CatalogSchema, struct {
		Sections []Section `json:"sections"`
	}{Sections: copySections})
	if err != nil {
		return Catalog{}, fmt.Errorf("catalog identity: %w", err)
	}
	return Catalog{Schema: CatalogSchema, SemanticID: semanticID, ID: fullID, sections: copySections}, nil
}

func (c Catalog) Sections() []Section {
	out := make([]Section, len(c.sections))
	for index, section := range c.sections {
		out[index] = cloneSection(section)
	}
	return out
}

func (c Catalog) Lookup(id VariableID) (VariableDescriptor, bool) {
	for _, section := range c.sections {
		for _, variable := range section.Variables {
			if variable.ID == id {
				return cloneVariableDescriptor(variable), true
			}
		}
	}
	return VariableDescriptor{}, false
}

func (c Catalog) Validate() error {
	rebuilt, err := NewCatalog(c.sections)
	if err != nil {
		return err
	}
	if c.Schema != rebuilt.Schema {
		return fmt.Errorf("catalog schema %s does not match %s", c.Schema, rebuilt.Schema)
	}
	if c.SemanticID != rebuilt.SemanticID {
		return fmt.Errorf("catalog semantic identity mismatch: got %s, want %s", c.SemanticID, rebuilt.SemanticID)
	}
	if c.ID != rebuilt.ID {
		return fmt.Errorf("catalog identity mismatch: got %s, want %s", c.ID, rebuilt.ID)
	}
	return nil
}

func (c Catalog) MarshalJSON() ([]byte, error) {
	type wireCatalog struct {
		Schema     record.SchemaID `json:"schema"`
		SemanticID record.Digest   `json:"semantic_id"`
		ID         record.Digest   `json:"id"`
		Sections   []Section       `json:"sections"`
	}
	return json.Marshal(wireCatalog{Schema: c.Schema, SemanticID: c.SemanticID, ID: c.ID, Sections: c.Sections()})
}

func (c *Catalog) UnmarshalJSON(data []byte) error {
	if c == nil {
		return fmt.Errorf("catalog target is nil")
	}
	type wireCatalog struct {
		Schema     record.SchemaID `json:"schema"`
		SemanticID record.Digest   `json:"semantic_id"`
		ID         record.Digest   `json:"id"`
		Sections   []Section       `json:"sections"`
	}
	var wire wireCatalog
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&wire); err != nil {
		return fmt.Errorf("decode catalog: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return fmt.Errorf("decode catalog: trailing JSON value")
		}
		return fmt.Errorf("decode catalog trailing data: %w", err)
	}
	rebuilt, err := NewCatalog(wire.Sections)
	if err != nil {
		return err
	}
	if wire.Schema != rebuilt.Schema || wire.SemanticID != rebuilt.SemanticID || wire.ID != rebuilt.ID {
		return fmt.Errorf("catalog wire identity does not match its sections")
	}
	*c = rebuilt
	return nil
}

type semanticSection struct {
	ID        SectionID          `json:"id"`
	Variables []semanticVariable `json:"variables"`
}

type semanticVariable struct {
	ID             VariableID      `json:"id"`
	Key            string          `json:"key"`
	ValueSchema    record.SchemaID `json:"value_schema"`
	Value          ValueSpec       `json:"value"`
	Default        json.RawMessage `json:"default,omitempty"`
	Sensitive      bool            `json:"sensitive"`
	CostHint       CostHint        `json:"cost_hint,omitempty"`
	BindingVersion string          `json:"binding_version"`
	Probes         []string        `json:"probes,omitempty"`
}

func semanticCatalogProjection(sections []Section) []semanticSection {
	out := make([]semanticSection, len(sections))
	for sectionIndex, section := range sections {
		variables := make([]semanticVariable, len(section.Variables))
		for variableIndex, variable := range section.Variables {
			variables[variableIndex] = semanticVariable{
				ID: variable.ID, Key: variable.Key, ValueSchema: variable.ValueSchema,
				Value: cloneValueSpec(variable.Value), Default: append([]byte(nil), variable.Default...),
				Sensitive: variable.Sensitive, CostHint: variable.CostHint,
				BindingVersion: variable.BindingVersion, Probes: append([]string(nil), variable.Probes...),
			}
		}
		out[sectionIndex] = semanticSection{ID: section.ID, Variables: variables}
	}
	return out
}
