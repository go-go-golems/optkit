package space

import (
	"testing"
	"time"

	"github.com/go-go-golems/optkit/record"
)

func baseCandidateIntent() CandidateIntent {
	return CandidateIntent{
		Proposer:   Proposer{Kind: ProposerHuman, Identity: "actor:demo"},
		Strategy:   "manual-coordinate/v1",
		Hypothesis: "Increasing the multiplier should improve exact accuracy.",
		ExpectedImprovement: ExpectedImprovement{
			Metric: "target.accuracy", Groups: []string{"development", "exact-target"},
		},
		Risks: []string{"may overfit", "may increase error elsewhere"},
		Motivation: Motivation{
			CaseIDs:          []string{"case-02", "case-01"},
			DiagnosticDigest: record.SumBytes([]byte("diagnostic")),
		},
	}
}

func newTestCandidate(t *testing.T, intent CandidateIntent, catalog record.Digest, now time.Time) Candidate {
	t.Helper()
	candidate, err := NewCandidate("snapshot:parent", "patch:change", "snapshot:child", intent, catalog, now)
	if err != nil {
		t.Fatal(err)
	}
	return candidate
}

func TestCandidateIdentityIncludesSemanticIntentAndExcludesTime(t *testing.T) {
	catalog := record.SumBytes([]byte("catalog"))
	baseIntent := baseCandidateIntent()
	base := newTestCandidate(t, baseIntent, catalog, time.Unix(100, 0))
	later := newTestCandidate(t, baseIntent, catalog, time.Unix(200, 0))
	if base.ID != later.ID {
		t.Fatalf("timestamp changed candidate identity: %s != %s", base.ID, later.ID)
	}
	if base.CreatedAt.Equal(later.CreatedAt) {
		t.Fatal("display timestamps unexpectedly equal")
	}

	tests := []struct {
		name   string
		mutate func(*CandidateIntent)
	}{
		{name: "proposer-kind", mutate: func(i *CandidateIntent) { i.Proposer.Kind = ProposerLLM }},
		{name: "proposer-identity", mutate: func(i *CandidateIntent) { i.Proposer.Identity = "actor:other" }},
		{name: "strategy", mutate: func(i *CandidateIntent) { i.Strategy = "search/v1" }},
		{name: "hypothesis", mutate: func(i *CandidateIntent) { i.Hypothesis += " More detail." }},
		{name: "metric", mutate: func(i *CandidateIntent) { i.ExpectedImprovement.Metric = "target.recall" }},
		{name: "groups", mutate: func(i *CandidateIntent) {
			i.ExpectedImprovement.Groups = append(i.ExpectedImprovement.Groups, "holdout")
		}},
		{name: "risk", mutate: func(i *CandidateIntent) { i.Risks[0] = "different risk" }},
		{name: "risk-order", mutate: func(i *CandidateIntent) { i.Risks[0], i.Risks[1] = i.Risks[1], i.Risks[0] }},
		{name: "case", mutate: func(i *CandidateIntent) { i.Motivation.CaseIDs = append(i.Motivation.CaseIDs, "case-03") }},
		{name: "diagnostic", mutate: func(i *CandidateIntent) { i.Motivation.DiagnosticDigest = record.SumBytes([]byte("other")) }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			intent := baseCandidateIntent()
			test.mutate(&intent)
			changed := newTestCandidate(t, intent, catalog, time.Unix(100, 0))
			if changed.ID == base.ID {
				t.Fatalf("%s did not change candidate identity", test.name)
			}
		})
	}
	otherCatalog := newTestCandidate(t, baseIntent, record.SumBytes([]byte("other catalog")), time.Unix(100, 0))
	if otherCatalog.ID == base.ID {
		t.Fatal("semantic catalog did not change candidate identity")
	}
}

func TestCandidateIntentValidatesWithoutDurableIdentities(t *testing.T) {
	intent := baseCandidateIntent()
	if err := intent.Validate(); err != nil {
		t.Fatalf("valid candidate intent rejected: %v", err)
	}
	normalized, err := NormalizeCandidateIntent(intent)
	if err != nil {
		t.Fatalf("normalize valid candidate intent: %v", err)
	}
	if normalized.Motivation.CaseIDs[0] != "case-01" || normalized.ExpectedImprovement.Groups[0] != "development" {
		t.Fatalf("standalone intent normalization failed: %+v", normalized)
	}
	intent.Hypothesis = ""
	if err := intent.Validate(); err == nil {
		t.Fatal("incomplete candidate intent accepted")
	}
}

func TestCandidateNormalizesOnlySetLikeCollections(t *testing.T) {
	catalog := record.SumBytes([]byte("catalog"))
	first := baseCandidateIntent()
	second := baseCandidateIntent()
	second.ExpectedImprovement.Groups = []string{"exact-target", "development", "development"}
	second.Motivation.CaseIDs = []string{"case-01", "case-02", "case-01"}
	one := newTestCandidate(t, first, catalog, time.Unix(100, 0))
	two := newTestCandidate(t, second, catalog, time.Unix(100, 0))
	if one.ID != two.ID {
		t.Fatalf("set order/duplicates changed identity: %s != %s", one.ID, two.ID)
	}
	if len(two.ExpectedImprovement.Groups) != 2 || two.ExpectedImprovement.Groups[0] != "development" {
		t.Fatalf("groups not normalized: %#v", two.ExpectedImprovement.Groups)
	}
	if len(two.Motivation.CaseIDs) != 2 || two.Motivation.CaseIDs[0] != "case-01" {
		t.Fatalf("case IDs not normalized: %#v", two.Motivation.CaseIDs)
	}
}

func TestCandidateRejectsIncompleteOrAmbiguousIntent(t *testing.T) {
	catalog := record.SumBytes([]byte("catalog"))
	tests := []struct {
		name   string
		mutate func(*CandidateIntent)
	}{
		{name: "proposer-kind", mutate: func(i *CandidateIntent) { i.Proposer.Kind = "robot" }},
		{name: "proposer-identity", mutate: func(i *CandidateIntent) { i.Proposer.Identity = "demo" }},
		{name: "strategy", mutate: func(i *CandidateIntent) { i.Strategy = " " }},
		{name: "hypothesis", mutate: func(i *CandidateIntent) { i.Hypothesis = "" }},
		{name: "metric", mutate: func(i *CandidateIntent) { i.ExpectedImprovement.Metric = "" }},
		{name: "blank-group", mutate: func(i *CandidateIntent) { i.ExpectedImprovement.Groups = []string{""} }},
		{name: "blank-case", mutate: func(i *CandidateIntent) { i.Motivation.CaseIDs = []string{" "} }},
		{name: "bad-digest", mutate: func(i *CandidateIntent) { i.Motivation.DiagnosticDigest = "bad" }},
		{name: "blank-risk", mutate: func(i *CandidateIntent) { i.Risks = []string{""} }},
		{name: "duplicate-risk", mutate: func(i *CandidateIntent) { i.Risks = []string{"same", "same"} }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			intent := baseCandidateIntent()
			test.mutate(&intent)
			if err := intent.Validate(); err == nil {
				t.Fatal("invalid standalone intent accepted")
			}
			if _, err := NewCandidate("snapshot:parent", "patch:change", "snapshot:child", intent, catalog, time.Now()); err == nil {
				t.Fatal("invalid candidate accepted")
			}
		})
	}
	if _, err := NewCandidate("bad", "patch:change", "snapshot:child", baseCandidateIntent(), catalog, time.Now()); err == nil {
		t.Fatal("invalid parent accepted")
	}
	if _, err := NewCandidate("snapshot:parent", "patch:change", "snapshot:child", baseCandidateIntent(), "bad", time.Now()); err == nil {
		t.Fatal("invalid semantic catalog accepted")
	}
}
