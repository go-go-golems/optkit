package experiment

import (
	"context"
	"testing"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/artifact/memory"
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
