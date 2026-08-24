package measure

import (
	"fmt"

	"github.com/go-go-golems/optkit/record"
)

type EpochDefinition struct {
	Construct      string `json:"construct"`
	Instrument     string `json:"instrument"`
	Protocol       string `json:"protocol"`
	Implementation string `json:"implementation"`
	Calibration    string `json:"calibration,omitempty"`
	Redaction      string `json:"redaction,omitempty"`
}

type Epoch struct {
	ID         record.EpochID  `json:"id"`
	Definition EpochDefinition `json:"definition"`
}

func NewEpoch(def EpochDefinition) (Epoch, error) {
	if def.Construct == "" || def.Instrument == "" || def.Protocol == "" || def.Implementation == "" {
		return Epoch{}, fmt.Errorf("epoch requires construct, instrument, protocol, and implementation")
	}
	digest, _, err := record.SemanticDigest("schema:optkit.measurement-epoch/v1", def)
	if err != nil {
		return Epoch{}, err
	}
	rawID, err := record.ContentID("epoch", digest)
	if err != nil {
		return Epoch{}, err
	}
	return Epoch{ID: record.EpochID(rawID), Definition: def}, nil
}
