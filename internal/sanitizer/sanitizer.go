// Package sanitizer provides utilities for cleaning and normalising
// .env file values before export or comparison.
package sanitizer

import (
	"strings"

	"github.com/envdiff/envdiff/internal/parser"
)

// Options controls which sanitisation passes are applied.
type Options struct {
	// TrimWhitespace removes leading/trailing whitespace from values.
	TrimWhitespace bool
	// RemoveNullBytes strips null bytes (\x00) from values.
	RemoveNullBytes bool
	// NormaliseLineEndings converts \r\n and bare \r to \n inside values.
	NormaliseLineEndings bool
	// CollapseBlankValues replaces values that consist solely of whitespace
	// with an empty string.
	CollapseBlankValues bool
}

// DefaultOptions returns a sensible default sanitisation configuration.
func DefaultOptions() Options {
	return Options{
		TrimWhitespace:       true,
		RemoveNullBytes:      true,
		NormaliseLineEndings: true,
		CollapseBlankValues:  true,
	}
}

// Apply runs the configured sanitisation passes over every entry in f and
// returns a new EnvFile with the cleaned values. The original file is not
// modified.
func Apply(f parser.EnvFile, opts Options) parser.EnvFile {
	out := parser.EnvFile{
		Path:    f.Path,
		Entries: make([]parser.Entry, 0, len(f.Entries)),
	}

	for _, e := range f.Entries {
		v := e.Value

		if opts.NormaliseLineEndings {
			v = strings.ReplaceAll(v, "\r\n", "\n")
			v = strings.ReplaceAll(v, "\r", "\n")
		}
		if opts.RemoveNullBytes {
			v = strings.ReplaceAll(v, "\x00", "")
		}
		if opts.TrimWhitespace {
			v = strings.TrimSpace(v)
		}
		if opts.CollapseBlankValues && strings.TrimSpace(v) == "" {
			v = ""
		}

		out.Entries = append(out.Entries, parser.Entry{
			Key:   e.Key,
			Value: v,
			Raw:   e.Raw,
		})
	}

	return out
}
