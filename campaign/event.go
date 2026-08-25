package campaign

import (
	"fmt"
	"time"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/record"
)

type EventKind string

const (
	CampaignCreated   EventKind = "CampaignCreated"
	PlanCompiled      EventKind = "PlanCompiled"
	CampaignStarted   EventKind = "CampaignStarted"
	CampaignPaused    EventKind = "CampaignPaused"
	CampaignResumed   EventKind = "CampaignResumed"
	CampaignStopping  EventKind = "CampaignStopping"
	CampaignStopped   EventKind = "CampaignStopped"
	CampaignFailed    EventKind = "CampaignFailed"
	CampaignCompleted EventKind = "CampaignCompleted"

	CandidateProposed     EventKind = "CandidateProposed"
	SnapshotMaterialized  EventKind = "SnapshotMaterialized"
	TrialPlanned          EventKind = "TrialPlanned"
	EpisodeScheduled      EventKind = "EpisodeScheduled"
	EpisodeLeaseGranted   EventKind = "EpisodeLeaseGranted"
	EpisodeAttemptStarted EventKind = "EpisodeAttemptStarted"
	EpisodeCompleted      EventKind = "EpisodeCompleted"
	EpisodeFailed         EventKind = "EpisodeFailed"
	ObservationRecorded   EventKind = "ObservationRecorded"
	EstimateRecorded      EventKind = "EstimateRecorded"
	DecisionRecorded      EventKind = "DecisionRecorded"
	BudgetReserved        EventKind = "BudgetReserved"
	UsageCommitted        EventKind = "UsageCommitted"
	BudgetReleased        EventKind = "BudgetReleased"
)

var GenesisDigest = record.SumBytes([]byte("optkit.campaign.genesis/v1"))

type NewEvent struct {
	Kind       EventKind
	Schema     record.SchemaID
	Subject    string
	OccurredAt time.Time
	Actor      record.ActorRef
	Command    *record.CommandID
	Payload    artifact.Ref
	Tags       map[string]string
}

type ControlEvent struct {
	ID             record.EventID    `json:"id"`
	Campaign       record.CampaignID `json:"campaign"`
	Seq            uint64            `json:"seq"`
	Kind           EventKind         `json:"kind"`
	Schema         record.SchemaID   `json:"schema"`
	Subject        string            `json:"subject,omitempty"`
	OccurredAt     time.Time         `json:"occurred_at"`
	RecordedAt     time.Time         `json:"recorded_at"`
	Actor          record.ActorRef   `json:"actor"`
	Command        *record.CommandID `json:"command,omitempty"`
	PreviousDigest record.Digest     `json:"previous_digest"`
	Payload        artifact.Ref      `json:"payload"`
	Tags           map[string]string `json:"tags,omitempty"`
	Digest         record.Digest     `json:"digest"`
}

func MaterializeEvent(campaignID record.CampaignID, seq uint64, recordedAt time.Time, previous record.Digest, input NewEvent) (ControlEvent, error) {
	if err := record.ValidateID("campaign", string(campaignID)); err != nil {
		return ControlEvent{}, err
	}
	if seq == 0 {
		return ControlEvent{}, fmt.Errorf("control event sequence must be positive")
	}
	if input.Kind == "" {
		return ControlEvent{}, fmt.Errorf("control event kind is required")
	}
	if err := input.Schema.Validate(); err != nil {
		return ControlEvent{}, err
	}
	if input.Actor == "" {
		return ControlEvent{}, fmt.Errorf("control event actor is required")
	}
	if err := input.Payload.Validate(); err != nil {
		return ControlEvent{}, fmt.Errorf("control event payload: %w", err)
	}
	if err := previous.Validate(); err != nil {
		return ControlEvent{}, fmt.Errorf("previous event digest: %w", err)
	}
	if input.OccurredAt.IsZero() {
		input.OccurredAt = recordedAt
	}
	event := ControlEvent{
		Campaign:       campaignID,
		Seq:            seq,
		Kind:           input.Kind,
		Schema:         input.Schema,
		Subject:        input.Subject,
		OccurredAt:     input.OccurredAt.UTC(),
		RecordedAt:     recordedAt.UTC(),
		Actor:          input.Actor,
		Command:        input.Command,
		PreviousDigest: previous,
		Payload:        input.Payload,
		Tags:           cloneTags(input.Tags),
	}
	identity := eventIdentity(event)
	digest, _, err := record.SemanticDigest("schema:optkit.control-event-identity/v1", identity)
	if err != nil {
		return ControlEvent{}, err
	}
	rawID, err := record.ContentID("event", digest)
	if err != nil {
		return ControlEvent{}, err
	}
	event.ID = record.EventID(rawID)
	event.Digest = digest
	return event, nil
}

func VerifyEventDigest(event ControlEvent) error {
	digest, _, err := record.SemanticDigest("schema:optkit.control-event-identity/v1", eventIdentity(event))
	if err != nil {
		return err
	}
	if digest != event.Digest {
		return fmt.Errorf("event %s digest mismatch: got %s, want %s", event.ID, digest, event.Digest)
	}
	rawID, err := record.ContentID("event", digest)
	if err != nil {
		return err
	}
	if record.EventID(rawID) != event.ID {
		return fmt.Errorf("event ID mismatch: got %s, want %s", rawID, event.ID)
	}
	return nil
}

func eventIdentity(event ControlEvent) any {
	return struct {
		Campaign       record.CampaignID `json:"campaign"`
		Seq            uint64            `json:"seq"`
		Kind           EventKind         `json:"kind"`
		Schema         record.SchemaID   `json:"schema"`
		Subject        string            `json:"subject,omitempty"`
		OccurredAt     time.Time         `json:"occurred_at"`
		RecordedAt     time.Time         `json:"recorded_at"`
		Actor          record.ActorRef   `json:"actor"`
		Command        *record.CommandID `json:"command,omitempty"`
		PreviousDigest record.Digest     `json:"previous_digest"`
		Payload        artifact.Ref      `json:"payload"`
		Tags           map[string]string `json:"tags,omitempty"`
	}{
		Campaign: event.Campaign, Seq: event.Seq, Kind: event.Kind, Schema: event.Schema,
		Subject: event.Subject, OccurredAt: event.OccurredAt, RecordedAt: event.RecordedAt,
		Actor: event.Actor, Command: event.Command, PreviousDigest: event.PreviousDigest,
		Payload: event.Payload, Tags: event.Tags,
	}
}

func cloneTags(tags map[string]string) map[string]string {
	if len(tags) == 0 {
		return nil
	}
	out := make(map[string]string, len(tags))
	for key, value := range tags {
		out[key] = value
	}
	return out
}
