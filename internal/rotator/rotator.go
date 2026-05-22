// Package rotator provides utilities for rotating keys in .env files,
// replacing old key names with new ones while preserving values and order.
package rotator

import (
	"fmt"

	"github.com/envdiff/envdiff/internal/parser"
)

// Mapping describes a single key rotation: OldKey -> NewKey.
type Mapping struct {
	OldKey string
	NewKey string
}

// Options configures the rotation behaviour.
type Options struct {
	// Mappings is the list of key rotations to apply.
	Mappings []Mapping
	// FailOnMissing returns an error if an OldKey is not found in the file.
	FailOnMissing bool
	// DropOld removes the old key entirely when it conflicts with an existing NewKey.
	DropOld bool
}

// Result holds the outcome of a rotation operation.
type Result struct {
	File    *parser.EnvFile
	Rotated []string // keys that were successfully rotated
	Missing []string // OldKey values not found in the source file
}

// Apply performs key rotation on src according to opts and returns a Result.
// The original file is not mutated.
func Apply(src *parser.EnvFile, opts Options) (*Result, error) {
	if src == nil {
		return &Result{File: &parser.EnvFile{}}, nil
	}

	// Build a lookup of old -> new for quick access.
	rotMap := make(map[string]string, len(opts.Mappings))
	for _, m := range opts.Mappings {
		rotMap[m.OldKey] = m.NewKey
	}

	// Detect which old keys are present.
	present := make(map[string]bool)
	for _, e := range src.Entries {
		if _, ok := rotMap[e.Key]; ok {
			present[e.Key] = true
		}
	}

	// Report missing.
	var missing []string
	for _, m := range opts.Mappings {
		if !present[m.OldKey] {
			missing = append(missing, m.OldKey)
			if opts.FailOnMissing {
				return nil, fmt.Errorf("rotator: key %q not found in source file", m.OldKey)
			}
		}
	}

	// Build the new entry list.
	newEntries := make([]parser.Entry, 0, len(src.Entries))
	var rotated []string
	for _, e := range src.Entries {
		if newKey, ok := rotMap[e.Key]; ok {
			e.Key = newKey
			rotated = append(rotated, newKey)
		}
		newEntries = append(newEntries, e)
	}

	return &Result{
		File:    &parser.EnvFile{Entries: newEntries},
		Rotated: rotated,
		Missing: missing,
	}, nil
}
