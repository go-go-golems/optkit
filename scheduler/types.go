package scheduler

import (
	"fmt"
	"time"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/record"
)

type WorkKind string

type WorkStatus string

const (
	WorkReady     WorkStatus = "ready"
	WorkLeased    WorkStatus = "leased"
	WorkCompleted WorkStatus = "completed"
	WorkFailed    WorkStatus = "failed"
)

type ResourceClaim struct {
	Resource string `json:"resource"`
	Units    int64  `json:"units"`
}

type WorkItem struct {
	ID             record.WorkID     `json:"id"`
	Campaign       record.CampaignID `json:"campaign"`
	Kind           WorkKind          `json:"kind"`
	SemanticKey    record.Digest     `json:"semantic_key"`
	Payload        artifact.Ref      `json:"payload"`
	Priority       int               `json:"priority"`
	EarliestStart  time.Time         `json:"earliest_start"`
	LeaseDuration  time.Duration     `json:"lease_duration"`
	ResourceClaims []ResourceClaim   `json:"resource_claims,omitempty"`
}

func NewWorkItem(campaign record.CampaignID, kind WorkKind, semanticKey record.Digest, payload artifact.Ref, priority int, earliest time.Time, leaseDuration time.Duration, claims []ResourceClaim) (WorkItem, error) {
	if err := record.ValidateID("campaign", string(campaign)); err != nil {
		return WorkItem{}, err
	}
	if kind == "" {
		return WorkItem{}, fmt.Errorf("work kind is required")
	}
	if err := semanticKey.Validate(); err != nil {
		return WorkItem{}, err
	}
	if err := payload.Validate(); err != nil {
		return WorkItem{}, err
	}
	if leaseDuration <= 0 {
		return WorkItem{}, fmt.Errorf("lease duration must be positive")
	}
	identity := struct {
		Campaign    record.CampaignID `json:"campaign"`
		Kind        WorkKind          `json:"kind"`
		SemanticKey record.Digest     `json:"semantic_key"`
	}{Campaign: campaign, Kind: kind, SemanticKey: semanticKey}
	digest, _, err := record.SemanticDigest("schema:optkit.work-identity/v1", identity)
	if err != nil {
		return WorkItem{}, err
	}
	rawID, err := record.ContentID("work", digest)
	if err != nil {
		return WorkItem{}, err
	}
	return WorkItem{
		ID: record.WorkID(rawID), Campaign: campaign, Kind: kind, SemanticKey: semanticKey,
		Payload: payload, Priority: priority, EarliestStart: earliest.UTC(), LeaseDuration: leaseDuration,
		ResourceClaims: append([]ResourceClaim(nil), claims...),
	}, nil
}

type LeaseRequest struct {
	Kinds []WorkKind
	Limit int
	Now   time.Time
}

type Lease struct {
	ID        record.LeaseID `json:"id"`
	Item      WorkItem       `json:"item"`
	Worker    string         `json:"worker"`
	Attempt   int            `json:"attempt"`
	ExpiresAt time.Time      `json:"expires_at"`
}

type WorkResult struct {
	Artifact artifact.Ref `json:"artifact"`
}

type WorkFailure struct {
	Code      string         `json:"code"`
	Message   string         `json:"message"`
	Retryable bool           `json:"retryable"`
	Backoff   time.Duration  `json:"backoff"`
	Evidence  []artifact.Ref `json:"evidence,omitempty"`
}

type WorkRecord struct {
	Item        WorkItem
	Status      WorkStatus
	Attempt     int
	LeaseID     record.LeaseID
	LeasedBy    string
	LeaseExpiry time.Time
	Result      *WorkResult
	Failure     *WorkFailure
}
