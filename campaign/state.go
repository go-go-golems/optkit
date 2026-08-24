package campaign

import "github.com/go-go-golems/optkit/record"

type Status string

const (
	StatusNone      Status = ""
	StatusDraft     Status = "draft"
	StatusReady     Status = "ready"
	StatusRunning   Status = "running"
	StatusPaused    Status = "paused"
	StatusStopping  Status = "stopping"
	StatusStopped   Status = "stopped"
	StatusFailed    Status = "failed"
	StatusCompleted Status = "completed"
)

type EpisodeStatus string

const (
	EpisodeQueued         EpisodeStatus = "queued"
	EpisodeLeased         EpisodeStatus = "leased"
	EpisodeRunning        EpisodeStatus = "running"
	EpisodeCompletedState EpisodeStatus = "completed"
	EpisodeFailedTerminal EpisodeStatus = "failed_terminal"
)

type State struct {
	Campaign   record.CampaignID
	Version    uint64
	LastDigest record.Digest
	Status     Status
	Episodes   map[string]EpisodeStatus
	Events     uint64
}

func Initial(campaignID record.CampaignID) State {
	return State{Campaign: campaignID, LastDigest: GenesisDigest, Episodes: make(map[string]EpisodeStatus)}
}

func (s State) Terminal() bool {
	return s.Status == StatusStopped || s.Status == StatusFailed || s.Status == StatusCompleted
}
