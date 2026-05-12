package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/envdiff/internal/parser"
)

// Snapshot captures the state of an env file at a point in time.
type Snapshot struct {
	Timestamp time.Time            `json:"timestamp"`
	Source    string               `json:"source"`
	Entries   []parser.Entry       `json:"entries"`
}

// Save writes a snapshot of the given env file to the destination path as JSON.
func Save(envFile *parser.EnvFile, dest string) error {
	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		Source:    envFile.Path,
		Entries:   envFile.Entries,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}

	if err := os.WriteFile(dest, data, 0600); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", dest, err)
	}

	return nil
}

// Load reads a previously saved snapshot from path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read %s: %w", path, err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}

	return &snap, nil
}

// ToEnvFile converts a Snapshot back into an EnvFile for diffing or reconciling.
func (s *Snapshot) ToEnvFile() *parser.EnvFile {
	return &parser.EnvFile{
		Path:    s.Source,
		Entries: s.Entries,
	}
}
