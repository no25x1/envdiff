// Package masker provides utilities for detecting and masking sensitive
// environment variable values before display or export.
package masker

import (
	"regexp"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// DefaultPlaceholder is used when a sensitive value is masked.
const DefaultPlaceholder = "***"

// DefaultPatterns are regex patterns matched against keys to detect secrets.
var DefaultPatterns = []string{
	`(?i)password`,
	`(?i)secret`,
	`(?i)token`,
	`(?i)api[_-]?key`,
	`(?i)private[_-]?key`,
	`(?i)auth`,
	`(?i)credential`,
}

// Masker masks sensitive values in env files.
type Masker struct {
	patterns    []*regexp.Regexp
	placeholder string
}

// New creates a Masker with default patterns and placeholder.
func New() *Masker {
	return NewWithOptions(DefaultPatterns, DefaultPlaceholder)
}

// NewWithOptions creates a Masker with custom patterns and placeholder.
func NewWithOptions(patterns []string, placeholder string) *Masker {
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		if re, err := regexp.Compile(p); err == nil {
			compiled = append(compiled, re)
		}
	}
	return &Masker{patterns: compiled, placeholder: placeholder}
}

// IsSensitive returns true if the key matches any sensitive pattern.
func (m *Masker) IsSensitive(key string) bool {
	for _, re := range m.patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

// MaskValue returns the placeholder if the key is sensitive, otherwise the value.
func (m *Masker) MaskValue(key, value string) string {
	if m.IsSensitive(key) {
		return m.placeholder
	}
	return value
}

// MaskEnv returns a copy of the env file with sensitive values replaced.
func (m *Masker) MaskEnv(f *parser.EnvFile) *parser.EnvFile {
	if f == nil {
		return nil
	}
	out := &parser.EnvFile{Path: f.Path}
	for _, e := range f.Entries {
		masked := e
		if m.IsSensitive(strings.TrimSpace(e.Key)) {
			masked.Value = m.placeholder
		}
		out.Entries = append(out.Entries, masked)
	}
	return out
}
