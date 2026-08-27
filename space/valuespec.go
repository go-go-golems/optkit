package space

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"regexp"
	"strings"

	"github.com/go-go-golems/optkit/record"
)

// ValueKind identifies the serialized value family accepted by a variable.
// A ValueSpec is a discriminated union: exactly the payload belonging to Kind
// may be present.
type ValueKind string

const (
	ValueInt      ValueKind = "int"
	ValueFloat    ValueKind = "float"
	ValueBool     ValueKind = "bool"
	ValueString   ValueKind = "string"
	ValueChoice   ValueKind = "choice"
	ValueArtifact ValueKind = "artifact_ref"
)

type IntegerRangeSpec struct {
	Minimum int64 `json:"minimum"`
	Maximum int64 `json:"maximum"`
}

type FloatRangeSpec struct {
	Minimum float64 `json:"minimum"`
	Maximum float64 `json:"maximum"`
}

type ChoiceSpec struct {
	Value json.RawMessage `json:"value"`
	Label string          `json:"label"`
}

type StringSpec struct {
	Pattern string `json:"pattern,omitempty"`
}

type ArtifactSpec struct {
	Schema record.SchemaID `json:"schema"`
}

type ValueSpec struct {
	Kind         ValueKind         `json:"kind"`
	IntegerRange *IntegerRangeSpec `json:"integer_range,omitempty"`
	FloatRange   *FloatRangeSpec   `json:"float_range,omitempty"`
	Choices      []ChoiceSpec      `json:"choices,omitempty"`
	String       *StringSpec       `json:"string,omitempty"`
	Artifact     *ArtifactSpec     `json:"artifact,omitempty"`
}

func (s ValueSpec) Validate() error {
	payloads := 0
	if s.IntegerRange != nil {
		payloads++
	}
	if s.FloatRange != nil {
		payloads++
	}
	if len(s.Choices) > 0 {
		payloads++
	}
	if s.String != nil {
		payloads++
	}
	if s.Artifact != nil {
		payloads++
	}

	switch s.Kind {
	case ValueInt:
		if payloads != 1 || s.IntegerRange == nil {
			return fmt.Errorf("int value spec requires only integer_range")
		}
		if s.IntegerRange.Minimum > s.IntegerRange.Maximum {
			return fmt.Errorf("invalid integer range [%d,%d]", s.IntegerRange.Minimum, s.IntegerRange.Maximum)
		}
	case ValueFloat:
		if payloads != 1 || s.FloatRange == nil {
			return fmt.Errorf("float value spec requires only float_range")
		}
		if !finite(s.FloatRange.Minimum) || !finite(s.FloatRange.Maximum) {
			return fmt.Errorf("float range bounds must be finite")
		}
		if s.FloatRange.Minimum > s.FloatRange.Maximum {
			return fmt.Errorf("invalid float range [%g,%g]", s.FloatRange.Minimum, s.FloatRange.Maximum)
		}
	case ValueBool:
		if payloads != 0 {
			return fmt.Errorf("bool value spec does not accept a payload")
		}
	case ValueString:
		if payloads != 1 || s.String == nil {
			return fmt.Errorf("string value spec requires only string")
		}
		if s.String.Pattern != "" {
			if _, err := regexp.Compile(s.String.Pattern); err != nil {
				return fmt.Errorf("invalid string pattern: %w", err)
			}
		}
	case ValueChoice:
		if payloads != 1 || len(s.Choices) == 0 {
			return fmt.Errorf("choice value spec requires only non-empty choices")
		}
		seenValues := make(map[string]struct{}, len(s.Choices))
		for index, choice := range s.Choices {
			if strings.TrimSpace(choice.Label) == "" {
				return fmt.Errorf("choice %d requires a label", index)
			}
			canonical, err := canonicalRawJSON(choice.Value)
			if err != nil {
				return fmt.Errorf("choice %d: %w", index, err)
			}
			if !bytes.Equal(choice.Value, canonical) {
				return fmt.Errorf("choice %d value must be canonical JSON", index)
			}
			key := string(canonical)
			if _, exists := seenValues[key]; exists {
				return fmt.Errorf("duplicate choice value %s", key)
			}
			seenValues[key] = struct{}{}
		}
	case ValueArtifact:
		if payloads != 1 || s.Artifact == nil {
			return fmt.Errorf("artifact_ref value spec requires only artifact")
		}
		if err := s.Artifact.Schema.Validate(); err != nil {
			return fmt.Errorf("artifact schema: %w", err)
		}
	default:
		return fmt.Errorf("unknown value kind %q", s.Kind)
	}
	return nil
}

func finite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}

func canonicalRawJSON(raw json.RawMessage) ([]byte, error) {
	var value any
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("decode JSON value: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("trailing JSON value")
		}
		return nil, fmt.Errorf("decode trailing JSON data: %w", err)
	}
	return record.CanonicalJSON(value)
}
