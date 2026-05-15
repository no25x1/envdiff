// Package comparator provides multi-file .env comparison across environments.
package comparator

import (
	"sort"

	"github.com/envdiff/internal/parser"
)

// Result holds the comparison matrix for a set of env files.
type Result struct {
	// Keys is the sorted union of all keys across all files.
	Keys []string
	// Files is the ordered list of file labels.
	Files []string
	// Matrix maps key -> file label -> value (empty string means absent).
	Matrix map[string]map[string]string
	// Absent tracks which keys are missing in which files.
	Absent map[string][]string
}

// Compare accepts a map of label -> EnvFile and returns a Result
// describing how values align (or diverge) across all files.
func Compare(files map[string]*parser.EnvFile) *Result {
	keySet := map[string]struct{}{}
	for _, ef := range files {
		for _, e := range ef.Entries {
			keySet[e.Key] = struct{}{}
		}
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	labels := make([]string, 0, len(files))
	for l := range files {
		labels = append(labels, l)
	}
	sort.Strings(labels)

	matrix := make(map[string]map[string]string, len(keys))
	absent := make(map[string][]string)

	for _, key := range keys {
		matrix[key] = make(map[string]string, len(labels))
		for _, label := range labels {
			ef := files[label]
			val, ok := lookup(ef, key)
			if ok {
				matrix[key][label] = val
			} else {
				matrix[key][label] = ""
				absent[key] = append(absent[key], label)
			}
		}
	}

	return &Result{
		Keys:   keys,
		Files:  labels,
		Matrix: matrix,
		Absent: absent,
	}
}

// Conflicts returns keys whose non-empty values differ across files.
func (r *Result) Conflicts() []string {
	var out []string
	for _, key := range r.Keys {
		vals := map[string]struct{}{}
		for _, v := range r.Matrix[key] {
			if v != "" {
				vals[v] = struct{}{}
			}
		}
		if len(vals) > 1 {
			out = append(out, key)
		}
	}
	return out
}

func lookup(ef *parser.EnvFile, key string) (string, bool) {
	for _, e := range ef.Entries {
		if e.Key == key {
			return e.Value, true
		}
	}
	return "", false
}
