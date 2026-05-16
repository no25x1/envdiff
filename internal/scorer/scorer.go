// Package scorer computes a quality score for an env file based on
// common hygiene rules: key naming, value presence, secret exposure, etc.
package scorer

import (
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Result holds the overall score and per-rule breakdown.
type Result struct {
	Score      int            // 0-100
	Max        int            // maximum possible score
	Breakdown  []RuleResult
}

// RuleResult captures the outcome of a single scoring rule.
type RuleResult struct {
	Rule    string
	Passed  int
	Failed  int
	Weight  int
	Comment string
}

// Score evaluates the given env file and returns a Result.
func Score(f *parser.EnvFile) Result {
	if f == nil || len(f.Entries) == 0 {
		return Result{Score: 0, Max: 100}
	}

	rules := []struct {
		name   string
		weight int
		fn     func(*parser.EnvFile) (int, int, string)
	}{
		{"no-empty-values", 30, checkEmptyValues},
		{"uppercase-keys", 25, checkUppercaseKeys},
		{"no-inline-secrets", 25, checkInlineSecrets},
		{"no-whitespace-in-keys", 20, checkWhitespaceKeys},
	}

	var totalWeight, weightedScore int
	var breakdown []RuleResult

	for _, r := range rules {
		passed, failed, comment := r.fn(f)
		total := passed + failed
		var contribution int
		if total > 0 {
			contribution = (passed * r.weight) / total
		}
		weightedScore += contribution
		totalWeight += r.weight
		breakdown = append(breakdown, RuleResult{
			Rule:    r.name,
			Passed:  passed,
			Failed:  failed,
			Weight:  r.weight,
			Comment: comment,
		})
	}

	score := 0
	if totalWeight > 0 {
		score = (weightedScore * 100) / totalWeight
	}
	return Result{Score: score, Max: 100, Breakdown: breakdown}
}

func checkEmptyValues(f *parser.EnvFile) (int, int, string) {
	passed, failed := 0, 0
	for _, e := range f.Entries {
		if strings.TrimSpace(e.Value) == "" {
			failed++
		} else {
			passed++
		}
	}
	return passed, failed, "keys with non-empty values"
}

func checkUppercaseKeys(f *parser.EnvFile) (int, int, string) {
	passed, failed := 0, 0
	for _, e := range f.Entries {
		if e.Key == strings.ToUpper(e.Key) {
			passed++
		} else {
			failed++
		}
	}
	return passed, failed, "keys that are fully uppercase"
}

var secretPatterns = []string{"password", "secret", "token", "apikey", "api_key", "passwd", "private_key"}

func checkInlineSecrets(f *parser.EnvFile) (int, int, string) {
	passed, failed := 0, 0
	for _, e := range f.Entries {
		lower := strings.ToLower(e.Key)
		sensitive := false
		for _, p := range secretPatterns {
			if strings.Contains(lower, p) {
				sensitive = true
				break
			}
		}
		if sensitive && strings.TrimSpace(e.Value) != "" && !strings.HasPrefix(e.Value, "${") {
			failed++
		} else {
			passed++
		}
	}
	return passed, failed, "sensitive keys without plain-text values"
}

func checkWhitespaceKeys(f *parser.EnvFile) (int, int, string) {
	passed, failed := 0, 0
	for _, e := range f.Entries {
		if strings.ContainsAny(e.Key, " \t") {
			failed++
		} else {
			passed++
		}
	}
	return passed, failed, "keys without whitespace"
}
