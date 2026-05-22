// Package highlighter applies ANSI colour highlighting to diff output,
// making added, removed, and modified entries visually distinct in terminals.
package highlighter

import (
	"fmt"
	"io"
	"strings"

	"github.com/envdiff/envdiff/internal/diff"
)

// ANSI escape codes.
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

// Options controls highlighting behaviour.
type Options struct {
	// NoColor disables ANSI codes (useful when output is not a TTY).
	NoColor bool
}

// DefaultOptions returns Options with colour enabled.
func DefaultOptions() Options {
	return Options{NoColor: false}
}

// Write renders diff results with ANSI colour highlighting to w.
func Write(w io.Writer, results []diff.Result, opts Options) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "No differences found.")
		return err
	}

	for _, r := range results {
		line := formatResult(r, opts)
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func formatResult(r diff.Result, opts Options) string {
	switch r.Status {
	case diff.Added:
		return colorize(opts, colorGreen, fmt.Sprintf("+ %s=%s", r.Key, r.NewValue))
	case diff.Removed:
		return colorize(opts, colorRed, fmt.Sprintf("- %s=%s", r.Key, r.OldValue))
	case diff.Modified:
		old := colorize(opts, colorRed, r.OldValue)
		new := colorize(opts, colorGreen, r.NewValue)
		key := colorize(opts, colorYellow, r.Key)
		return fmt.Sprintf("~ %s: %s → %s", key, old, new)
	case diff.Unchanged:
		return colorize(opts, colorCyan, fmt.Sprintf("  %s=%s", r.Key, r.OldValue))
	default:
		return fmt.Sprintf("  %s", r.Key)
	}
}

func colorize(opts Options, code, text string) string {
	if opts.NoColor || strings.TrimSpace(text) == "" {
		return text
	}
	return code + text + colorReset
}
