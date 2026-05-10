package diff

import (
	"sort"

	"github.com/user/envdiff/internal/parser"
)

// Status represents the kind of change for a key.
type Status string

const (
	Added     Status = "added"
	Removed   Status = "removed"
	Modified  Status = "modified"
	Unchanged Status = "unchanged"
)

// Result holds the diff outcome for a single key.
type Result struct {
	Key      string
	Status   Status
	OldValue string
	NewValue string
}

// Diff compares two EnvFiles and returns an ordered slice of Results.
// base is the reference environment; target is the environment being compared.
func Diff(base, target *parser.EnvFile) []Result {
	results := make([]Result, 0)
	seen := make(map[string]bool)

	// Collect all keys from both files in sorted order.
	keySet := make(map[string]struct{})
	for _, e := range base.Entries {
		keySet[e.Key] = struct{}{}
	}
	for _, e := range target.Entries {
		keySet[e.Key] = struct{}{}
	}
	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	baseMap := base.Map()
	targetMap := target.Map()

	for _, key := range keys {
		if seen[key] {
			continue
		}
		seen[key] = true

		oldVal, inBase := baseMap[key]
		newVal, inTarget := targetMap[key]

		switch {
		case inBase && inTarget && oldVal == newVal:
			results = append(results, Result{Key: key, Status: Unchanged, OldValue: oldVal, NewValue: newVal})
		case inBase && inTarget && oldVal != newVal:
			results = append(results, Result{Key: key, Status: Modified, OldValue: oldVal, NewValue: newVal})
		case !inBase && inTarget:
			results = append(results, Result{Key: key, Status: Added, OldValue: "", NewValue: newVal})
		case inBase && !inTarget:
			results = append(results, Result{Key: key, Status: Removed, OldValue: oldVal, NewValue: ""})
		}
	}
	return results
}
