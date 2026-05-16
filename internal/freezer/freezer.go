// Package freezer provides functionality to freeze an env file by locking
// key-value pairs so that future diffs can detect unexpected mutations.
package freezer

import (
	"crypto/sha256"
	"fmt"
	"sort"

	"envdiff/internal/parser"
)

// FrozenEntry holds a key and the hash of its value at freeze time.
type FrozenEntry struct {
	Key  string `json:"key"`
	Hash string `json:"hash"`
}

// Violation describes a key whose value changed since it was frozen.
type Violation struct {
	Key      string
	Expected string // hash at freeze time
	Actual   string // hash of current value
}

// Freeze computes a deterministic SHA-256 hash for every entry in f and
// returns the slice sorted by key.
func Freeze(f *parser.EnvFile) []FrozenEntry {
	if f == nil {
		return nil
	}
	entries := make([]FrozenEntry, 0, len(f.Entries))
	for _, e := range f.Entries {
		entries = append(entries, FrozenEntry{
			Key:  e.Key,
			Hash: hashValue(e.Value),
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	return entries
}

// Check compares the current env file against a frozen snapshot and returns
// any violations where values have changed.
func Check(f *parser.EnvFile, frozen []FrozenEntry) []Violation {
	if f == nil || len(frozen) == 0 {
		return nil
	}
	index := make(map[string]string, len(frozen))
	for _, fe := range frozen {
		index[fe.Key] = fe.Hash
	}
	var violations []Violation
	for _, e := range f.Entries {
		expected, ok := index[e.Key]
		if !ok {
			continue
		}
		actual := hashValue(e.Value)
		if actual != expected {
			violations = append(violations, Violation{
				Key:      e.Key,
				Expected: expected,
				Actual:   actual,
			})
		}
	}
	return violations
}

func hashValue(v string) string {
	sum := sha256.Sum256([]byte(v))
	return fmt.Sprintf("%x", sum)
}
