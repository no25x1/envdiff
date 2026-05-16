// Package defaulter fills in missing keys in a target env file from a
// defaults env file, leaving existing values untouched.
package defaulter

import (
	"fmt"

	"github.com/user/envdiff/internal/parser"
)

// Options controls how defaults are applied.
type Options struct {
	// Keys restricts which keys are considered. If empty, all keys are used.
	Keys []string
	// Overwrite forces the default value even when the key already exists.
	Overwrite bool
}

// Result describes a single key that was acted upon.
type Result struct {
	Key     string
	Value   string
	Skipped bool // true when key existed and Overwrite was false
}

// Apply merges defaults into dst, returning the updated file and a report of
// every key that was evaluated.
func Apply(defaults, dst *parser.EnvFile, opts Options) (*parser.EnvFile, []Result, error) {
	if defaults == nil {
		return nil, nil, fmt.Errorf("defaulter: defaults file must not be nil")
	}
	if dst == nil {
		dst = &parser.EnvFile{}
	}

	allowSet := toSet(opts.Keys)

	// Build a mutable copy of dst entries.
	out := &parser.EnvFile{Path: dst.Path}
	out.Entries = make([]parser.Entry, len(dst.Entries))
	copy(out.Entries, dst.Entries)

	// Index existing keys in dst for fast lookup.
	existing := make(map[string]int, len(out.Entries))
	for i, e := range out.Entries {
		existing[e.Key] = i
	}

	var results []Result

	for _, de := range defaults.Entries {
		if len(allowSet) > 0 && !allowSet[de.Key] {
			continue
		}
		if idx, found := existing[de.Key]; found {
			if opts.Overwrite {
				out.Entries[idx].Value = de.Value
				results = append(results, Result{Key: de.Key, Value: de.Value})
			} else {
				results = append(results, Result{Key: de.Key, Value: out.Entries[idx].Value, Skipped: true})
			}
		} else {
			out.Entries = append(out.Entries, parser.Entry{Key: de.Key, Value: de.Value})
			existing[de.Key] = len(out.Entries) - 1
			results = append(results, Result{Key: de.Key, Value: de.Value})
		}
	}

	return out, results, nil
}

func toSet(keys []string) map[string]bool {
	if len(keys) == 0 {
		return nil
	}
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
