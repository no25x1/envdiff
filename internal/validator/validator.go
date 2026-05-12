package validator

import (
	"fmt"
	"strings"

	"github.com/envdiff/internal/parser"
)

// Rule defines a validation rule applied to env file entries.
type Rule struct {
	Name    string
	Check   func(key, value string) error
}

// Result holds the outcome of validating a single entry.
type Result struct {
	Key     string
	Rule    string
	Message string
}

// Validator runs a set of rules against an EnvFile.
type Validator struct {
	rules []Rule
}

// New creates a Validator with the default built-in rules.
func New() *Validator {
	return &Validator{
		rules: defaultRules(),
	}
}

// WithRules creates a Validator with the provided rules appended to defaults.
func WithRules(extra ...Rule) *Validator {
	v := New()
	v.rules = append(v.rules, extra...)
	return v
}

// Validate checks every entry in the file against all rules.
// It returns a slice of Result for every violation found.
func (v *Validator) Validate(f *parser.EnvFile) []Result {
	var results []Result
	for _, entry := range f.Entries {
		for _, rule := range v.rules {
			if err := rule.Check(entry.Key, entry.Value); err != nil {
				results = append(results, Result{
					Key:     entry.Key,
					Rule:    rule.Name,
					Message: err.Error(),
				})
			}
		}
	}
	return results
}

func defaultRules() []Rule {
	return []Rule{
		{
			Name: "no-empty-key",
			Check: func(key, _ string) error {
				if strings.TrimSpace(key) == "" {
					return fmt.Errorf("key must not be empty")
				}
				return nil
			},
		},
		{
			Name: "no-spaces-in-key",
			Check: func(key, _ string) error {
				if strings.ContainsAny(key, " \t") {
					return fmt.Errorf("key %q contains whitespace", key)
				}
				return nil
			},
		},
		{
			Name: "uppercase-key",
			Check: func(key, _ string) error {
				if key != strings.ToUpper(key) {
					return fmt.Errorf("key %q should be uppercase", key)
				}
				return nil
			},
		},
		{
			Name: "no-empty-value",
			Check: func(key, value string) error {
				if strings.TrimSpace(value) == "" {
					return fmt.Errorf("key %q has an empty value", key)
				}
				return nil
			},
		},
	}
}
