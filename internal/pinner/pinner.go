// Package pinner provides functionality to pin (freeze) env file values,
// producing a snapshot of current values that can be committed or compared later.
package pinner

import (
	"fmt"
	"sort"
	"strings"

	"envdiff/internal/parser"
)

// Options controls pinning behaviour.
type Options struct {
	// Keys restricts pinning to the listed keys. Empty means all keys.
	Keys []string
	// Placeholder is the value written for empty entries when PinEmpty is false.
	Placeholder string
	// PinEmpty includes entries with empty values when true.
	PinEmpty bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Placeholder: "CHANGE_ME",
		PinEmpty:    false,
	}
}

// Result holds a single pinned entry.
type Result struct {
	Key      string
	Value    string
	Pinned   bool
	Skipped  bool
	Reason   string
}

// Apply pins values from src according to opts.
// It returns a new EnvFile containing only pinned entries and a slice of Results.
func Apply(src *parser.EnvFile, opts Options) (*parser.EnvFile, []Result, error) {
	if src == nil {
		return nil, nil, fmt.Errorf("pinner: source file is nil")
	}

	allowSet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		allowSet[strings.TrimSpace(k)] = struct{}{}
	}

	var results []Result
	var entries []parser.Entry

	for _, e := range src.Entries {
		if len(allowSet) > 0 {
			if _, ok := allowSet[e.Key]; !ok {
				results = append(results, Result{Key: e.Key, Skipped: true, Reason: "not in key list"})
				continue
			}
		}

		val := e.Value
		if val == "" && !opts.PinEmpty {
			if opts.Placeholder != "" {
				val = opts.Placeholder
			} else {
				results = append(results, Result{Key: e.Key, Skipped: true, Reason: "empty value"})
				continue
			}
		}

		pinned := parser.Entry{
			Key:     e.Key,
			Value:   val,
			Comment: e.Comment,
		}
		entries = append(entries, pinned)
		results = append(results, Result{Key: e.Key, Value: val, Pinned: true})
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].Key < entries[j].Key })

	return &parser.EnvFile{Entries: entries}, results, nil
}
