// Package aliaser provides functionality to define and resolve key aliases
// within an env file, allowing one key to mirror the value of another.
package aliaser

import (
	"errors"
	"fmt"

	"envdiff/internal/parser"
)

// Mapping represents a single alias definition: Alias is the new key name,
// Source is the existing key whose value should be copied.
type Mapping struct {
	Alias  string
	Source string
}

// Options controls the behaviour of Apply.
type Options struct {
	// Mappings is the list of alias definitions to apply.
	Mappings []Mapping
	// Overwrite controls whether an existing key with the alias name is replaced.
	Overwrite bool
}

// Apply returns a new EnvFile that contains all original entries plus any
// alias entries resolved from opts.Mappings. Keys that cannot be resolved
// (source not found) are reported as errors but do not abort processing;
// all resolvable aliases are still applied.
func Apply(file *parser.EnvFile, opts Options) (*parser.EnvFile, []error) {
	if file == nil {
		return &parser.EnvFile{}, nil
	}
	if len(opts.Mappings) == 0 {
		return file.Clone(), nil
	}

	// Build a lookup map of existing entries.
	lookup := make(map[string]string, len(file.Entries))
	for _, e := range file.Entries {
		lookup[e.Key] = e.Value
	}

	// Clone original entries.
	out := file.Clone()

	// Track alias keys already present in output to honour Overwrite.
	existing := make(map[string]int, len(out.Entries))
	for i, e := range out.Entries {
		existing[e.Key] = i
	}

	var errs []error

	for _, m := range opts.Mappings {
		if m.Alias == "" || m.Source == "" {
			errs = append(errs, errors.New("aliaser: alias and source must not be empty"))
			continue
		}

		srcVal, found := lookup[m.Source]
		if !found {
			errs = append(errs, fmt.Errorf("aliaser: source key %q not found", m.Source))
			continue
		}

		newEntry := parser.Entry{
			Key:     m.Alias,
			Value:   srcVal,
			Comment: fmt.Sprintf("alias of %s", m.Source),
		}

		if idx, exists := existing[m.Alias]; exists {
			if opts.Overwrite {
				out.Entries[idx] = newEntry
			}
			continue
		}

		out.Entries = append(out.Entries, newEntry)
		existing[m.Alias] = len(out.Entries) - 1
	}

	return out, errs
}
