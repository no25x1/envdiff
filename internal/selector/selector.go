// Package selector provides key-based selection from env files,
// supporting explicit key lists, prefix patterns, and regex matching.
package selector

import (
	"fmt"
	"regexp"
	"strings"

	"envdiff/internal/parser"
)

// Options controls how entries are selected.
type Options struct {
	// Keys is an explicit list of keys to select.
	Keys []string
	// Prefix selects all keys with the given prefix.
	Prefix string
	// Regex selects keys matching the regular expression.
	Regex string
	// Invert returns entries that do NOT match the criteria.
	Invert bool
}

// Apply returns a new EnvFile containing only the entries that match
// the given options. At least one of Keys, Prefix, or Regex must be set.
func Apply(src *parser.EnvFile, opts Options) (*parser.EnvFile, error) {
	if src == nil {
		return &parser.EnvFile{}, nil
	}
	if len(opts.Keys) == 0 && opts.Prefix == "" && opts.Regex == "" {
		return nil, fmt.Errorf("selector: at least one of Keys, Prefix, or Regex must be provided")
	}

	var re *regexp.Regexp
	if opts.Regex != "" {
		var err error
		re, err = regexp.Compile(opts.Regex)
		if err != nil {
			return nil, fmt.Errorf("selector: invalid regex %q: %w", opts.Regex, err)
		}
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	out := &parser.EnvFile{}
	for _, e := range src.Entries {
		matched := matchesEntry(e, keySet, opts.Prefix, re)
		if opts.Invert {
			matched = !matched
		}
		if matched {
			out.Entries = append(out.Entries, e)
		}
	}
	return out, nil
}

func matchesEntry(e parser.Entry, keys map[string]struct{}, prefix string, re *regexp.Regexp) bool {
	if _, ok := keys[e.Key]; ok {
		return true
	}
	if prefix != "" && strings.HasPrefix(e.Key, prefix) {
		return true
	}
	if re != nil && re.MatchString(e.Key) {
		return true
	}
	return false
}
