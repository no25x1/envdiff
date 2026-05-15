package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/exporter"
	"github.com/user/envdiff/internal/interpolator"
	"github.com/user/envdiff/internal/parser"
)

// runInterpolate handles the `interpolate` sub-command.
// Usage: envdiff interpolate [--override KEY=VAL] [--fail-on-missing] [--format dotenv|json] <file>
func runInterpolate(args []string) error {
	fs := flag.NewFlagSet("interpolate", flag.ContinueOnError)
	overrides := fs.String("override", "", "comma-separated KEY=VALUE pairs that override file values during expansion")
	failOnMissing := fs.Bool("fail-on-missing", false, "return an error when a referenced variable cannot be resolved")
	format := fs.String("format", "dotenv", "output format: dotenv or json")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("interpolate: missing required argument <file>")
	}

	filePath := fs.Arg(0)
	f, err := parser.Parse(filePath)
	if err != nil {
		return fmt.Errorf("interpolate: parse %q: %w", filePath, err)
	}

	overrideMap, err := parseOverrides(*overrides)
	if err != nil {
		return fmt.Errorf("interpolate: invalid --override value: %w", err)
	}

	opts := interpolator.Options{
		FailOnMissing: *failOnMissing,
		Overrides:     overrideMap,
	}

	result, err := interpolator.Apply(f, opts)
	if err != nil {
		return fmt.Errorf("interpolate: %w", err)
	}

	exp, err := exporter.New(*format)
	if err != nil {
		return fmt.Errorf("interpolate: %w", err)
	}
	return exp.Write(os.Stdout, result)
}

// parseOverrides converts "KEY1=VAL1,KEY2=VAL2" into a map.
// Returns an error if a pair is missing the '=' separator.
func parseOverrides(raw string) (map[string]string, error) {
	m := make(map[string]string)
	if raw == "" {
		return m, nil
	}
	for _, pair := range strings.Split(raw, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed pair %q: expected KEY=VALUE", pair)
		}
		m[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return m, nil
}
