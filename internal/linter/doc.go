// Package linter provides a rule-based linting engine for .env files.
//
// A Linter is constructed with a set of Rules. Each Rule inspects a single
// parser.Entry and returns zero or more Findings. Findings carry a Severity
// (error, warning, or info) and a human-readable message.
//
// Usage:
//
//	l := linter.New()               // uses default rules
//	findings := l.Lint(envFile)
//	for _, f := range findings {
//		fmt.Println(f)
//	}
//
// Custom rules can be composed with WithRules:
//
//	l := linter.WithRules(linter.RuleUppercaseKey, myCustomRule)
package linter
