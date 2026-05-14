// Package redactor provides functionality to redact sensitive values
// from env files before sharing or exporting.
package redactor

import (
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Options controls redaction behaviour.
type Options struct {
	// Placeholder replaces redacted values. Defaults to "[REDACTED]".
	Placeholder string
	// Patterns is a list of key substrings (case-insensitive) that trigger redaction.
	Patterns []string
}

// defaultPatterns mirrors the masker's sensitive key heuristics.
var defaultPatterns = []string{
	"secret", "password", "passwd", "token", "apikey", "api_key",
	"private", "credential", "auth", "cert", "key",
}

// Apply returns a copy of f with sensitive values replaced by the placeholder.
func Apply(f parser.EnvFile, opts Options) parser.EnvFile {
	if opts.Placeholder == "" {
		opts.Placeholder = "[REDACTED]"
	}
	patterns := opts.Patterns
	if len(patterns) == 0 {
		patterns = defaultPatterns
	}

	result := parser.EnvFile{
		Path:    f.Path,
		Entries: make([]parser.Entry, len(f.Entries)),
	}
	for i, e := range f.Entries {
		if isSensitive(e.Key, patterns) {
			e.Value = opts.Placeholder
			e.Masked = true
		}
		result.Entries[i] = e
	}
	return result
}

// isSensitive returns true when key contains any of the given patterns.
func isSensitive(key string, patterns []string) bool {
	lower := strings.ToLower(key)
	for _, p := range patterns {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}
