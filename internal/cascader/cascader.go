// Package cascader merges a sequence of env files in cascade order,
// where later files override earlier ones only for keys that are present.
package cascader

import (
	"errors"
	"fmt"

	"envdiff/internal/parser"
)

// Options controls cascade behaviour.
type Options struct {
	// AllowEmpty permits source files that contain no entries.
	AllowEmpty bool
	// StopOnMissing returns an error when a key in the base file is absent
	// from every subsequent layer.
	StopOnMissing bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		AllowEmpty:    true,
		StopOnMissing: false,
	}
}

// Apply cascades layers on top of base. The first file is the lowest-priority
// base; each subsequent file overrides matching keys from the layer below.
func Apply(files []*parser.EnvFile, opts Options) (*parser.EnvFile, error) {
	if len(files) == 0 {
		return nil, errors.New("cascader: at least one file is required")
	}
	for i, f := range files {
		if f == nil {
			return nil, fmt.Errorf("cascader: file at index %d is nil", i)
		}
		if !opts.AllowEmpty && len(f.Entries) == 0 {
			return nil, fmt.Errorf("cascader: file at index %d is empty", i)
		}
	}

	// Build index from base.
	result := &parser.EnvFile{}
	order := make([]string, 0, len(files[0].Entries))
	index := make(map[string]parser.Entry)

	for _, e := range files[0].Entries {
		if _, seen := index[e.Key]; !seen {
			order = append(order, e.Key)
		}
		index[e.Key] = e
	}

	// Apply each subsequent layer.
	for _, layer := range files[1:] {
		for _, e := range layer.Entries {
			if _, exists := index[e.Key]; exists {
				index[e.Key] = e
			}
			// Keys not in base are intentionally ignored.
		}
	}

	if opts.StopOnMissing {
		for _, key := range order {
			found := false
			for _, layer := range files[1:] {
				for _, e := range layer.Entries {
					if e.Key == key {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found && len(files) > 1 {
				return nil, fmt.Errorf("cascader: key %q missing from all overlay layers", key)
			}
		}
	}

	for _, key := range order {
		result.Entries = append(result.Entries, index[key])
	}
	return result, nil
}
