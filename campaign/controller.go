package campaign

import (
	"context"
	"fmt"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/record"
)

type CommandKind string

const (
	CreateCampaign   CommandKind = "CreateCampaign"
	CompilePlan      CommandKind = "CompilePlan"
	StartCampaign    CommandKind = "StartCampaign"
	PauseCampaign    CommandKind = "PauseCampaign"
	ResumeCampaign   CommandKind = "ResumeCampaign"
	StopCampaign     CommandKind = "StopCampaign"
	ConfirmStopped   CommandKind = "ConfirmStopped"
	FailCampaign     CommandKind = "FailCampaign"
	CompleteCampaign CommandKind = "CompleteCampaign"
)

type Command struct {
	ID              record.CommandID
	Campaign        record.CampaignID
	ExpectedVersion uint64
	Kind            CommandKind
	Actor           record.ActorRef
	Schema          record.SchemaID
	Payload         any
	Tags            map[string]string
}

type Controller struct {
	Journal   Journal
	Artifacts artifact.Store
}

func (c Controller) Handle(ctx context.Context, command Command) (AppendResult, error) {
	if c.Journal == nil || c.Artifacts == nil {
		return AppendResult{}, fmt.Errorf("controller requires journal and artifact store")
	}
	if err := record.ValidateID("command", string(command.ID)); err != nil {
		return AppendResult{}, err
	}
	if lookup, ok := c.Journal.(CommandLookup); ok {
		prior, found, err := lookup.LookupCommand(ctx, command.Campaign, command.ID)
		if err != nil {
			return AppendResult{}, err
		}
		if found {
			prior.Duplicate = true
			return prior, nil
		}
	}
	events, err := c.Journal.Read(ctx, command.Campaign, 0)
	if err != nil {
		return AppendResult{}, err
	}
	state, err := Fold(command.Campaign, events)
	if err != nil {
		return AppendResult{}, err
	}
	kind, err := commandEventKind(command.Kind)
	if err != nil {
		return AppendResult{}, err
	}
	if err := ValidateTransition(state, kind, "", command.Tags); err != nil {
		return AppendResult{}, err
	}
	ref, err := artifact.PutCanonical(ctx, c.Artifacts, command.Schema, artifact.SensitivityInternal, command.Payload)
	if err != nil {
		return AppendResult{}, err
	}
	return c.Journal.Append(ctx, command.Campaign, command.ExpectedVersion, []NewEvent{{
		Kind: kind, Schema: command.Schema, Actor: command.Actor, Command: &command.ID,
		Payload: ref, Tags: command.Tags,
	}})
}

func commandEventKind(kind CommandKind) (EventKind, error) {
	switch kind {
	case CreateCampaign:
		return CampaignCreated, nil
	case CompilePlan:
		return PlanCompiled, nil
	case StartCampaign:
		return CampaignStarted, nil
	case PauseCampaign:
		return CampaignPaused, nil
	case ResumeCampaign:
		return CampaignResumed, nil
	case StopCampaign:
		return CampaignStopping, nil
	case ConfirmStopped:
		return CampaignStopped, nil
	case FailCampaign:
		return CampaignFailed, nil
	case CompleteCampaign:
		return CampaignCompleted, nil
	default:
		return "", fmt.Errorf("unknown command kind %s", kind)
	}
}
