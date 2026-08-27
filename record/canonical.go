package record

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// CanonicalJSON returns the repository's v0 canonical JSON representation.
//
// Go's encoding/json sorts map keys and emits a stable compact representation.
// Schemas are responsible for defining omitted-vs-zero semantics. HTML escaping
// is disabled because it is not semantic and makes prompt/document bytes harder
// to inspect.
func CanonicalJSON(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, fmt.Errorf("encode canonical JSON: %w", err)
	}
	return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), nil
}

func SemanticDigest(schema SchemaID, v any) (Digest, []byte, error) {
	if err := schema.Validate(); err != nil {
		return "", nil, err
	}
	payload, err := CanonicalJSON(v)
	if err != nil {
		return "", nil, err
	}
	digest := DigestParts([]byte(schema), []byte{0}, payload)
	return digest, payload, nil
}
