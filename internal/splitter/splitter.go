// Package splitter splits an EnvFile into multiple files based on key prefix groups.
package splitter

import (
	"fmt"
	"strings"

	"envdiff/internal/parser"
)

// Options controls how the split is performed.
type Options struct {
	// Prefixes maps a group name to the key prefix it captures.
	// e.g. {"db": "DB_", "aws": "AWS_"}
	Prefixes map[string]string

	// IncludeUnmatched places keys that match no prefix into a group named "other".
	IncludeUnmatched bool
}

// Result holds the split output.
type Result map[string]*parser.EnvFile

// Apply splits src into groups according to opts.
// Keys are matched case-insensitively against prefixes.
// The first matching prefix wins.
func Apply(src *parser.EnvFile, opts Options) (Result, error) {
	if src == nil {
		return nil, fmt.Errorf("splitter: source file is nil")
	}
	if len(opts.Prefixes) == 0 {
		return nil, fmt.Errorf("splitter: no prefixes defined")
	}

	// Build ordered group names for deterministic first-match.
	result := make(Result)
	for name := range opts.Prefixes {
		result[name] = &parser.EnvFile{}
	}
	if opts.IncludeUnmatched {
		result["other"] = &parser.EnvFile{}
	}

	for _, entry := range src.Entries {
		matched := false
		for name, prefix := range opts.Prefixes {
			if strings.HasPrefix(strings.ToUpper(entry.Key), strings.ToUpper(prefix)) {
				result[name].Entries = append(result[name].Entries, entry)
				matched = true
				break
			}
		}
		if !matched && opts.IncludeUnmatched {
			result["other"].Entries = append(result["other"].Entries, entry)
		}
	}

	return result, nil
}
