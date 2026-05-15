// Package schema provides validation of .env files against a schema definition
// that specifies required keys, optional keys, and allowed value patterns.
package schema

import (
	"fmt"
	"regexp"

	"github.com/envdiff/envdiff/internal/parser"
)

// FieldRule describes constraints for a single env key.
type FieldRule struct {
	Key      string
	Required bool
	Pattern  string // optional regex pattern the value must match
}

// Schema holds the set of rules for an env file.
type Schema struct {
	Rules []FieldRule
}

// Violation represents a single schema violation.
type Violation struct {
	Key     string
	Message string
}

func (v Violation) Error() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// Validate checks the given EnvFile against the schema and returns any violations.
func (s *Schema) Validate(f *parser.EnvFile) []Violation {
	var violations []Violation

	present := make(map[string]string, len(f.Entries))
	for _, e := range f.Entries {
		present[e.Key] = e.Value
	}

	for _, rule := range s.Rules {
		val, ok := present[rule.Key]
		if !ok {
			if rule.Required {
				violations = append(violations, Violation{Key: rule.Key, Message: "required key is missing"})
			}
			continue
		}
		if rule.Pattern != "" {
			re, err := regexp.Compile(rule.Pattern)
			if err != nil {
				violations = append(violations, Violation{Key: rule.Key, Message: fmt.Sprintf("invalid pattern %q: %v", rule.Pattern, err)})
				continue
			}
			if !re.MatchString(val) {
				violations = append(violations, Violation{Key: rule.Key, Message: fmt.Sprintf("value %q does not match pattern %q", val, rule.Pattern)})
			}
		}
	}
	return violations
}
