package campaign

import (
	"fmt"
	"strconv"

	"github.com/go-go-golems/optkit/record"
)

func Fold(campaignID record.CampaignID, events []ControlEvent) (State, error) {
	state := Initial(campaignID)
	for _, event := range events {
		var err error
		state, err = Apply(state, event)
		if err != nil {
			return State{}, err
		}
	}
	return state, nil
}

func Apply(state State, event ControlEvent) (State, error) {
	if event.Campaign != state.Campaign {
		return State{}, fmt.Errorf("event campaign %s does not match state campaign %s", event.Campaign, state.Campaign)
	}
	if event.Seq != state.Version+1 {
		return State{}, fmt.Errorf("event sequence %d does not follow state version %d", event.Seq, state.Version)
	}
	if event.PreviousDigest != state.LastDigest {
		return State{}, fmt.Errorf("event previous digest %s does not match state digest %s", event.PreviousDigest, state.LastDigest)
	}
	if err := VerifyEventDigest(event); err != nil {
		return State{}, err
	}
	if err := ValidateTransition(state, event.Kind, event.Subject, event.Tags); err != nil {
		return State{}, err
	}

	next := state
	next.Episodes = cloneEpisodes(state.Episodes)
	switch event.Kind {
	case CampaignCreated:
		next.Status = StatusDraft
	case PlanCompiled:
		next.Status = StatusReady
	case CampaignStarted, CampaignResumed:
		next.Status = StatusRunning
	case CampaignPaused:
		next.Status = StatusPaused
	case CampaignStopping:
		next.Status = StatusStopping
	case CampaignStopped:
		next.Status = StatusStopped
	case CampaignFailed:
		next.Status = StatusFailed
	case CampaignCompleted:
		next.Status = StatusCompleted
	case EpisodeScheduled:
		next.Episodes[event.Subject] = EpisodeQueued
	case EpisodeLeaseGranted:
		next.Episodes[event.Subject] = EpisodeLeased
	case EpisodeAttemptStarted:
		next.Episodes[event.Subject] = EpisodeRunning
	case EpisodeCompleted:
		next.Episodes[event.Subject] = EpisodeCompletedState
	case EpisodeFailed:
		retryable, _ := strconv.ParseBool(event.Tags["retryable"])
		if retryable {
			next.Episodes[event.Subject] = EpisodeQueued
		} else {
			next.Episodes[event.Subject] = EpisodeFailedTerminal
		}
	case CandidateProposed, SnapshotMaterialized, TrialPlanned,
		ObservationRecorded, EstimateRecorded, DecisionRecorded,
		BudgetReserved, UsageCommitted, BudgetReleased:
		// These evidence events advance journal custody without changing the
		// lifecycle or per-episode status projection.
	}
	next.Version = event.Seq
	next.LastDigest = event.Digest
	next.Events++
	return next, nil
}

func ValidateTransition(state State, kind EventKind, subject string, tags map[string]string) error {
	requireStatus := func(allowed ...Status) error {
		for _, status := range allowed {
			if state.Status == status {
				return nil
			}
		}
		return fmt.Errorf("event %s is invalid while campaign is %s", kind, state.Status)
	}
	requireEpisode := func(expected ...EpisodeStatus) error {
		if subject == "" {
			return fmt.Errorf("event %s requires an episode subject", kind)
		}
		current, ok := state.Episodes[subject]
		if !ok {
			return fmt.Errorf("event %s references unknown episode %s", kind, subject)
		}
		for _, status := range expected {
			if current == status {
				return nil
			}
		}
		return fmt.Errorf("event %s is invalid for episode %s in state %s", kind, subject, current)
	}

	switch kind {
	case CampaignCreated:
		return requireStatus(StatusNone)
	case PlanCompiled:
		return requireStatus(StatusDraft)
	case CampaignStarted:
		return requireStatus(StatusReady)
	case CampaignPaused:
		return requireStatus(StatusRunning)
	case CampaignResumed:
		return requireStatus(StatusPaused)
	case CampaignStopping:
		return requireStatus(StatusRunning, StatusPaused)
	case CampaignStopped:
		return requireStatus(StatusStopping)
	case CampaignFailed:
		if state.Terminal() || state.Status == StatusNone {
			return fmt.Errorf("event %s is invalid while campaign is %s", kind, state.Status)
		}
		return nil
	case CampaignCompleted:
		if err := requireStatus(StatusRunning); err != nil {
			return err
		}
		for id, status := range state.Episodes {
			if status != EpisodeCompletedState && status != EpisodeFailedTerminal {
				return fmt.Errorf("campaign cannot complete while episode %s is %s", id, status)
			}
		}
		return nil
	case CandidateProposed, SnapshotMaterialized, TrialPlanned, ObservationRecorded, EstimateRecorded, DecisionRecorded, BudgetReserved:
		return requireStatus(StatusRunning, StatusPaused)
	case UsageCommitted, BudgetReleased:
		if state.Status == StatusNone {
			return fmt.Errorf("event %s is invalid before campaign creation", kind)
		}
		return nil
	case EpisodeScheduled:
		if err := requireStatus(StatusRunning); err != nil {
			return err
		}
		if subject == "" {
			return fmt.Errorf("event %s requires an episode subject", kind)
		}
		if _, exists := state.Episodes[subject]; exists {
			return fmt.Errorf("episode %s is already scheduled", subject)
		}
		return nil
	case EpisodeLeaseGranted:
		return requireEpisode(EpisodeQueued)
	case EpisodeAttemptStarted:
		return requireEpisode(EpisodeLeased)
	case EpisodeCompleted:
		return requireEpisode(EpisodeRunning, EpisodeLeased)
	case EpisodeFailed:
		if _, err := strconv.ParseBool(tags["retryable"]); err != nil {
			return fmt.Errorf("event %s requires boolean retryable tag", kind)
		}
		return requireEpisode(EpisodeRunning, EpisodeLeased)
	default:
		return fmt.Errorf("unknown event kind %s", kind)
	}
}

func cloneEpisodes(in map[string]EpisodeStatus) map[string]EpisodeStatus {
	out := make(map[string]EpisodeStatus, len(in))
	for id, status := range in {
		out[id] = status
	}
	return out
}
