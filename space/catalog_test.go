package space

import (
	"bytes"
	"encoding/json"
	"testing"
)

func testSections() []Section {
	return []Section{
		{
			ID: "math", Label: "Math", Short: "Arithmetic controls.", Long: "Controls deterministic arithmetic behavior.",
			Variables: []VariableDescriptor{multiplierVariable().Descriptor},
		},
		{
			ID: "noise", Label: "Noise", Short: "Noise controls.", Long: "Controls deterministic seeded-noise behavior.",
			Variables: []VariableDescriptor{modeVariable().Descriptor},
		},
	}
}

func TestCatalogIdentitySeparatesSemanticsFromDocumentation(t *testing.T) {
	base, err := NewCatalog(testSections())
	if err != nil {
		t.Fatal(err)
	}
	docEdit := testSections()
	docEdit[0].Label = "Arithmetic"
	docEdit[0].Short = "Edited section copy."
	docEdit[0].Long = "Edited long section documentation."
	docEdit[0].Variables[0].Label = "Multiplication factor"
	docEdit[0].Variables[0].Short = "Edited variable copy."
	docEdit[0].Variables[0].Long = "Edited long variable documentation."
	edited, err := NewCatalog(docEdit)
	if err != nil {
		t.Fatal(err)
	}
	if base.SemanticID != edited.SemanticID {
		t.Fatalf("documentation changed semantic identity: %s != %s", base.SemanticID, edited.SemanticID)
	}
	if base.ID == edited.ID {
		t.Fatal("documentation did not change full catalog identity")
	}

	semanticEdit := testSections()
	semanticEdit[0].Variables[0].Value.IntegerRange.Maximum = 9
	changed, err := NewCatalog(semanticEdit)
	if err != nil {
		t.Fatal(err)
	}
	if base.SemanticID == changed.SemanticID || base.ID == changed.ID {
		t.Fatal("domain edit did not change both catalog identities")
	}
}

func TestCatalogPreservesOrderAndReturnsDetachedCopies(t *testing.T) {
	sections := testSections()
	catalog, err := NewCatalog(sections)
	if err != nil {
		t.Fatal(err)
	}
	got := catalog.Sections()
	if got[0].ID != "math" || got[1].ID != "noise" {
		t.Fatalf("section order changed: %+v", got)
	}
	got[0].Label = "mutated"
	got[0].Variables[0].Label = "mutated"
	got[0].Variables[0].Default[0] = '9'
	gotAgain := catalog.Sections()
	if gotAgain[0].Label == "mutated" || gotAgain[0].Variables[0].Label == "mutated" || string(gotAgain[0].Variables[0].Default) != "2" {
		t.Fatal("catalog exposed mutable internal data")
	}

	reordered := []Section{sections[1], sections[0]}
	other, err := NewCatalog(reordered)
	if err != nil {
		t.Fatal(err)
	}
	if catalog.SemanticID == other.SemanticID || catalog.ID == other.ID {
		t.Fatal("intentional section order did not affect identity")
	}
}

func TestCatalogJSONIsDeterministicAndVerifiesIdentity(t *testing.T) {
	catalog, err := NewCatalog(testSections())
	if err != nil {
		t.Fatal(err)
	}
	first, err := json.Marshal(catalog)
	if err != nil {
		t.Fatal(err)
	}
	for range 10 {
		repeated, err := json.Marshal(catalog)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(first, repeated) {
			t.Fatalf("catalog JSON changed: %s != %s", first, repeated)
		}
	}
	var decoded Catalog
	if err := json.Unmarshal(first, &decoded); err != nil {
		t.Fatal(err)
	}
	if err := decoded.Validate(); err != nil {
		t.Fatal(err)
	}
	if decoded.ID != catalog.ID || decoded.SemanticID != catalog.SemanticID {
		t.Fatal("catalog identities changed during round trip")
	}

	var wire map[string]any
	if err := json.Unmarshal(first, &wire); err != nil {
		t.Fatal(err)
	}
	wire["id"] = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	tampered, _ := json.Marshal(wire)
	if err := json.Unmarshal(tampered, &decoded); err == nil {
		t.Fatal("tampered catalog identity accepted")
	}
}

func TestCatalogRejectsDuplicateAndMalformedEntries(t *testing.T) {
	sections := testSections()
	sections = append(sections, cloneSection(sections[0]))
	if _, err := NewCatalog(sections); err == nil {
		t.Fatal("duplicate section accepted")
	}

	sections = testSections()
	sections[1].Variables = append(sections[1].Variables, multiplierVariable().Descriptor)
	if _, err := NewCatalog(sections); err == nil {
		t.Fatal("duplicate variable accepted")
	}

	sections = testSections()
	sections[0].Variables = nil
	if _, err := NewCatalog(sections); err == nil {
		t.Fatal("empty section accepted")
	}
}
