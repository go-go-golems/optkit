package filesystem

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/record"
)

type Store struct {
	root string
}

func New(root string) (*Store, error) {
	if root == "" {
		return nil, errors.New("artifact root is required")
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve artifact root: %w", err)
	}
	for _, dir := range []string{filepath.Join(root, "sha256"), filepath.Join(root, "tmp")} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create artifact directory %s: %w", dir, err)
		}
	}
	return &Store{root: root}, nil
}

func (s *Store) Root() string { return s.root }

func (s *Store) Put(ctx context.Context, request artifact.PutRequest, r io.Reader) (artifact.Ref, error) {
	if err := ctx.Err(); err != nil {
		return artifact.Ref{}, err
	}
	if request.MediaType == "" {
		return artifact.Ref{}, artifact.ErrMissingMediaType
	}

	tmp, err := os.CreateTemp(filepath.Join(s.root, "tmp"), "put-*")
	if err != nil {
		return artifact.Ref{}, fmt.Errorf("create artifact temp file: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	digest, size, copyErr := copyAndHash(ctx, tmp, r)
	syncErr := tmp.Sync()
	closeErr := tmp.Close()
	if copyErr != nil {
		return artifact.Ref{}, copyErr
	}
	if syncErr != nil {
		return artifact.Ref{}, fmt.Errorf("sync artifact temp file: %w", syncErr)
	}
	if closeErr != nil {
		return artifact.Ref{}, fmt.Errorf("close artifact temp file: %w", closeErr)
	}
	if request.ExpectedDigest != nil && digest != *request.ExpectedDigest {
		return artifact.Ref{}, fmt.Errorf("%w: expected %s, got %s", artifact.ErrDigestMismatch, *request.ExpectedDigest, digest)
	}

	ref := artifact.Ref{
		Digest:      digest,
		MediaType:   request.MediaType,
		Schema:      request.Schema,
		Size:        size,
		Sensitivity: request.Sensitivity,
	}
	target, err := s.path(ref.Digest)
	if err != nil {
		return artifact.Ref{}, err
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return artifact.Ref{}, fmt.Errorf("create artifact shard directory: %w", err)
	}

	// Link is atomic and fails if a concurrent writer already installed the
	// content. Unlike Rename on Unix, it cannot replace a corrupt incumbent.
	if err := os.Link(tmpPath, target); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return artifact.Ref{}, fmt.Errorf("publish artifact: %w", err)
		}
		if err := s.Verify(ctx, ref); err != nil {
			return artifact.Ref{}, fmt.Errorf("existing artifact failed verification: %w", err)
		}
		return ref, nil
	}
	if err := os.Chmod(target, 0o444); err != nil {
		return artifact.Ref{}, fmt.Errorf("make artifact immutable: %w", err)
	}
	if err := syncDir(filepath.Dir(target)); err != nil {
		return artifact.Ref{}, err
	}
	return ref, nil
}

func (s *Store) Open(ctx context.Context, ref artifact.Ref) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := ref.Validate(); err != nil {
		return nil, err
	}
	path, err := s.path(ref.Digest)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, artifact.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("open artifact: %w", err)
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, fmt.Errorf("stat artifact: %w", err)
	}
	if info.Size() != ref.Size {
		f.Close()
		return nil, fmt.Errorf("%w: size %d, expected %d", artifact.ErrCorrupt, info.Size(), ref.Size)
	}
	return f, nil
}

func (s *Store) Stat(ctx context.Context, ref artifact.Ref) (artifact.Info, error) {
	if err := ctx.Err(); err != nil {
		return artifact.Info{}, err
	}
	path, err := s.path(ref.Digest)
	if err != nil {
		return artifact.Info{}, err
	}
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return artifact.Info{}, artifact.ErrNotFound
	}
	if err != nil {
		return artifact.Info{}, fmt.Errorf("stat artifact: %w", err)
	}
	if info.Size() != ref.Size {
		return artifact.Info{}, fmt.Errorf("%w: size %d, expected %d", artifact.ErrCorrupt, info.Size(), ref.Size)
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

func (s *Store) path(digest record.Digest) (string, error) {
	hexPart, err := digest.Hex()
	if err != nil {
		return "", err
	}
	return filepath.Join(s.root, "sha256", hexPart[:2], hexPart[2:4], hexPart), nil
}

func copyAndHash(ctx context.Context, dst io.Writer, src io.Reader) (record.Digest, int64, error) {
	pr, pw := io.Pipe()
	result := make(chan struct {
		digest record.Digest
		size   int64
		err    error
	}, 1)
	go func() {
		digest, size, err := record.SumReader(pr)
		result <- struct {
			digest record.Digest
			size   int64
			err    error
		}{digest: digest, size: size, err: err}
	}()

	written, err := io.Copy(io.MultiWriter(dst, pw), &contextReader{ctx: ctx, r: src})
	closeErr := pw.CloseWithError(err)
	hashed := <-result
	if err != nil {
		return "", written, fmt.Errorf("write artifact: %w", err)
	}
	if closeErr != nil {
		return "", written, fmt.Errorf("close artifact hash stream: %w", closeErr)
	}
	if hashed.err != nil {
		return "", written, hashed.err
	}
	if hashed.size != written {
		return "", written, fmt.Errorf("hash stream size %d differs from written size %d", hashed.size, written)
	}
	return hashed.digest, written, nil
}

type contextReader struct {
	ctx context.Context
	r   io.Reader
}

func (r *contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.r.Read(p)
}

func syncDir(path string) error {
	dir, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open artifact directory for sync: %w", err)
	}
	defer dir.Close()
	if err := dir.Sync(); err != nil {
		return fmt.Errorf("sync artifact directory: %w", err)
	}
	return nil
}
