// Package capitalizer provides utilities for normalizing the case of
// environment variable keys and values within an EnvFile.
package capitalizer

import (
	"strings"

	"github.com/envdiff/envdiff/internal/parser"
)

// Options controls the behaviour of Apply.
type Options struct {
	// KeysToUpper converts all keys to UPPER_CASE.
	KeysToUpper bool
	// KeysToLower converts all keys to lower_case.
	KeysToLower bool
	// ValuesToUpper converts all values to upper-case.
	ValuesToUpper bool
	// ValuesToLower converts all values to lower-case.
	ValuesToLower bool
	// OnlyKeys restricts key transformations to this set (empty = all keys).
	OnlyKeys []string
}

// Apply returns a new EnvFile with key/value capitalisation applied
// according to opts. The original file is never modified.
func Apply(f *parser.EnvFile, opts Options) *parser.EnvFile {
	if f == nil {
		return &parser.EnvFile{}
	}

	allowed := toSet(opts.OnlyKeys)

	out := &parser.EnvFile{}
	for _, e := range f.Entries {
		entry := e

		if len(allowed) == 0 || allowed[e.Key] {
			if opts.KeysToUpper {
				entry.Key = strings.ToUpper(e.Key)
			} else if opts.KeysToLower {
				entry.Key = strings.ToLower(e.Key)
			}

			if opts.ValuesToUpper {
				entry.Value = strings.ToUpper(e.Value)
			} else if opts.ValuesToLower {
				entry.Value = strings.ToLower(e.Value)
			}
		}

		out.Entries = append(out.Entries, entry)
	}
	return out
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
