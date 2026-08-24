package local

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-go-golems/optkit/artifact/filesystem"
	sqlitestore "github.com/go-go-golems/optkit/store/sqlite"
)

type Profile struct {
	Root      string
	Metadata  *sqlitestore.Store
	Artifacts *filesystem.Store
}

func Open(root string) (*Profile, error) {
	return OpenWithClock(root, func() time.Time { return time.Now().UTC() })
}

func OpenWithClock(root string, now func() time.Time) (*Profile, error) {
	if root == "" {
		return nil, fmt.Errorf("local store root is required")
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve local store root: %w", err)
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return nil, fmt.Errorf("create local store root: %w", err)
	}
	artifacts, err := filesystem.New(filepath.Join(abs, "artifacts"))
	if err != nil {
		return nil, err
	}
	metadata, err := sqlitestore.OpenWithClock(filepath.Join(abs, "optkit.db"), now)
	if err != nil {
		return nil, err
	}
	return &Profile{Root: abs, Metadata: metadata, Artifacts: artifacts}, nil
}

func (p *Profile) Close() error {
	if p == nil || p.Metadata == nil {
		return nil
	}
	return p.Metadata.Close()
}
