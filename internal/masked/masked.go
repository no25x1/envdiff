// Package masked provides utilities for comparing two env files while
// redacting sensitive values before any output is produced.
package masked

import (
	"fmt"
	"strings"

	"envdiff/internal/diff"
	"envdiff/internal/masker"
	"envdiff/internal/parser"
)

// Result holds a single compared entry with its (possibly masked) values.
type Result struct {
	Key      string
	Status   diff.Status
	BaseVal  string
	NextVal  string
	Masked   bool
}

// Options controls how the masked diff is produced.
type Options struct {
	// ExtraPatterns are additional key patterns to treat as sensitive.
	ExtraPatterns []string
}

// Compare diffs base against next and masks sensitive values in the results.
func Compare(base, next *parser.EnvFile, opts Options) ([]Result, error) {
	if base == nil || next == nil {
		return nil, fmt.Errorf("masked: base and next files must not be nil")
	}

	m := masker.New()
	if len(opts.ExtraPatterns) > 0 {
		m = masker.NewWithOptions(masker.Options{ExtraPatterns: opts.ExtraPatterns})
	}

	results := diff.Diff(base, next)
	out := make([]Result, 0, len(results))

	for _, r := range results {
		sensitive := m.IsSensitive(r.Key)
		bv := r.BaseVal
		nv := r.NextVal
		if sensitive {
			bv = m.MaskValue(bv)
			nv = m.MaskValue(nv)
		}
		out = append(out, Result{
			Key:     r.Key,
			Status:  r.Status,
			BaseVal: bv,
			NextVal: nv,
			Masked:  sensitive,
		})
	}
	return out, nil
}

// Format returns a human-readable summary line for a Result.
func Format(r Result) string {
	var sb strings.Builder
	switch r.Status {
	case diff.StatusAdded:
		fmt.Fprintf(&sb, "+ %s=%s", r.Key, r.NextVal)
	case diff.StatusRemoved:
		fmt.Fprintf(&sb, "- %s=%s", r.Key, r.BaseVal)
	case diff.StatusModified:
		fmt.Fprintf(&sb, "~ %s: %s -> %s", r.Key, r.BaseVal, r.NextVal)
	default:
		fmt.Fprintf(&sb, "  %s=%s", r.Key, r.BaseVal)
	}
	if r.Masked {
		sb.WriteString(" [masked]")
	}
	return sb.String()
}
