// Package deduplicator removes duplicate keys from an EnvFile,
// keeping either the first or last occurrence based on options.
package deduplicator

import "github.com/user/envdiff/internal/parser"

// Strategy controls which occurrence of a duplicate key is kept.
type Strategy string

const (
	// KeepFirst retains the first occurrence of each key.
	KeepFirst Strategy = "first"
	// KeepLast retains the last occurrence of each key.
	KeepLast Strategy = "last"
)

// Options configures deduplication behaviour.
type Options struct {
	Strategy Strategy
}

// DefaultOptions returns sensible defaults (keep last, matching shell semantics).
func DefaultOptions() Options {
	return Options{Strategy: KeepLast}
}

// Result holds the deduplicated file and a report of removed entries.
type Result struct {
	File    *parser.EnvFile
	Removed []Duplicate
}

// Duplicate records a key that appeared more than once.
type Duplicate struct {
	Key   string
	Count int
}

// Apply deduplicates src according to opts and returns a Result.
func Apply(src *parser.EnvFile, opts Options) Result {
	if opts.Strategy == "" {
		opts = DefaultOptions()
	}

	// Count occurrences.
	counts := make(map[string]int)
	for _, e := range src.Entries {
		counts[e.Key]++
	}

	seen := make(map[string]bool)
	var kept []parser.Entry

	entries := src.Entries
	if opts.Strategy == KeepLast {
		// Iterate in reverse so we mark the last occurrence first.
		reversed := make([]parser.Entry, len(entries))
		copy(reversed, entries)
		for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
			reversed[i], reversed[j] = reversed[j], reversed[i]
		}
		for _, e := range reversed {
			if !seen[e.Key] {
				seen[e.Key] = true
				kept = append([]parser.Entry{e}, kept...)
			}
		}
	} else {
		for _, e := range entries {
			if !seen[e.Key] {
				seen[e.Key] = true
				kept = append(kept, e)
			}
		}
	}

	var dups []Duplicate
	for k, c := range counts {
		if c > 1 {
			dups = append(dups, Duplicate{Key: k, Count: c})
		}
	}

	return Result{
		File:    &parser.EnvFile{Path: src.Path, Entries: kept},
		Removed: dups,
	}
}
