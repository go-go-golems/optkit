// Package query builds read-only campaign views from authoritative journal
// facts and immutable artifacts. It intentionally exposes no mutation API.
package query

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/budget"
	"github.com/go-go-golems/optkit/campaign"
	"github.com/go-go-golems/optkit/projection"
	"github.com/go-go-golems/optkit/record"
)

const (
	APIVersion          = "optkit.query/v1"
	DefaultEventLimit   = 200
	MaximumEventLimit   = 500
	MaximumPayloadBytes = 256 << 10
)

type CampaignCatalog interface {
	ListCampaigns(context.Context) ([]record.CampaignID, error)
}

type CampaignReader interface {
	Read(context.Context, record.CampaignID, uint64) ([]campaign.ControlEvent, error)
	Head(context.Context, record.CampaignID) (campaign.Head, error)
	Verify(context.Context, record.CampaignID) error
}

type BudgetReader interface {
	Snapshot(context.Context, record.CampaignID) (budget.Snapshot, error)
}

type MetadataStore interface {
	CampaignCatalog
	CampaignReader
	BudgetReader
}

type Service struct {
	Metadata  MetadataStore
	Artifacts artifact.Store
	Now       func() time.Time
}

type Integrity struct {
	Verified bool   `json:"verified"`
	Error    string `json:"error,omitempty"`
}

type CampaignSummary struct {
	APIVersion string              `json:"api_version"`
	ID         record.CampaignID   `json:"id"`
	Version    uint64              `json:"version"`
	Status     campaign.Status     `json:"status"`
	System     record.SystemID     `json:"system,omitempty"`
	Dataset    string              `json:"dataset,omitempty"`
	Trial      record.TrialID      `json:"trial,omitempty"`
	Overview   projection.Overview `json:"overview"`
	Budget     *budget.Snapshot    `json:"budget,omitempty"`
	EventKinds map[string]int      `json:"event_kinds"`
	Integrity  Integrity           `json:"integrity"`
	Generated  time.Time           `json:"generated_at"`
}

type CampaignPage struct {
	APIVersion string            `json:"api_version"`
	Campaigns  []CampaignSummary `json:"campaigns"`
	Count      int               `json:"count"`
	Generated  time.Time         `json:"generated_at"`
}

type EventView struct {
	Seq              uint64                `json:"seq"`
	ID               record.EventID        `json:"id"`
	Kind             campaign.EventKind    `json:"kind"`
	Schema           record.SchemaID       `json:"schema"`
	Subject          string                `json:"subject,omitempty"`
	Actor            record.ActorRef       `json:"actor"`
	OccurredAt       time.Time             `json:"occurred_at"`
	RecordedAt       time.Time             `json:"recorded_at"`
	Command          *record.CommandID     `json:"command,omitempty"`
	Causation        *record.EventID       `json:"causation,omitempty"`
	Correlation      *record.CorrelationID `json:"correlation,omitempty"`
	Tags             map[string]string     `json:"tags,omitempty"`
	Payload          artifact.Ref          `json:"payload"`
	PayloadJSON      json.RawMessage       `json:"payload_json,omitempty"`
	PayloadAvailable bool                  `json:"payload_available"`
	PayloadReason    string                `json:"payload_reason,omitempty"`
	Digest           record.Digest         `json:"digest"`
}

type EventPage struct {
	APIVersion string            `json:"api_version"`
	Campaign   record.CampaignID `json:"campaign"`
	After      uint64            `json:"after"`
	Through    uint64            `json:"through"`
	Head       uint64            `json:"head"`
	Events     []EventView       `json:"events"`
	HasMore    bool              `json:"has_more"`
	Generated  time.Time         `json:"generated_at"`
}

func (s Service) validate() error {
	if s.Metadata == nil || s.Artifacts == nil {
		return fmt.Errorf("query service requires metadata and artifact stores")
	}
	return nil
}

func (s Service) now() time.Time {
	if s.Now != nil {
		return s.Now().UTC()
	}
	return time.Now().UTC()
}

func (s Service) ListCampaigns(ctx context.Context) (CampaignPage, error) {
	if err := s.validate(); err != nil {
		return CampaignPage{}, err
	}
	ids, err := s.Metadata.ListCampaigns(ctx)
	if err != nil {
		return CampaignPage{}, err
	}
	page := CampaignPage{APIVersion: APIVersion, Campaigns: make([]CampaignSummary, 0, len(ids)), Generated: s.now()}
	for _, id := range ids {
		summary, err := s.Campaign(ctx, id)
		if err != nil {
			return CampaignPage{}, fmt.Errorf("project campaign %s: %w", id, err)
		}
		page.Campaigns = append(page.Campaigns, summary)
	}
	page.Count = len(page.Campaigns)
	return page, nil
}

func (s Service) Campaign(ctx context.Context, id record.CampaignID) (CampaignSummary, error) {
	if err := s.validate(); err != nil {
		return CampaignSummary{}, err
	}
	if err := record.ValidateID("campaign", string(id)); err != nil {
		return CampaignSummary{}, err
	}
	events, err := s.Metadata.Read(ctx, id, 0)
	if err != nil {
		return CampaignSummary{}, err
	}
	if len(events) == 0 {
		return CampaignSummary{}, fmt.Errorf("campaign %s not found", id)
	}
	overview, err := projection.RebuildOverview(id, events)
	if err != nil {
		return CampaignSummary{}, err
	}
	summary := CampaignSummary{
		APIVersion: APIVersion, ID: id, Version: overview.Version,
		Status: overview.Status, Overview: overview, EventKinds: make(map[string]int),
		Generated: s.now(),
	}
	for _, event := range events {
		summary.EventKinds[string(event.Kind)]++
		if event.Kind == campaign.CampaignCreated {
			summary.System, summary.Dataset, summary.Trial = s.campaignIdentity(ctx, event.Payload)
		}
	}
	if snapshot, budgetErr := s.Metadata.Snapshot(ctx, id); budgetErr == nil {
		summary.Budget = &snapshot
	} else if !errors.Is(budgetErr, budget.ErrNotDefined) {
		return CampaignSummary{}, budgetErr
	}
	if verifyErr := s.Metadata.Verify(ctx, id); verifyErr != nil {
		summary.Integrity = Integrity{Error: verifyErr.Error()}
	} else {
		summary.Integrity = Integrity{Verified: true}
	}
	return summary, nil
}

func (s Service) campaignIdentity(ctx context.Context, ref artifact.Ref) (record.SystemID, string, record.TrialID) {
	if ref.Size > MaximumPayloadBytes || ref.Sensitivity == artifact.SensitivityConfidential || ref.Sensitivity == artifact.SensitivityRestricted {
		return "", "", ""
	}
	data, err := artifact.ReadAll(ctx, s.Artifacts, ref)
	if err != nil {
		return "", "", ""
	}
	var value struct {
		System  record.SystemID `json:"system"`
		Dataset string          `json:"dataset"`
		Trial   record.TrialID  `json:"trial"`
	}
	if json.Unmarshal(data, &value) != nil {
		return "", "", ""
	}
	return value.System, value.Dataset, value.Trial
}

func (s Service) Events(ctx context.Context, id record.CampaignID, after uint64, limit int) (EventPage, error) {
	if err := s.validate(); err != nil {
		return EventPage{}, err
	}
	if err := record.ValidateID("campaign", string(id)); err != nil {
		return EventPage{}, err
	}
	if limit == 0 {
		limit = DefaultEventLimit
	}
	if limit < 1 || limit > MaximumEventLimit {
		return EventPage{}, fmt.Errorf("event limit must be in [1,%d]", MaximumEventLimit)
	}
	head, err := s.Metadata.Head(ctx, id)
	if err != nil {
		return EventPage{}, err
	}
	if head.Version == 0 {
		return EventPage{}, fmt.Errorf("campaign %s not found", id)
	}
	events, err := s.Metadata.Read(ctx, id, after)
	if err != nil {
		return EventPage{}, err
	}
	hasMore := len(events) > limit
	if hasMore {
		events = events[:limit]
	}
	views := make([]EventView, 0, len(events))
	through := after
	for _, event := range events {
		view := EventView{
			Seq: event.Seq, ID: event.ID, Kind: event.Kind, Schema: event.Schema,
			Subject: event.Subject, Actor: event.Actor, OccurredAt: event.OccurredAt,
			RecordedAt: event.RecordedAt, Command: event.Command, Causation: event.Causation,
			Correlation: event.Correlation, Tags: event.Tags, Payload: event.Payload,
			Digest: event.Digest,
		}
		s.populatePayload(ctx, &view)
		views = append(views, view)
		through = event.Seq
	}
	return EventPage{
		APIVersion: APIVersion, Campaign: id, After: after, Through: through,
		Head: head.Version, Events: views, HasMore: hasMore, Generated: s.now(),
	}, nil
}

func (s Service) populatePayload(ctx context.Context, view *EventView) {
	ref := view.Payload
	switch ref.Sensitivity {
	case artifact.SensitivityPublic, artifact.SensitivityInternal:
		// Bounded JSON previews are allowed for the local read-only explorer.
	case artifact.SensitivityConfidential, artifact.SensitivityRestricted:
		view.PayloadReason = "sensitivity policy"
		return
	default:
		view.PayloadReason = "unknown sensitivity"
		return
	}
	if ref.Size > MaximumPayloadBytes {
		view.PayloadReason = "payload exceeds preview limit"
		return
	}
	data, err := artifact.ReadAll(ctx, s.Artifacts, ref)
	if err != nil {
		view.PayloadReason = err.Error()
		return
	}
	if !json.Valid(data) {
		view.PayloadReason = "payload is not JSON"
		return
	}
	view.PayloadJSON = append(json.RawMessage(nil), data...)
	view.PayloadAvailable = true
}

func (s Service) Head(ctx context.Context, id record.CampaignID) (campaign.Head, error) {
	if err := s.validate(); err != nil {
		return campaign.Head{}, err
	}
	return s.Metadata.Head(ctx, id)
}
