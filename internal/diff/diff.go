package diff

import (
	"github.com/user/envdiff/internal/parser"
)

// ChangeType describes the kind of difference between two env files.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Change represents a single diffed entry between two env files.
type Change struct {
	Key        string
	Type       ChangeType
	OldValue   string
	NewValue   string
}

// Result holds the complete diff between two env files.
type Result struct {
	BaseFile   string
	TargetFile string
	Changes    []Change
}

// Diff computes the difference between base and target env files.
func Diff(base, target *parser.EnvFile) *Result {
	result := &Result{
		BaseFile:   base.Path,
		TargetFile: target.Path,
	}

	// Check for removed or modified keys.
	for _, entry := range base.Entries {
		if t, ok := target.Index[entry.Key]; ok {
			if t.Value != entry.Value {
				result.Changes = append(result.Changes, Change{
					Key:      entry.Key,
					Type:     Modified,
					OldValue: entry.Value,
					NewValue: t.Value,
				})
			} else {
				result.Changes = append(result.Changes, Change{
					Key:  entry.Key,
					Type: Unchanged,
				})
			}
		} else {
			result.Changes = append(result.Changes, Change{
				Key:      entry.Key,
				Type:     Removed,
				OldValue: entry.Value,
			})
		}
	}

	// Check for added keys.
	for _, entry := range target.Entries {
		if _, ok := base.Index[entry.Key]; !ok {
			result.Changes = append(result.Changes, Change{
				Key:      entry.Key,
				Type:     Added,
				NewValue: entry.Value,
			})
		}
	}

	return result
}

// Summary returns counts of each change type.
func (r *Result) Summary() (added, removed, modified, unchanged int) {
	for _, c := range r.Changes {
		switch c.Type {
		case Added:
			added++
		case Removed:
			removed++
		case Modified:
			modified++
		case Unchanged:
			unchanged++
		}
	}
	return
}
