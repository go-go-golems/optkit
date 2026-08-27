//go:build cgo

package numbergame

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/campaign"
	"github.com/go-go-golems/optkit/episode"
	"github.com/go-go-golems/optkit/local"
)

func TestRunDemoPersistsAndReplaysCompleteCampaign(t *testing.T) {
	root := filepath.Join(t.TempDir(), "local-store")
	clock := NewSequenceClock(time.Date(2026, 8, 24, 12, 0, 0, 0, time.UTC), time.Millisecond)
	summary, err := RunDemo(context.Background(), DemoOptions{Root: root, Reset: true, Clock: clock})
	if err != nil {
		t.Fatal(err)
	}
	if summary.Overview.Status != campaign.StatusCompleted {
		t.Fatalf("campaign status = %s", summary.Overview.Status)
	}
	if summary.Overview.Episodes != 8 || summary.Overview.Completed != 8 || summary.Overview.Active != 0 || summary.Overview.FailedTerminal != 0 {
		t.Fatalf("unexpected overview: %+v", summary.Overview)
	}
	if summary.Estimate.SampleSize != 4 || summary.Estimate.Missing != 0 || summary.Estimate.Value <= 0 {
		t.Fatalf("unexpected estimate: %+v", summary.Estimate)
	}
	if summary.Decision.Status != "eligible" || !summary.Decision.AllInterventionsExercised {
		t.Fatalf("unexpected decision: %+v", summary.Decision)
	}
	if summary.Budget.Violated || len(summary.Budget.Resources) != 2 {
		t.Fatalf("unexpected budget: %+v", summary.Budget)
	}
	for _, resource := range summary.Budget.Resources {
		if resource.Reserved != 0 || resource.Committed != 8 || resource.Available != 0 {
			t.Fatalf("unexpected budget resource: %+v", resource)
		}
	}

	databaseBytes, err := os.ReadFile(summary.DatabasePath)
	if err != nil {
		t.Fatal(err)
	}
	if len(databaseBytes) < 16 || string(databaseBytes[:15]) != "SQLite format 3" {
		t.Fatalf("%s is not a SQLite database", summary.DatabasePath)
	}

	profile, err := local.OpenWithClock(summary.StoreRoot, clock.Now)
	if err != nil {
		t.Fatal(err)
	}
	defer profile.Close()
	if err := profile.Metadata.Verify(context.Background(), summary.Campaign); err != nil {
		t.Fatal(err)
	}
	events, err := profile.Metadata.Read(context.Background(), summary.Campaign, 0)
	if err != nil {
		t.Fatal(err)
	}
	var proposal CandidateProposal
	for _, event := range events {
		if event.Kind == campaign.CandidateProposed {
			if err := artifact.DecodeJSON(context.Background(), profile.Artifacts, event.Payload, &proposal); err != nil {
				t.Fatal(err)
			}
			break
		}
	}
	if proposal.Candidate.ID == "" {
		t.Fatal("persisted candidate proposal not found")
	}
	if err := proposal.Catalog.Validate(); err != nil {
		t.Fatalf("persisted catalog invalid: %v", err)
	}
	if proposal.Candidate.SemanticCatalogID != proposal.Catalog.SemanticID {
		t.Fatalf("candidate/catalog provenance mismatch: %s != %s", proposal.Candidate.SemanticCatalogID, proposal.Catalog.SemanticID)
	}
	trajectory, err := episode.LoadTrajectory(context.Background(), profile.Artifacts, summary.SampleTrajectory)
	if err != nil {
		t.Fatal(err)
	}
	if len(trajectory.Events) != 5 {
		t.Fatalf("trajectory has %d events, want 5", len(trajectory.Events))
	}
	if trajectory.Events[0].Kind != "input.received" || trajectory.Events[len(trajectory.Events)-1].Kind != "episode.completed" {
		t.Fatalf("unexpected trajectory endpoints: %s -> %s", trajectory.Events[0].Kind, trajectory.Events[len(trajectory.Events)-1].Kind)
	}
}
