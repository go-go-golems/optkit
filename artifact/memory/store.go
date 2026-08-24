package memory

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/record"
)

type Store struct {
	mu    sync.RWMutex
	blobs map[record.Digest][]byte
}

func New() *Store {
	return &Store{blobs: make(map[record.Digest][]byte)}
}

func (s *Store) Put(ctx context.Context, request artifact.PutRequest, r io.Reader) (artifact.Ref, error) {
	if err := ctx.Err(); err != nil {
		return artifact.Ref{}, err
	}
	if request.MediaType == "" {
		return artifact.Ref{}, artifact.ErrMissingMediaType
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return artifact.Ref{}, fmt.Errorf("read artifact: %w", err)
	}
	digest := record.SumBytes(data)
	if request.ExpectedDigest != nil && digest != *request.ExpectedDigest {
		return artifact.Ref{}, fmt.Errorf("%w: expected %s, got %s", artifact.ErrDigestMismatch, *request.ExpectedDigest, digest)
	}

	s.mu.Lock()
	if existing, ok := s.blobs[digest]; ok {
		if !bytes.Equal(existing, data) {
			s.mu.Unlock()
			return artifact.Ref{}, fmt.Errorf("%w: digest collision for %s", artifact.ErrCorrupt, digest)
		}
	} else {
		s.blobs[digest] = append([]byte(nil), data...)
	}
	s.mu.Unlock()

	return artifact.Ref{
		Digest:      digest,
		MediaType:   request.MediaType,
		Schema:      request.Schema,
		Size:        int64(len(data)),
		Sensitivity: request.Sensitivity,
	}, nil
}

func (s *Store) Open(ctx context.Context, ref artifact.Ref) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := ref.Validate(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	data, ok := s.blobs[ref.Digest]
	copyData := append([]byte(nil), data...)
	s.mu.RUnlock()
	if !ok {
		return nil, artifact.ErrNotFound
	}
	if int64(len(copyData)) != ref.Size {
		return nil, fmt.Errorf("%w: size %d, expected %d", artifact.ErrCorrupt, len(copyData), ref.Size)
	}
	return io.NopCloser(bytes.NewReader(copyData)), nil
}

func (s *Store) Stat(ctx context.Context, ref artifact.Ref) (artifact.Info, error) {
	if err := ctx.Err(); err != nil {
		return artifact.Info{}, err
	}
	s.mu.RLock()
	data, ok := s.blobs[ref.Digest]
	s.mu.RUnlock()
	if !ok {
		return artifact.Info{}, artifact.ErrNotFound
	}
	if int64(len(data)) != ref.Size {
		return artifact.Info{}, fmt.Errorf("%w: size %d, expected %d", artifact.ErrCorrupt, len(data), ref.Size)
	}
	return artifact.Info{Ref: ref}, nil
}

func (s *Store) Verify(ctx context.Context, ref artifact.Ref) error {
	r, err := s.Open(ctx, ref)
	if err != nil {
		return err
	}
	defer r.Close()
	digest, size, err := record.SumReader(r)
	if err != nil {
		return err
	}
	if digest != ref.Digest || size != ref.Size {
		return fmt.Errorf("%w: got %s/%d, expected %s/%d", artifact.ErrCorrupt, digest, size, ref.Digest, ref.Size)
	}
	return nil
}
