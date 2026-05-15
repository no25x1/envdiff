// Package renamer provides utilities for bulk-renaming keys in an env file.
package renamer

import (
	"fmt"
	"strings"

	"github.com/your-org/envdiff/internal/parser"
)

// Options controls how keys are renamed.
type Options struct {
	// Mappings is an explicit old->new key map.
	Mappings map[string]string
	// AddPrefix prepends a string to every key.
	AddPrefix string
	// StripPrefix removes a leading prefix from every key.
	StripPrefix string
	// ToUpper converts all keys to upper-case after other transforms.
	ToUpper bool
}

// Result holds the renamed file and a log of every change made.
type Result struct {
	File    *parser.EnvFile
	Changes []Change
}

// Change records a single key rename.
type Change struct {
	OldKey string
	NewKey string
}

// Apply renames keys in f according to opts and returns a new EnvFile plus a
// list of changes. The original file is never mutated.
func Apply(f *parser.EnvFile, opts Options) (*Result, error) {
	out := &parser.EnvFile{}
	var changes []Change
	seen := make(map[string]bool)

	for _, e := range f.Entries {
		newKey, err := rename(e.Key, opts)
		if err != nil {
			return nil, err
		}
		if seen[newKey] {
			return nil, fmt.Errorf("renamer: rename would produce duplicate key %q", newKey)
		}
		seen[newKey] = true
		if newKey != e.Key {
			changes = append(changes, Change{OldKey: e.Key, NewKey: newKey})
		}
		out.Entries = append(out.Entries, parser.Entry{
			Key:     newKey,
			Value:   e.Value,
			Comment: e.Comment,
		})
	}
	return &Result{File: out, Changes: changes}, nil
}

func rename(key string, opts Options) (string, error) {
	// Explicit mapping takes priority.
	if opts.Mappings != nil {
		if v, ok := opts.Mappings[key]; ok {
			return v, nil
		}
	}
	result := key
	if opts.StripPrefix != "" {
		result = strings.TrimPrefix(result, opts.StripPrefix)
	}
	if opts.AddPrefix != "" {
		result = opts.AddPrefix + result
	}
	if opts.ToUpper {
		result = strings.ToUpper(result)
	}
	return result, nil
}
