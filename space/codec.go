package space

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-go-golems/optkit/record"
)

type Codec[V any] interface {
	EncodeCanonical(V) ([]byte, error)
	Decode([]byte) (V, error)
	Schema() record.SchemaID
}

type JSONCodec[V any] struct {
	SchemaID record.SchemaID
	Strict   bool
}

func NewJSONCodec[V any](schema record.SchemaID) JSONCodec[V] {
	return JSONCodec[V]{SchemaID: schema, Strict: true}
}

func (c JSONCodec[V]) EncodeCanonical(v V) ([]byte, error) {
	if err := c.SchemaID.Validate(); err != nil {
		return nil, err
	}
	return record.CanonicalJSON(v)
}

func (c JSONCodec[V]) Decode(data []byte) (V, error) {
	var out V
	dec := json.NewDecoder(bytes.NewReader(data))
	if c.Strict {
		dec.DisallowUnknownFields()
	}
	if err := dec.Decode(&out); err != nil {
		return out, fmt.Errorf("decode %s: %w", c.SchemaID, err)
	}
	var trailing any
	if err := dec.Decode(&trailing); err != io.EOF {
		if err == nil {
			return out, fmt.Errorf("decode %s: trailing JSON value", c.SchemaID)
		}
		return out, fmt.Errorf("decode %s trailing data: %w", c.SchemaID, err)
	}
	return out, nil
}

func (c JSONCodec[V]) Schema() record.SchemaID { return c.SchemaID }
