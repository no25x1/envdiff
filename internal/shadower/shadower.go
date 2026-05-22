// Package shadower detects keys in a target env file that shadow (override)
// keys defined in a base env file, optionally filtering by prefix.
package shadower

import "github.com/your-org/envdiff/internal/parser"

// Shadow represents a key that exists in both base and target with differing values.
type Shadow struct {
	Key       string
	BaseValue string
	NextValue string
}

// Options controls shadower behaviour.
type Options struct {
	// Prefix restricts detection to keys with the given prefix. Empty means all keys.
	Prefix string
	// IgnoreEqual when true skips keys whose values are identical.
	IgnoreEqual bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{IgnoreEqual: true}
}

// Apply compares base and next env files and returns keys in next that shadow
// keys in base. Returns an error if either file is nil.
func Apply(base, next *parser.EnvFile, opts Options) ([]Shadow, error) {
	if base == nil || next == nil {
		return nil, fmt.Errorf("shadower: base and next files must not be nil")
	}

	baseMap := make(map[string]string, len(base.Entries))
	for _, e := range base.Entries {
		baseMap[e.Key] = e.Value
	}

	var shadows []Shadow
	for _, e := range next.Entries {
		if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
			continue
		}
		baseVal, exists := baseMap[e.Key]
		if !exists {
			continue
		}
		if opts.IgnoreEqual && baseVal == e.Value {
			continue
		}
		shadows = append(shadows, Shadow{
			Key:       e.Key,
			BaseValue: baseVal,
			NextValue: e.Value,
		})
	}
	return shadows, nil
}
