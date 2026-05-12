// Package filter provides functionality to filter env file entries
// by key prefix, suffix, or pattern matching.
package filter

import (
	"regexp"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Options holds the filtering criteria.
type Options struct {
	// Prefix filters keys that start with the given prefix.
	Prefix string
	// Suffix filters keys that end with the given suffix.
	Suffix string
	// Pattern filters keys matching the given regular expression.
	Pattern string
	// Exclude inverts the filter, keeping entries that do NOT match.
	Exclude bool
}

// Apply returns a new EnvFile containing only the entries that match
// the given Options. If no criteria are set, the original entries are
// returned unchanged.
func Apply(file parser.EnvFile, opts Options) (parser.EnvFile, error) {
	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return parser.EnvFile{}, err
		}
	}

	filtered := make([]parser.Entry, 0, len(file.Entries))
	for _, entry := range file.Entries {
		matched := matches(entry.Key, opts.Prefix, opts.Suffix, re)
		if opts.Exclude {
			matched = !matched
		}
		if matched {
			filtered = append(filtered, entry)
		}
	}

	return parser.EnvFile{
		Path:    file.Path,
		Entries: filtered,
	}, nil
}

// matches returns true if key satisfies the provided criteria.
// When no criteria are specified (all zero values), every key matches.
func matches(key, prefix, suffix string, re *regexp.Regexp) bool {
	if prefix == "" && suffix == "" && re == nil {
		return true
	}
	if prefix != "" && !strings.HasPrefix(key, prefix) {
		return false
	}
	if suffix != "" && !strings.HasSuffix(key, suffix) {
		return false
	}
	if re != nil && !re.MatchString(key) {
		return false
	}
	return true
}
