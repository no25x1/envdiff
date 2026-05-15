// Package stripper removes entries from an EnvFile based on key patterns or predicates.
package stripper

import (
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Options controls which keys are stripped.
type Options struct {
	// Keys is an explicit list of key names to remove.
	Keys []string
	// Prefixes removes any key that starts with one of these prefixes.
	Prefixes []string
	// Suffixes removes any key that ends with one of these suffixes.
	Suffixes []string
	// EmptyValues removes keys whose value is the empty string.
	EmptyValues bool
}

// Apply returns a new EnvFile with matching entries removed.
func Apply(src *parser.EnvFile, opts Options) *parser.EnvFile {
	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	out := &parser.EnvFile{}
	for _, e := range src.Entries {
		if shouldStrip(e, keySet, opts) {
			continue
		}
		out.Entries = append(out.Entries, e)
	}
	return out
}

func shouldStrip(e parser.Entry, keySet map[string]struct{}, opts Options) bool {
	if _, ok := keySet[e.Key]; ok {
		return true
	}
	for _, p := range opts.Prefixes {
		if strings.HasPrefix(e.Key, p) {
			return true
		}
	}
	for _, s := range opts.Suffixes {
		if strings.HasSuffix(e.Key, s) {
			return true
		}
	}
	if opts.EmptyValues && e.Value == "" {
		return true
	}
	return false
}
