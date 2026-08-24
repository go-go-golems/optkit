package space

import (
	"fmt"
	"time"

	"github.com/go-go-golems/optkit/record"
)

type Candidate struct {
	ID         record.CandidateID `json:"id"`
	Parent     record.SnapshotID  `json:"parent"`
	Patch      record.PatchID     `json:"patch"`
	Child      record.SnapshotID  `json:"child"`
	Proposer   record.ActorRef    `json:"proposer"`
	Strategy   string             `json:"strategy"`
	Hypothesis string             `json:"hypothesis"`
	Targets    []string           `json:"targets,omitempty"`
	Risks      []string           `json:"risks,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
}

func NewCandidate(parent record.SnapshotID, patch record.PatchID, child record.SnapshotID, proposer record.ActorRef, strategy, hypothesis string, targets, risks []string, createdAt time.Time) (Candidate, error) {
	if hypothesis == "" {
		return Candidate{}, fmt.Errorf("candidate hypothesis is required")
	}
	identity := struct {
		Parent     record.SnapshotID `json:"parent"`
		Patch      record.PatchID    `json:"patch"`
		Child      record.SnapshotID `json:"child"`
		Proposer   record.ActorRef   `json:"proposer"`
		Strategy   string            `json:"strategy"`
		Hypothesis string            `json:"hypothesis"`
		Targets    []string          `json:"targets,omitempty"`
		Risks      []string          `json:"risks,omitempty"`
	}{parent, patch, child, proposer, strategy, hypothesis, targets, risks}
	digest, _, err := record.SemanticDigest("schema:optkit.candidate-identity/v1", identity)
	if err != nil {
		return Candidate{}, err
	}
	rawID, err := record.ContentID("candidate", digest)
	if err != nil {
		return Candidate{}, err
	}
	return Candidate{
		ID:         record.CandidateID(rawID),
		Parent:     parent,
		Patch:      patch,
		Child:      child,
		Proposer:   proposer,
		Strategy:   strategy,
		Hypothesis: hypothesis,
		Targets:    append([]string(nil), targets...),
		Risks:      append([]string(nil), risks...),
		CreatedAt:  createdAt.UTC(),
	}, nil
}
