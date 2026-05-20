// Package summarizer provides a high-level statistical summary of an env file.
package summarizer

import (
	"sort"
	"strings"

	"envdiff/internal/parser"
)

// Summary holds aggregate statistics about an env file.
type Summary struct {
	TotalKeys      int            `json:"total_keys"`
	EmptyValues    int            `json:"empty_values"`
	SensitiveKeys  int            `json:"sensitive_keys"`
	Uniqueпрефixes int            `json:"unique_prefixes"`
	TopPrefixes    []PrefixCount  `json:"top_prefixes"`
	LongestKey     string         `json:"longest_key"`
	ShortestKey    string         `json:"shortest_key"`
}

// PrefixCount pairs a prefix with its occurrence count.
type PrefixCount struct {
	Prefix string `json:"prefix"`
	Count  int    `json:"count"`
}

var sensitivePatterns = []string{
	"secret", "password", "passwd", "token", "apikey", "api_key",
	"private", "auth", "credential", "key",
}

// Analyse returns a Summary for the given env file.
func Analyse(f *parser.EnvFile) Summary {
	if f == nil || len(f.Entries) == 0 {
		return Summary{}
	}

	prefixCounts := map[string]int{}
	longest, shortest := "", ""

	var s Summary
	s.TotalKeys = len(f.Entries)

	for _, e := range f.Entries {
		if e.Value == "" {
			s.EmptyValues++
		}
		if isSensitive(e.Key) {
			s.SensitiveKeys++
		}
		if longest == "" || len(e.Key) > len(longest) {
			longest = e.Key
		}
		if shortest == "" || len(e.Key) < len(shortest) {
			shortest = e.Key
		}
		if idx := strings.Index(e.Key, "_"); idx > 0 {
			prefixCounts[e.Key[:idx]]++
		}
	}

	s.LongestKey = longest
	s.ShortestKey = shortest
	s.UniquePrefix = len(prefixCounts)
	s.TopPrefixes = topN(prefixCounts, 5)
	return s
}

func isSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

func topN(m map[string]int, n int) []PrefixCount {
	result := make([]PrefixCount, 0, len(m))
	for k, v := range m {
		result = append(result, PrefixCount{Prefix: k, Count: v})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		return result[i].Prefix < result[j].Prefix
	})
	if len(result) > n {
		return result[:n]
	}
	return result
}
