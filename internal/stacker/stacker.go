// Package stacker layers multiple .env files, applying each one on top of
// the previous result so that later files override earlier ones — similar
// to Docker's layer model or docker-compose env_file stacking.
package stacker

import (
	"errors"
	"fmt"

	"envdiff/internal/parser"
)

// Options controls stacking behaviour.
type Options struct {
	// AllowEmpty permits source files that contain no entries.
	AllowEmpty bool
	// Prefix is prepended to every key in the final result.
	Prefix string
}

// Apply stacks files in order; later files override earlier ones.
// At least one file must be provided.
func Apply(files []*parser.EnvFile, opts Options) (*parser.EnvFile, error) {
	if len(files) == 0 {
		return nil, errors.New("stacker: at least one file is required")
	}

	index := map[string]int{} // key -> position in entries
	var entries []parser.Entry

	for fileIdx, f := range files {
		if f == nil {
			return nil, fmt.Errorf("stacker: file at index %d is nil", fileIdx)
		}
		if !opts.AllowEmpty && len(f.Entries) == 0 {
			return nil, fmt.Errorf("stacker: file at index %d contains no entries", fileIdx)
		}
		for _, e := range f.Entries {
			key := e.Key
			if opts.Prefix != "" {
				key = opts.Prefix + key
			}
			out := parser.Entry{Key: key, Value: e.Value, Comment: e.Comment}
			if pos, exists := index[key]; exists {
				entries[pos] = out
			} else {
				index[key] = len(entries)
				entries = append(entries, out)
			}
		}
	}

	return &parser.EnvFile{Entries: entries}, nil
}
