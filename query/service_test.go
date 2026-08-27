//go:build cgo

package query_test

import (
	"context"
	"testing"

	"github.com/go-go-golems/optkit/examples/numbergame"
	"github.com/go-go-golems/optkit/local"
	"github.com/go-go-golems/optkit/query"
)

func TestServiceProjectsNumbergameCampaignWithoutMutation(t *testing.T) {
	root := t.TempDir()
	summary, err := numbergame.RunDemo(context.Background(), numbergame.DemoOptions{Root: root})
	if err != nil {
		t.Fatalf("run Numbergame: %v", err)
	}
	profile, err := local.Open(root)
	if err != nil {
		t.Fatalf("open profile: %v", err)
	}
	t.Cleanup(func() { _ = profile.Close() })
	service := query.Service{Metadata: profile.Metadata, Artifacts: profile.Artifacts}

	before, err := profile.Metadata.Head(context.Background(), summary.Campaign)
	if err != nil {
		t.Fatalf("head before query: %v", err)
	}
	campaigns, err := service.ListCampaigns(context.Background())
	if err != nil {
		t.Fatalf("list campaigns: %v", err)
	}
	if campaigns.Count != 1 || campaigns.Campaigns[0].ID != summary.Campaign {
		t.Fatalf("campaign page = %+v", campaigns)
	}
	view, err := service.Campaign(context.Background(), summary.Campaign)
	if err != nil {
		t.Fatalf("campaign view: %v", err)
	}
	if !view.Integrity.Verified || view.Version != 81 || view.Overview.Completed != 8 {
		t.Fatalf("campaign view = %+v", view)
	}
	if view.System != numbergame.SystemID || view.Trial != summary.Trial || view.Budget == nil {
		t.Fatalf("campaign identity/budget missing: %+v", view)
	}
	page, err := service.Events(context.Background(), summary.Campaign, 0, 2)
	if err != nil {
		t.Fatalf("event page: %v", err)
	}
	if len(page.Events) != 2 || !page.HasMore || page.Through != 2 || !page.Events[0].PayloadAvailable {
		t.Fatalf("event page = %+v", page)
	}
	if _, err := service.Events(context.Background(), summary.Campaign, 0, query.MaximumEventLimit+1); err == nil {
		t.Fatal("event query accepted a limit above the server ceiling")
	}
	after, err := profile.Metadata.Head(context.Background(), summary.Campaign)
	if err != nil {
		t.Fatalf("head after query: %v", err)
	}
	if before != after {
		t.Fatalf("query changed journal head: before=%+v after=%+v", before, after)
	}
}
