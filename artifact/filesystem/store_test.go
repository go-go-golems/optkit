package filesystem

import (
	"bytes"
	"context"
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/go-go-golems/optkit/artifact"
)

func TestConcurrentIdenticalPut(t *testing.T) {
	store, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	const workers = 16
	refs := make(chan artifact.Ref, workers)
	errs := make(chan error, workers)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ref, err := store.Put(ctx, artifact.PutRequest{MediaType: "application/octet-stream"}, bytes.NewReader([]byte("same bytes")))
			if err != nil {
				errs <- err
				return
			}
			refs <- ref
		}()
	}
	wg.Wait()
	close(refs)
	close(errs)
	for err := range errs {
		t.Fatal(err)
	}
	var first artifact.Ref
	for ref := range refs {
		if first.Digest == "" {
			first = ref
			continue
		}
		if ref.Digest != first.Digest {
			t.Fatalf("digests differ: %s and %s", ref.Digest, first.Digest)
		}
	}
	if err := store.Verify(ctx, first); err != nil {
		t.Fatal(err)
	}
}

func TestCorruptionFailsClosed(t *testing.T) {
	store, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	ref, err := store.Put(ctx, artifact.PutRequest{MediaType: "text/plain"}, bytes.NewReader([]byte("correct")))
	if err != nil {
		t.Fatal(err)
	}
	path, err := store.path(ref.Digest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("corrupt"), 0o444); err != nil {
		t.Fatal(err)
	}
	if err := store.Verify(ctx, ref); !errors.Is(err, artifact.ErrCorrupt) {
		t.Fatalf("Verify error = %v, want ErrCorrupt", err)
	}
	_, err = store.Put(ctx, artifact.PutRequest{MediaType: "text/plain"}, bytes.NewReader([]byte("correct")))
	if err == nil {
		t.Fatal("put unexpectedly replaced corrupt existing content")
	}
}
