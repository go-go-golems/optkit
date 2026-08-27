package campaign

import (
	"context"
	"testing"
	"time"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/artifact/memory"
	"github.com/go-go-golems/optkit/record"
)

func TestReducerLifecycleAndEpisode(t *testing.T) {
	ctx := context.Background()
	store := memory.New()
	campaignID := record.CampaignID("campaign:test")
	state := Initial(campaignID)
	previous := GenesisDigest
	seq := uint64(0)
	appendEvent := func(kind EventKind, subject string, tags map[string]string) {
		t.Helper()
		seq++
		ref, err := artifact.PutCanonical(ctx, store, "schema:test.event/v1", artifact.SensitivityInternal, map[string]any{"kind": kind})
		if err != nil {
			t.Fatal(err)
		}
		event, err := MaterializeEvent(campaignID, seq, time.Unix(int64(seq), 0), previous, NewEvent{Kind: kind, Schema: "schema:test.event/v1", Subject: subject, Actor: "actor:test", Payload: ref, Tags: tags})
		if err != nil {
			t.Fatal(err)
		}
		state, err = Apply(state, event)
		if err != nil {
			t.Fatal(err)
		}
		previous = event.Digest
	}
	appendEvent(CampaignCreated, "", nil)
	appendEvent(PlanCompiled, "", nil)
	appendEvent(CampaignStarted, "", nil)
	appendEvent(EpisodeScheduled, "episode:e1", nil)
	appendEvent(EpisodeLeaseGranted, "episode:e1", nil)
	appendEvent(EpisodeAttemptStarted, "episode:e1", nil)
	appendEvent(EpisodeCompleted, "episode:e1", nil)
	appendEvent(CampaignCompleted, "", nil)
	if state.Status != StatusCompleted || state.Episodes["episode:e1"] != EpisodeCompletedState {
		t.Fatalf("unexpected state: %+v", state)
	}
}

func TestReducerRejectsPrematureCompletion(t *testing.T) {
	state := State{Campaign: "campaign:test", Status: StatusRunning, LastDigest: GenesisDigest, Episodes: map[string]EpisodeStatus{"episode:e1": EpisodeQueued}}
	if err := ValidateTransition(state, CampaignCompleted, "", nil); err == nil {
		t.Fatal("premature completion unexpectedly accepted")
	}
}

func TestUsageCanBeCommittedAfterTerminalCampaign(t *testing.T) {
	state := State{Campaign: "campaign:test", Status: StatusCompleted, LastDigest: GenesisDigest, Episodes: map[string]EpisodeStatus{}}
	if err := ValidateTransition(state, UsageCommitted, "work:test", nil); err != nil {
		t.Fatalf("terminal usage custody was rejected: %v", err)
	}
	if err := ValidateTransition(state, BudgetReserved, "work:test", nil); err == nil {
		t.Fatal("new reservation unexpectedly accepted after completion")
	}
}
