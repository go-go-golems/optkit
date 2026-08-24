//go:build cgo

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/optkit/examples/numbergame"
)

func TestDemoAndCampaignCommands(t *testing.T) {
	root := filepath.Join(t.TempDir(), "store")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := run(context.Background(), []string{"demo", "--store", root, "--reset"}, &stdout, &stderr); err != nil {
		t.Fatalf("demo failed: %v\nstderr: %s", err, stderr.String())
	}
	var summary numbergame.DemoSummary
	if err := json.Unmarshal(stdout.Bytes(), &summary); err != nil {
		t.Fatal(err)
	}
	if summary.Campaign == "" || summary.Decision.Status != "eligible" {
		t.Fatalf("unexpected summary: %+v", summary)
	}

	stdout.Reset()
	stderr.Reset()
	if err := run(context.Background(), []string{"campaign", "inspect", "--store", root, "--id", string(summary.Campaign), "--tail", "3"}, &stdout, &stderr); err != nil {
		t.Fatalf("inspect failed: %v\nstderr: %s", err, stderr.String())
	}
	var inspection campaignInspection
	if err := json.Unmarshal(stdout.Bytes(), &inspection); err != nil {
		t.Fatal(err)
	}
	if inspection.Overview.Completed != 8 || inspection.Budget == nil || inspection.Budget.Violated || len(inspection.Tail) != 3 {
		t.Fatalf("unexpected inspection: %+v", inspection)
	}

	stdout.Reset()
	stderr.Reset()
	if err := run(context.Background(), []string{"campaign", "verify", "--store", root, "--id", string(summary.Campaign)}, &stdout, &stderr); err != nil {
		t.Fatalf("verify failed: %v\nstderr: %s", err, stderr.String())
	}
}
