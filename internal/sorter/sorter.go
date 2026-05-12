// Package sorter provides utilities for sorting and grouping .env file entries.
package sorter

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Options controls how entries are sorted.
type Options struct {
	// Alphabetical sorts keys A→Z when true.
	Alphabetical bool
	// GroupByPrefix groups keys sharing the same prefix (up to the first '_').
	GroupByPrefix bool
	// Reverse inverts the final order.
	Reverse bool
}

// Apply returns a new EnvFile whose entries are sorted according to opts.
func Apply(f parser.EnvFile, opts Options) parser.EnvFile {
	entries := make([]parser.Entry, len(f.Entries))
	copy(entries, f.Entries)

	if opts.Alphabetical {
		sort.SliceStable(entries, func(i, j int) bool {
			return entries[i].Key < entries[j].Key
		})
	}

	if opts.GroupByPrefix {
		entries = groupByPrefix(entries)
	}

	if opts.Reverse {
		for l, r := 0, len(entries)-1; l < r; l, r = l+1, r-1 {
			entries[l], entries[r] = entries[r], entries[l]
		}
	}

	return parser.EnvFile{Path: f.Path, Entries: entries}
}

// groupByPrefix stable-sorts entries so that keys sharing the same prefix
// (the segment before the first '_') are contiguous.
func groupByPrefix(entries []parser.Entry) []parser.Entry {
	prefixOf := func(key string) string {
		if idx := strings.Index(key, "_"); idx > 0 {
			return key[:idx]
		}
		return key
	}

	// Build insertion-order prefix list to preserve relative group order.
	seen := make(map[string]int)
	order := []string{}
	for _, e := range entries {
		p := prefixOf(e.Key)
		if _, ok := seen[p]; !ok {
			seen[p] = len(order)
			order = append(order, p)
		}
	}

	sort.SliceStable(entries, func(i, j int) bool {
		pi := seen[prefixOf(entries[i].Key)]
		pj := seen[prefixOf(entries[j].Key)]
		return pi < pj
	})

	return entries
}
