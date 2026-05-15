// Package promote provides utilities for promoting .env values
// from one environment to another, with optional key filtering and masking.
package promote

import (
	"fmt"

	"github.com/user/envdiff/internal/parser"
)

// Options controls which keys are promoted and how.
type Options struct {
	// Keys is an explicit allowlist of keys to promote.
	// If empty, all keys from Source are candidates.
	Keys []string

	// Overwrite controls whether existing keys in Target are replaced.
	Overwrite bool
}

// Result describes what happened to a single key during promotion.
type Result struct {
	Key    string
	Action string // "promoted", "skipped_exists", "skipped_filter"
}

// Apply promotes entries from src into dst according to opts.
// It returns a new EnvFile representing the merged result and a
// slice of Result records describing each decision.
func Apply(src, dst parser.EnvFile, opts Options) (parser.EnvFile, []Result, error) {
	if src.Path == "" {
		return parser.EnvFile{}, nil, fmt.Errorf("promote: source path must not be empty")
	}

	allowSet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		allowSet[k] = true
	}

	// Build a mutable copy of dst entries indexed by key.
	dstIndex := make(map[string]int, len(dst.Entries))
	out := make([]parser.Entry, len(dst.Entries))
	copy(out, dst.Entries)
	for i, e := range out {
		dstIndex[e.Key] = i
	}

	var results []Result

	for _, e := range src.Entries {
		if len(allowSet) > 0 && !allowSet[e.Key] {
			results = append(results, Result{Key: e.Key, Action: "skipped_filter"})
			continue
		}
		if idx, exists := dstIndex[e.Key]; exists {
			if !opts.Overwrite {
				results = append(results, Result{Key: e.Key, Action: "skipped_exists"})
				continue
			}
			out[idx] = e
		} else {
			dstIndex[e.Key] = len(out)
			out = append(out, e)
		}
		results = append(results, Result{Key: e.Key, Action: "promoted"})
	}

	return parser.EnvFile{Path: dst.Path, Entries: out}, results, nil
}
