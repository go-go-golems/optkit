package artifact

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

func ReadAll(ctx context.Context, store Store, ref Ref) ([]byte, error) {
	r, err := store.Open(ctx, ref)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read artifact %s: %w", ref.Digest, err)
	}
	return data, nil
}

func DecodeJSON(ctx context.Context, store Store, ref Ref, dst any) error {
	data, err := ReadAll(ctx, store, ref)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, dst); err != nil {
		return fmt.Errorf("decode artifact %s: %w", ref.Digest, err)
	}
	return nil
}
