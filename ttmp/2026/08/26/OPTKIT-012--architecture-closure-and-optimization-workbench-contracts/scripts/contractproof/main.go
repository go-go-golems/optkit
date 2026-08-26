// Command contractproof is an isolated compile-and-runtime proof for OPTKIT-012.
// It demonstrates that a typed Variable[C,V] can be erased behind Binding[C]
// while preserving one codec/domain/lens path for pure drafts and durable seals.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-go-golems/optkit/artifact/memory"
	"github.com/go-go-golems/optkit/record"
	"github.com/go-go-golems/optkit/space"
)

type Config struct {
	Count int `json:"count"`
}

type Binding[C any] interface {
	ID() space.VariableID
	Normalize(json.RawMessage) (json.RawMessage, error)
	ApplyPure(C, json.RawMessage) (C, error)
	Assign(*space.PatchBuilder[C], json.RawMessage) error
}

type typedBinding[C, V any] struct {
	variable space.Variable[C, V]
}

func (b typedBinding[C, V]) ID() space.VariableID { return b.variable.Descriptor.ID }

func (b typedBinding[C, V]) decode(raw json.RawMessage) (V, json.RawMessage, error) {
	value, err := b.variable.Codec.Decode(raw)
	if err != nil {
		var zero V
		return zero, nil, err
	}
	if err := b.variable.Domain.Validate(value); err != nil {
		var zero V
		return zero, nil, err
	}
	canonical, err := b.variable.Codec.EncodeCanonical(value)
	if err != nil {
		var zero V
		return zero, nil, err
	}
	return value, canonical, nil
}

func (b typedBinding[C, V]) Normalize(raw json.RawMessage) (json.RawMessage, error) {
	_, canonical, err := b.decode(raw)
	return canonical, err
}

func (b typedBinding[C, V]) ApplyPure(config C, raw json.RawMessage) (C, error) {
	value, _, err := b.decode(raw)
	if err != nil {
		return config, err
	}
	return b.variable.Lens.Put(config, value)
}

func (b typedBinding[C, V]) Assign(builder *space.PatchBuilder[C], raw json.RawMessage) error {
	value, _, err := b.decode(raw)
	if err != nil {
		return err
	}
	return space.Set(builder, b.variable, value)
}

func main() {
	ctx := context.Background()
	codec := space.NewJSONCodec[Config](record.SchemaID("schema:optkit.contract-proof-config/v1"))
	intCodec := space.NewJSONCodec[int](record.SchemaID("schema:optkit.contract-proof-count/v1"))
	variable := space.Variable[Config, int]{
		Descriptor: space.VariableDescriptor{
			ID:          "proof.count",
			Name:        "Count",
			Description: "Contract proof count.",
			ValueSchema: intCodec.Schema(),
			Domain:      space.IntRange(1, 10).Descriptor(),
		},
		Lens: space.Lens[Config, int]{
			Get: func(config Config) int { return config.Count },
			Put: func(config Config, value int) (Config, error) {
				config.Count = value
				return config, nil
			},
		},
		Domain: space.IntRange(1, 10),
		Codec:  intCodec,
	}
	if err := variable.Validate(); err != nil {
		panic(err)
	}

	var binding Binding[Config] = typedBinding[Config, int]{variable: variable}
	raw := json.RawMessage(" 3 ")
	normalized, err := binding.Normalize(raw)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(normalized, []byte("3")) {
		panic(fmt.Sprintf("normalization mismatch: %q", normalized))
	}
	pure, err := binding.ApplyPure(Config{Count: 2}, raw)
	if err != nil {
		panic(err)
	}

	store := memory.New()
	base, err := space.MaterializeSnapshot(ctx, store, record.SystemID("system:optkit-contract-proof/v1"), codec, Config{Count: 2})
	if err != nil {
		panic(err)
	}
	builder := space.NewPatchBuilder(base, store, codec)
	if err := binding.Assign(builder, raw); err != nil {
		panic(err)
	}
	_, child, err := builder.Build(ctx)
	if err != nil {
		panic(err)
	}
	if pure != child.Value {
		panic(fmt.Sprintf("pure/durable mismatch: pure=%+v durable=%+v", pure, child.Value))
	}
	fmt.Printf("contract proof OK: binding=%s normalized=%s child=%s\n", binding.ID(), normalized, child.ID)
}
