// Package placeholder detects and reports keys whose values are placeholder
// strings (e.g. "changeme", "TODO", "<your-value-here>").
package placeholder

import (
	"regexp"
	"strings"

	"envdiff/internal/parser"
)

// Finding describes a key whose value looks like a placeholder.
type Finding struct {
	Key     string
	Value   string
	Pattern string
}

// DefaultPatterns is the set of patterns used when none are provided.
var DefaultPatterns = []string{
	`(?i)^changeme$`,
	`(?i)^todo$`,
	`(?i)^fixme$`,
	`(?i)^replace.?me$`,
	`(?i)^your.?`,
	`(?i)^<[^>]+>$`,
	`(?i)^\[.*\]$`,
	`(?i)^example$`,
	`(?i)^placeholder$`,
	`(?i)^dummy$`,
}

// Options configures the placeholder detector.
type Options struct {
	// Patterns overrides DefaultPatterns when non-nil.
	Patterns []string
}

// Apply scans f and returns all entries whose values match a placeholder
// pattern. Returns an error if any pattern fails to compile.
func Apply(f *parser.EnvFile, opts Options) ([]Finding, error) {
	if f == nil {
		return nil, nil
	}

	pats := opts.Patterns
	if len(pats) == 0 {
		pats = DefaultPatterns
	}

	compiled := make([]*regexp.Regexp, 0, len(pats))
	for _, p := range pats {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}

	var findings []Finding
	for _, e := range f.Entries {
		v := strings.TrimSpace(e.Value)
		for _, re := range compiled {
			if re.MatchString(v) {
				findings = append(findings, Finding{
					Key:     e.Key,
					Value:   v,
					Pattern: re.String(),
				})
				break
			}
		}
	}
	return findings, nil
}
