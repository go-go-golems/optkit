package space

import (
	"encoding/json"
	"fmt"
	"strings"
)

// BindingErrorKind classifies failures at the type-erased binding boundary so
// application services can expose stable diagnostics without parsing messages.
type BindingErrorKind string

const (
	BindingErrorDecode BindingErrorKind = "decode"
	BindingErrorDomain BindingErrorKind = "domain"
	BindingErrorEncode BindingErrorKind = "encode"
	BindingErrorApply  BindingErrorKind = "apply"
)

// BindingError preserves the underlying codec, domain, or lens error while
// naming the variable and operation that failed.
type BindingError struct {
	Variable VariableID
	Kind     BindingErrorKind
	Err      error
}

func (e *BindingError) Error() string {
	if e == nil {
		return "binding error"
	}
	return fmt.Sprintf("variable %s: %s: %v", e.Variable, e.Kind, e.Err)
}

func (e *BindingError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func bindingError(id VariableID, kind BindingErrorKind, err error) error {
	return &BindingError{Variable: id, Kind: kind, Err: err}
}

// Binding erases a variable's value type while retaining its real typed codec,
// domain, and lens. Drafting uses ApplyPure; sealing uses Assign.
type Binding[C any] interface {
	ID() VariableID
	Descriptor() VariableDescriptor
	ReadCanonical(C) (json.RawMessage, error)
	Normalize(json.RawMessage) (json.RawMessage, error)
	ApplyPure(C, json.RawMessage) (C, error)
	Assign(*PatchBuilder[C], json.RawMessage) error
}

type typedBinding[C, V any] struct {
	variable Variable[C, V]
}

func (b typedBinding[C, V]) ID() VariableID { return b.variable.Descriptor.ID }

func (b typedBinding[C, V]) Descriptor() VariableDescriptor {
	return cloneVariableDescriptor(b.variable.Descriptor)
}

func (b typedBinding[C, V]) decode(raw json.RawMessage) (V, json.RawMessage, error) {
	value, err := b.variable.Codec.Decode(raw)
	if err != nil {
		var zero V
		return zero, nil, bindingError(b.ID(), BindingErrorDecode, err)
	}
	if err := b.variable.Domain.Validate(value); err != nil {
		var zero V
		return zero, nil, bindingError(b.ID(), BindingErrorDomain, err)
	}
	canonical, err := b.variable.Codec.EncodeCanonical(value)
	if err != nil {
		var zero V
		return zero, nil, bindingError(b.ID(), BindingErrorEncode, err)
	}
	return value, append([]byte(nil), canonical...), nil
}

func (b typedBinding[C, V]) ReadCanonical(config C) (json.RawMessage, error) {
	value := b.variable.Lens.Get(config)
	if err := b.variable.Domain.Validate(value); err != nil {
		return nil, bindingError(b.ID(), BindingErrorDomain, fmt.Errorf("current value: %w", err))
	}
	canonical, err := b.variable.Codec.EncodeCanonical(value)
	if err != nil {
		return nil, bindingError(b.ID(), BindingErrorEncode, fmt.Errorf("current value: %w", err))
	}
	return append([]byte(nil), canonical...), nil
}

func (b typedBinding[C, V]) Normalize(raw json.RawMessage) (json.RawMessage, error) {
	_, canonical, err := b.decode(raw)
	return canonical, err
}

func (b typedBinding[C, V]) ApplyPure(config C, raw json.RawMessage) (C, error) {
	value, _, err := b.decode(raw)
	if err != nil {
		return config, err
	}
	updated, err := b.variable.Lens.Put(config, value)
	if err != nil {
		return config, bindingError(b.ID(), BindingErrorApply, err)
	}
	return updated, nil
}

func (b typedBinding[C, V]) Assign(builder *PatchBuilder[C], raw json.RawMessage) error {
	value, _, err := b.decode(raw)
	if err != nil {
		return err
	}
	return Set(builder, b.variable, value)
}

type SectionMetadata struct {
	ID    SectionID
	Label string
	Short string
	Long  string
}

func (m SectionMetadata) validate() error {
	if err := m.ID.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(m.Label) == "" || strings.TrimSpace(m.Short) == "" || strings.TrimSpace(m.Long) == "" {
		return fmt.Errorf("section %s requires label, short, and long documentation", m.ID)
	}
	return nil
}

type registrySection struct {
	metadata  SectionMetadata
	variables []VariableDescriptor
	keys      map[string]struct{}
}

type RegistryBuilder[C any] struct {
	sections     []registrySection
	sectionIndex map[SectionID]int
	bindings     map[VariableID]Binding[C]
}

func NewRegistryBuilder[C any]() *RegistryBuilder[C] {
	return &RegistryBuilder[C]{sectionIndex: make(map[SectionID]int), bindings: make(map[VariableID]Binding[C])}
}

func (b *RegistryBuilder[C]) AddSection(metadata SectionMetadata) error {
	if b == nil {
		return fmt.Errorf("registry builder is nil")
	}
	if err := metadata.validate(); err != nil {
		return err
	}
	if _, exists := b.sectionIndex[metadata.ID]; exists {
		return fmt.Errorf("duplicate section %s", metadata.ID)
	}
	b.sectionIndex[metadata.ID] = len(b.sections)
	b.sections = append(b.sections, registrySection{metadata: metadata, keys: make(map[string]struct{})})
	return nil
}

// Register captures V in a private adapter and adds both descriptor and binding
// in one operation. Go does not permit generic methods, so Register is a
// package function rather than a RegistryBuilder method.
func Register[C, V any](builder *RegistryBuilder[C], section SectionID, variable Variable[C, V]) error {
	if builder == nil {
		return fmt.Errorf("registry builder is nil")
	}
	if err := variable.Validate(); err != nil {
		return err
	}
	sectionIndex, exists := builder.sectionIndex[section]
	if !exists {
		return fmt.Errorf("unknown section %s", section)
	}
	if _, exists := builder.bindings[variable.Descriptor.ID]; exists {
		return fmt.Errorf("duplicate variable %s", variable.Descriptor.ID)
	}
	registrySection := &builder.sections[sectionIndex]
	if _, exists := registrySection.keys[variable.Descriptor.Key]; exists {
		return fmt.Errorf("duplicate variable key %q in section %s", variable.Descriptor.Key, section)
	}
	binding := typedBinding[C, V]{variable: variable}
	builder.bindings[variable.Descriptor.ID] = binding
	registrySection.keys[variable.Descriptor.Key] = struct{}{}
	registrySection.variables = append(registrySection.variables, binding.Descriptor())
	return nil
}

type Registry[C any] struct {
	catalog  Catalog
	bindings map[VariableID]Binding[C]
}

func (b *RegistryBuilder[C]) Build() (*Registry[C], error) {
	if b == nil {
		return nil, fmt.Errorf("registry builder is nil")
	}
	sections := make([]Section, len(b.sections))
	for index, registered := range b.sections {
		sections[index] = Section{
			ID: registered.metadata.ID, Label: registered.metadata.Label,
			Short: registered.metadata.Short, Long: registered.metadata.Long,
			Variables: append([]VariableDescriptor(nil), registered.variables...),
		}
	}
	catalog, err := NewCatalog(sections)
	if err != nil {
		return nil, err
	}
	bindings := make(map[VariableID]Binding[C], len(b.bindings))
	for id, binding := range b.bindings {
		bindings[id] = binding
	}
	return &Registry[C]{catalog: catalog, bindings: bindings}, nil
}

func (r *Registry[C]) Catalog() Catalog {
	if r == nil {
		return Catalog{}
	}
	return r.catalog
}

func (r *Registry[C]) Lookup(id VariableID) (Binding[C], bool) {
	if r == nil {
		return nil, false
	}
	binding, ok := r.bindings[id]
	return binding, ok
}

func (r *Registry[C]) Binding(id VariableID) (Binding[C], error) {
	binding, ok := r.Lookup(id)
	if !ok {
		return nil, fmt.Errorf("unknown variable %s", id)
	}
	return binding, nil
}
