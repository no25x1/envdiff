package masker

import (
	"regexp"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

const defaultPlaceholder = "***"

var defaultSensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)password`),
	regexp.MustCompile(`(?i)secret`),
	regexp.MustCompile(`(?i)token`),
	regexp.MustCompile(`(?i)api_key`),
	regexp.MustCompile(`(?i)private_key`),
	regexp.MustCompile(`(?i)auth`),
	regexp.MustCompile(`(?i)credential`),
}

// Masker masks sensitive values in env files.
type Masker struct {
	patterns    []*regexp.Regexp
	placeholder string
}

// New returns a Masker with default sensitive patterns.
func New() *Masker {
	return &Masker{
		patterns:    defaultSensitivePatterns,
		placeholder: defaultPlaceholder,
	}
}

// NewWithOptions returns a Masker with custom patterns and placeholder.
func NewWithOptions(patterns []*regexp.Regexp, placeholder string) *Masker {
	if len(patterns) == 0 {
		patterns = defaultSensitivePatterns
	}
	if placeholder == "" {
		placeholder = defaultPlaceholder
	}
	return &Masker{patterns: patterns, placeholder: placeholder}
}

// IsSensitive reports whether the given key matches any sensitive pattern.
func (m *Masker) IsSensitive(key string) bool {
	for _, p := range m.patterns {
		if p.MatchString(strings.ToLower(key)) {
			return true
		}
	}
	return false
}

// MaskValue returns the placeholder if the key is sensitive, otherwise the original value.
func (m *Masker) MaskValue(key, value string) string {
	if m.IsSensitive(key) {
		return m.placeholder
	}
	return value
}

// MaskEnv returns a copy of the env file with sensitive values replaced.
func (m *Masker) MaskEnv(f *parser.EnvFile) *parser.EnvFile {
	if f == nil {
		return &parser.EnvFile{}
	}
	out := &parser.EnvFile{Path: f.Path}
	for _, e := range f.Entries {
		masked := e
		masked.Value = m.MaskValue(e.Key, e.Value)
		out.Entries = append(out.Entries, masked)
	}
	return out
}
