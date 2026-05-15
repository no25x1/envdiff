// Package differ provides line-level diffing of .env file values,
// showing before/after context for changed entries.
package differ

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// ChangeType indicates the kind of change for an entry.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Line represents a single diff line with context.
type Line struct {
	Key    string
	Before string
	After  string
	Change ChangeType
}

// Result holds the full line-level diff between two env files.
type Result struct {
	Lines []Line
}

// HasChanges returns true if any line is not Unchanged.
func (r *Result) HasChanges() bool {
	for _, l := range r.Lines {
		if l.Change != Unchanged {
			return true
		}
	}
	return false
}

// Summary returns a brief count string.
func (r *Result) Summary() string {
	var added, removed, modified int
	for _, l := range r.Lines {
		switch l.Change {
		case Added:
			added++
		case Removed:
			removed++
		case Modified:
			modified++
		}
	}
	return fmt.Sprintf("+%d -%d ~%d", added, removed, modified)
}

// Compare produces a line-level diff between base and target env files.
func Compare(base, target *parser.EnvFile) *Result {
	baseMap := make(map[string]string, len(base.Entries))
	for _, e := range base.Entries {
		baseMap[e.Key] = e.Value
	}
	targetMap := make(map[string]string, len(target.Entries))
	for _, e := range target.Entries {
		targetMap[e.Key] = e.Value
	}

	seen := make(map[string]bool)
	var lines []Line

	for _, e := range base.Entries {
		seen[e.Key] = true
		if tv, ok := targetMap[e.Key]; ok {
			if e.Value == tv {
				lines = append(lines, Line{Key: e.Key, Before: e.Value, After: tv, Change: Unchanged})
			} else {
				lines = append(lines, Line{Key: e.Key, Before: e.Value, After: tv, Change: Modified})
			}
		} else {
			lines = append(lines, Line{Key: e.Key, Before: e.Value, After: "", Change: Removed})
		}
	}

	for _, e := range target.Entries {
		if !seen[e.Key] {
			lines = append(lines, Line{Key: e.Key, Before: "", After: e.Value, Change: Added})
		}
	}

	return &Result{Lines: lines}
}

// Format renders a unified-style diff string.
func Format(r *Result) string {
	var sb strings.Builder
	for _, l := range r.Lines {
		switch l.Change {
		case Added:
			fmt.Fprintf(&sb, "+ %s=%s\n", l.Key, l.After)
		case Removed:
			fmt.Fprintf(&sb, "- %s=%s\n", l.Key, l.Before)
		case Modified:
			fmt.Fprintf(&sb, "~ %s: %s -> %s\n", l.Key, l.Before, l.After)
		case Unchanged:
			fmt.Fprintf(&sb, "  %s=%s\n", l.Key, l.Before)
		}
	}
	return sb.String()
}
