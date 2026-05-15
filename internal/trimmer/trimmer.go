// Package trimmer removes trailing whitespace and normalises blank lines
// in parsed env files.
package trimmer

import (
	"strings"

	"github.com/your-org/envdiff/internal/parser"
)

// Options controls trimmer behaviour.
type Options struct {
	// TrimValues removes leading and trailing whitespace from entry values.
	TrimValues bool
	// TrimKeys removes leading and trailing whitespace from entry keys.
	TrimKeys bool
	// RemoveEmptyValues drops entries whose value is empty after trimming.
	RemoveEmptyValues bool
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		TrimValues:        true,
		TrimKeys:          true,
		RemoveEmptyValues: false,
	}
}

// Apply returns a new EnvFile with whitespace normalised according to opts.
// The original file is not modified.
func Apply(file parser.EnvFile, opts Options) parser.EnvFile {
	out := make([]parser.Entry, 0, len(file.Entries))

	for _, e := range file.Entries {
		if opts.TrimKeys {
			e.Key = strings.TrimSpace(e.Key)
		}
		if opts.TrimValues {
			e.Value = strings.TrimSpace(e.Value)
		}
		if opts.RemoveEmptyValues && e.Value == "" {
			continue
		}
		out = append(out, e)
	}

	return parser.EnvFile{
		Path:    file.Path,
		Entries: out,
	}
}
