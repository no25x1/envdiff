// Package counter provides utilities for counting and summarising
// key/value statistics across an EnvFile.
package counter

import (
	"sort"

	"github.com/envdiff/envdiff/internal/parser"
)

// Result holds aggregated counts for an EnvFile.
type Result struct {
	Total      int
	Empty      int
	NonEmpty   int
	Unique     int
	Duplicated int
	// TopPrefixes lists the top prefix groups by count (prefix -> count).
	TopPrefixes map[string]int
}

// Apply counts statistics for the given EnvFile.
// It returns an error only when src is nil.
func Apply(src *parser.EnvFile) (*Result, error) {
	if src == nil {
		return nil, fmt.Errorf("counter: nil file")
	}

	valueSeen := make(map[string]int)
	prefixCount := make(map[string]int)

	r := &Result{}

	for _, e := range src.Entries {
		r.Total++
		if e.Value == "" {
			r.Empty++
		} else {
			r.NonEmpty++
		}
		valueSeen[e.Value]++

		if p := prefix(e.Key); p != "" {
			prefixCount[p]++
		}
	}

	for _, cnt := range valueSeen {
		if cnt == 1 {
			r.Unique++
		} else {
			r.Duplicated++
		}
	}

	r.TopPrefixes = topN(prefixCount, 5)
	return r, nil
}

// prefix returns the portion of key before the first '_', or "" if none.
func prefix(key string) string {
	for i, c := range key {
		if c == '_' && i > 0 {
			return key[:i]
		}
	}
	return ""
}

// topN returns up to n entries with the highest counts.
func topN(m map[string]int, n int) map[string]int {
	type kv struct {
		k string
		v int
	}
	var pairs []kv
	for k, v := range m {
		pairs = append(pairs, kv{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].v != pairs[j].v {
			return pairs[i].v > pairs[j].v
		}
		return pairs[i].k < pairs[j].k
	})
	out := make(map[string]int)
	for i, p := range pairs {
		if i >= n {
			break
		}
		out[p.k] = p.v
	}
	return out
}
