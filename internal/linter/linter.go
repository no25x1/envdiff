package linter

import (
	"fmt"
	"strings"

	"github.com/envdiff/internal/parser"
)

// Severity represents the level of a lint finding.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Finding represents a single lint result for a key.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s] %s: %s", f.Severity, f.Key, f.Message)
}

// Rule is a function that inspects an entry and returns findings.
type Rule func(entry parser.Entry) []Finding

// Linter applies rules to env file entries.
type Linter struct {
	rules []Rule
}

// New creates a Linter with the default rule set.
func New() *Linter {
	return &Linter{rules: defaultRules()}
}

// WithRules creates a Linter with the provided rules only.
func WithRules(rules ...Rule) *Linter {
	return &Linter{rules: rules}
}

// Lint runs all rules against every entry in the file.
func (l *Linter) Lint(file *parser.EnvFile) []Finding {
	var findings []Finding
	for _, entry := range file.Entries {
		for _, rule := range l.rules {
			findings = append(findings, rule(entry)...)
		}
	}
	return findings
}

func defaultRules() []Rule {
	return []Rule{
		RuleNoEmptyValue,
		RuleUppercaseKey,
		RuleNoWhitespaceInKey,
		RuleNoLeadingUnderscore,
	}
}

// RuleNoEmptyValue warns when a key has no value.
func RuleNoEmptyValue(e parser.Entry) []Finding {
	if strings.TrimSpace(e.Value) == "" {
		return []Finding{{Key: e.Key, Message: "value is empty", Severity: SeverityWarning}}
	}
	return nil
}

// RuleUppercaseKey errors when a key contains lowercase letters.
func RuleUppercaseKey(e parser.Entry) []Finding {
	if e.Key != strings.ToUpper(e.Key) {
		return []Finding{{Key: e.Key, Message: "key should be uppercase", Severity: SeverityError}}
	}
	return nil
}

// RuleNoWhitespaceInKey errors when a key contains whitespace.
func RuleNoWhitespaceInKey(e parser.Entry) []Finding {
	if strings.ContainsAny(e.Key, " \t") {
		return []Finding{{Key: e.Key, Message: "key must not contain whitespace", Severity: SeverityError}}
	}
	return nil
}

// RuleNoLeadingUnderscore warns when a key starts with an underscore.
func RuleNoLeadingUnderscore(e parser.Entry) []Finding {
	if strings.HasPrefix(e.Key, "_") {
		return []Finding{{Key: e.Key, Message: "key starts with underscore", Severity: SeverityWarning}}
	}
	return nil
}
