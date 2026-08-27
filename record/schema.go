package record

import (
	"fmt"
	"sort"
	"sync"
)

type Schema struct {
	ID          SchemaID `json:"id"`
	Description string   `json:"description"`
}

type SchemaRegistry struct {
	mu      sync.RWMutex
	schemas map[SchemaID]Schema
}

func NewSchemaRegistry() *SchemaRegistry {
	return &SchemaRegistry{schemas: make(map[SchemaID]Schema)}
}

func (r *SchemaRegistry) Register(schema Schema) error {
	if err := schema.ID.Validate(); err != nil {
		return err
	}
	if schema.Description == "" {
		return fmt.Errorf("schema %s requires a description", schema.ID)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if existing, ok := r.schemas[schema.ID]; ok {
		if existing == schema {
			return nil
		}
		return fmt.Errorf("schema %s already registered with different definition", schema.ID)
	}
	r.schemas[schema.ID] = schema
	return nil
}

func (r *SchemaRegistry) Lookup(id SchemaID) (Schema, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	schema, ok := r.schemas[id]
	return schema, ok
}

func (r *SchemaRegistry) List() []Schema {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Schema, 0, len(r.schemas))
	for _, schema := range r.schemas {
		out = append(out, schema)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
