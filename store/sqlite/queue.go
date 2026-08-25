package sqlite

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sqlitedb "github.com/go-go-golems/optkit/internal/sqlite"
	"github.com/go-go-golems/optkit/record"
	"github.com/go-go-golems/optkit/scheduler"
)

func (s *Store) Enqueue(ctx context.Context, items []scheduler.WorkItem) error {
	if len(items) == 0 {
		return nil
	}
	return s.db.InTx(ctx, func(conn *sqlitedb.Conn) error {
		now := formatTime(s.now())
		for _, item := range items {
			if err := validateWorkItem(item); err != nil {
				return err
			}
			var payloadSchema any
			if item.Payload.Schema != nil {
				payloadSchema = string(*item.Payload.Schema)
			}
			changes, err := conn.Exec(ctx, `
INSERT OR IGNORE INTO work_items(
    id, campaign_id, kind, semantic_key,
    payload_digest, payload_media_type, payload_schema, payload_size, payload_sensitivity,
    priority, earliest_start, lease_duration_ns,
    status, attempt, created_at, updated_at
) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?)`,
				string(item.ID), string(item.Campaign), string(item.Kind), item.SemanticKey.String(),
				item.Payload.Digest.String(), item.Payload.MediaType, payloadSchema, item.Payload.Size, string(item.Payload.Sensitivity),
				item.Priority, formatTime(item.EarliestStart), int64(item.LeaseDuration),
				string(scheduler.WorkReady), now, now)
			if err != nil {
				return err
			}
			if changes == 0 {
				rows, err := conn.Query(ctx, workSelect+` WHERE campaign_id = ? AND kind = ? AND semantic_key = ?`, string(item.Campaign), string(item.Kind), item.SemanticKey.String())
				if err != nil {
					return err
				}
				if len(rows) != 1 {
					return scheduler.ErrWorkConflict
				}
				existing, err := decodeWork(rows[0])
				if err != nil {
					return err
				}
				if existing.Item.ID != item.ID || existing.Item.Payload.Digest != item.Payload.Digest {
					return fmt.Errorf("%w: existing %s differs from %s", scheduler.ErrWorkConflict, existing.Item.ID, item.ID)
				}
			}
		}
		return nil
	})
}

func (s *Store) Lease(ctx context.Context, worker string, request scheduler.LeaseRequest) ([]scheduler.Lease, error) {
	if worker == "" {
		return nil, fmt.Errorf("worker ID is required")
	}
	if request.Limit <= 0 {
		return nil, fmt.Errorf("lease limit must be positive")
	}
	now := request.Now.UTC()
	if now.IsZero() {
		now = s.now().UTC()
	}
	var leases []scheduler.Lease
	err := s.db.InTx(ctx, func(conn *sqlitedb.Conn) error {
		if _, err := conn.Exec(ctx, `UPDATE work_items SET status = ?, lease_id = NULL, leased_by = NULL, lease_expires_at = NULL, updated_at = ? WHERE status = ? AND lease_expires_at <= ?`, string(scheduler.WorkReady), formatTime(now), string(scheduler.WorkLeased), formatTime(now)); err != nil {
			return err
		}
		query := workSelect + ` WHERE status = ? AND earliest_start <= ?`
		args := []any{string(scheduler.WorkReady), formatTime(now)}
		if len(request.Kinds) > 0 {
			query += ` AND kind IN (` + strings.TrimRight(strings.Repeat("?,", len(request.Kinds)), ",") + `)`
			for _, kind := range request.Kinds {
				args = append(args, string(kind))
			}
		}
		query += ` ORDER BY priority DESC, earliest_start, created_at, id LIMIT ?`
		args = append(args, int64(request.Limit))
		rows, err := conn.Query(ctx, query, args...)
		if err != nil {
			return err
		}
		for _, row := range rows {
			rec, err := decodeWork(row)
			if err != nil {
				return err
			}
			rawLease, err := record.NewID("lease")
			if err != nil {
				return err
			}
			leaseID := record.LeaseID(rawLease)
			expires := now.Add(rec.Item.LeaseDuration)
			attempt := rec.Attempt + 1
			changes, err := conn.Exec(ctx, `UPDATE work_items SET status = ?, attempt = ?, lease_id = ?, leased_by = ?, lease_expires_at = ?, updated_at = ? WHERE id = ? AND status = ?`,
				string(scheduler.WorkLeased), int64(attempt), string(leaseID), worker, formatTime(expires), formatTime(now), string(rec.Item.ID), string(scheduler.WorkReady))
			if err != nil {
				return err
			}
			if changes != 1 {
				return fmt.Errorf("work item %s was not claimable", rec.Item.ID)
			}
			leases = append(leases, scheduler.Lease{ID: leaseID, Item: rec.Item, Worker: worker, Attempt: attempt, ExpiresAt: expires})
		}
		return nil
	})
	return leases, err
}

func (s *Store) Complete(ctx context.Context, leaseID record.LeaseID, result scheduler.WorkResult) error {
	if err := result.Artifact.Validate(); err != nil {
		return err
	}
	return s.db.InTx(ctx, func(conn *sqlitedb.Conn) error {
		rows, err := conn.Query(ctx, workSelect+` WHERE lease_id = ?`, string(leaseID))
		if err != nil {
			return err
		}
		if len(rows) != 1 {
			return scheduler.ErrLeaseNotFound
		}
		rec, err := decodeWork(rows[0])
		if err != nil {
			return err
		}
		if rec.Status == scheduler.WorkCompleted {
			if rec.Result != nil && rec.Result.Artifact.Digest == result.Artifact.Digest {
				return nil
			}
			return scheduler.ErrTerminalConflict
		}
		if rec.Status != scheduler.WorkLeased {
			return scheduler.ErrLeaseNotFound
		}
		var schema any
		if result.Artifact.Schema != nil {
			schema = string(*result.Artifact.Schema)
		}
		changes, err := conn.Exec(ctx, `UPDATE work_items SET status = ?, result_digest = ?, result_media_type = ?, result_schema = ?, result_size = ?, result_sensitivity = ?, updated_at = ? WHERE id = ? AND lease_id = ? AND status = ?`,
			string(scheduler.WorkCompleted), result.Artifact.Digest.String(), result.Artifact.MediaType, schema, result.Artifact.Size, string(result.Artifact.Sensitivity), formatTime(s.now()), string(rec.Item.ID), string(leaseID), string(scheduler.WorkLeased))
		if err != nil {
			return err
		}
		if changes != 1 {
			return scheduler.ErrLeaseNotFound
		}
		return nil
	})
}

func (s *Store) Fail(ctx context.Context, leaseID record.LeaseID, failure scheduler.WorkFailure) error {
	if failure.Code == "" {
		return fmt.Errorf("work failure code is required")
	}
	data, err := record.CanonicalJSON(failure)
	if err != nil {
		return err
	}
	return s.db.InTx(ctx, func(conn *sqlitedb.Conn) error {
		rows, err := conn.Query(ctx, workSelect+` WHERE lease_id = ?`, string(leaseID))
		if err != nil {
			return err
		}
		if len(rows) != 1 {
			return scheduler.ErrLeaseNotFound
		}
		rec, err := decodeWork(rows[0])
		if err != nil {
			return err
		}
		if rec.Status == scheduler.WorkCompleted || rec.Status == scheduler.WorkFailed {
			return scheduler.ErrTerminalConflict
		}
		if rec.Status != scheduler.WorkLeased {
			return scheduler.ErrLeaseNotFound
		}
		if failure.Retryable {
			_, err = conn.Exec(ctx, `UPDATE work_items SET status = ?, earliest_start = ?, lease_id = NULL, leased_by = NULL, lease_expires_at = NULL, failure_json = ?, updated_at = ? WHERE id = ? AND lease_id = ? AND status = ?`,
				string(scheduler.WorkReady), formatTime(s.now().Add(failure.Backoff)), data, formatTime(s.now()), string(rec.Item.ID), string(leaseID), string(scheduler.WorkLeased))
			return err
		}
		_, err = conn.Exec(ctx, `UPDATE work_items SET status = ?, failure_json = ?, updated_at = ? WHERE id = ? AND lease_id = ? AND status = ?`,
			string(scheduler.WorkFailed), data, formatTime(s.now()), string(rec.Item.ID), string(leaseID), string(scheduler.WorkLeased))
		return err
	})
}

func (s *Store) ReclaimExpired(ctx context.Context, now time.Time) (int64, error) {
	if now.IsZero() {
		now = s.now()
	}
	return s.db.Exec(ctx, `UPDATE work_items SET status = ?, lease_id = NULL, leased_by = NULL, lease_expires_at = NULL, updated_at = ? WHERE status = ? AND lease_expires_at <= ?`, string(scheduler.WorkReady), formatTime(now), string(scheduler.WorkLeased), formatTime(now))
}

func (s *Store) Get(ctx context.Context, id record.WorkID) (scheduler.WorkRecord, error) {
	rows, err := s.db.Query(ctx, workSelect+` WHERE id = ?`, string(id))
	if err != nil {
		return scheduler.WorkRecord{}, err
	}
	if len(rows) == 0 {
		return scheduler.WorkRecord{}, scheduler.ErrWorkNotFound
	}
	return decodeWork(rows[0])
}

const workSelect = `SELECT
    id, campaign_id, kind, semantic_key,
    payload_digest, payload_media_type, payload_schema, payload_size, payload_sensitivity,
    priority, earliest_start, lease_duration_ns,
    status, attempt, lease_id, leased_by, lease_expires_at,
    result_digest, result_media_type, result_schema, result_size, result_sensitivity,
    failure_json, created_at, updated_at
FROM work_items`

func decodeWork(row sqlitedb.Row) (scheduler.WorkRecord, error) {
	id, err := text(row, "id")
	if err != nil {
		return scheduler.WorkRecord{}, err
	}
	campaignID, err := text(row, "campaign_id")
	if err != nil {
		return scheduler.WorkRecord{}, err
	}
	kind, err := text(row, "kind")
	if err != nil {
		return scheduler.WorkRecord{}, err
	}
	semanticRaw, err := text(row, "semantic_key")
	if err != nil {
		return scheduler.WorkRecord{}, err
	}
	semantic, err := record.ParseDigest(semanticRaw)
	if err != nil {
		return scheduler.WorkRecord{}, err
	}
	payload, err := decodeRef(row, "payload")
	if err != nil {
		return scheduler.WorkRecord{}, err
	}
	priority, err := integer(row, "priority")
	if err != nil {
		return scheduler.WorkRecord{}, err
	}
	earliestRaw, err := text(row, "earliest_start")
	if err != nil {
		return scheduler.WorkRecord{}, err
	}
	earliest, err := parseTime(earliestRaw)
	if err != nil {
		return scheduler.WorkRecord{}, err
	}
	leaseDuration, err := integer(row, "lease_duration_ns")
	if err != nil {
		return scheduler.WorkRecord{}, err
	}
	statusRaw, err := text(row, "status")
	if err != nil {
		return scheduler.WorkRecord{}, err
	}
	attempt, err := integer(row, "attempt")
	if err != nil {
		return scheduler.WorkRecord{}, err
	}
	leaseRaw, err := nullableText(row, "lease_id")
	if err != nil {
		return scheduler.WorkRecord{}, err
	}
	workerRaw, err := nullableText(row, "leased_by")
	if err != nil {
		return scheduler.WorkRecord{}, err
	}
	expiryRaw, err := nullableText(row, "lease_expires_at")
	if err != nil {
		return scheduler.WorkRecord{}, err
	}
	rec := scheduler.WorkRecord{Item: scheduler.WorkItem{
		ID: record.WorkID(id), Campaign: record.CampaignID(campaignID), Kind: scheduler.WorkKind(kind), SemanticKey: semantic, Payload: payload,
		Priority: int(priority), EarliestStart: earliest, LeaseDuration: time.Duration(leaseDuration),
	}, Status: scheduler.WorkStatus(statusRaw), Attempt: int(attempt)}
	if leaseRaw != nil {
		rec.LeaseID = record.LeaseID(*leaseRaw)
	}
	if workerRaw != nil {
		rec.LeasedBy = *workerRaw
	}
	if expiryRaw != nil {
		rec.LeaseExpiry, err = parseTime(*expiryRaw)
		if err != nil {
			return scheduler.WorkRecord{}, err
		}
	}
	resultDigest, err := nullableText(row, "result_digest")
	if err != nil {
		return scheduler.WorkRecord{}, err
	}
	if resultDigest != nil {
		result, err := decodeRef(row, "result")
		if err != nil {
			return scheduler.WorkRecord{}, err
		}
		rec.Result = &scheduler.WorkResult{Artifact: result}
	}
	failureRaw, err := bytesValue(row, "failure_json")
	if err != nil {
		return scheduler.WorkRecord{}, err
	}
	if len(failureRaw) > 0 {
		var failure scheduler.WorkFailure
		if err := json.Unmarshal(failureRaw, &failure); err != nil {
			return scheduler.WorkRecord{}, err
		}
		rec.Failure = &failure
	}
	return rec, nil
}

func validateWorkItem(item scheduler.WorkItem) error {
	if err := record.ValidateID("work", string(item.ID)); err != nil {
		return err
	}
	if err := record.ValidateID("campaign", string(item.Campaign)); err != nil {
		return err
	}
	if item.Kind == "" {
		return fmt.Errorf("work kind is required")
	}
	if err := item.SemanticKey.Validate(); err != nil {
		return err
	}
	if err := item.Payload.Validate(); err != nil {
		return err
	}
	if item.LeaseDuration <= 0 {
		return fmt.Errorf("work lease duration must be positive")
	}
	return nil
}

var _ scheduler.Queue = (*Store)(nil)
