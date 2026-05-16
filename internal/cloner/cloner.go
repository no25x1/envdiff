// Package cloner copies entries from one env file into another,
// optionally renaming keys via a prefix swap and filtering by key pattern.
package cloner

import (
	"fmt"
	"strings"

	"envdiff/internal/parser"
)

// Options controls how entries are cloned.
type Options struct {
	// StripPrefix removes this prefix from source keys before writing.
	StripPrefix string
	// AddPrefix prepends this prefix to destination keys.
	AddPrefix string
	// Keys is an explicit allow-list of source key names (pre-strip).
	// When empty, all keys are cloned.
	Keys []string
	// Overwrite replaces existing keys in the destination.
	Overwrite bool
}

// Apply clones entries from src into dst according to opts.
// It returns a new EnvFile containing the merged result.
func Apply(src, dst *parser.EnvFile, opts Options) (*parser.EnvFile, error) {
	if src == nil {
		return nil, fmt.Errorf("cloner: source file must not be nil")
	}
	if dst == nil {
		dst = &parser.EnvFile{}
	}

	allowSet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		allowSet[k] = true
	}

	// Build a lookup of existing destination keys.
	existing := make(map[string]bool, len(dst.Entries))
	for _, e := range dst.Entries {
		existing[e.Key] = true
	}

	out := &parser.EnvFile{}
	out.Entries = append(out.Entries, dst.Entries...)

	for _, e := range src.Entries {
		if len(allowSet) > 0 && !allowSet[e.Key] {
			continue
		}

		destKey := e.Key
		if opts.StripPrefix != "" {
			if !strings.HasPrefix(destKey, opts.StripPrefix) {
				continue
			}
			destKey = strings.TrimPrefix(destKey, opts.StripPrefix)
		}
		destKey = opts.AddPrefix + destKey

		if existing[destKey] && !opts.Overwrite {
			continue
		}

		cloned := parser.Entry{
			Key:     destKey,
			Value:   e.Value,
			Comment: e.Comment,
		}

		if existing[destKey] {
			// Replace in-place.
			for i, oe := range out.Entries {
				if oe.Key == destKey {
					out.Entries[i] = cloned
					break
				}
			}
		} else {
			out.Entries = append(out.Entries, cloned)
			existing[destKey] = true
		}
	}

	return out, nil
}
