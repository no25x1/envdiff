// Package truncator shortens env variable values to a maximum length,
// optionally appending a suffix such as "..." to indicate truncation.
package truncator

import "github.com/your-org/envdiff/internal/parser"

// DefaultOptions returns a safe default configuration.
func DefaultOptions() Options {
	return Options{
		MaxLen:  64,
		Suffix:  "...",
		KeysOnly: nil,
	}
}

// Options controls truncation behaviour.
type Options struct {
	// MaxLen is the maximum rune length of a value after truncation.
	MaxLen int
	// Suffix is appended when a value is shortened. Counts toward MaxLen.
	Suffix string
	// KeysOnly restricts truncation to the listed keys. nil means all keys.
	KeysOnly []string
}

// Apply returns a copy of src with values truncated according to opts.
// If src is nil an empty EnvFile is returned.
func Apply(src *parser.EnvFile, opts Options) *parser.EnvFile {
	out := &parser.EnvFile{}
	if src == nil {
		return out
	}

	allowSet := toSet(opts.KeysOnly)

	for _, e := range src.Entries {
		copy := e
		if len(allowSet) == 0 || allowSet[e.Key] {
			copy.Value = truncate(e.Value, opts.MaxLen, opts.Suffix)
		}
		out.Entries = append(out.Entries, copy)
	}
	return out
}

func truncate(val string, maxLen int, suffix string) string {
	runes := []rune(val)
	if len(runes) <= maxLen {
		return val
	}
	if len([]rune(suffix)) >= maxLen {
		return string([]rune(suffix)[:maxLen])
	}
	cutAt := maxLen - len([]rune(suffix))
	return string(runes[:cutAt]) + suffix
}

func toSet(keys []string) map[string]bool {
	if len(keys) == 0 {
		return nil
	}
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
