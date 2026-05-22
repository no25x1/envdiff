// Package blocker prevents forbidden keys or values from appearing in env files.
package blocker

import (
	"fmt"
	"strings"

	"envdiff/internal/parser"
)

// Violation describes a blocked entry.
type Violation struct {
	Key    string
	Reason string
}

// Options controls which keys and values are forbidden.
type Options struct {
	// ForbiddenKeys is a list of exact key names that must not appear.
	ForbiddenKeys []string
	// ForbiddenPrefixes blocks any key whose name starts with one of these.
	ForbiddenPrefixes []string
	// ForbiddenValues blocks any entry whose value matches one of these exactly.
	ForbiddenValues []string
}

// Apply checks f against the options and returns any violations found.
// A nil file returns nil, nil.
func Apply(f *parser.EnvFile, opts Options) ([]Violation, error) {
	if f == nil {
		return nil, nil
	}

	keySet := toSet(opts.ForbiddenKeys)
	valSet := toSet(opts.ForbiddenValues)

	var violations []Violation

	for _, e := range f.Entries {
		// Exact key match.
		if keySet[e.Key] {
			violations = append(violations, Violation{
				Key:    e.Key,
				Reason: fmt.Sprintf("key %q is forbidden", e.Key),
			})
			continue
		}

		// Prefix match.
		for _, pfx := range opts.ForbiddenPrefixes {
			if strings.HasPrefix(e.Key, pfx) {
				violations = append(violations, Violation{
					Key:    e.Key,
					Reason: fmt.Sprintf("key %q matches forbidden prefix %q", e.Key, pfx),
				})
				break
			}
		}

		// Forbidden value match.
		if valSet[e.Value] {
			violations = append(violations, Violation{
				Key:    e.Key,
				Reason: fmt.Sprintf("value of key %q is forbidden", e.Key),
			})
		}
	}

	return violations, nil
}

func toSet(items []string) map[string]bool {
	s := make(map[string]bool, len(items))
	for _, v := range items {
		s[v] = true
	}
	return s
}
