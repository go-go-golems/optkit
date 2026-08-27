package artifact

import (
	"context"
	"errors"
	"io"

	"github.com/go-go-golems/optkit/record"
)

var (
	ErrNotFound         = errors.New("artifact not found")
	ErrCorrupt          = errors.New("artifact corrupt")
	ErrDigestMismatch   = errors.New("artifact digest mismatch")
	ErrMissingMediaType = errors.New("artifact media type is required")
	ErrInvalidSize      = errors.New("artifact size is invalid")
)

type PutRequest struct {
	MediaType      string
	Schema         *record.SchemaID
	Sensitivity    Sensitivity
	ExpectedDigest *record.Digest
}

type Info struct {
	Ref Ref
}

type Store interface {
	Put(context.Context, PutRequest, io.Reader) (Ref, error)
	Open(context.Context, Ref) (io.ReadCloser, error)
	Stat(context.Context, Ref) (Info, error)
	Verify(context.Context, Ref) error
}
