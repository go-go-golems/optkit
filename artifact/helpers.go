package artifact

import (
	"bytes"
	"context"
	"fmt"

	"github.com/go-go-golems/optkit/record"
)

func PutBytes(ctx context.Context, store Store, request PutRequest, data []byte) (Ref, error) {
	return store.Put(ctx, request, bytes.NewReader(data))
}

func PutCanonical(ctx context.Context, store Store, schema record.SchemaID, sensitivity Sensitivity, value any) (Ref, error) {
	data, err := record.CanonicalJSON(value)
	if err != nil {
		return Ref{}, err
	}
	ref, err := store.Put(ctx, PutRequest{
		MediaType:   "application/vnd.optkit.record+json",
		Schema:      &schema,
		Sensitivity: sensitivity,
	}, bytes.NewReader(data))
	if err != nil {
		return Ref{}, fmt.Errorf("put canonical artifact: %w", err)
	}
	return ref, nil
}
