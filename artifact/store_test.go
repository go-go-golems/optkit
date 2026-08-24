package artifact_test

import (
	"context"
	"io"
	"testing"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/artifact/filesystem"
	"github.com/go-go-golems/optkit/artifact/memory"
)

func TestStores(t *testing.T) {
	stores := map[string]artifact.Store{
		"memory": memory.New(),
	}
	fs, err := filesystem.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	stores["filesystem"] = fs

	for name, store := range stores {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			ref, err := artifact.PutBytes(ctx, store, artifact.PutRequest{
				MediaType:   "text/plain",
				Sensitivity: artifact.SensitivityInternal,
			}, []byte("immutable evidence"))
			if err != nil {
				t.Fatal(err)
			}
			if err := store.Verify(ctx, ref); err != nil {
				t.Fatal(err)
			}
			r, err := store.Open(ctx, ref)
			if err != nil {
				t.Fatal(err)
			}
			got, err := io.ReadAll(r)
			r.Close()
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != "immutable evidence" {
				t.Fatalf("content = %q", got)
			}
			ref2, err := artifact.PutBytes(ctx, store, artifact.PutRequest{MediaType: "text/plain"}, got)
			if err != nil {
				t.Fatal(err)
			}
			if ref2.Digest != ref.Digest {
				t.Fatalf("idempotent put changed digest: %s != %s", ref2.Digest, ref.Digest)
			}
		})
	}
}
