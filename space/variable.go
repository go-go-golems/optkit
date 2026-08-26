package space

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-go-golems/optkit/record"
)

type VariableID string

var (
	variablePattern = regexp.MustCompile(`^[a-z][a-z0-9_.-]*$`)
	keyPattern      = regexp.MustCompile(`^[a-z][a-z0-9_-]*$`)
)

func (id VariableID) Validate() error {
	if !variablePattern.MatchString(string(id)) {
		return fmt.Errorf("invalid variable ID %q", id)
	}
	return nil
}

type CostHint string

const (
	CostLow    CostHint = "low"
	CostMedium CostHint = "medium"
	CostHigh   CostHint = "high"
)

func (c CostHint) Validate() error {
	switch c {
	case "", CostLow, CostMedium, CostHigh:
		return nil
	default:
		return fmt.Errorf("unknown cost hint %q", c)
	}
}

type VariableDescriptor struct {
	ID             VariableID      `json:"id"`
	Key            string          `json:"key"`
	Label          string          `json:"label"`
	Short          string          `json:"short"`
	Long           string          `json:"long"`
	ValueSchema    record.SchemaID `json:"value_schema"`
	Value          ValueSpec       `json:"value"`
	Default        json.RawMessage `json:"default,omitempty"`
	Sensitive      bool            `json:"sensitive"`
	CostHint       CostHint        `json:"cost_hint,omitempty"`
	BindingVersion string          `json:"binding_version"`
	Probes         []string        `json:"probes,omitempty"`
}

type VariableMetadata struct {
	ID             VariableID
	Key            string
	Label          string
	Short          string
	Long           string
	Sensitive      bool
	CostHint       CostHint
	BindingVersion string
	Probes         []string
}

type Variable[C, V any] struct {
	Descriptor VariableDescriptor
	Lens       Lens[C, V]
	Domain     Domain[V]
	Codec      Codec[V]
}

// NewVariable derives serialized legality and default metadata from the same
// typed domain and codec used by the executable variable.
func NewVariable[C, V any](metadata VariableMetadata, lens Lens[C, V], domain Domain[V], codec Codec[V], defaultValue *V) (Variable[C, V], error) {
	if domain == nil {
		return Variable[C, V]{}, fmt.Errorf("variable %s has nil domain", metadata.ID)
	}
	if codec == nil {
		return Variable[C, V]{}, fmt.Errorf("variable %s has nil codec", metadata.ID)
	}
	valueSpec, err := domain.Descriptor(codec)
	if err != nil {
		return Variable[C, V]{}, fmt.Errorf("describe variable %s: %w", metadata.ID, err)
	}
	var defaultRaw json.RawMessage
	if defaultValue != nil {
		if err := domain.Validate(*defaultValue); err != nil {
			return Variable[C, V]{}, fmt.Errorf("variable %s default: %w", metadata.ID, err)
		}
		encoded, err := codec.EncodeCanonical(*defaultValue)
		if err != nil {
			return Variable[C, V]{}, fmt.Errorf("encode variable %s default: %w", metadata.ID, err)
		}
		defaultRaw = append([]byte(nil), encoded...)
	}
	variable := Variable[C, V]{
		Descriptor: VariableDescriptor{
			ID:             metadata.ID,
			Key:            metadata.Key,
			Label:          metadata.Label,
			Short:          metadata.Short,
			Long:           metadata.Long,
			ValueSchema:    codec.Schema(),
			Value:          valueSpec,
			Default:        defaultRaw,
			Sensitive:      metadata.Sensitive,
			CostHint:       metadata.CostHint,
			BindingVersion: metadata.BindingVersion,
			Probes:         append([]string(nil), metadata.Probes...),
		},
		Lens:   lens,
		Domain: domain,
		Codec:  codec,
	}
	if err := variable.Validate(); err != nil {
		return Variable[C, V]{}, err
	}
	return variable, nil
}

func (v Variable[C, V]) Validate() error {
	d := v.Descriptor
	if err := d.ID.Validate(); err != nil {
		return err
	}
	if !keyPattern.MatchString(d.Key) {
		return fmt.Errorf("variable %s has invalid key %q", d.ID, d.Key)
	}
	if strings.TrimSpace(d.Label) == "" {
		return fmt.Errorf("variable %s requires a label", d.ID)
	}
	if strings.TrimSpace(d.Short) == "" {
		return fmt.Errorf("variable %s requires short documentation", d.ID)
	}
	if strings.TrimSpace(d.Long) == "" {
		return fmt.Errorf("variable %s requires long documentation", d.ID)
	}
	if strings.TrimSpace(d.BindingVersion) == "" {
		return fmt.Errorf("variable %s requires a binding version", d.ID)
	}
	if err := d.CostHint.Validate(); err != nil {
		return fmt.Errorf("variable %s: %w", d.ID, err)
	}
	if err := v.Lens.Validate(); err != nil {
		return fmt.Errorf("variable %s: %w", d.ID, err)
	}
	if v.Domain == nil {
		return fmt.Errorf("variable %s has nil domain", d.ID)
	}
	if v.Codec == nil {
		return fmt.Errorf("variable %s has nil codec", d.ID)
	}
	if d.ValueSchema != v.Codec.Schema() {
		return fmt.Errorf("variable %s descriptor schema %s differs from codec schema %s", d.ID, d.ValueSchema, v.Codec.Schema())
	}
	if err := d.Value.Validate(); err != nil {
		return fmt.Errorf("variable %s value spec: %w", d.ID, err)
	}
	derived, err := v.Domain.Descriptor(v.Codec)
	if err != nil {
		return fmt.Errorf("variable %s domain descriptor: %w", d.ID, err)
	}
	actualJSON, err := record.CanonicalJSON(d.Value)
	if err != nil {
		return fmt.Errorf("variable %s value spec: %w", d.ID, err)
	}
	derivedJSON, err := record.CanonicalJSON(derived)
	if err != nil {
		return fmt.Errorf("variable %s derived value spec: %w", d.ID, err)
	}
	if !bytes.Equal(actualJSON, derivedJSON) {
		return fmt.Errorf("variable %s value spec differs from typed domain", d.ID)
	}
	if len(d.Default) > 0 {
		value, err := v.Codec.Decode(d.Default)
		if err != nil {
			return fmt.Errorf("variable %s default: %w", d.ID, err)
		}
		if err := v.Domain.Validate(value); err != nil {
			return fmt.Errorf("variable %s default: %w", d.ID, err)
		}
		canonical, err := v.Codec.EncodeCanonical(value)
		if err != nil {
			return fmt.Errorf("variable %s default: %w", d.ID, err)
		}
		if !bytes.Equal(d.Default, canonical) {
			return fmt.Errorf("variable %s default must be canonical JSON", d.ID)
		}
	}
	return nil
}
