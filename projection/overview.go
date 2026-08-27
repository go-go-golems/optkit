package projection

import (
	"fmt"

	"github.com/go-go-golems/optkit/campaign"
	"github.com/go-go-golems/optkit/record"
)

type Overview struct {
	Campaign       record.CampaignID `json:"campaign"`
	Version        uint64            `json:"version"`
	Status         campaign.Status   `json:"status"`
	Events         uint64            `json:"events"`
	Episodes       int               `json:"episodes"`
	Queued         int               `json:"queued"`
	Active         int               `json:"active"`
	Completed      int               `json:"completed"`
	FailedTerminal int               `json:"failed_terminal"`
	LastDigest     record.Digest     `json:"last_digest"`
}

func RebuildOverview(campaignID record.CampaignID, events []campaign.ControlEvent) (Overview, error) {
	state, err := campaign.Fold(campaignID, events)
	if err != nil {
		return Overview{}, fmt.Errorf("fold campaign overview: %w", err)
	}
	overview := Overview{
		Campaign: campaignID, Version: state.Version, Status: state.Status,
		Events: state.Events, Episodes: len(state.Episodes), LastDigest: state.LastDigest,
	}
	for _, status := range state.Episodes {
		switch status {
		case campaign.EpisodeQueued:
			overview.Queued++
		case campaign.EpisodeLeased, campaign.EpisodeRunning:
			overview.Active++
		case campaign.EpisodeCompletedState:
			overview.Completed++
		case campaign.EpisodeFailedTerminal:
			overview.FailedTerminal++
		}
	}
	return overview, nil
}
