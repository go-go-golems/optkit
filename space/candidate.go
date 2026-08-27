package space

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-go-golems/optkit/record"
)

type ProposerKind string

const (
	ProposerHuman  ProposerKind = "human"
	ProposerLLM    ProposerKind = "llm"
	ProposerSearch ProposerKind = "search"
)

type Proposer struct {
	Kind     ProposerKind    `json:"kind" yaml:"kind"`
	Identity record.ActorRef `json:"identity" yaml:"identity"`
}

func (p Proposer) Validate() error {
	switch p.Kind {
	case ProposerHuman, ProposerLLM, ProposerSearch:
	default:
		return fmt.Errorf("unknown proposer kind %q", p.Kind)
	}
	if err := record.ValidateID("actor", string(p.Identity)); err != nil {
		return fmt.Errorf("proposer identity: %w", err)
	}
	return nil
}

type ExpectedImprovement struct {
	Metric string   `json:"metric" yaml:"metric"`
	Groups []string `json:"groups,omitempty" yaml:"groups,omitempty"`
}

type Motivation struct {
	CaseIDs          []string      `json:"case_ids,omitempty" yaml:"case_ids,omitempty"`
	DiagnosticDigest record.Digest `json:"diagnostic_digest,omitempty" yaml:"diagnostic_digest,omitempty"`
}

type CandidateIntent struct {
	Proposer            Proposer            `json:"proposer" yaml:"proposer"`
	Strategy            string              `json:"strategy" yaml:"strategy"`
	Hypothesis          string              `json:"hypothesis" yaml:"hypothesis"`
	ExpectedImprovement ExpectedImprovement `json:"expected_improvement" yaml:"expected_improvement"`
	Risks               []string            `json:"risks,omitempty" yaml:"risks,omitempty"`
	Motivation          Motivation          `json:"motivation,omitempty" yaml:"motivation,omitempty"`
}

type Candidate struct {
	ID                  record.CandidateID  `json:"id"`
	Parent              record.SnapshotID   `json:"parent"`
	Patch               record.PatchID      `json:"patch"`
	Child               record.SnapshotID   `json:"child"`
	Proposer            Proposer            `json:"proposer"`
	Strategy            string              `json:"strategy"`
	Hypothesis          string              `json:"hypothesis"`
	ExpectedImprovement ExpectedImprovement `json:"expected_improvement"`
	Risks               []string            `json:"risks,omitempty"`
	Motivation          Motivation          `json:"motivation,omitempty"`
	SemanticCatalogID   record.Digest       `json:"semantic_catalog_id"`
	CreatedAt           time.Time           `json:"created_at"`
}

// Validate checks the complete structured authoring intent independently from
// snapshot, patch, and catalog identities. Manifest resolution uses this before
// any durable materialization occurs.
func (intent CandidateIntent) Validate() error {
	_, err := NormalizeCandidateIntent(intent)
	return err
}

// NormalizeCandidateIntent validates intent and canonicalizes set-like groups
// and motivating case IDs for semantic request identity and candidate creation.
func NormalizeCandidateIntent(intent CandidateIntent) (CandidateIntent, error) {
	if err := intent.Proposer.Validate(); err != nil {
		return CandidateIntent{}, err
	}
	if strings.TrimSpace(intent.Strategy) == "" {
		return CandidateIntent{}, fmt.Errorf("candidate strategy is required")
	}
	if strings.TrimSpace(intent.Hypothesis) == "" {
		return CandidateIntent{}, fmt.Errorf("candidate hypothesis is required")
	}
	if strings.TrimSpace(intent.ExpectedImprovement.Metric) == "" {
		return CandidateIntent{}, fmt.Errorf("candidate expected improvement metric is required")
	}
	groups, err := normalizeSet("expected improvement group", intent.ExpectedImprovement.Groups)
	if err != nil {
		return CandidateIntent{}, err
	}
	caseIDs, err := normalizeSet("motivation case ID", intent.Motivation.CaseIDs)
	if err != nil {
		return CandidateIntent{}, err
	}
	if intent.Motivation.DiagnosticDigest != "" {
		if err := intent.Motivation.DiagnosticDigest.Validate(); err != nil {
			return CandidateIntent{}, fmt.Errorf("candidate diagnostic digest: %w", err)
		}
	}
	seenRisks := make(map[string]struct{}, len(intent.Risks))
	risks := make([]string, len(intent.Risks))
	for index, risk := range intent.Risks {
		if strings.TrimSpace(risk) == "" {
			return CandidateIntent{}, fmt.Errorf("candidate risk %d is blank", index)
		}
		if _, exists := seenRisks[risk]; exists {
			return CandidateIntent{}, fmt.Errorf("duplicate candidate risk %q", risk)
		}
		seenRisks[risk] = struct{}{}
		risks[index] = risk
	}
	intent.ExpectedImprovement.Groups = groups
	intent.Risks = risks
	intent.Motivation.CaseIDs = caseIDs
	return intent, nil
}

func NewCandidate(parent record.SnapshotID, patch record.PatchID, child record.SnapshotID, intent CandidateIntent, semanticCatalogID record.Digest, createdAt time.Time) (Candidate, error) {
	if err := record.ValidateID("snapshot", string(parent)); err != nil {
		return Candidate{}, fmt.Errorf("candidate parent: %w", err)
	}
	if err := record.ValidateID("patch", string(patch)); err != nil {
		return Candidate{}, fmt.Errorf("candidate patch: %w", err)
	}
	if err := record.ValidateID("snapshot", string(child)); err != nil {
		return Candidate{}, fmt.Errorf("candidate child: %w", err)
	}
	normalizedIntent, err := NormalizeCandidateIntent(intent)
	if err != nil {
		return Candidate{}, err
	}
	if err := semanticCatalogID.Validate(); err != nil {
		return Candidate{}, fmt.Errorf("candidate semantic catalog identity: %w", err)
	}

	identity := struct {
		Parent              record.SnapshotID   `json:"parent"`
		Patch               record.PatchID      `json:"patch"`
		Child               record.SnapshotID   `json:"child"`
		Proposer            Proposer            `json:"proposer"`
		Strategy            string              `json:"strategy"`
		Hypothesis          string              `json:"hypothesis"`
		ExpectedImprovement ExpectedImprovement `json:"expected_improvement"`
		Risks               []string            `json:"risks,omitempty"`
		Motivation          Motivation          `json:"motivation,omitempty"`
		SemanticCatalogID   record.Digest       `json:"semantic_catalog_id"`
	}{parent, patch, child, normalizedIntent.Proposer, normalizedIntent.Strategy, normalizedIntent.Hypothesis, normalizedIntent.ExpectedImprovement, normalizedIntent.Risks, normalizedIntent.Motivation, semanticCatalogID}
	digest, _, err := record.SemanticDigest("schema:optkit.candidate-identity/v2", identity)
	if err != nil {
		return Candidate{}, err
	}
	rawID, err := record.ContentID("candidate", digest)
	if err != nil {
		return Candidate{}, err
	}
	return Candidate{
		ID: record.CandidateID(rawID), Parent: parent, Patch: patch, Child: child,
		Proposer: normalizedIntent.Proposer, Strategy: normalizedIntent.Strategy, Hypothesis: normalizedIntent.Hypothesis,
		ExpectedImprovement: normalizedIntent.ExpectedImprovement, Risks: normalizedIntent.Risks, Motivation: normalizedIntent.Motivation,
		SemanticCatalogID: semanticCatalogID, CreatedAt: createdAt.UTC(),
	}, nil
}

func normalizeSet(name string, values []string) ([]string, error) {
	if len(values) == 0 {
		return nil, nil
	}
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for index, value := range values {
		if strings.TrimSpace(value) == "" {
			return nil, fmt.Errorf("%s %d is blank", name, index)
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	sort.Strings(out)
	return out, nil
}
