// Package obscurer partially masks env values for safe display,
// preserving a configurable number of leading/trailing characters.
package obscurer

import (
	"strings"

	"github.com/envdiff/envdiff/internal/parser"
)

// Options controls how values are obscured.
type Options struct {
	// ShowPrefix is the number of leading characters to reveal (default 2).
	ShowPrefix int
	// ShowSuffix is the number of trailing characters to reveal (default 2).
	ShowSuffix int
	// Mask is the string used to replace hidden characters (default "****").
	Mask string
	// Keys limits obscuring to these specific keys; empty means all keys.
	Keys []string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		ShowPrefix: 2,
		ShowSuffix: 2,
		Mask:       "****",
	}
}

// Apply returns a copy of src with values obscured according to opts.
// Entries whose values are empty are left unchanged.
func Apply(src *parser.EnvFile, opts Options) *parser.EnvFile {
	if src == nil {
		return &parser.EnvFile{}
	}

	allow := toSet(opts.Keys)

	out := &parser.EnvFile{}
	for _, e := range src.Entries {
		copy := e
		if e.Value != "" && (len(allow) == 0 || allow[e.Key]) {
			copy.Value = obscure(e.Value, opts)
		}
		out.Entries = append(out.Entries, copy)
	}
	return out
}

func obscure(val string, opts Options) string {
	n := len(val)
	pre := opts.ShowPrefix
	suf := opts.ShowSuffix

	// If the value is too short to reveal both ends, mask entirely.
	if pre+suf >= n {
		return opts.Mask
	}

	var sb strings.Builder
	sb.WriteString(val[:pre])
	sb.WriteString(opts.Mask)
	sb.WriteString(val[n-suf:])
	return sb.String()
}

func toSet(keys []string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}
