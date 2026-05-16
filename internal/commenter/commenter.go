// Package commenter provides utilities for adding, removing, and updating
// inline comments on .env file entries.
package commenter

import (
	"strings"

	"github.com/envdiff/envdiff/internal/parser"
)

// Options controls how comments are applied.
type Options struct {
	// Set maps key names to comment strings to add or replace.
	Set map[string]string
	// Remove lists key names whose comments should be stripped.
	Remove []string
	// Overwrite replaces existing comments when true; otherwise only adds to
	// entries that have no comment.
	Overwrite bool
}

// Apply returns a new EnvFile with comments modified according to opts.
// Entries not referenced by opts are returned unchanged.
func Apply(src *parser.EnvFile, opts Options) (*parser.EnvFile, error) {
	if src == nil {
		return &parser.EnvFile{}, nil
	}

	removeSet := make(map[string]bool, len(opts.Remove))
	for _, k := range opts.Remove {
		removeSet[k] = true
	}

	out := &parser.EnvFile{}
	for _, e := range src.Entries {
		entry := e // copy

		if removeSet[entry.Key] {
			entry.Comment = ""
			out.Entries = append(out.Entries, entry)
			continue
		}

		if comment, ok := opts.Set[entry.Key]; ok {
			if entry.Comment == "" || opts.Overwrite {
				entry.Comment = normalizeComment(comment)
			}
		}

		out.Entries = append(out.Entries, entry)
	}
	return out, nil
}

// normalizeComment ensures a comment string starts with "# ".
func normalizeComment(c string) string {
	c = strings.TrimSpace(c)
	if c == "" {
		return ""
	}
	if !strings.HasPrefix(c, "#") {
		return "# " + c
	}
	// Ensure space after #
	if len(c) > 1 && c[1] != ' ' {
		return "# " + strings.TrimPrefix(c, "#")
	}
	return c
}
