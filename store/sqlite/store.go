package sqlite

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	sqlitedb "github.com/go-go-golems/optkit/internal/sqlite"
)

type Store struct {
	db  *sqlitedb.DB
	now func() time.Time
}

func Open(path string) (*Store, error) {
	return OpenWithClock(path, func() time.Time { return time.Now().UTC() })
}

func OpenWithClock(path string, now func() time.Time) (*Store, error) {
	if now == nil {
		return nil, fmt.Errorf("sqlite store clock is required")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve sqlite store path: %w", err)
	}
	db, err := sqlitedb.Open(abs)
	if err != nil {
		return nil, err
	}
	store := &Store{db: db, now: now}
	if err := store.migrate(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) migrate(ctx context.Context) error {
	return s.db.ExecScript(ctx, `
CREATE TABLE IF NOT EXISTS campaign_heads (
    campaign_id TEXT PRIMARY KEY,
    version INTEGER NOT NULL,
    last_digest TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS campaign_events (
    campaign_id TEXT NOT NULL,
    seq INTEGER NOT NULL,
    event_id TEXT NOT NULL UNIQUE,
    kind TEXT NOT NULL,
    schema_id TEXT NOT NULL,
    subject TEXT NOT NULL DEFAULT '',
    occurred_at TEXT NOT NULL,
    recorded_at TEXT NOT NULL,
    actor TEXT NOT NULL,
    command_id TEXT,
    causation_id TEXT,
    correlation_id TEXT,
    previous_digest TEXT NOT NULL,
    payload_digest TEXT NOT NULL,
    payload_media_type TEXT NOT NULL,
    payload_schema TEXT,
    payload_size INTEGER NOT NULL,
    payload_sensitivity TEXT NOT NULL,
    tags_json BLOB NOT NULL,
    digest TEXT NOT NULL,
    PRIMARY KEY (campaign_id, seq),
    UNIQUE (campaign_id, digest),
    FOREIGN KEY (campaign_id) REFERENCES campaign_heads(campaign_id)
);

CREATE INDEX IF NOT EXISTS campaign_events_kind_idx
    ON campaign_events(campaign_id, kind, seq);

CREATE TABLE IF NOT EXISTS campaign_commands (
    campaign_id TEXT NOT NULL,
    command_id TEXT NOT NULL,
    first_seq INTEGER NOT NULL,
    last_seq INTEGER NOT NULL,
    PRIMARY KEY (campaign_id, command_id),
    FOREIGN KEY (campaign_id) REFERENCES campaign_heads(campaign_id)
);

CREATE TABLE IF NOT EXISTS work_items (
    id TEXT PRIMARY KEY,
    campaign_id TEXT NOT NULL,
    kind TEXT NOT NULL,
    semantic_key TEXT NOT NULL,
    payload_digest TEXT NOT NULL,
    payload_media_type TEXT NOT NULL,
    payload_schema TEXT,
    payload_size INTEGER NOT NULL,
    payload_sensitivity TEXT NOT NULL,
    priority INTEGER NOT NULL,
    earliest_start TEXT NOT NULL,
    lease_duration_ns INTEGER NOT NULL,
    resource_claims_json BLOB NOT NULL,
    status TEXT NOT NULL,
    attempt INTEGER NOT NULL DEFAULT 0,
    lease_id TEXT,
    leased_by TEXT,
    lease_expires_at TEXT,
    result_digest TEXT,
    result_media_type TEXT,
    result_schema TEXT,
    result_size INTEGER,
    result_sensitivity TEXT,
    failure_json BLOB,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    UNIQUE (campaign_id, kind, semantic_key),
    FOREIGN KEY (campaign_id) REFERENCES campaign_heads(campaign_id)
);

CREATE INDEX IF NOT EXISTS work_ready_idx
    ON work_items(status, earliest_start, priority DESC, created_at, id);
CREATE INDEX IF NOT EXISTS work_lease_idx
    ON work_items(lease_id);

CREATE TABLE IF NOT EXISTS budget_limits (
    campaign_id TEXT NOT NULL,
    resource TEXT NOT NULL,
    limit_units INTEGER NOT NULL,
    PRIMARY KEY (campaign_id, resource),
    FOREIGN KEY (campaign_id) REFERENCES campaign_heads(campaign_id)
);

CREATE TABLE IF NOT EXISTS budget_reservations (
    id TEXT PRIMARY KEY,
    campaign_id TEXT NOT NULL,
    work_id TEXT NOT NULL UNIQUE,
    requested_json BLOB NOT NULL,
    actual_json BLOB NOT NULL,
    status TEXT NOT NULL,
    overage INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (campaign_id) REFERENCES campaign_heads(campaign_id)
);

CREATE INDEX IF NOT EXISTS budget_reservations_campaign_idx
    ON budget_reservations(campaign_id, status, id);
`)
}

const timeLayout = time.RFC3339Nano

func formatTime(value time.Time) string { return value.UTC().Format(timeLayout) }

func parseTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	parsed, err := time.Parse(timeLayout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse SQLite time %q: %w", value, err)
	}
	return parsed.UTC(), nil
}
