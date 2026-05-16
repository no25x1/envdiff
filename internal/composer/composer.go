// Package composer merges multiple .env files into a single output,
// applying optional prefix namespacing per source file.
package composer

import (
	"fmt"
	"strings"

	"envdiff/internal/parser"
)

// Source describes one input file and an optional namespace prefix.
type Source struct {
	Path      string
	Namespace string // if non-empty, keys are prefixed with NAMESPACE_
}

// Options controls composition behaviour.
type Options struct {
	Sources         []Source
	OverwriteOnConflict bool // last-writer wins when true; first-writer wins when false
}

// Apply composes all source files into a single EnvFile.
// Order in Sources determines precedence when OverwriteOnConflict is false.
func Apply(opts Options) (*parser.EnvFile, error) {
	if len(opts.Sources) == 0 {
		return nil, fmt.Errorf("composer: no source files provided")
	}

	seen := make(map[string]bool)
	var entries []parser.Entry

	for _, src := range opts.Sources {
		f, err := parser.Parse(src.Path)
		if err != nil {
			return nil, fmt.Errorf("composer: parsing %s: %w", src.Path, err)
		}

		for _, e := range f.Entries {
			key := namespacedKey(src.Namespace, e.Key)
			e.Key = key

			if seen[key] {
				if !opts.OverwriteOnConflict {
					continue
				}
				// replace existing entry
				for i, existing := range entries {
					if existing.Key == key {
						entries[i] = e
						break
					}
				}
				continue
			}

			seen[key] = true
			entries = append(entries, e)
		}
	}

	return &parser.EnvFile{Entries: entries}, nil
}

func namespacedKey(ns, key string) string {
	if ns == "" {
		return key
	}
	return strings.ToUpper(ns) + "_" + key
}
