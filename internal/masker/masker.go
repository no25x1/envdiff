// Package masker provides secret masking for .env file values.
package masker

import (
	"regexp"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

const DefaultPlaceholder = "***"

// defaultSensitivePatterns are case-insensitive key patterns that indicate secrets.
var defaultSensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)password`),
	regexp.MustCompile(`(?i)secret`),
	regexp.MustCompile(`(?i)token`),
	regexp.MustCompile(`(?i)api_key`),
	regexp.MustCompile(`(?i)private_key`),
	regexp.MustCompile(`(?i)auth`),
	regexp.MustCompile(`(?i)credential`),
}

// Options configures masking behaviour.
type Options struct {
	Patterns    []*regexp.Regexp
	Placeholder string
}

// Masker masks sensitive values in env files.
type Masker struct {
	opts Options
}

// New returns a Masker with default sensitive patterns.
func New() *Masker {
	return NewWithOptions(Options{
		Patterns:    defaultSensitivePatterns,
		Placeholder: DefaultPlaceholder,
	})
}

// NewWithOptions returns a Masker configured with the supplied options.
func NewWithOptions(opts Options) *Masker {
	if opts.Placeholder == "" {
		opts.Placeholder = DefaultPlaceholder
	}
	if len(opts.Patterns) == 0 {
		opts.Patterns = defaultSensitivePatterns
	}
	return &Masker{opts: opts}
}

// IsSensitive reports whether key matches any sensitive pattern.
func (m *Masker) IsSensitive(key string) bool {
	for _, p := range m.opts.Patterns {
		if p.MatchString(strings.ToLower(key)) {
			return true
		}
	}
	return false
}

// MaskValue returns the placeholder if key is sensitive, otherwise value.
func (m *Masker) MaskValue(key, value string) string {
	if m.IsSensitive(key) {
		return m.opts.Placeholder
	}
	return value
}

// MaskEnv returns a copy of f with sensitive values replaced by the placeholder.
func (m *Masker) MaskEnv(f *parser.EnvFile) *parser.EnvFile {
	if f == nil {
		return nil
	}
	out := &parser.EnvFile{Path: f.Path}
	for _, e := range f.Entries {
		masked := e
		masked.Value = m.MaskValue(e.Key, e.Value)
		out.Entries = append(out.Entries, masked)
	}
	return out
}
