// Package annotator provides functionality to attach inline comments
// (annotations) to .env file entries based on configurable rules.
package annotator

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Rule defines a condition and the comment to attach when it matches.
type Rule struct {
	// KeyPrefix, if non-empty, matches keys that start with this prefix.
	KeyPrefix string
	// KeySuffix, if non-empty, matches keys that end with this suffix.
	KeySuffix string
	// Comment is the annotation text to attach (without leading "# ").
	Comment string
}

// Options controls annotation behaviour.
type Options struct {
	Rules []Rule
	// OverwriteExisting replaces existing inline comments when true.
	OverwriteExisting bool
}

// DefaultOptions returns a sensible set of built-in annotation rules.
func DefaultOptions() Options {
	return Options{
		Rules: []Rule{
			{KeySuffix: "_SECRET", Comment: "sensitive – do not commit"},
			{KeySuffix: "_PASSWORD", Comment: "sensitive – do not commit"},
			{KeySuffix: "_TOKEN", Comment: "sensitive – do not commit"},
			{KeySuffix: "_KEY", Comment: "sensitive – do not commit"},
			{KeyPrefix: "DEBUG_", Comment: "debug flag – disable in production"},
		},
		OverwriteExisting: false,
	}
}

// Apply returns a new EnvFile with annotations added to matching entries.
func Apply(src *parser.EnvFile, opts Options) *parser.EnvFile {
	if src == nil {
		return &parser.EnvFile{}
	}

	out := &parser.EnvFile{}
	for _, e := range src.Entries {
		annotated := e
		comment := matchComment(e.Key, opts.Rules)
		if comment != "" {
			if annotated.Comment == "" || opts.OverwriteExisting {
				annotated.Comment = fmt.Sprintf("# %s", comment)
			}
		}
		out.Entries = append(out.Entries, annotated)
	}
	return out
}

// matchComment returns the first matching rule comment for the given key,
// or an empty string if no rule matches.
func matchComment(key string, rules []Rule) string {
	upper := strings.ToUpper(key)
	for _, r := range rules {
		prefixMatch := r.KeyPrefix == "" || strings.HasPrefix(upper, strings.ToUpper(r.KeyPrefix))
		suffixMatch := r.KeySuffix == "" || strings.HasSuffix(upper, strings.ToUpper(r.KeySuffix))
		if prefixMatch && suffixMatch {
			return r.Comment
		}
	}
	return ""
}
