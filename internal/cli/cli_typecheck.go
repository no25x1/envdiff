package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/envdiff/envdiff/internal/parser"
	"github.com/envdiff/envdiff/internal/typecheck"
)

// runTypecheck implements the `typecheck` sub-command.
// Usage: envdiff typecheck <file> <KEY=type,...>
// Example: envdiff typecheck .env PORT=int,ENABLED=bool,API_URL=url
func runTypecheck(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff typecheck <file> <KEY=type,...>")
	}

	filePath := args[0]
	rawRules := args[1]

	file, err := parser.Parse(filePath)
	if err != nil {
		return fmt.Errorf("parse %s: %w", filePath, err)
	}

	rules, err := parseTypeRules(rawRules)
	if err != nil {
		return err
	}

	violations := typecheck.Check(file, rules)

	format := flagValue(args[2:], "--format", "text")

	if format == "json" {
		return writeTypecheckJSON(violations)
	}
	return writeTypecheckText(violations)
}

func parseTypeRules(raw string) ([]typecheck.Rule, error) {
	var rules []typecheck.Rule
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		idx := strings.Index(part, "=")
		if idx < 1 {
			return nil, fmt.Errorf("invalid rule %q: expected KEY=type", part)
		}
		key := part[:idx]
		typ := typecheck.Type(strings.ToLower(part[idx+1:]))
		rules = append(rules, typecheck.Rule{Pattern: "^" + key + "$", Type: typ})
	}
	if len(rules) == 0 {
		return nil, fmt.Errorf("no valid rules provided")
	}
	return rules, nil
}

func writeTypecheckText(violations []typecheck.Violation) error {
	if len(violations) == 0 {
		fmt.Fprintln(os.Stdout, "typecheck: all values OK")
		return nil
	}
	for _, v := range violations {
		fmt.Fprintf(os.Stdout, "FAIL  %s\n", v.Error())
	}
	return fmt.Errorf("typecheck: %d violation(s) found", len(violations))
}

func writeTypecheckJSON(violations []typecheck.Violation) error {
	out := struct {
		Violations []typecheck.Violation `json:"violations"`
		Count      int                   `json:"count"`
	}{
		Violations: violations,
		Count:      len(violations),
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return err
	}
	if len(violations) > 0 {
		return fmt.Errorf("typecheck: %d violation(s) found", len(violations))
	}
	return nil
}

// flagValue is a helper to extract --flag=value or --flag value from args.
func flagValue(args []string, flag, def string) string {
	prefix := flag + "="
	for i, a := range args {
		if strings.HasPrefix(a, prefix) {
			return strings.TrimPrefix(a, prefix)
		}
		if a == flag && i+1 < len(args) {
			return args[i+1]
		}
	}
	return def
}
