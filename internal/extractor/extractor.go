// Package extractor extracts a subset of keys from an EnvFile into a new EnvFile.
package extractor

import (
	"fmt"
	"strings"

	"github.com/your-org/envdiff/internal/parser"
)

// Options controls which keys are extracted.
type Options struct {
	// Keys is an explicit list of keys to extract.
	Keys []string
	// Prefix retains only keys that start with the given prefix.
	Prefix string
	// StripPrefix removes the prefix from the key name in the output.
	StripPrefix bool
}

// Apply extracts matching entries from src and returns a new EnvFile.
// Returns an error if src is nil or no keys are specified.
func Apply(src *parser.EnvFile, opts Options) (*parser.EnvFile, error) {
	if src == nil {
		return nil, fmt.Errorf("extractor: source file is nil")
	}
	if len(opts.Keys) == 0 && opts.Prefix == "" {
		return nil, fmt.Errorf("extractor: at least one of Keys or Prefix must be set")
	}

	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = true
	}

	out := &parser.EnvFile{}
	for _, e := range src.Entries {
		if !matches(e.Key, keySet, opts.Prefix) {
			continue
		}
		key := e.Key
		if opts.StripPrefix && opts.Prefix != "" {
			key = strings.TrimPrefix(key, opts.Prefix)
		}
		out.Entries = append(out.Entries, parser.Entry{
			Key:     key,
			Value:   e.Value,
			Comment: e.Comment,
		})
	}
	return out, nil
}

func matches(key string, keySet map[string]bool, prefix string) bool {
	if keySet[key] {
		return true
	}
	if prefix != "" && strings.HasPrefix(key, prefix) {
		return true
	}
	return false
}
