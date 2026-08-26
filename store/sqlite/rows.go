package sqlite

import (
	"fmt"
	"math"

	sqlitedb "github.com/go-go-golems/optkit/internal/sqlite"
)

func text(row sqlitedb.Row, name string) (string, error) {
	value, ok := row[name]
	if !ok {
		return "", fmt.Errorf("SQLite row has no column %q", name)
	}
	if value.Type != sqlitedb.Text {
		return "", fmt.Errorf("SQLite column %q is not text", name)
	}
	return value.Text, nil
}

func nullableText(row sqlitedb.Row, name string) (*string, error) {
	value, ok := row[name]
	if !ok {
		return nil, fmt.Errorf("SQLite row has no column %q", name)
	}
	if value.Type == sqlitedb.Null {
		return nil, nil
	}
	if value.Type != sqlitedb.Text {
		return nil, fmt.Errorf("SQLite column %q is not nullable text", name)
	}
	out := value.Text
	return &out, nil
}

func sqliteSequence(value uint64) (int64, error) {
	if value > math.MaxInt64 {
		return 0, fmt.Errorf("sequence %d exceeds SQLite's signed integer range", value)
	}
	return int64(value), nil
}

func integer(row sqlitedb.Row, name string) (int64, error) {
	value, ok := row[name]
	if !ok {
		return 0, fmt.Errorf("SQLite row has no column %q", name)
	}
	if value.Type != sqlitedb.Integer {
		return 0, fmt.Errorf("SQLite column %q is not integer", name)
	}
	return value.Integer, nil
}

func bytesValue(row sqlitedb.Row, name string) ([]byte, error) {
	value, ok := row[name]
	if !ok {
		return nil, fmt.Errorf("SQLite row has no column %q", name)
	}
	if value.Type == sqlitedb.Null {
		return nil, nil
	}
	if value.Type != sqlitedb.Blob {
		return nil, fmt.Errorf("SQLite column %q is not blob", name)
	}
	return append([]byte(nil), value.Blob...), nil
}
