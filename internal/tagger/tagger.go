// Package tagger provides functionality for adding, removing, and querying
// inline tags on .env entries using comment annotations (e.g. # @tag:production).
package tagger

import (
	"fmt"
	"strings"

	"envdiff/internal/parser"
)

// Options controls tagging behaviour.
type Options struct {
	// Add appends these tags to matching keys.
	Add []string
	// Remove strips these tags from matching keys.
	Remove []string
	// Keys restricts operations to the given key names. Empty means all keys.
	Keys []string
	// FilterTag returns only entries that carry this tag when non-empty.
	FilterTag string
}

// Apply adds or removes tags on entries in src and returns a new EnvFile.
// Tags are stored in the entry comment as "# @tag:<name>".
func Apply(src *parser.EnvFile, opts Options) (*parser.EnvFile, error) {
	if src == nil {
		return &parser.EnvFile{}, nil
	}

	result := &parser.EnvFile{}
	for _, e := range src.Entries {
		entry := e
		if len(opts.Keys) == 0 || containsKey(opts.Keys, entry.Key) {
			for _, t := range opts.Add {
				entry.Comment = addTag(entry.Comment, t)
			}
			for _, t := range opts.Remove {
				entry.Comment = removeTag(entry.Comment, t)
			}
		}
		if opts.FilterTag != "" && !hasTag(entry.Comment, opts.FilterTag) {
			continue
		}
		result.Entries = append(result.Entries, entry)
	}
	return result, nil
}

// Tags returns the list of tags present in the given comment string.
func Tags(comment string) []string {
	var tags []string
	for _, part := range strings.Fields(comment) {
		part = strings.TrimPrefix(part, "#")
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "@tag:") {
			tags = append(tags, strings.TrimPrefix(part, "@tag:"))
		}
	}
	return tags
}

func addTag(comment, tag string) string {
	marker := fmt.Sprintf("@tag:%s", tag)
	if strings.Contains(comment, marker) {
		return comment
	}
	if comment == "" {
		return "# " + marker
	}
	return strings.TrimRight(comment, " ") + " " + marker
}

func removeTag(comment, tag string) string {
	marker := fmt.Sprintf("@tag:%s", tag)
	parts := strings.Fields(comment)
	var kept []string
	for _, p := range parts {
		if p != marker {
			kept = append(kept, p)
		}
	}
	result := strings.Join(kept, " ")
	if result == "#" {
		return ""
	}
	return result
}

func hasTag(comment, tag string) bool {
	return strings.Contains(comment, fmt.Sprintf("@tag:%s", tag))
}

func containsKey(keys []string, key string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}
