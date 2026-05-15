package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/exporter"
	"github.com/user/envdiff/internal/flattener"
	"github.com/user/envdiff/internal/parser"
)

func runFlatten(args []string) error {
	fs := flag.NewFlagSet("flatten", flag.ContinueOnError)
	strategy := fs.String("strategy", "last", "conflict resolution strategy: first | last | error")
	output := fs.String("output", "", "output file path (default: stdout)")
	format := fs.String("format", "dotenv", "output format: dotenv | json")

	if err := fs.Parse(args); err != nil {
		return err
	}

	paths := fs.Args()
	if len(paths) < 2 {
		return fmt.Errorf("flatten requires at least two input files")
	}

	var files []*parser.EnvFile
	for _, p := range paths {
		f, err := parser.Parse(p)
		if err != nil {
			return fmt.Errorf("flatten: cannot parse %q: %w", p, err)
		}
		files = append(files, f)
	}

	strat, err := parseStrategy(*strategy)
	if err != nil {
		return fmt.Errorf("flatten: %w", err)
	}

	res, err := flattener.Apply(files, flattener.Options{Strategy: strat})
	if err != nil {
		return fmt.Errorf("flatten: %w", err)
	}

	if len(res.Conflicts) > 0 {
		for _, c := range res.Conflicts {
			fmt.Fprintf(os.Stderr, "conflict: key %q found in [%s] — kept %q\n",
				c.Key, strings.Join(c.Files, ", "), c.Chosen)
		}
	}

	exp, err := exporter.New(*format)
	if err != nil {
		return fmt.Errorf("flatten: %w", err)
	}

	w := os.Stdout
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			return fmt.Errorf("flatten: cannot create output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	return exp.Write(w, res.File)
}

// parseStrategy converts a strategy name string into a flattener.Strategy value.
// It returns an error for unrecognised strategy names instead of silently
// falling back to a default, so callers get explicit feedback on bad input.
func parseStrategy(s string) (flattener.Strategy, error) {
	switch strings.ToLower(s) {
	case "first":
		return flattener.StrategyFirst, nil
	case "last":
		return flattener.StrategyLast, nil
	case "error":
		return flattener.StrategyError, nil
	default:
		return flattener.StrategyLast, fmt.Errorf("unknown strategy %q: must be first, last, or error", s)
	}
}
