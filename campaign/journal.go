package campaign

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-go-golems/optkit/record"
)

var (
	ErrVersionConflict = errors.New("campaign version conflict")
	ErrJournalCorrupt  = errors.New("campaign journal corrupt")
)

type Head struct {
	Version    uint64
	LastDigest record.Digest
}

type AppendResult struct {
	Version   uint64
	Events    []ControlEvent
	Duplicate bool
}

type CommandLookup interface {
	LookupCommand(context.Context, record.CampaignID, record.CommandID) (AppendResult, bool, error)
}

type Journal interface {
	Append(context.Context, record.CampaignID, uint64, []NewEvent) (AppendResult, error)
	Read(context.Context, record.CampaignID, uint64) ([]ControlEvent, error)
	Head(context.Context, record.CampaignID) (Head, error)
	Verify(context.Context, record.CampaignID) error
}

func VerifyChain(events []ControlEvent) error {
	previous := GenesisDigest
	var seq uint64
	for _, event := range events {
		seq++
		if event.Seq != seq {
			return fmt.Errorf("%w: expected sequence %d, got %d", ErrJournalCorrupt, seq, event.Seq)
		}
		if event.PreviousDigest != previous {
			return fmt.Errorf("%w: event %s previous digest %s, expected %s", ErrJournalCorrupt, event.ID, event.PreviousDigest, previous)
		}
		if err := VerifyEventDigest(event); err != nil {
			return fmt.Errorf("%w: %v", ErrJournalCorrupt, err)
		}
		previous = event.Digest
	}
	return nil
}
