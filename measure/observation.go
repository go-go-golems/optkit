package measure

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/record"
)

type Role string

const (
	RoleFact         Role = "fact"
	RoleIntervention Role = "intervention"
	RoleConstraint   Role = "constraint"
	RoleMeasurement  Role = "measurement"
)

type Status string

const (
	StatusMeasured     Status = "measured"
	StatusSatisfied    Status = "satisfied"
	StatusViolated     Status = "violated"
	StatusInapplicable Status = "inapplicable"
	StatusUnknown      Status = "unknown"
	StatusFailed       Status = "failed"
)

type SubjectRef struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
}

type Value struct {
	Kind   string  `json:"kind"`
	Number *string `json:"number,omitempty"`
	Bool   *bool   `json:"bool,omitempty"`
	Text   *string `json:"text,omitempty"`
}

func NumberValue(value string) Value { return Value{Kind: "decimal", Number: &value} }
func BoolValue(value bool) Value     { return Value{Kind: "boolean", Bool: &value} }
func TextValue(value string) Value   { return Value{Kind: "text", Text: &value} }

type Draft struct {
	Role        Role
	Subject     SubjectRef
	Construct   string
	Instrument  string
	Protocol    string
	Epoch       Epoch
	Status      Status
	Value       Value
	Evidence    []artifact.Ref
	Diagnostics map[string]any
	Repeat      int
}

type Observation struct {
	ID          record.ObservationID `json:"id"`
	Role        Role                 `json:"role"`
	Subject     SubjectRef           `json:"subject"`
	Construct   string               `json:"construct"`
	Instrument  string               `json:"instrument"`
	Protocol    string               `json:"protocol"`
	Epoch       record.EpochID       `json:"epoch"`
	Status      Status               `json:"status"`
	Value       Value                `json:"value"`
	Evidence    []artifact.Ref       `json:"evidence,omitempty"`
	Diagnostics json.RawMessage      `json:"diagnostics,omitempty"`
	Repeat      int                  `json:"repeat"`
	CreatedAt   time.Time            `json:"created_at"`
}

func NewObservation(draft Draft, createdAt time.Time) (Observation, error) {
	if draft.Construct == "" || draft.Instrument == "" || draft.Protocol == "" {
		return Observation{}, fmt.Errorf("observation construct, instrument, and protocol are required")
	}
	if draft.Epoch.Definition.Construct != draft.Construct || draft.Epoch.Definition.Instrument != draft.Instrument || draft.Epoch.Definition.Protocol != draft.Protocol {
		return Observation{}, fmt.Errorf("observation fields do not match epoch definition")
	}
	diagnostics, err := record.CanonicalJSON(draft.Diagnostics)
	if err != nil {
		return Observation{}, err
	}
	identity := struct {
		Role        Role            `json:"role"`
		Subject     SubjectRef      `json:"subject"`
		Construct   string          `json:"construct"`
		Epoch       record.EpochID  `json:"epoch"`
		Status      Status          `json:"status"`
		Value       Value           `json:"value"`
		Evidence    []record.Digest `json:"evidence,omitempty"`
		Diagnostics json.RawMessage `json:"diagnostics,omitempty"`
		Repeat      int             `json:"repeat"`
	}{
		Role:        draft.Role,
		Subject:     draft.Subject,
		Construct:   draft.Construct,
		Epoch:       draft.Epoch.ID,
		Status:      draft.Status,
		Value:       draft.Value,
		Diagnostics: diagnostics,
		Repeat:      draft.Repeat,
	}
	for _, ref := range draft.Evidence {
		identity.Evidence = append(identity.Evidence, ref.Digest)
	}
	digest, _, err := record.SemanticDigest("schema:optkit.observation-identity/v1", identity)
	if err != nil {
		return Observation{}, err
	}
	rawID, err := record.ContentID("observation", digest)
	if err != nil {
		return Observation{}, err
	}
	return Observation{
		ID:          record.ObservationID(rawID),
		Role:        draft.Role,
		Subject:     draft.Subject,
		Construct:   draft.Construct,
		Instrument:  draft.Instrument,
		Protocol:    draft.Protocol,
		Epoch:       draft.Epoch.ID,
		Status:      draft.Status,
		Value:       draft.Value,
		Evidence:    append([]artifact.Ref(nil), draft.Evidence...),
		Diagnostics: diagnostics,
		Repeat:      draft.Repeat,
		CreatedAt:   createdAt.UTC(),
	}, nil
}
