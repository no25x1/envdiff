// Package interpolator resolves variable references within .env file values.
// It expands expressions like ${OTHER_KEY} or $OTHER_KEY using values defined
// in the same file or provided via an override map.
package interpolator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

var refPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

// Options controls interpolation behaviour.
type Options struct {
	// Overrides are resolved before the file's own keys.
	Overrides map[string]string
	// FailOnMissing returns an error if a referenced key cannot be resolved.
	FailOnMissing bool
}

// Apply resolves variable references in every entry of f and returns a new
// EnvFile with the expanded values. Entries whose values contain no references
// are returned unchanged.
func Apply(f parser.EnvFile, opts Options) (parser.EnvFile, error) {
	// Build a lookup map: file values first, then overrides win.
	lookup := make(map[string]string, len(f.Entries))
	for _, e := range f.Entries {
		lookup[e.Key] = e.Value
	}
	for k, v := range opts.Overrides {
		lookup[k] = v
	}

	out := parser.EnvFile{Path: f.Path}
	for _, e := range f.Entries {
		expanded, err := expand(e.Value, lookup, opts.FailOnMissing)
		if err != nil {
			return parser.EnvFile{}, fmt.Errorf("interpolator: key %q: %w", e.Key, err)
		}
		out.Entries = append(out.Entries, parser.Entry{
			Key:     e.Key,
			Value:   expanded,
			Comment: e.Comment,
		})
	}
	return out, nil
}

func expand(value string, lookup map[string]string, failOnMissing bool) (string, error) {
	var expandErr error
	result := refPattern.ReplaceAllStringFunc(value, func(match string) string {
		if expandErr != nil {
			return match
		}
		name := extractName(match)
		if v, ok := lookup[name]; ok {
			return v
		}
		if failOnMissing {
			expandErr = fmt.Errorf("unresolved reference %q", match)
			return match
		}
		return ""
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}

// extractName strips the sigil and optional braces from a matched reference.
// For example, "${FOO}" and "$FOO" both return "FOO".
func extractName(match string) string {
	match = strings.TrimPrefix(match, "$")
	match = strings.TrimPrefix(match, "{")
	match = strings.TrimSuffix(match, "}")
	return match
}

// References returns the list of unique variable names referenced in value,
// in the order they first appear. This can be used to inspect dependencies
// between entries without performing a full expansion.
func References(value string) []string {
	seen := make(map[string]bool)
	var refs []string
	for _, match := range refPattern.FindAllString(value, -1) {
		name := extractName(match)
		if !seen[name] {
			seen[name] = true
			refs = append(refs, name)
		}
	}
	return refs
}
