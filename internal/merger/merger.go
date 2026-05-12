// Package merger provides functionality to merge multiple .env files
// with configurable precedence and conflict resolution strategies.
package merger

import (
	"fmt"

	"github.com/user/envdiff/internal/parser"
)

// Strategy defines how conflicts are resolved when merging env files.
type Strategy int

const (
	// StrategyFirst keeps the value from the first file that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast keeps the value from the last file that defines the key.
	StrategyLast
	// StrategyError returns an error if the same key is defined in multiple files.
	StrategyError
)

// Result holds the merged output and metadata about conflicts.
type Result struct {
	File      *parser.EnvFile
	Conflicts []Conflict
}

// Conflict describes a key that appeared in more than one source file.
type Conflict struct {
	Key    string
	Values []string // values in order of source files
	Sources []string // file paths corresponding to each value
}

// Merge combines multiple EnvFiles according to the given strategy.
// Files are processed in the order provided; index 0 has lowest precedence
// for StrategyLast and highest precedence for StrategyFirst.
func Merge(files []*parser.EnvFile, strategy Strategy) (*Result, error) {
	if len(files) == 0 {
		return &Result{File: &parser.EnvFile{}}, nil
	}

	merged := &parser.EnvFile{
		Entries: make(map[string]parser.Entry),
	}
	result := &Result{File: merged}

	// Track insertion order for deterministic output.
	seen := make(map[string]bool)
	conflictMap := make(map[string]*Conflict)

	for _, f := range files {
		if f == nil {
			continue
		}
		for _, key := range f.Order {
			entry := f.Entries[key]
			if existing, exists := merged.Entries[key]; exists {
				// Record conflict metadata.
				if _, ok := conflictMap[key]; !ok {
					conflictMap[key] = &Conflict{
						Key:     key,
						Values:  []string{existing.Value},
						Sources: []string{existing.Source},
					}
				}
				conflictMap[key].Values = append(conflictMap[key].Values, entry.Value)
				conflictMap[key].Sources = append(conflictMap[key].Sources, entry.Source)

				switch strategy {
				case StrategyError:
					return nil, fmt.Errorf("merger: conflict on key %q (sources: %v)", key, conflictMap[key].Sources)
				case StrategyFirst:
					// Keep existing; do nothing.
				case StrategyLast:
					merged.Entries[key] = entry
				}
			} else {
				merged.Entries[key] = entry
				if !seen[key] {
					merged.Order = append(merged.Order, key)
					seen[key] = true
				}
			}
		}
	}

	for _, c := range conflictMap {
		result.Conflicts = append(result.Conflicts, *c)
	}
	return result, nil
}
