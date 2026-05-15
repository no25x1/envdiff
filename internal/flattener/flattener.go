// Package flattener merges multiple env files into a single flat key=value map,
// resolving conflicts according to a configurable priority strategy.
package flattener

import (
	"fmt"

	"github.com/user/envdiff/internal/parser"
)

// Strategy controls how key conflicts are resolved across files.
type Strategy string

const (
	// StrategyFirst keeps the value from the first file that defines the key.
	StrategyFirst Strategy = "first"
	// StrategyLast keeps the value from the last file that defines the key.
	StrategyLast Strategy = "last"
	// StrategyError returns an error when the same key appears in multiple files.
	StrategyError Strategy = "error"
)

// Options configures the flattening behaviour.
type Options struct {
	Strategy Strategy
}

// DefaultOptions returns sensible defaults (last-wins).
func DefaultOptions() Options {
	return Options{Strategy: StrategyLast}
}

// Result holds the flattened output together with metadata about conflicts.
type Result struct {
	File      *parser.EnvFile
	Conflicts []Conflict
}

// Conflict records a key that was defined in more than one source file.
type Conflict struct {
	Key    string
	Files  []string
	Chosen string
}

// Apply flattens the provided env files into a single EnvFile according to opts.
// Files are processed in the order they are supplied.
func Apply(files []*parser.EnvFile, opts Options) (*Result, error) {
	if len(files) == 0 {
		return &Result{File: &parser.EnvFile{}}, nil
	}

	if opts.Strategy == "" {
		opts = DefaultOptions()
	}

	seen := make(map[string]string)   // key -> origin filename
	index := make(map[string]int)     // key -> position in entries
	conflicts := []Conflict{}
	var entries []parser.Entry

	for _, f := range files {
		if f == nil {
			continue
		}
		for _, e := range f.Entries {
			origin, exists := seen[e.Key]
			if !exists {
				seen[e.Key] = f.Path
				index[e.Key] = len(entries)
				entries = append(entries, e)
				continue
			}
			// conflict
			switch opts.Strategy {
			case StrategyError:
				return nil, fmt.Errorf("flattener: key %q defined in both %q and %q", e.Key, origin, f.Path)
			case StrategyFirst:
				// keep existing — record conflict but do not update
				conflicts = append(conflicts, Conflict{Key: e.Key, Files: []string{origin, f.Path}, Chosen: entries[index[e.Key]].Value})
			default: // StrategyLast
				conflicts = append(conflicts, Conflict{Key: e.Key, Files: []string{origin, f.Path}, Chosen: e.Value})
				entries[index[e.Key]] = e
				seen[e.Key] = f.Path
			}
		}
	}

	return &Result{
		File:      &parser.EnvFile{Entries: entries},
		Conflicts: conflicts,
	}, nil
}
