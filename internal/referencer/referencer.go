// Package referencer finds all keys in an env file that are referenced
// by other keys via variable interpolation (e.g. ${VAR} or $VAR).
package referencer

import (
	"fmt"
	"regexp"
	"sort"

	"envdiff/internal/parser"
)

var refPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

// Result holds the reference graph for a single env file.
type Result struct {
	// Referrers maps a key to the list of keys whose values reference it.
	Referrers map[string][]string
	// Undefined contains keys that are referenced but not defined in the file.
	Undefined []string
	// Unused contains keys that are defined but never referenced by any other key.
	Unused []string
}

// Analyse inspects f and returns a Result describing cross-key references.
func Analyse(f *parser.EnvFile) (*Result, error) {
	if f == nil {
		return nil, fmt.Errorf("referencer: nil file")
	}

	defined := make(map[string]struct{}, len(f.Entries))
	for _, e := range f.Entries {
		defined[e.Key] = struct{}{}
	}

	referrers := make(map[string][]string)
	referenced := make(map[string]struct{})

	for _, e := range f.Entries {
		matches := refPattern.FindAllStringSubmatch(e.Value, -1)
		for _, m := range matches {
			name := m[1]
			if name == "" {
				name = m[2]
			}
			if name == e.Key {
				continue // skip self-references
			}
			referrers[name] = append(referrers[name], e.Key)
			referenced[name] = struct{}{}
		}
	}

	var undefined []string
	for name := range referenced {
		if _, ok := defined[name]; !ok {
			undefined = append(undefined, name)
		}
	}

	var unused []string
	for _, e := range f.Entries {
		if _, ok := referenced[e.Key]; !ok {
			unused = append(unused, e.Key)
		}
	}

	sort.Strings(undefined)
	sort.Strings(unused)

	return &Result{
		Referrers: referrers,
		Undefined: undefined,
		Unused:    unused,
	}, nil
}
