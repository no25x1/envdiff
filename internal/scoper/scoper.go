// Package scoper restricts an EnvFile to a named scope (prefix group),
// stripping the scope prefix from keys so callers work with clean names.
package scoper

import (
	"fmt"
	"strings"

	"envdiff/internal/parser"
)

// Options controls scoping behaviour.
type Options struct {
	// Scope is the prefix that identifies the target group, e.g. "PROD".
	// Keys must match "<Scope>_" to be included.
	Scope string

	// StripPrefix removes the "<Scope>_" prefix from keys in the result.
	StripPrefix bool

	// KeepUnscoped includes keys that carry no recognised scope prefix.
	KeepUnscoped bool
}

// Apply returns a new EnvFile containing only the entries that belong to
// the given scope. If opts.StripPrefix is true the scope prefix is removed
// from each key in the returned file.
func Apply(src *parser.EnvFile, opts Options) (*parser.EnvFile, error) {
	if src == nil {
		return nil, fmt.Errorf("scoper: source file is nil")
	}
	if strings.TrimSpace(opts.Scope) == "" {
		return nil, fmt.Errorf("scoper: scope must not be empty")
	}

	prefix := strings.ToUpper(strings.TrimSpace(opts.Scope)) + "_"
	out := &parser.EnvFile{}

	for _, e := range src.Entries {
		upper := strings.ToUpper(e.Key)
		if strings.HasPrefix(upper, prefix) {
			key := e.Key
			if opts.StripPrefix {
				key = e.Key[len(prefix):]
			}
			out.Entries = append(out.Entries, parser.Entry{
				Key:     key,
				Value:   e.Value,
				Comment: e.Comment,
			})
			continue
		}
		if opts.KeepUnscoped && !hasAnyScope(e.Key, src) {
			out.Entries = append(out.Entries, e)
		}
	}
	return out, nil
}

// hasAnyScope returns true when key contains an underscore-separated uppercase
// prefix that appears as a scope in other keys of the file.
func hasAnyScope(key string, src *parser.EnvFile) bool {
	idx := strings.Index(key, "_")
	if idx < 0 {
		return false
	}
	candidate := strings.ToUpper(key[:idx]) + "_"
	for _, e := range src.Entries {
		if e.Key == key {
			continue
		}
		if strings.HasPrefix(strings.ToUpper(e.Key), candidate) {
			return true
		}
	}
	return false
}
