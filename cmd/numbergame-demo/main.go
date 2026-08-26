// Command numbergame-demo runs the complete numbergame campaign into a local
// store and exports one JSON bundle for the specialist Experiment Lab UI:
// the candidate proposal (hypothesis, expected improvement, risks), the configurations,
// the cases, every episode's measurements, the estimate, the decision, and
// the journal timeline. The bundle is a teaching record — everything in it
// was reconstructed from the sealed store, never from in-memory state.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/campaign"
	"github.com/go-go-golems/optkit/examples/numbergame"
	"github.com/go-go-golems/optkit/local"
	"github.com/go-go-golems/optkit/measure"
	"github.com/go-go-golems/optkit/space"
)

type bundleCandidate struct {
	ID                  string                    `json:"id"`
	Parent              string                    `json:"parent"`
	Patch               string                    `json:"patch"`
	Child               string                    `json:"child"`
	Proposer            space.Proposer            `json:"proposer"`
	Strategy            string                    `json:"strategy"`
	Hypothesis          string                    `json:"hypothesis"`
	ExpectedImprovement space.ExpectedImprovement `json:"expected_improvement"`
	Risks               []string                  `json:"risks,omitempty"`
	Motivation          space.Motivation          `json:"motivation,omitempty"`
	SemanticCatalogID   string                    `json:"semantic_catalog_id"`
	CreatedAt           time.Time                 `json:"created_at"`
}

type bundleEpisode struct {
	Episode               string  `json:"episode"`
	Arm                   string  `json:"arm"`
	Case                  string  `json:"case"`
	Repeat                int     `json:"repeat"`
	Score                 float64 `json:"score"`
	AbsoluteError         int     `json:"absolute_error"`
	InterventionSatisfied bool    `json:"intervention_satisfied"`
	JudgeLabel            string  `json:"judge_label,omitempty"`
	JudgeRationale        string  `json:"judge_rationale,omitempty"`
}

type bundleEvent struct {
	Seq        uint64    `json:"seq"`
	Kind       string    `json:"kind"`
	Subject    string    `json:"subject,omitempty"`
	Actor      string    `json:"actor,omitempty"`
	OccurredAt time.Time `json:"occurred_at"`
}

type bundle struct {
	Schema           string                 `json:"schema"`
	Campaign         string                 `json:"campaign"`
	Candidate        bundleCandidate        `json:"candidate"`
	BaselineConfig   numbergame.Config      `json:"baseline_config"`
	ChallengerConfig numbergame.Config      `json:"challenger_config"`
	Cases            []bundleCase           `json:"cases"`
	Episodes         []bundleEpisode        `json:"episodes"`
	Estimate         bundleEstimate         `json:"estimate"`
	Decision         numbergame.Decision    `json:"decision"`
	Budget           []bundleBudgetResource `json:"budget"`
	Journal          []bundleEvent          `json:"journal"`
	StoreRoot        string                 `json:"store_root"`
}

type bundleCase struct {
	ID     string `json:"id"`
	Input  int    `json:"input"`
	Target int    `json:"target"`
}

type bundleEstimate struct {
	Value      float64 `json:"value"`
	SampleSize int     `json:"sample_size"`
	Missing    int     `json:"missing"`
}

type bundleBudgetResource struct {
	Resource  string `json:"resource"`
	Limit     int64  `json:"limit"`
	Committed int64  `json:"committed"`
	Available int64  `json:"available"`
}

func main() {
	root := flag.String("root", "", "store root (default: temp dir)")
	out := flag.String("out", "numbergame-demo.json", "bundle output path")
	flag.Parse()
	if err := run(*root, *out); err != nil {
		fmt.Fprintln(os.Stderr, "numbergame-demo:", err)
		os.Exit(1)
	}
}

func run(root, out string) error {
	ctx := context.Background()
	if root == "" {
		dir, err := os.MkdirTemp("", "numbergame-store-")
		if err != nil {
			return err
		}
		root = filepath.Join(dir, "store")
	}
	clock := numbergame.NewSequenceClock(time.Date(2026, 8, 26, 9, 0, 0, 0, time.UTC), time.Millisecond)
	summary, err := numbergame.RunDemo(ctx, numbergame.DemoOptions{Root: root, Reset: true, Clock: clock})
	if err != nil {
		return err
	}

	// Reconstruct everything for the bundle from the sealed store.
	profile, err := local.Open(root)
	if err != nil {
		return err
	}
	defer func() { _ = profile.Close() }()
	events, err := profile.Metadata.Read(ctx, summary.Campaign, 0)
	if err != nil {
		return err
	}

	result := bundle{
		Schema:   "numbergame.lab-bundle/v2",
		Campaign: string(summary.Campaign),
		Decision: summary.Decision,
		Estimate: bundleEstimate{
			Value: summary.Estimate.Value, SampleSize: summary.Estimate.SampleSize, Missing: summary.Estimate.Missing,
		},
		StoreRoot: root,
	}
	for _, resource := range summary.Budget.Resources {
		result.Budget = append(result.Budget, bundleBudgetResource{
			Resource: string(resource.Resource), Limit: resource.Limit,
			Committed: resource.Committed, Available: resource.Available,
		})
	}

	var baselineRecord, challengerRecord space.SnapshotRecord
	for _, event := range events {
		result.Journal = append(result.Journal, bundleEvent{
			Seq: event.Seq, Kind: string(event.Kind), Subject: event.Subject,
			Actor: string(event.Actor), OccurredAt: event.OccurredAt,
		})
		switch event.Kind {
		case campaign.CampaignCreated:
			var spec numbergame.CampaignSpec
			if err := artifact.DecodeJSON(ctx, profile.Artifacts, event.Payload, &spec); err != nil {
				return fmt.Errorf("decode campaign spec: %w", err)
			}
			baselineRecord = spec.Baseline
		case campaign.SnapshotMaterialized:
			if err := artifact.DecodeJSON(ctx, profile.Artifacts, event.Payload, &challengerRecord); err != nil {
				return fmt.Errorf("decode challenger snapshot: %w", err)
			}
		case campaign.CandidateProposed:
			var proposal numbergame.CandidateProposal
			if err := artifact.DecodeJSON(ctx, profile.Artifacts, event.Payload, &proposal); err != nil {
				return fmt.Errorf("decode candidate proposal: %w", err)
			}
			result.Candidate = bundleCandidate{
				ID: string(proposal.Candidate.ID), Parent: string(proposal.Candidate.Parent),
				Patch: string(proposal.Candidate.Patch), Child: string(proposal.Candidate.Child),
				Proposer: proposal.Candidate.Proposer, Strategy: proposal.Candidate.Strategy,
				Hypothesis: proposal.Candidate.Hypothesis, ExpectedImprovement: proposal.Candidate.ExpectedImprovement,
				Risks: proposal.Candidate.Risks, Motivation: proposal.Candidate.Motivation,
				SemanticCatalogID: string(proposal.Candidate.SemanticCatalogID), CreatedAt: proposal.Candidate.CreatedAt,
			}
		case campaign.EpisodeCompleted:
			var completion numbergame.EpisodeCompletion
			if err := artifact.DecodeJSON(ctx, profile.Artifacts, event.Payload, &completion); err != nil {
				return fmt.Errorf("decode episode completion: %w", err)
			}
			row, err := decodeEpisode(ctx, profile, completion)
			if err != nil {
				return err
			}
			result.Episodes = append(result.Episodes, row)
		}
	}

	codec := numbergame.ConfigCodec()
	baseline, err := space.LoadSnapshot(ctx, profile.Artifacts, baselineRecord, codec)
	if err != nil {
		return fmt.Errorf("load baseline config: %w", err)
	}
	challenger, err := space.LoadSnapshot(ctx, profile.Artifacts, challengerRecord, codec)
	if err != nil {
		return fmt.Errorf("load challenger config: %w", err)
	}
	result.BaselineConfig, result.ChallengerConfig = baseline.Value, challenger.Value
	result.Cases = demoCases()

	payload, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(out, append(payload, '\n'), 0o644); err != nil {
		return err
	}
	fmt.Printf("wrote %s (campaign %s, %d episodes, decision %s)\n", out, result.Campaign, len(result.Episodes), result.Decision.Status)
	return nil
}

func decodeEpisode(ctx context.Context, profile *local.Profile, completion numbergame.EpisodeCompletion) (bundleEpisode, error) {
	row := bundleEpisode{
		Episode: string(completion.Episode), Arm: completion.Arm,
		Case: completion.Case, Repeat: completion.Repeat,
	}
	var judge measure.Observation
	if err := artifact.DecodeJSON(ctx, profile.Artifacts, completion.Judge, &judge); err != nil {
		return bundleEpisode{}, fmt.Errorf("decode judge observation: %w", err)
	}
	if judge.Value.Number != nil {
		score, err := strconv.ParseFloat(*judge.Value.Number, 64)
		if err != nil {
			return bundleEpisode{}, err
		}
		row.Score = score
	}
	var assessment numbergame.Assessment
	if err := artifact.DecodeJSON(ctx, profile.Artifacts, completion.Assessment, &assessment); err == nil {
		row.JudgeLabel, row.JudgeRationale = assessment.Label, assessment.Rationale
	}
	var absErr measure.Observation
	if err := artifact.DecodeJSON(ctx, profile.Artifacts, completion.AbsoluteError, &absErr); err != nil {
		return bundleEpisode{}, fmt.Errorf("decode error observation: %w", err)
	}
	if absErr.Value.Number != nil {
		value, err := strconv.Atoi(*absErr.Value.Number)
		if err != nil {
			return bundleEpisode{}, err
		}
		row.AbsoluteError = value
	}
	var intervention measure.Observation
	if err := artifact.DecodeJSON(ctx, profile.Artifacts, completion.Intervention, &intervention); err != nil {
		return bundleEpisode{}, fmt.Errorf("decode intervention observation: %w", err)
	}
	row.InterventionSatisfied = intervention.Status == measure.StatusSatisfied
	return row, nil
}

func demoCases() []bundleCase {
	return []bundleCase{
		{ID: "case-01", Input: 1, Target: 3},
		{ID: "case-02", Input: 2, Target: 6},
		{ID: "case-03", Input: 4, Target: 12},
		{ID: "case-04", Input: 7, Target: 21},
	}
}
