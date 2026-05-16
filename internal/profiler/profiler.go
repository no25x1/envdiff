// Package profiler analyses an EnvFile and produces statistics about its
// contents: key count, empty values, sensitive key ratio, prefix distribution.
package profiler

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Profile holds statistical information about an EnvFile.
type Profile struct {
	TotalKeys      int
	EmptyValues    int
	SensitiveKeys  int
	CommentedLines int
	PrefixCounts   map[string]int // top-level prefix (before first "_") -> count
}

// sensitivePatterns mirrors the default patterns used by the masker/redactor.
var sensitivePatterns = []string{
	"PASSWORD", "SECRET", "TOKEN", "KEY", "PRIVATE", "CREDENTIAL", "AUTH",
}

// Analyse returns a Profile for the given EnvFile.
func Analyse(f *parser.EnvFile) Profile {
	if f == nil {
		return Profile{PrefixCounts: map[string]int{}}
	}

	p := Profile{PrefixCounts: map[string]int{}}

	for _, e := range f.Entries {
		if e.Comment != "" && e.Key == "" {
			p.CommentedLines++
			continue
		}
		p.TotalKeys++
		if e.Value == "" {
			p.EmptyValues++
		}
		if isSensitive(e.Key) {
			p.SensitiveKeys++
		}
		prefix := topPrefix(e.Key)
		if prefix != "" {
			p.PrefixCounts[prefix]++
		}
	}
	return p
}

// TopPrefixes returns prefix names sorted by descending count, limited to n.
func (p Profile) TopPrefixes(n int) []string {
	type kv struct {
		key   string
		count int
	}
	var pairs []kv
	for k, v := range p.PrefixCounts {
		pairs = append(pairs, kv{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].key < pairs[j].key
	})
	result := make([]string, 0, n)
	for i, p := range pairs {
		if i >= n {
			break
		}
		result = append(result, p.key)
	}
	return result
}

func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, pat := range sensitivePatterns {
		if strings.Contains(upper, pat) {
			return true
		}
	}
	return false
}

func topPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return ""
}
