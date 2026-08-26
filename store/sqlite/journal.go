package sqlite

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-go-golems/optkit/artifact"
	"github.com/go-go-golems/optkit/campaign"
	sqlitedb "github.com/go-go-golems/optkit/internal/sqlite"
	"github.com/go-go-golems/optkit/record"
)

func (s *Store) Append(ctx context.Context, campaignID record.CampaignID, expected uint64, inputs []campaign.NewEvent) (campaign.AppendResult, error) {
	if err := record.ValidateID("campaign", string(campaignID)); err != nil {
		return campaign.AppendResult{}, err
	}
	if len(inputs) == 0 {
		return campaign.AppendResult{}, fmt.Errorf("journal append requires at least one event")
	}
	var result campaign.AppendResult
	err := s.db.InTx(ctx, func(conn *sqlitedb.Conn) error {
		if _, err := conn.Exec(ctx, `INSERT OR IGNORE INTO campaign_heads(campaign_id, version, last_digest) VALUES(?, 0, ?)`, string(campaignID), campaign.GenesisDigest.String()); err != nil {
			return err
		}

		if commandID, ok := sharedCommand(inputs); ok {
			rows, err := conn.Query(ctx, `SELECT first_seq, last_seq FROM campaign_commands WHERE campaign_id = ? AND command_id = ?`, string(campaignID), string(commandID))
			if err != nil {
				return err
			}
			if len(rows) == 1 {
				first, err := integer(rows[0], "first_seq")
				if err != nil {
					return err
				}
				last, err := integer(rows[0], "last_seq")
				if err != nil {
					return err
				}
				events, err := readEventsConn(ctx, conn, campaignID, uint64(first-1), uint64(last))
				if err != nil {
					return err
				}
				result = campaign.AppendResult{Version: uint64(last), Events: events, Duplicate: true}
				return nil
			}
		}

		headRows, err := conn.Query(ctx, `SELECT version, last_digest FROM campaign_heads WHERE campaign_id = ?`, string(campaignID))
		if err != nil {
			return err
		}
		if len(headRows) != 1 {
			return fmt.Errorf("campaign head %s was not created", campaignID)
		}
		versionValue, err := integer(headRows[0], "version")
		if err != nil {
			return err
		}
		lastRaw, err := text(headRows[0], "last_digest")
		if err != nil {
			return err
		}
		lastDigest, err := record.ParseDigest(lastRaw)
		if err != nil {
			return err
		}
		version := uint64(versionValue)
		if version != expected {
			return fmt.Errorf("%w: expected %d, current %d", campaign.ErrVersionConflict, expected, version)
		}

		materialized := make([]campaign.ControlEvent, 0, len(inputs))
		previous := lastDigest
		for i, input := range inputs {
			seq := version + uint64(i) + 1
			event, err := campaign.MaterializeEvent(campaignID, seq, s.now().UTC(), previous, input)
			if err != nil {
				return err
			}
			tags, err := record.CanonicalJSON(event.Tags)
			if err != nil {
				return err
			}
			var payloadSchema any
			if event.Payload.Schema != nil {
				payloadSchema = string(*event.Payload.Schema)
			}
			if _, err := conn.Exec(ctx, `
INSERT INTO campaign_events(
    campaign_id, seq, event_id, kind, schema_id, subject,
    occurred_at, recorded_at, actor, command_id,
    previous_digest, payload_digest, payload_media_type, payload_schema,
    payload_size, payload_sensitivity, tags_json, digest
) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				string(event.Campaign), int64(event.Seq), string(event.ID), string(event.Kind), string(event.Schema), event.Subject,
				formatTime(event.OccurredAt), formatTime(event.RecordedAt), string(event.Actor), optionalCommand(event.Command),
				event.PreviousDigest.String(), event.Payload.Digest.String(), event.Payload.MediaType, payloadSchema,
				event.Payload.Size, string(event.Payload.Sensitivity), tags, event.Digest.String()); err != nil {
				return err
			}
			materialized = append(materialized, event)
			previous = event.Digest
		}
		newVersion := version + uint64(len(materialized))
		changes, err := conn.Exec(ctx, `UPDATE campaign_heads SET version = ?, last_digest = ? WHERE campaign_id = ? AND version = ?`, int64(newVersion), previous.String(), string(campaignID), int64(version))
		if err != nil {
			return err
		}
		if changes != 1 {
			return fmt.Errorf("%w: campaign head changed during append", campaign.ErrVersionConflict)
		}
		if commandID, ok := sharedCommand(inputs); ok {
			if _, err := conn.Exec(ctx, `INSERT INTO campaign_commands(campaign_id, command_id, first_seq, last_seq) VALUES(?, ?, ?, ?)`, string(campaignID), string(commandID), int64(version+1), int64(newVersion)); err != nil {
				return err
			}
		}
		result = campaign.AppendResult{Version: newVersion, Events: materialized}
		return nil
	})
	if err != nil {
		return campaign.AppendResult{}, err
	}
	return result, nil
}

func (s *Store) Read(ctx context.Context, campaignID record.CampaignID, from uint64) ([]campaign.ControlEvent, error) {
	sequence, err := sqliteSequence(from)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query(ctx, eventSelect+` WHERE campaign_id = ? AND seq > ? ORDER BY seq`, string(campaignID), sequence)
	if err != nil {
		return nil, err
	}
	return decodeEvents(rows)
}

func (s *Store) Head(ctx context.Context, campaignID record.CampaignID) (campaign.Head, error) {
	rows, err := s.db.Query(ctx, `SELECT version, last_digest FROM campaign_heads WHERE campaign_id = ?`, string(campaignID))
	if err != nil {
		return campaign.Head{}, err
	}
	if len(rows) == 0 {
		return campaign.Head{LastDigest: campaign.GenesisDigest}, nil
	}
	version, err := integer(rows[0], "version")
	if err != nil {
		return campaign.Head{}, err
	}
	digestRaw, err := text(rows[0], "last_digest")
	if err != nil {
		return campaign.Head{}, err
	}
	digest, err := record.ParseDigest(digestRaw)
	if err != nil {
		return campaign.Head{}, err
	}
	return campaign.Head{Version: uint64(version), LastDigest: digest}, nil
}

func (s *Store) Verify(ctx context.Context, campaignID record.CampaignID) error {
	events, err := s.Read(ctx, campaignID, 0)
	if err != nil {
		return err
	}
	if err := campaign.VerifyChain(events); err != nil {
		return err
	}
	head, err := s.Head(ctx, campaignID)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		if head.Version != 0 || head.LastDigest != campaign.GenesisDigest {
			return fmt.Errorf("%w: empty journal has non-genesis head", campaign.ErrJournalCorrupt)
		}
		return nil
	}
	last := events[len(events)-1]
	if head.Version != last.Seq || head.LastDigest != last.Digest {
		return fmt.Errorf("%w: head %d/%s differs from last event %d/%s", campaign.ErrJournalCorrupt, head.Version, head.LastDigest, last.Seq, last.Digest)
	}
	return nil
}

const eventSelect = `SELECT
    campaign_id, seq, event_id, kind, schema_id, subject,
    occurred_at, recorded_at, actor, command_id,
    previous_digest, payload_digest, payload_media_type, payload_schema,
    payload_size, payload_sensitivity, tags_json, digest
FROM campaign_events`

func readEventsConn(ctx context.Context, conn *sqlitedb.Conn, campaignID record.CampaignID, from, through uint64) ([]campaign.ControlEvent, error) {
	fromSequence, err := sqliteSequence(from)
	if err != nil {
		return nil, err
	}
	throughSequence, err := sqliteSequence(through)
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(ctx, eventSelect+` WHERE campaign_id = ? AND seq > ? AND seq <= ? ORDER BY seq`, string(campaignID), fromSequence, throughSequence)
	if err != nil {
		return nil, err
	}
	return decodeEvents(rows)
}

func decodeEvents(rows []sqlitedb.Row) ([]campaign.ControlEvent, error) {
	events := make([]campaign.ControlEvent, 0, len(rows))
	for _, row := range rows {
		event, err := decodeEvent(row)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func decodeEvent(row sqlitedb.Row) (campaign.ControlEvent, error) {
	campaignRaw, err := text(row, "campaign_id")
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	seq, err := integer(row, "seq")
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	id, err := text(row, "event_id")
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	kind, err := text(row, "kind")
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	schema, err := text(row, "schema_id")
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	subject, err := text(row, "subject")
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	occurredRaw, err := text(row, "occurred_at")
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	recordedRaw, err := text(row, "recorded_at")
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	occurred, err := parseTime(occurredRaw)
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	recorded, err := parseTime(recordedRaw)
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	actor, err := text(row, "actor")
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	commandRaw, err := nullableText(row, "command_id")
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	previousRaw, err := text(row, "previous_digest")
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	previous, err := record.ParseDigest(previousRaw)
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	payload, err := decodeRef(row, "payload")
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	tagsRaw, err := bytesValue(row, "tags_json")
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	var tags map[string]string
	if err := json.Unmarshal(tagsRaw, &tags); err != nil {
		return campaign.ControlEvent{}, fmt.Errorf("decode event tags: %w", err)
	}
	digestRaw, err := text(row, "digest")
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	digest, err := record.ParseDigest(digestRaw)
	if err != nil {
		return campaign.ControlEvent{}, err
	}
	event := campaign.ControlEvent{
		ID: record.EventID(id), Campaign: record.CampaignID(campaignRaw), Seq: uint64(seq), Kind: campaign.EventKind(kind), Schema: record.SchemaID(schema), Subject: subject,
		OccurredAt: occurred, RecordedAt: recorded, Actor: record.ActorRef(actor), PreviousDigest: previous, Payload: payload, Tags: tags, Digest: digest,
	}
	if commandRaw != nil {
		value := record.CommandID(*commandRaw)
		event.Command = &value
	}
	return event, nil
}

func decodeRef(row sqlitedb.Row, prefix string) (artifact.Ref, error) {
	digestRaw, err := text(row, prefix+"_digest")
	if err != nil {
		return artifact.Ref{}, err
	}
	digest, err := record.ParseDigest(digestRaw)
	if err != nil {
		return artifact.Ref{}, err
	}
	mediaType, err := text(row, prefix+"_media_type")
	if err != nil {
		return artifact.Ref{}, err
	}
	schemaRaw, err := nullableText(row, prefix+"_schema")
	if err != nil {
		return artifact.Ref{}, err
	}
	size, err := integer(row, prefix+"_size")
	if err != nil {
		return artifact.Ref{}, err
	}
	sensitivity, err := text(row, prefix+"_sensitivity")
	if err != nil {
		return artifact.Ref{}, err
	}
	ref := artifact.Ref{Digest: digest, MediaType: mediaType, Size: size, Sensitivity: artifact.Sensitivity(sensitivity)}
	if schemaRaw != nil {
		schema := record.SchemaID(*schemaRaw)
		ref.Schema = &schema
	}
	return ref, nil
}

func sharedCommand(events []campaign.NewEvent) (record.CommandID, bool) {
	if len(events) == 0 || events[0].Command == nil {
		return "", false
	}
	command := *events[0].Command
	for _, event := range events[1:] {
		if event.Command == nil || *event.Command != command {
			return "", false
		}
	}
	return command, true
}

func optionalCommand(value *record.CommandID) any {
	if value == nil {
		return nil
	}
	return string(*value)
}

var _ campaign.Journal = (*Store)(nil)

func (s *Store) LookupCommand(ctx context.Context, campaignID record.CampaignID, commandID record.CommandID) (campaign.AppendResult, bool, error) {
	rows, err := s.db.Query(ctx, `SELECT first_seq, last_seq FROM campaign_commands WHERE campaign_id = ? AND command_id = ?`, string(campaignID), string(commandID))
	if err != nil {
		return campaign.AppendResult{}, false, err
	}
	if len(rows) == 0 {
		return campaign.AppendResult{}, false, nil
	}
	first, err := integer(rows[0], "first_seq")
	if err != nil {
		return campaign.AppendResult{}, false, err
	}
	last, err := integer(rows[0], "last_seq")
	if err != nil {
		return campaign.AppendResult{}, false, err
	}
	events, err := s.db.Query(ctx, eventSelect+` WHERE campaign_id = ? AND seq >= ? AND seq <= ? ORDER BY seq`, string(campaignID), first, last)
	if err != nil {
		return campaign.AppendResult{}, false, err
	}
	decoded, err := decodeEvents(events)
	if err != nil {
		return campaign.AppendResult{}, false, err
	}
	return campaign.AppendResult{Version: uint64(last), Events: decoded, Duplicate: true}, true, nil
}

var _ campaign.CommandLookup = (*Store)(nil)
