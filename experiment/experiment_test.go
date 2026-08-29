package experiment

import (
	"context"
	"testing"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/artifact/memory"
	"github.com/go-go-golems/optkit/record"
)

func TestCompleteBlockExpansionIsDeterministic(t *testing.T) {
	store := memory.New()
	ctx := context.Background()
	ref1, err := artifact.PutCanonical(ctx, store, "schema:test.case/v1", artifact.SensitivityInternal, map[string]int{"input": 1})
	if err != nil {
		t.Fatal(err)
	}
	ref2, err := artifact.PutCanonical(ctx, store, "schema:test.case/v1", artifact.SensitivityInternal, map[string]int{"input": 2})
	if err != nil {
		t.Fatal(err)
	}
	dataset, err := NewDatasetManifest("development", []Case{{ID: "c1", Input: ref1}, {ID: "c2", Input: ref2}})
	if err != nil {
		t.Fatal(err)
	}
	plan, err := NewCompleteBlockTrial([]Arm{{ID: "base", Snapshot: "snapshot:base"}, {ID: "candidate", Snapshot: "snapshot:candidate"}}, dataset, 2, "numbergame/v1")
	if err != nil {
		t.Fatal(err)
	}
	first, err := Expand(plan)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Expand(plan)
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 8 || len(second) != 8 {
		t.Fatalf("got %d and %d specs", len(first), len(second))
	}
	for i := range first {
		if first[i].ID != second[i].ID || first[i].SemanticKey != second[i].SemanticKey {
			t.Fatalf("spec %d is not deterministic", i)
		}
	}
}

func TestDatasetManifestOwnsNestedCaseData(t *testing.T) {
	schema := record.SchemaID("schema:test.case/v1")
	cases := []Case{{
		ID: "c1",
		Input: artifact.Ref{
			Digest: record.SumBytes([]byte("case")), MediaType: "application/json",
			Schema: &schema, Size: 4, Sensitivity: artifact.SensitivityInternal,
		},
		Groups:   []string{"development"},
		Metadata: map[string]string{"source": "fixture"},
	}}
	manifest, err := NewDatasetManifest("development", cases)
	if err != nil {
		t.Fatal(err)
	}
	originalDigest := manifest.Digest
	originalID := manifest.ID

	cases[0].Groups[0] = "mutated"
	cases[0].Metadata["source"] = "mutated"
	*cases[0].Input.Schema = "schema:mutated/v1"

	if manifest.Cases[0].Groups[0] != "development" || manifest.Cases[0].Metadata["source"] != "fixture" || *manifest.Cases[0].Input.Schema != "schema:test.case/v1" {
		t.Fatalf("manifest aliases caller-owned case data: %+v", manifest.Cases[0])
	}
	rebuilt, err := NewDatasetManifest(manifest.Role, manifest.Cases)
	if err != nil {
		t.Fatal(err)
	}
	if rebuilt.Digest != originalDigest || rebuilt.ID != originalID {
		t.Fatalf("manifest identity changed after caller mutation: got %s/%s want %s/%s", rebuilt.ID, rebuilt.Digest, originalID, originalDigest)
	}

	plan, err := NewCompleteBlockTrial([]Arm{{ID: "base", Snapshot: "snapshot:base"}, {ID: "candidate", Snapshot: "snapshot:candidate"}}, manifest, 1, "test/v1")
	if err != nil {
		t.Fatal(err)
	}
	manifest.Cases[0].Metadata["source"] = "changed-after-plan"
	if plan.Dataset.Cases[0].Metadata["source"] != "fixture" {
		t.Fatalf("trial plan aliases dataset case metadata: %+v", plan.Dataset.Cases[0])
	}
}

func TestPairedMeanRejectsMissingPair(t *testing.T) {
	rows := []NumericObservation{{CaseID: "c1", ArmID: "base", Value: 1, Valid: true}, {CaseID: "c1", ArmID: "candidate", Value: 2, Valid: true}, {CaseID: "c2", ArmID: "base", Value: 3, Valid: true}}
	if _, err := PairedMean(rows, "base", "candidate"); err == nil {
		t.Fatal("missing pair unexpectedly accepted")
	}
}

func TestPairedMean(t *testing.T) {
	rows := []NumericObservation{{CaseID: "c1", ArmID: "base", Value: 1, Valid: true}, {CaseID: "c1", ArmID: "candidate", Value: 2, Valid: true}, {CaseID: "c2", ArmID: "base", Value: 4, Valid: true}, {CaseID: "c2", ArmID: "candidate", Value: 7, Valid: true}}
	estimate, err := PairedMean(rows, "base", "candidate")
	if err != nil {
		t.Fatal(err)
	}
	if estimate.Value != 2 || estimate.SampleSize != 2 {
		t.Fatalf("unexpected estimate %+v", estimate)
	}
}
