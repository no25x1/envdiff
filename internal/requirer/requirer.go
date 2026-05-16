// Package requirer checks that all required keys are present and non-empty
// in an env file, returning a list of violations for any that are missing.
package requirer

import "github.com/user/envdiff/internal/parser"

// Violation describes a single missing or empty required key.
type Violation struct {
	Key    string
	Reason string
}

// Options controls which keys are required and how strictness is applied.
type Options struct {
	// Keys is the explicit list of keys that must be present.
	Keys []string
	// AllowEmpty permits keys to exist with an empty value.
	AllowEmpty bool
}

// Apply checks src against the required keys defined in opts.
// It returns a slice of Violation for every key that is absent or (when
// AllowEmpty is false) present but empty.
func Apply(src *parser.EnvFile, opts Options) ([]Violation, error) {
	if src == nil {
		return nil, nil
	}
	if len(opts.Keys) == 0 {
		return nil, nil
	}

	// Build a lookup map for O(1) access.
	index := make(map[string]string, len(src.Entries))
	for _, e := range src.Entries {
		index[e.Key] = e.Value
	}

	var violations []Violation
	for _, key := range opts.Keys {
		val, found := index[key]
		if !found {
			violations = append(violations, Violation{
				Key:    key,
				Reason: "key is missing",
			})
			continue
		}
		if !opts.AllowEmpty && val == "" {
			violations = append(violations, Violation{
				Key:    key,
				Reason: "key is present but empty",
			})
		}
	}
	return violations, nil
}
