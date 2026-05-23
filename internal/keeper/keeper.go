// Package keeper provides functionality to selectively retain only the
// specified keys from an EnvFile, discarding all others.
package keeper

import (
	"errors"
	"fmt"

	"envdiff/internal/parser"
)

// Options controls which keys are retained.
type Options struct {
	// Keys is the explicit list of keys to keep.
	Keys []string
	// Prefix retains all keys that start with the given prefix.
	Prefix string
}

// Apply returns a new EnvFile containing only the entries matched by opts.
// At least one of Keys or Prefix must be non-empty.
func Apply(src *parser.EnvFile, opts Options) (*parser.EnvFile, error) {
	if src == nil {
		return &parser.EnvFile{}, nil
	}
	if len(opts.Keys) == 0 && opts.Prefix == "" {
		return nil, errors.New("keeper: at least one of Keys or Prefix must be set")
	}

	allowSet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		allowSet[k] = struct{}{}
	}

	out := &parser.EnvFile{}
	for _, e := range src.Entries {
		if _, ok := allowSet[e.Key]; ok {
			out.Entries = append(out.Entries, e)
			continue
		}
		if opts.Prefix != "" && len(e.Key) >= len(opts.Prefix) && e.Key[:len(opts.Prefix)] == opts.Prefix {
			out.Entries = append(out.Entries, e)
		}
	}
	return out, nil
}

// Count returns the number of entries that would be retained.
func Count(src *parser.EnvFile, opts Options) (int, error) {
	result, err := Apply(src, opts)
	if err != nil {
		return 0, fmt.Errorf("keeper: %w", err)
	}
	return len(result.Entries), nil
}
