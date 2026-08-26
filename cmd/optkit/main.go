package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/budget"
	"github.com/go-go-golems/optkit/campaign"
	"github.com/go-go-golems/optkit/examples/numbergame"
	webserver "github.com/go-go-golems/optkit/internal/web"
	"github.com/go-go-golems/optkit/local"
	"github.com/go-go-golems/optkit/projection"
	"github.com/go-go-golems/optkit/query"
	"github.com/go-go-golems/optkit/record"
)

func main() {
	if err := run(context.Background(), os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "optkit:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		writeUsage(stderr)
		return errors.New("a command is required")
	}
	switch args[0] {
	case "demo":
		return runDemo(ctx, args[1:], stdout, stderr)
	case "serve":
		return runServe(ctx, args[1:], stdout, stderr)
	case "campaign":
		if len(args) < 2 {
			writeCampaignUsage(stderr)
			return errors.New("a campaign subcommand is required")
		}
		switch args[1] {
		case "inspect":
			return runCampaignInspect(ctx, args[2:], stdout, stderr)
		case "verify":
			return runCampaignVerify(ctx, args[2:], stdout, stderr)
		default:
			writeCampaignUsage(stderr)
			return fmt.Errorf("unknown campaign subcommand %q", args[1])
		}
	case "artifact":
		if len(args) < 2 || args[1] != "verify" {
			writeArtifactUsage(stderr)
			return errors.New("supported artifact subcommand: verify")
		}
		return runArtifactVerify(ctx, args[2:], stdout, stderr)
	case "help", "-h", "--help":
		writeUsage(stdout)
		return nil
	default:
		writeUsage(stderr)
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func runServe(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("serve", flag.ContinueOnError)
	flags.SetOutput(stderr)
	store := flags.String("store", "./tmp/demo", "local store root")
	listen := flags.String("listen", "127.0.0.1:8080", "HTTP listen address")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("serve does not accept positional arguments")
	}
	profile, err := local.Open(*store)
	if err != nil {
		return err
	}
	defer func() { _ = profile.Close() }()
	handler, err := (webserver.Server{Query: query.Service{Metadata: profile.Metadata, Artifacts: profile.Artifacts}}).Handler()
	if err != nil {
		return err
	}
	server := &http.Server{Addr: *listen, Handler: handler, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()
	_, _ = fmt.Fprintf(stdout, "optkit explorer: http://%s/ (read-only)\n", *listen)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func runDemo(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("demo", flag.ContinueOnError)
	flags.SetOutput(stderr)
	store := flags.String("store", "./tmp/demo", "local SQLite + artifact store root")
	reset := flags.Bool("reset", false, "delete the selected store before running")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("demo does not accept positional arguments")
	}
	summary, err := numbergame.RunDemo(ctx, numbergame.DemoOptions{Root: *store, Reset: *reset})
	if err != nil {
		return err
	}
	return writeJSON(stdout, summary)
}

type inspectEvent struct {
	Seq           uint64             `json:"seq"`
	Kind          campaign.EventKind `json:"kind"`
	Subject       string             `json:"subject,omitempty"`
	OccurredAt    string             `json:"occurred_at"`
	PayloadDigest record.Digest      `json:"payload_digest"`
}

type campaignInspection struct {
	Overview   projection.Overview `json:"overview"`
	Budget     *budget.Snapshot    `json:"budget,omitempty"`
	EventKinds map[string]int      `json:"event_kinds"`
	Tail       []inspectEvent      `json:"tail"`
}

func runCampaignInspect(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("campaign inspect", flag.ContinueOnError)
	flags.SetOutput(stderr)
	store := flags.String("store", "./tmp/demo", "local store root")
	id := flags.String("id", "", "campaign ID")
	tailSize := flags.Int("tail", 10, "number of latest events to include")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("campaign inspect does not accept positional arguments")
	}
	campaignID, err := parseCampaignID(*id)
	if err != nil {
		return err
	}
	if *tailSize < 0 {
		return fmt.Errorf("tail must be non-negative")
	}
	profile, err := local.Open(*store)
	if err != nil {
		return err
	}
	defer func() { _ = profile.Close() }()
	if err := profile.Metadata.Verify(ctx, campaignID); err != nil {
		return err
	}
	events, err := profile.Metadata.Read(ctx, campaignID, 0)
	if err != nil {
		return err
	}
	overview, err := projection.RebuildOverview(campaignID, events)
	if err != nil {
		return err
	}
	var budgetState *budget.Snapshot
	value, err := profile.Metadata.Snapshot(ctx, campaignID)
	if err == nil {
		budgetState = &value
	} else if !errors.Is(err, budget.ErrNotDefined) {
		return err
	}
	kinds := make(map[string]int)
	for _, event := range events {
		kinds[string(event.Kind)]++
	}
	start := len(events) - *tailSize
	if start < 0 {
		start = 0
	}
	tail := make([]inspectEvent, 0, len(events)-start)
	for _, event := range events[start:] {
		tail = append(tail, inspectEvent{
			Seq: event.Seq, Kind: event.Kind, Subject: event.Subject,
			OccurredAt:    event.OccurredAt.UTC().Format("2006-01-02T15:04:05.999999999Z07:00"),
			PayloadDigest: event.Payload.Digest,
		})
	}
	return writeJSON(stdout, campaignInspection{Overview: overview, Budget: budgetState, EventKinds: kinds, Tail: tail})
}

func runCampaignVerify(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("campaign verify", flag.ContinueOnError)
	flags.SetOutput(stderr)
	store := flags.String("store", "./tmp/demo", "local store root")
	id := flags.String("id", "", "campaign ID")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("campaign verify does not accept positional arguments")
	}
	campaignID, err := parseCampaignID(*id)
	if err != nil {
		return err
	}
	profile, err := local.Open(*store)
	if err != nil {
		return err
	}
	defer func() { _ = profile.Close() }()
	if err := profile.Metadata.Verify(ctx, campaignID); err != nil {
		return err
	}
	events, err := profile.Metadata.Read(ctx, campaignID, 0)
	if err != nil {
		return err
	}
	refs := make(map[record.Digest]artifact.Ref)
	for _, event := range events {
		refs[event.Payload.Digest] = event.Payload
	}
	digests := make([]string, 0, len(refs))
	for digest := range refs {
		digests = append(digests, digest.String())
	}
	sort.Strings(digests)
	for _, digest := range digests {
		parsed, err := record.ParseDigest(digest)
		if err != nil {
			return err
		}
		if err := profile.Artifacts.Verify(ctx, refs[parsed]); err != nil {
			return fmt.Errorf("verify event payload %s: %w", digest, err)
		}
	}
	return writeJSON(stdout, map[string]any{
		"campaign": campaignID, "journal": "verified", "events": len(events),
		"unique_event_payloads": len(refs),
	})
}

func runArtifactVerify(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("artifact verify", flag.ContinueOnError)
	flags.SetOutput(stderr)
	store := flags.String("store", "./tmp/demo", "local store root")
	digestRaw := flags.String("digest", "", "sha256:<hex> digest")
	mediaType := flags.String("media-type", "application/vnd.optkit.record+json", "artifact media type")
	schemaRaw := flags.String("schema", "", "optional schema ID")
	size := flags.Int64("size", -1, "exact artifact byte size")
	sensitivityRaw := flags.String("sensitivity", string(artifact.SensitivityInternal), "public|internal|confidential|restricted")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("artifact verify does not accept positional arguments")
	}
	digest, err := record.ParseDigest(*digestRaw)
	if err != nil {
		return err
	}
	if *size < 0 {
		return fmt.Errorf("size is required and must be non-negative")
	}
	sensitivity, err := parseSensitivity(*sensitivityRaw)
	if err != nil {
		return err
	}
	ref := artifact.Ref{Digest: digest, MediaType: *mediaType, Size: *size, Sensitivity: sensitivity}
	if *schemaRaw != "" {
		schema := record.SchemaID(*schemaRaw)
		if err := schema.Validate(); err != nil {
			return err
		}
		ref.Schema = &schema
	}
	if err := ref.Validate(); err != nil {
		return err
	}
	profile, err := local.Open(*store)
	if err != nil {
		return err
	}
	defer func() { _ = profile.Close() }()
	if err := profile.Artifacts.Verify(ctx, ref); err != nil {
		return err
	}
	return writeJSON(stdout, map[string]any{"artifact": ref, "status": "verified"})
}

func parseCampaignID(raw string) (record.CampaignID, error) {
	if strings.TrimSpace(raw) == "" {
		return "", fmt.Errorf("campaign ID is required")
	}
	if err := record.ValidateID("campaign", raw); err != nil {
		return "", err
	}
	return record.CampaignID(raw), nil
}

func parseSensitivity(raw string) (artifact.Sensitivity, error) {
	value := artifact.Sensitivity(raw)
	switch value {
	case artifact.SensitivityPublic, artifact.SensitivityInternal, artifact.SensitivityConfidential, artifact.SensitivityRestricted:
		return value, nil
	default:
		return "", fmt.Errorf("invalid sensitivity %q", raw)
	}
}

func writeJSON(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(value)
}

func writeUsage(writer io.Writer) {
	fmt.Fprintln(writer, `Usage:
  optkit demo [--store PATH] [--reset]
  optkit serve [--store PATH] [--listen 127.0.0.1:8080]
  optkit campaign inspect --store PATH --id CAMPAIGN [--tail N]
  optkit campaign verify --store PATH --id CAMPAIGN
  optkit artifact verify --store PATH --digest DIGEST --size BYTES [options]

All commands emit stable JSON to stdout. The local profile is SQLite in WAL mode
plus a filesystem content-addressed artifact store.`)
}

func writeCampaignUsage(writer io.Writer) {
	fmt.Fprintln(writer, "Usage: optkit campaign <inspect|verify> --store PATH --id CAMPAIGN")
}

func writeArtifactUsage(writer io.Writer) {
	fmt.Fprintln(writer, "Usage: optkit artifact verify --store PATH --digest sha256:<hex> --size BYTES [--media-type TYPE] [--schema ID] [--sensitivity LEVEL]")
}
