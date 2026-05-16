package freezer

import (
	"encoding/json"
	"fmt"
	"os"
)

// Save writes frozen entries to path as JSON.
func Save(path string, entries []FrozenEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("freezer: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("freezer: write %s: %w", path, err)
	}
	return nil
}

// Load reads frozen entries from a JSON file at path.
func Load(path string) ([]FrozenEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("freezer: freeze file not found: %s", path)
		}
		return nil, fmt.Errorf("freezer: read %s: %w", path, err)
	}
	var entries []FrozenEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("freezer: unmarshal: %w", err)
	}
	return entries, nil
}
