// Package normalizer provides utilities for normalizing .env file entries,
// including key casing, value trimming, and quote stripping.
package normalizer

import (
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Options controls which normalization steps are applied.
type Options struct {
	// KeyToUpper converts all keys to UPPER_CASE.
	KeyToUpper bool
	// StripQuotes removes surrounding single or double quotes from values.
	StripQuotes bool
	// TrimValues trims leading and trailing whitespace from values.
	TrimValues bool
	// CollapseEmptyValues replaces whitespace-only values with an empty string.
	CollapseEmptyValues bool
}

// DefaultOptions returns an Options with all normalization steps enabled.
func DefaultOptions() Options {
	return Options{
		KeyToUpper:          true,
		StripQuotes:         true,
		TrimValues:          true,
		CollapseEmptyValues: true,
	}
}

// Apply normalizes the entries in src according to opts and returns a new EnvFile.
// The original file is not modified.
func Apply(src *parser.EnvFile, opts Options) *parser.EnvFile {
	if src == nil {
		return &parser.EnvFile{}
	}

	out := &parser.EnvFile{}
	for _, e := range src.Entries {
		normalized := normalize(e, opts)
		out.Entries = append(out.Entries, normalized)
	}
	return out
}

func normalize(e parser.Entry, opts Options) parser.Entry {
	key := e.Key
	val := e.Value

	if opts.KeyToUpper {
		key = strings.ToUpper(key)
	}

	if opts.TrimValues {
		val = strings.TrimSpace(val)
	}

	if opts.StripQuotes {
		val = stripQuotes(val)
	}

	if opts.CollapseEmptyValues && strings.TrimSpace(val) == "" {
		val = ""
	}

	return parser.Entry{Key: key, Value: val}
}

func stripQuotes(s string) string {
	if len(s) < 2 {
		return s
	}
	if (s[0] == '"' && s[len(s)-1] == '"') ||
		(s[0] == '\'' && s[len(s)-1] == '\'') {
		return s[1 : len(s)-1]
	}
	return s
}
