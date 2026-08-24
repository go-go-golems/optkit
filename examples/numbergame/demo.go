package numbergame

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/budget"
	"github.com/go-go-golems/optkit/campaign"
	"github.com/go-go-golems/optkit/episode"
	"github.com/go-go-golems/optkit/experiment"
	"github.com/go-go-golems/optkit/local"
	"github.com/go-go-golems/optkit/measure"
	"github.com/go-go-golems/optkit/projection"
	"github.com/go-go-golems/optkit/record"
	"github.com/go-go-golems/optkit/scheduler"
	"github.com/go-go-golems/optkit/space"
)

type DemoOptions struct {
	Root  string
	Reset bool
	Clock episode.Clock
}

type DemoSummary struct {
	Campaign         record.CampaignID   `json:"campaign"`
	Candidate        record.CandidateID  `json:"candidate"`
	Baseline         record.SnapshotID   `json:"baseline"`
	Challenger       record.SnapshotID   `json:"challenger"`
	Patch            record.PatchID      `json:"patch"`
	Trial            record.TrialID      `json:"trial"`
	Estimate         experiment.Estimate `json:"estimate"`
	Decision         Decision            `json:"decision"`
	Overview         projection.Overview `json:"overview"`
	Budget           budget.Snapshot     `json:"budget"`
	SampleTrajectory artifact.Ref        `json:"sample_trajectory"`
	StoreRoot        string              `json:"store_root"`
	DatabasePath     string              `json:"database_path"`
	ArtifactsPath    string              `json:"artifacts_path"`
}

type CampaignSpec struct {
	System     record.SystemID      `json:"system"`
	Baseline   space.SnapshotRecord `json:"baseline"`
	SearchVars []space.VariableID   `json:"search_variables"`
	Dataset    string               `json:"dataset"`
	Trial      record.TrialID       `json:"trial"`
	Budget     []budget.Limit       `json:"budget"`
}

type CandidateProposal struct {
	Candidate space.Candidate   `json:"candidate"`
	Patch     space.PatchRecord `json:"patch"`
}

type EpisodeWork struct {
	Spec     experiment.EpisodeSpec `json:"spec"`
	Snapshot space.SnapshotRecord   `json:"snapshot"`
}

type EpisodeCompletion struct {
	Episode       record.EpisodeID `json:"episode"`
	Arm           string           `json:"arm"`
	Case          string           `json:"case"`
	Repeat        int              `json:"repeat"`
	Result        artifact.Ref     `json:"result"`
	Intervention  artifact.Ref     `json:"intervention"`
	AbsoluteError artifact.Ref     `json:"absolute_error"`
	Judge         artifact.Ref     `json:"judge"`
	Assessment    artifact.Ref     `json:"assessment"`
}

type Decision struct {
	Policy                    string  `json:"policy"`
	Status                    string  `json:"status"`
	PairedAccuracyDelta       float64 `json:"paired_accuracy_delta"`
	AllInterventionsExercised bool    `json:"all_interventions_exercised"`
	BudgetViolated            bool    `json:"budget_violated"`
	Reason                    string  `json:"reason"`
}

func RunDemo(ctx context.Context, options DemoOptions) (DemoSummary, error) {
	if options.Root == "" {
		return DemoSummary{}, fmt.Errorf("demo store root is required")
	}
	if options.Reset {
		if err := resetRoot(options.Root); err != nil {
			return DemoSummary{}, err
		}
	}
	clock := options.Clock
	if clock == nil {
		clock = episode.RealClock{}
	}
	profile, err := local.OpenWithClock(options.Root, clock.Now)
	if err != nil {
		return DemoSummary{}, err
	}
	defer func() {
		if profile != nil {
			_ = profile.Close()
		}
	}()

	configCodec := ConfigCodec()
	baseline, err := space.MaterializeSnapshot(ctx, profile.Artifacts, SystemID, configCodec, Config{Multiplier: 2, Noise: NoiseNone})
	if err != nil {
		return DemoSummary{}, err
	}
	patchBuilder := space.NewPatchBuilder(baseline, profile.Artifacts, configCodec)
	if err := space.Set(patchBuilder, MultiplierVariable(), 3); err != nil {
		return DemoSummary{}, err
	}
	patch, challenger, err := patchBuilder.Build(ctx)
	if err != nil {
		return DemoSummary{}, err
	}
	candidate, err := space.NewCandidate(
		baseline.ID, patch.ID, challenger.ID, "actor:demo", "manual-coordinate/v1",
		"Set the multiplier to three so outputs match the target relation on all development cases.",
		[]string{"target.accuracy"}, []string{"a fixed multiplier may overfit other target relations"}, clock.Now(),
	)
	if err != nil {
		return DemoSummary{}, err
	}

	cases := []Case{{Input: 1, Target: 3}, {Input: 2, Target: 6}, {Input: 4, Target: 12}, {Input: 7, Target: 21}}
	experimentCases := make([]experiment.Case, 0, len(cases))
	for index, value := range cases {
		ref, err := artifact.PutCanonical(ctx, profile.Artifacts, "schema:numbergame.case/v1", artifact.SensitivityInternal, value)
		if err != nil {
			return DemoSummary{}, err
		}
		experimentCases = append(experimentCases, experiment.Case{
			ID: fmt.Sprintf("case-%02d", index+1), Input: ref,
			Groups: []string{"development", "exact-target"},
		})
	}
	dataset, err := experiment.NewDatasetManifest("development", experimentCases)
	if err != nil {
		return DemoSummary{}, err
	}
	trial, err := experiment.NewCompleteBlockTrial([]experiment.Arm{
		{ID: "baseline", Snapshot: baseline.ID}, {ID: "challenger", Snapshot: challenger.ID},
	}, dataset, 1, "numbergame.execution/v1")
	if err != nil {
		return DemoSummary{}, err
	}

	rawCampaign, err := record.NewID("campaign")
	if err != nil {
		return DemoSummary{}, err
	}
	campaignID := record.CampaignID(rawCampaign)
	controller := campaign.Controller{Journal: profile.Metadata, Artifacts: profile.Artifacts}
	budgetLimits := []budget.Limit{
		{Resource: "episodes", Units: int64(len(experimentCases) * 2)},
		{Resource: "numbergame.operations", Units: int64(len(experimentCases) * 2)},
	}
	spec := CampaignSpec{
		System: SystemID, Baseline: baseline.SnapshotRecord,
		SearchVars: []space.VariableID{MultiplierVariable().Descriptor.ID},
		Dataset:    dataset.ID, Trial: trial.ID, Budget: budgetLimits,
	}
	if _, err := handleLifecycle(ctx, controller, profile.Metadata, campaignID, campaign.CreateCampaign, "schema:numbergame.campaign-spec/v1", spec, "actor:demo"); err != nil {
		return DemoSummary{}, err
	}
	if err := profile.Metadata.Define(ctx, campaignID, budgetLimits); err != nil {
		return DemoSummary{}, err
	}
	if _, err := handleLifecycle(ctx, controller, profile.Metadata, campaignID, campaign.CompilePlan, "schema:numbergame.trial-plan/v1", trial, "actor:demo"); err != nil {
		return DemoSummary{}, err
	}
	if _, err := handleLifecycle(ctx, controller, profile.Metadata, campaignID, campaign.StartCampaign, "schema:numbergame.lifecycle/v1", map[string]string{"reason": "demo start"}, "actor:demo"); err != nil {
		return DemoSummary{}, err
	}

	candidateRef, err := artifact.PutCanonical(ctx, profile.Artifacts, "schema:numbergame.candidate-proposal/v1", artifact.SensitivityInternal, CandidateProposal{
		Candidate: candidate, Patch: patch.PatchRecord,
	})
	if err != nil {
		return DemoSummary{}, err
	}
	if _, err := appendArtifactFact(ctx, profile, campaignID, campaign.CandidateProposed, "schema:numbergame.candidate-proposal/v1", string(candidate.ID), candidateRef, nil); err != nil {
		return DemoSummary{}, err
	}
	childRef, err := artifact.PutCanonical(ctx, profile.Artifacts, "schema:optkit.snapshot/v1", artifact.SensitivityInternal, challenger.SnapshotRecord)
	if err != nil {
		return DemoSummary{}, err
	}
	if _, err := appendArtifactFact(ctx, profile, campaignID, campaign.SnapshotMaterialized, "schema:optkit.snapshot/v1", string(challenger.ID), childRef, nil); err != nil {
		return DemoSummary{}, err
	}
	trialRef, err := artifact.PutCanonical(ctx, profile.Artifacts, "schema:numbergame.trial-plan/v1", artifact.SensitivityInternal, trial)
	if err != nil {
		return DemoSummary{}, err
	}
	if _, err := appendArtifactFact(ctx, profile, campaignID, campaign.TrialPlanned, "schema:numbergame.trial-plan/v1", string(trial.ID), trialRef, nil); err != nil {
		return DemoSummary{}, err
	}

	specs, err := experiment.Expand(trial)
	if err != nil {
		return DemoSummary{}, err
	}
	snapshots := map[record.SnapshotID]space.SnapshotRecord{
		baseline.ID: baseline.SnapshotRecord, challenger.ID: challenger.SnapshotRecord,
	}
	workItems := make([]scheduler.WorkItem, 0, len(specs))
	reservations := make([]budget.Reservation, 0, len(specs))
	for _, specValue := range specs {
		work := EpisodeWork{Spec: specValue, Snapshot: snapshots[specValue.Arm.Snapshot]}
		workRef, err := artifact.PutCanonical(ctx, profile.Artifacts, "schema:numbergame.episode-work/v1", artifact.SensitivityInternal, work)
		if err != nil {
			return DemoSummary{}, err
		}
		claims := []budget.Quantity{{Resource: "episodes", Units: 1}, {Resource: "numbergame.operations", Units: 1}}
		item, err := scheduler.NewWorkItem(campaignID, "episode", specValue.SemanticKey, workRef, 0, clock.Now(), time.Minute, []scheduler.ResourceClaim{
			{Resource: "episodes", Units: 1}, {Resource: "numbergame.operations", Units: 1},
		})
		if err != nil {
			return DemoSummary{}, err
		}
		reservation, err := profile.Metadata.Reserve(ctx, budget.ReserveRequest{Campaign: campaignID, Work: item.ID, Claims: claims})
		if err != nil {
			return DemoSummary{}, err
		}
		reservationRef, err := artifact.PutCanonical(ctx, profile.Artifacts, "schema:optkit.budget-reservation/v1", artifact.SensitivityInternal, reservation)
		if err != nil {
			return DemoSummary{}, err
		}
		if _, err := appendArtifactFact(ctx, profile, campaignID, campaign.BudgetReserved, "schema:optkit.budget-reservation/v1", string(item.ID), reservationRef, nil); err != nil {
			return DemoSummary{}, err
		}
		if _, err := appendArtifactFact(ctx, profile, campaignID, campaign.EpisodeScheduled, "schema:numbergame.episode-work/v1", string(specValue.ID), workRef, nil); err != nil {
			return DemoSummary{}, err
		}
		workItems = append(workItems, item)
		reservations = append(reservations, reservation)
	}
	if err := profile.Metadata.Enqueue(ctx, workItems); err != nil {
		for _, reservation := range reservations {
			_, _ = releaseBudgetReservation(ctx, profile, campaignID, reservation.ID)
		}
		return DemoSummary{}, err
	}

	for {
		leases, err := profile.Metadata.Lease(ctx, "numbergame-worker", scheduler.LeaseRequest{Kinds: []scheduler.WorkKind{"episode"}, Limit: 2, Now: clock.Now()})
		if err != nil {
			return DemoSummary{}, err
		}
		if len(leases) == 0 {
			break
		}
		for _, lease := range leases {
			if err := executeEpisodeLease(ctx, profile, campaignID, lease, clock); err != nil {
				return DemoSummary{}, err
			}
		}
	}
	if err := profile.Metadata.Verify(ctx, campaignID); err != nil {
		return DemoSummary{}, err
	}

	// Simulate controller and worker process restart before analysis. All later
	// decisions are reconstructed from the SQLite journal and sealed artifacts.
	if err := profile.Close(); err != nil {
		return DemoSummary{}, err
	}
	profile = nil
	profile, err = local.OpenWithClock(options.Root, clock.Now)
	if err != nil {
		return DemoSummary{}, err
	}
	controller = campaign.Controller{Journal: profile.Metadata, Artifacts: profile.Artifacts}
	if err := profile.Metadata.Verify(ctx, campaignID); err != nil {
		return DemoSummary{}, err
	}

	events, err := profile.Metadata.Read(ctx, campaignID, 0)
	if err != nil {
		return DemoSummary{}, err
	}
	rows := make([]experiment.NumericObservation, 0, len(specs))
	allInterventions := true
	var sampleTrajectory artifact.Ref
	for _, event := range events {
		if event.Kind != campaign.EpisodeCompleted {
			continue
		}
		var completion EpisodeCompletion
		if err := artifact.DecodeJSON(ctx, profile.Artifacts, event.Payload, &completion); err != nil {
			return DemoSummary{}, err
		}
		var result episode.Result
		if err := artifact.DecodeJSON(ctx, profile.Artifacts, completion.Result, &result); err != nil {
			return DemoSummary{}, err
		}
		if sampleTrajectory.Digest == "" {
			sampleTrajectory = result.Trajectory
		}
		var judge measure.Observation
		if err := artifact.DecodeJSON(ctx, profile.Artifacts, completion.Judge, &judge); err != nil {
			return DemoSummary{}, err
		}
		score, err := parseScore(judge)
		if err != nil {
			return DemoSummary{}, err
		}
		var intervention measure.Observation
		if err := artifact.DecodeJSON(ctx, profile.Artifacts, completion.Intervention, &intervention); err != nil {
			return DemoSummary{}, err
		}
		if intervention.Status != measure.StatusSatisfied {
			allInterventions = false
		}
		rows = append(rows, experiment.NumericObservation{
			CaseID: completion.Case, ArmID: completion.Arm, Repeat: completion.Repeat,
			Value: score, Valid: judge.Status == measure.StatusMeasured,
		})
	}
	if len(rows) != len(specs) {
		return DemoSummary{}, fmt.Errorf("analysis found %d episode observations, expected %d", len(rows), len(specs))
	}
	estimate, err := experiment.PairedMean(rows, "baseline", "challenger")
	if err != nil {
		return DemoSummary{}, err
	}
	estimateRef, err := artifact.PutCanonical(ctx, profile.Artifacts, "schema:optkit.estimate/v1", artifact.SensitivityInternal, estimate)
	if err != nil {
		return DemoSummary{}, err
	}
	if _, err := appendArtifactFact(ctx, profile, campaignID, campaign.EstimateRecorded, "schema:optkit.estimate/v1", string(trial.ID), estimateRef, nil); err != nil {
		return DemoSummary{}, err
	}
	budgetState, err := profile.Metadata.Snapshot(ctx, campaignID)
	if err != nil {
		return DemoSummary{}, err
	}

	decision := Decision{
		Policy: "numbergame.lexicographic/v1", PairedAccuracyDelta: estimate.Value,
		AllInterventionsExercised: allInterventions, BudgetViolated: budgetState.Violated,
	}
	switch {
	case budgetState.Violated:
		decision.Status = "invalid"
		decision.Reason = "Committed resource use violated the campaign budget."
	case !allInterventions:
		decision.Status = "invalid"
		decision.Reason = "At least one configured multiplier was not observed in its episode trajectory."
	case estimate.Value <= 0:
		decision.Status = "ineligible"
		decision.Reason = "The challenger did not improve paired target accuracy."
	default:
		decision.Status = "eligible"
		decision.Reason = "Every intervention was exercised and paired target accuracy improved."
	}
	decisionRef, err := artifact.PutCanonical(ctx, profile.Artifacts, "schema:numbergame.decision/v1", artifact.SensitivityInternal, decision)
	if err != nil {
		return DemoSummary{}, err
	}
	if _, err := appendArtifactFact(ctx, profile, campaignID, campaign.DecisionRecorded, "schema:numbergame.decision/v1", string(candidate.ID), decisionRef, nil); err != nil {
		return DemoSummary{}, err
	}
	if _, err := handleLifecycle(ctx, controller, profile.Metadata, campaignID, campaign.CompleteCampaign, "schema:numbergame.lifecycle/v1", map[string]any{
		"decision": decision, "estimate": estimate, "budget": budgetState,
	}, "actor:demo"); err != nil {
		return DemoSummary{}, err
	}
	if err := profile.Metadata.Verify(ctx, campaignID); err != nil {
		return DemoSummary{}, err
	}

	// Reopen once more after terminal commit to prove the summary is a replayed
	// projection, not an in-memory object retained by RunDemo.
	if err := profile.Close(); err != nil {
		return DemoSummary{}, err
	}
	profile = nil
	profile, err = local.OpenWithClock(options.Root, clock.Now)
	if err != nil {
		return DemoSummary{}, err
	}
	if err := profile.Metadata.Verify(ctx, campaignID); err != nil {
		return DemoSummary{}, err
	}
	finalEvents, err := profile.Metadata.Read(ctx, campaignID, 0)
	if err != nil {
		return DemoSummary{}, err
	}
	overview, err := projection.RebuildOverview(campaignID, finalEvents)
	if err != nil {
		return DemoSummary{}, err
	}
	if overview.Status != campaign.StatusCompleted || overview.Completed != len(specs) {
		return DemoSummary{}, fmt.Errorf("unexpected terminal overview: %+v", overview)
	}
	budgetState, err = profile.Metadata.Snapshot(ctx, campaignID)
	if err != nil {
		return DemoSummary{}, err
	}

	return DemoSummary{
		Campaign: campaignID, Candidate: candidate.ID, Baseline: baseline.ID,
		Challenger: challenger.ID, Patch: patch.ID, Trial: trial.ID,
		Estimate: estimate, Decision: decision, Overview: overview, Budget: budgetState,
		SampleTrajectory: sampleTrajectory,
		StoreRoot:        profile.Root, DatabasePath: filepath.Join(profile.Root, "optkit.db"),
		ArtifactsPath: filepath.Join(profile.Root, "artifacts"),
	}, nil
}

func executeEpisodeLease(ctx context.Context, profile *local.Profile, campaignID record.CampaignID, lease scheduler.Lease, clock episode.Clock) error {
	var work EpisodeWork
	if err := artifact.DecodeJSON(ctx, profile.Artifacts, lease.Item.Payload, &work); err != nil {
		return err
	}
	subject := string(work.Spec.ID)

	leaseRef, err := artifact.PutCanonical(ctx, profile.Artifacts, "schema:optkit.work-lease/v1", artifact.SensitivityInternal, lease)
	if err != nil {
		return err
	}
	if _, err := appendArtifactFact(ctx, profile, campaignID, campaign.EpisodeLeaseGranted, "schema:optkit.work-lease/v1", subject, leaseRef, nil); err != nil {
		return err
	}
	attemptRef, err := artifact.PutCanonical(ctx, profile.Artifacts, "schema:numbergame.episode-attempt/v1", artifact.SensitivityInternal, map[string]any{
		"episode": work.Spec.ID, "lease": lease.ID, "worker": lease.Worker, "attempt": lease.Attempt,
	})
	if err != nil {
		return err
	}
	if _, err := appendArtifactFact(ctx, profile, campaignID, campaign.EpisodeAttemptStarted, "schema:numbergame.episode-attempt/v1", subject, attemptRef, nil); err != nil {
		return err
	}

	snapshot, err := space.LoadSnapshot(ctx, profile.Artifacts, work.Snapshot, ConfigCodec())
	if err != nil {
		return failEpisodeLease(ctx, profile, campaignID, lease, work.Spec.ID, "snapshot_load", err)
	}
	var caseValue Case
	if err := artifact.DecodeJSON(ctx, profile.Artifacts, work.Spec.Case.Input, &caseValue); err != nil {
		return failEpisodeLease(ctx, profile, campaignID, lease, work.Spec.ID, "case_decode", err)
	}
	writer, err := episode.NewWriter(work.Spec.ID, profile.Artifacts, clock)
	if err != nil {
		return failEpisodeLease(ctx, profile, campaignID, lease, work.Spec.ID, "trajectory_writer", err)
	}
	run, err := (Executable{Config: snapshot.Value, Clock: clock}).Run(ctx, Invocation{Case: caseValue, Seed: work.Spec.Seed}, writer)
	if err != nil {
		return failEpisodeLease(ctx, profile, campaignID, lease, work.Spec.ID, "system_run", err)
	}
	trajectoryRef, err := writer.Seal(ctx)
	if err != nil {
		return failEpisodeLease(ctx, profile, campaignID, lease, work.Spec.ID, "trajectory_seal", err)
	}
	result := run.WithTrajectory(trajectoryRef)
	resultRef, err := artifact.PutCanonical(ctx, profile.Artifacts, "schema:optkit.episode-result/v1", artifact.SensitivityInternal, result)
	if err != nil {
		return failEpisodeLease(ctx, profile, campaignID, lease, work.Spec.ID, "result_store", err)
	}
	measurements, err := Measure(ctx, profile.Artifacts, result, snapshot.Value, work.Spec.Repeat, clock.Now())
	if err != nil {
		return failEpisodeLease(ctx, profile, campaignID, lease, work.Spec.ID, "measurement", err)
	}
	interventionRef, err := artifact.PutCanonical(ctx, profile.Artifacts, "schema:optkit.observation/v1", artifact.SensitivityInternal, measurements.Intervention)
	if err != nil {
		return err
	}
	absoluteErrorRef, err := artifact.PutCanonical(ctx, profile.Artifacts, "schema:optkit.observation/v1", artifact.SensitivityInternal, measurements.AbsoluteError)
	if err != nil {
		return err
	}
	judgeRef, err := artifact.PutCanonical(ctx, profile.Artifacts, "schema:optkit.observation/v1", artifact.SensitivityInternal, measurements.Judge)
	if err != nil {
		return err
	}
	completion := EpisodeCompletion{
		Episode: work.Spec.ID, Arm: work.Spec.Arm.ID, Case: work.Spec.Case.ID, Repeat: work.Spec.Repeat,
		Result: resultRef, Intervention: interventionRef, AbsoluteError: absoluteErrorRef,
		Judge: judgeRef, Assessment: measurements.Assessment,
	}
	completionRef, err := artifact.PutCanonical(ctx, profile.Artifacts, "schema:numbergame.episode-completion/v1", artifact.SensitivityInternal, completion)
	if err != nil {
		return err
	}
	if err := profile.Metadata.Complete(ctx, lease.ID, scheduler.WorkResult{Artifact: completionRef}); err != nil {
		return err
	}
	reservation, err := profile.Metadata.ReservationForWork(ctx, lease.Item.ID)
	if err != nil {
		return err
	}
	actual, err := budgetQuantities(result.Usage)
	if err != nil {
		return err
	}
	committed, err := profile.Metadata.Commit(ctx, reservation.ID, actual)
	if err != nil {
		return err
	}
	usageRef, err := artifact.PutCanonical(ctx, profile.Artifacts, "schema:optkit.budget-usage/v1", artifact.SensitivityInternal, committed)
	if err != nil {
		return err
	}
	if _, err := appendArtifactFact(ctx, profile, campaignID, campaign.UsageCommitted, "schema:optkit.budget-usage/v1", string(lease.Item.ID), usageRef, nil); err != nil {
		return err
	}
	if _, err := appendArtifactFact(ctx, profile, campaignID, campaign.EpisodeCompleted, "schema:numbergame.episode-completion/v1", subject, completionRef, nil); err != nil {
		return err
	}
	for _, observation := range []struct {
		value measure.Observation
		ref   artifact.Ref
	}{
		{measurements.Intervention, interventionRef},
		{measurements.AbsoluteError, absoluteErrorRef},
		{measurements.Judge, judgeRef},
	} {
		if _, err := appendArtifactFact(ctx, profile, campaignID, campaign.ObservationRecorded, "schema:optkit.observation/v1", string(observation.value.ID), observation.ref, nil); err != nil {
			return err
		}
	}
	return nil
}

func failEpisodeLease(ctx context.Context, profile *local.Profile, campaignID record.CampaignID, lease scheduler.Lease, episodeID record.EpisodeID, code string, cause error) error {
	failure := scheduler.WorkFailure{Code: code, Message: cause.Error(), Retryable: false}
	failureRef, storeErr := artifact.PutCanonical(ctx, profile.Artifacts, "schema:optkit.work-failure/v1", artifact.SensitivityInternal, failure)
	if storeErr != nil {
		return fmt.Errorf("%s: %w (store failure evidence: %v)", code, cause, storeErr)
	}
	if err := profile.Metadata.Fail(ctx, lease.ID, failure); err != nil {
		return fmt.Errorf("%s: %w (commit work failure: %v)", code, cause, err)
	}
	if reservation, err := profile.Metadata.ReservationForWork(ctx, lease.Item.ID); err == nil && reservation.Status == budget.StatusReserved {
		if _, releaseErr := releaseBudgetReservation(ctx, profile, campaignID, reservation.ID); releaseErr != nil {
			return fmt.Errorf("%s: %w (release budget reservation: %v)", code, cause, releaseErr)
		}
	} else if err != nil && !errors.Is(err, budget.ErrReservationNotFound) {
		return fmt.Errorf("%s: %w (lookup budget reservation: %v)", code, cause, err)
	}
	if _, err := appendArtifactFact(ctx, profile, campaignID, campaign.EpisodeFailed, "schema:optkit.work-failure/v1", string(episodeID), failureRef, map[string]string{"retryable": "false"}); err != nil {
		return fmt.Errorf("%s: %w (append episode failure: %v)", code, cause, err)
	}
	return fmt.Errorf("episode %s failed at %s: %w", episodeID, code, cause)
}

func releaseBudgetReservation(ctx context.Context, profile *local.Profile, campaignID record.CampaignID, reservationID record.ReservationID) (budget.Reservation, error) {
	released, err := profile.Metadata.Release(ctx, reservationID)
	if err != nil {
		return budget.Reservation{}, err
	}
	releasedRef, err := artifact.PutCanonical(ctx, profile.Artifacts, "schema:optkit.budget-release/v1", artifact.SensitivityInternal, released)
	if err != nil {
		return budget.Reservation{}, err
	}
	if _, err := appendArtifactFact(ctx, profile, campaignID, campaign.BudgetReleased, "schema:optkit.budget-release/v1", string(released.Work), releasedRef, nil); err != nil {
		return budget.Reservation{}, err
	}
	return released, nil
}

func budgetQuantities(usage []episode.ResourceUsage) ([]budget.Quantity, error) {
	totals := make(map[budget.Resource]int64)
	for _, value := range usage {
		resource := budget.Resource(value.Resource)
		if err := resource.Validate(); err != nil {
			return nil, err
		}
		if value.Units < 0 {
			return nil, fmt.Errorf("negative usage for %s", value.Resource)
		}
		totals[resource] += value.Units
	}
	values := make([]budget.Quantity, 0, len(totals))
	for resource, units := range totals {
		values = append(values, budget.Quantity{Resource: resource, Units: units})
	}
	return budget.NormalizeQuantities(values, true)
}

func handleLifecycle(ctx context.Context, controller campaign.Controller, journal campaign.Journal, campaignID record.CampaignID, kind campaign.CommandKind, schema record.SchemaID, payload any, actor record.ActorRef) (campaign.AppendResult, error) {
	head, err := journal.Head(ctx, campaignID)
	if err != nil {
		return campaign.AppendResult{}, err
	}
	rawCommand, err := record.NewID("command")
	if err != nil {
		return campaign.AppendResult{}, err
	}
	return controller.Handle(ctx, campaign.Command{
		ID: record.CommandID(rawCommand), Campaign: campaignID, ExpectedVersion: head.Version,
		Kind: kind, Actor: actor, Schema: schema, Payload: payload,
	})
}

func appendArtifactFact(ctx context.Context, profile *local.Profile, campaignID record.CampaignID, kind campaign.EventKind, schema record.SchemaID, subject string, payload artifact.Ref, tags map[string]string) (campaign.ControlEvent, error) {
	events, err := profile.Metadata.Read(ctx, campaignID, 0)
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	state, err := campaign.Fold(campaignID, events)
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	if err := campaign.ValidateTransition(state, kind, subject, tags); err != nil {
		return campaign.ControlEvent{}, err
	}
	head, err := profile.Metadata.Head(ctx, campaignID)
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	result, err := profile.Metadata.Append(ctx, campaignID, head.Version, []campaign.NewEvent{
		{Kind: kind, Schema: schema, Subject: subject, Actor: "actor:demo", Payload: payload, Tags: tags},
	})
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	return result.Events[0], nil
}

func resetRoot(root string) error {
	abs, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	if abs == string(filepath.Separator) || abs == "." || len(abs) < 4 {
		return fmt.Errorf("refusing to reset unsafe path %q", abs)
	}
	if err := os.RemoveAll(abs); err != nil {
		return fmt.Errorf("reset demo store: %w", err)
	}
	return nil
}

func parseScore(observation measure.Observation) (float64, error) {
	if observation.Value.Number == nil {
		return 0, fmt.Errorf("observation %s has no numeric value", observation.ID)
	}
	value, err := strconv.ParseFloat(*observation.Value.Number, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}
