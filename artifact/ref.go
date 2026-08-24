package artifact

import "github.com/go-go-golems/optkit/record"

type Sensitivity string

const (
	SensitivityPublic       Sensitivity = "public"
	SensitivityInternal     Sensitivity = "internal"
	SensitivityConfidential Sensitivity = "confidential"
	SensitivityRestricted   Sensitivity = "restricted"
)

type Ref struct {
	Digest      record.Digest    `json:"digest"`
	MediaType   string           `json:"media_type"`
	Schema      *record.SchemaID `json:"schema,omitempty"`
	Size        int64            `json:"size"`
	Sensitivity Sensitivity      `json:"sensitivity"`
}

func (r Ref) Validate() error {
	if err := r.Digest.Validate(); err != nil {
		return err
	}
	if r.MediaType == "" {
		return ErrMissingMediaType
	}
	if r.Size < 0 {
		return ErrInvalidSize
	}
	if r.Schema != nil {
		if err := r.Schema.Validate(); err != nil {
			return err
		}
	}
	return nil
}
