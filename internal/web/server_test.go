//go:build cgo

package web

import (
	"bufio"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-go-golems/optkit/examples/numbergame"
	"github.com/go-go-golems/optkit/local"
	"github.com/go-go-golems/optkit/query"
)

func testServer(t *testing.T) (*httptest.Server, numbergame.DemoSummary) {
	t.Helper()
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
	handler, err := (Server{
		Query:        query.Service{Metadata: profile.Metadata, Artifacts: profile.Artifacts},
		PollInterval: 5 * time.Millisecond,
	}).Handler()
	if err != nil {
		t.Fatalf("construct handler: %v", err)
	}
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return server, summary
}

func TestServerServesReadOnlyAPIAndStaticExplorer(t *testing.T) {
	server, summary := testServer(t)
	for _, path := range []string{"/", "/static/styles.css", "/static/app.js", "/api/v1/health", "/api/v1/campaigns", "/api/v1/campaigns/" + string(summary.Campaign), "/api/v1/campaigns/" + string(summary.Campaign) + "/events?limit=10"} {
		response, err := http.Get(server.URL + path)
		if err != nil {
			t.Fatalf("GET %s: %v", path, err)
		}
		_ = response.Body.Close()
		if response.StatusCode != http.StatusOK {
			t.Fatalf("GET %s status = %d", path, response.StatusCode)
		}
		if response.Header.Get("X-Content-Type-Options") != "nosniff" {
			t.Fatalf("GET %s missing security headers", path)
		}
	}

	request, err := http.NewRequest(http.MethodPost, server.URL+"/api/v1/campaigns", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("POST campaigns: %v", err)
	}
	_ = response.Body.Close()
	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("POST campaigns status = %d, want 405", response.StatusCode)
	}
}

func TestCampaignAPIReportsVerifiedNumbergame(t *testing.T) {
	server, summary := testServer(t)
	response, err := http.Get(server.URL + "/api/v1/campaigns/" + string(summary.Campaign))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var view query.CampaignSummary
	if err := json.NewDecoder(response.Body).Decode(&view); err != nil {
		t.Fatalf("decode campaign: %v", err)
	}
	if view.Version != 81 || view.Overview.Completed != 8 || !view.Integrity.Verified {
		t.Fatalf("campaign view = %+v", view)
	}
}

func TestSSEReplaysAfterCursor(t *testing.T) {
	server, summary := testServer(t)
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL+"/api/v1/campaigns/"+string(summary.Campaign)+"/stream?after=80", nil)
	if err != nil {
		t.Fatal(err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	scanner := bufio.NewScanner(response.Body)
	var id, event string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "id: ") {
			id = strings.TrimPrefix(line, "id: ")
		}
		if strings.HasPrefix(line, "event: ") {
			event = strings.TrimPrefix(line, "event: ")
		}
		if line == "" && id != "" {
			break
		}
	}
	if id != "81" || event != "campaign-event" {
		t.Fatalf("SSE replay id/event = %q/%q", id, event)
	}
}
