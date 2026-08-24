//go:build cgo

package sqlite

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
)

func TestOpenUsesWALAndTransactions(t *testing.T) {
	ctx := context.Background()
	db, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.ExecScript(ctx, `CREATE TABLE values_table (id INTEGER PRIMARY KEY, name TEXT NOT NULL);`); err != nil {
		t.Fatal(err)
	}
	if err := db.InTx(ctx, func(conn *Conn) error {
		if _, err := conn.Exec(ctx, `INSERT INTO values_table(id, name) VALUES(?, ?)`, int64(1), "one"); err != nil {
			return err
		}
		return errors.New("force rollback")
	}); err == nil {
		t.Fatal("transaction unexpectedly committed")
	}
	rows, err := db.Query(ctx, `SELECT COUNT(*) AS count FROM values_table`)
	if err != nil {
		t.Fatal(err)
	}
	if got := rows[0]["count"].Integer; got != 0 {
		t.Fatalf("rollback left %d rows", got)
	}
	if _, err := db.Exec(ctx, `INSERT INTO values_table(id, name) VALUES(?, ?)`, int64(2), "two"); err != nil {
		t.Fatal(err)
	}
	rows, err = db.Query(ctx, `SELECT id, name FROM values_table`)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0]["name"].Text != "two" {
		t.Fatalf("unexpected rows: %#v", rows)
	}
	pragma, err := db.Query(ctx, `PRAGMA journal_mode`)
	if err != nil {
		t.Fatal(err)
	}
	if pragma[0]["journal_mode"].Text != "wal" {
		t.Fatalf("journal_mode = %q", pragma[0]["journal_mode"].Text)
	}
}
