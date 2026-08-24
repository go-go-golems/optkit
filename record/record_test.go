package record

import (
	"testing"
)

func TestCanonicalJSONSortsMapKeys(t *testing.T) {
	a, err := CanonicalJSON(map[string]any{"z": 1, "a": 2})
	if err != nil {
		t.Fatal(err)
	}
	b, err := CanonicalJSON(map[string]any{"a": 2, "z": 1})
	if err != nil {
		t.Fatal(err)
	}
	if string(a) != string(b) {
		t.Fatalf("canonical bytes differ:\n%s\n%s", a, b)
	}
	if got, want := string(a), `{"a":2,"z":1}`; got != want {
		t.Fatalf("canonical JSON = %s, want %s", got, want)
	}
}

func TestSemanticDigestIncludesSchema(t *testing.T) {
	v := map[string]any{"value": 7}
	d1, _, err := SemanticDigest("schema:test.one/v1", v)
	if err != nil {
		t.Fatal(err)
	}
	d2, _, err := SemanticDigest("schema:test.two/v1", v)
	if err != nil {
		t.Fatal(err)
	}
	if d1 == d2 {
		t.Fatal("different schema IDs produced equal semantic digest")
	}
}

func TestSchemaRegistryRejectsConflictingRegistration(t *testing.T) {
	r := NewSchemaRegistry()
	if err := r.Register(Schema{ID: "schema:test/v1", Description: "first"}); err != nil {
		t.Fatal(err)
	}
	if err := r.Register(Schema{ID: "schema:test/v1", Description: "first"}); err != nil {
		t.Fatalf("idempotent registration failed: %v", err)
	}
	if err := r.Register(Schema{ID: "schema:test/v1", Description: "changed"}); err == nil {
		t.Fatal("conflicting registration unexpectedly succeeded")
	}
}
