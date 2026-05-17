// Package inspector provides key-level inspection of env files,
// returning metadata about each entry such as type hints, length,
// and whether the value appears sensitive.
package inspector

import (
	"strings"

	"envdiff/internal/parser"
)

// Result holds inspection metadata for a single env entry.
type Result struct {
	Key       string
	Value     string
	Length    int
	IsEmpty   bool
	IsSensitive bool
	TypeHint  string
	Comment   string
}

var sensitivePatterns = []string{
	"secret", "password", "passwd", "token", "key", "api_key",
	"auth", "credential", "private", "cert",
}

// Inspect analyses every entry in f and returns a slice of Results.
func Inspect(f *parser.EnvFile) ([]Result, error) {
	if f == nil {
		return nil, nil
	}
	results := make([]Result, 0, len(f.Entries))
	for _, e := range f.Entries {
		results = append(results, Result{
			Key:         e.Key,
			Value:       e.Value,
			Length:      len(e.Value),
			IsEmpty:     e.Value == "",
			IsSensitive: isSensitive(e.Key),
			TypeHint:    typeHint(e.Value),
			Comment:     e.Comment,
		})
	}
	return results, nil
}

func isSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

func typeHint(v string) string {
	if v == "" {
		return "empty"
	}
	if v == "true" || v == "false" {
		return "bool"
	}
	if isNumeric(v) {
		return "number"
	}
	if strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") {
		return "url"
	}
	return "string"
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for i, c := range s {
		if c == '-' && i == 0 {
			continue
		}
		if c == '.' {
			continue
		}
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
