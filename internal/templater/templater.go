// Package templater provides functionality to render .env files
// from a template with variable substitution and missing key detection.
package templater

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// varPattern matches ${VAR_NAME} or $VAR_NAME style placeholders.
var varPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

// Result holds the output of a template render operation.
type Result struct {
	Entries []parser.Entry
	Missing []string
}

// Render takes a template EnvFile and a values EnvFile, substitutes
// placeholders in template values with values from the values file,
// and returns a Result. Entries whose placeholders cannot be resolved
// are included with their raw template value; missing keys are collected.
func Render(tmpl parser.EnvFile, values parser.EnvFile) Result {
	valMap := make(map[string]string, len(values.Entries))
	for _, e := range values.Entries {
		valMap[e.Key] = e.Value
	}

	missingSet := make(map[string]struct{})
	var entries []parser.Entry

	for _, e := range tmpl.Entries {
		resolved, missing := substitute(e.Value, valMap)
		for _, m := range missing {
			missingSet[m] = struct{}{}
		}
		entries = append(entries, parser.Entry{Key: e.Key, Value: resolved})
	}

	var missing []string
	for k := range missingSet {
		missing = append(missing, k)
	}

	return Result{Entries: entries, Missing: missing}
}

// substitute replaces all variable references in s using vals.
// Returns the substituted string and a slice of unresolved variable names.
func substitute(s string, vals map[string]string) (string, []string) {
	var missing []string
	result := varPattern.ReplaceAllStringFunc(s, func(match string) string {
		name := extractName(match)
		if v, ok := vals[name]; ok {
			return v
		}
		missing = append(missing, name)
		return match
	})
	return result, missing
}

// extractName pulls the variable name out of a $VAR or ${VAR} token.
func extractName(token string) string {
	token = strings.TrimPrefix(token, "$")
	token = strings.Trim(token, "{}")
	return token
}

// Validate checks that every placeholder in the template has a corresponding
// key in the values file and returns a list of human-readable error strings.
func Validate(tmpl parser.EnvFile, values parser.EnvFile) []string {
	res := Render(tmpl, values)
	if len(res.Missing) == 0 {
		return nil
	}
	var errs []string
	for _, m := range res.Missing {
		errs = append(errs, fmt.Sprintf("missing value for placeholder: %s", m))
	}
	return errs
}
