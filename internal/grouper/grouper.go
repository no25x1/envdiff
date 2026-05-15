// Package grouper provides functionality to group env file entries
// by a common prefix delimiter, returning named sections.
package grouper

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Group represents a named collection of env entries sharing a common prefix.
type Group struct {
	Name    string
	Entries []parser.Entry
}

// Options controls how grouping is performed.
type Options struct {
	// Delimiter separates the prefix from the rest of the key (default: "_").
	Delimiter string
	// IncludeUngrouped places keys with no delimiter into a group named "".
	IncludeUngrouped bool
}

// Apply groups the entries in f by the leading prefix before the first
// occurrence of Delimiter. Entries that contain no delimiter are placed
// in the ungrouped bucket when IncludeUngrouped is true, otherwise skipped.
func Apply(f *parser.EnvFile, opts Options) []Group {
	if opts.Delimiter == "" {
		opts.Delimiter = "_"
	}

	buckets := make(map[string][]parser.Entry)

	for _, e := range f.Entries {
		idx := strings.Index(e.Key, opts.Delimiter)
		if idx < 0 {
			if opts.IncludeUngrouped {
				buckets[""] = append(buckets[""], e)
			}
			continue
		}
		prefix := e.Key[:idx]
		buckets[prefix] = append(buckets[prefix], e)
	}

	names := make([]string, 0, len(buckets))
	for k := range buckets {
		names = append(names, k)
	}
	sort.Strings(names)

	groups := make([]Group, 0, len(names))
	for _, name := range names {
		groups = append(groups, Group{
			Name:    name,
			Entries: buckets[name],
		})
	}
	return groups
}
