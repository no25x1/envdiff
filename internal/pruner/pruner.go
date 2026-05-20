// Package pruner removes keys from an EnvFile that have no effect —
// specifically keys whose values are empty and that are not explicitly
// required by the caller.
package pruner

import (
	"errors"

	"github.com/yourusername/envdiff/internal/parser"
)

// Options controls which keys are pruned.
type Options struct {
	// RemoveEmpty removes keys whose value is the empty string.
	RemoveEmpty bool

	// RemoveCommented removes entries that carry only a comment and no key.
	RemoveCommented bool

	// Allowlist is an optional set of keys that must never be removed,
	// even if they would otherwise qualify for pruning.
	Allowlist []string
}

// DefaultOptions returns a safe default: only empty-value keys are pruned.
func DefaultOptions() Options {
	return Options{RemoveEmpty: true}
}

// Apply returns a new EnvFile with qualifying entries removed.
// The original file is never modified.
func Apply(f *parser.EnvFile, opts Options) (*parser.EnvFile, error) {
	if f == nil {
		return &parser.EnvFile{}, nil
	}

	if !opts.RemoveEmpty && !opts.RemoveCommented {
		return nil, errors.New("pruner: no pruning options enabled")
	}

	allow := toSet(opts.Allowlist)

	out := &parser.EnvFile{}
	for _, e := range f.Entries {
		if allow[e.Key] {
			out.Entries = append(out.Entries, e)
			continue
		}

		if opts.RemoveCommented && e.Key == "" {
			// entry is a blank line or pure comment — skip it
			continue
		}

		if opts.RemoveEmpty && e.Key != "" && e.Value == "" {
			continue
		}

		out.Entries = append(out.Entries, e)
	}

	return out, nil
}

// Count returns the number of entries that would be removed by Apply without
// actually modifying anything. It returns an error under the same conditions
// as Apply.
func Count(f *parser.EnvFile, opts Options) (int, error) {
	if f == nil {
		return 0, nil
	}

	if !opts.RemoveEmpty && !opts.RemoveCommented {
		return 0, errors.New("pruner: no pruning options enabled")
	}

	allow := toSet(opts.Allowlist)
	removed := 0

	for _, e := range f.Entries {
		if allow[e.Key] {
			continue
		}
		if opts.RemoveCommented && e.Key == "" {
			removed++
			continue
		}
		if opts.RemoveEmpty && e.Key != "" && e.Value == "" {
			removed++
		}
	}

	return removed, nil
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
