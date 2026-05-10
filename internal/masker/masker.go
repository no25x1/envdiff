package masker

import (
	"strings"
)

// DefaultSecretPatterns holds common key substrings that indicate sensitive values.
var DefaultSecretPatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"APIKEY",
	"PRIVATE",
	"CREDENTIALS",
	"AUTH",
	"ACCESS_KEY",
	"SIGNING_KEY",
}

const MaskedValue = "***MASKED***"

// Masker masks sensitive values in env maps based on key patterns.
type Masker struct {
	patterns []string
}

// New creates a Masker with the provided secret key patterns.
// If no patterns are provided, DefaultSecretPatterns are used.
func New(patterns ...string) *Masker {
	if len(patterns) == 0 {
		patterns = DefaultSecretPatterns
	}
	return &Masker{patterns: patterns}
}

// IsSensitive returns true if the key matches any known secret pattern.
func (m *Masker) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range m.patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

// MaskEnv returns a copy of the env map with sensitive values replaced.
func (m *Masker) MaskEnv(env map[string]string) map[string]string {
	masked := make(map[string]string, len(env))
	for k, v := range env {
		if m.IsSensitive(k) {
			masked[k] = MaskedValue
		} else {
			masked[k] = v
		}
	}
	return masked
}

// MaskValue returns the masked constant if the key is sensitive, otherwise the original value.
func (m *Masker) MaskValue(key, value string) string {
	if m.IsSensitive(key) {
		return MaskedValue
	}
	return value
}
