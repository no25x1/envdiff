// Package transformer provides key/value transformation utilities for env files.
package transformer

import (
	"strings"

	"github.com/envdiff/envdiff/internal/parser"
)

// TransformFunc is a function that transforms a key-value pair.
// Returning ok=false drops the entry from the output.
type TransformFunc func(key, value string) (newKey, newValue string, ok bool)

// Options controls which transformations are applied.
type Options struct {
	// PrefixAdd prepends a prefix to every key.
	PrefixAdd string
	// PrefixStrip removes a prefix from every key (entries without the prefix are dropped).
	PrefixStrip string
	// KeyToUpper converts all keys to upper-case.
	KeyToUpper bool
	// KeyToLower converts all keys to lower-case.
	KeyToLower bool
	// Custom holds user-supplied transform functions applied in order.
	Custom []TransformFunc
}

// Apply runs all configured transformations on f and returns a new EnvFile.
func Apply(f parser.EnvFile, opts Options) parser.EnvFile {
	out := parser.EnvFile{
		Path:    f.Path,
		Entries: make([]parser.Entry, 0, len(f.Entries)),
	}

	for _, e := range f.Entries {
		key, value, ok := transform(e.Key, e.Value, opts)
		if !ok {
			continue
		}
		out.Entries = append(out.Entries, parser.Entry{
			Key:   key,
			Value: value,
		})
	}
	return out
}

func transform(key, value string, opts Options) (string, string, bool) {
	if opts.PrefixStrip != "" {
		if !strings.HasPrefix(key, opts.PrefixStrip) {
			return "", "", false
		}
		key = strings.TrimPrefix(key, opts.PrefixStrip)
	}

	if opts.PrefixAdd != "" {
		key = opts.PrefixAdd + key
	}

	if opts.KeyToUpper {
		key = strings.ToUpper(key)
	} else if opts.KeyToLower {
		key = strings.ToLower(key)
	}

	for _, fn := range opts.Custom {
		var ok bool
		key, value, ok = fn(key, value)
		if !ok {
			return "", "", false
		}
	}

	return key, value, true
}
