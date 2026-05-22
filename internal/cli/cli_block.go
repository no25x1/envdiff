package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"envdiff/internal/blocker"
	"envdiff/internal/parser"
)

func runBlock(args []string) error {
	fs := flag.NewFlagSet("block", flag.ContinueOnError)
	keys := fs.String("keys", "", "comma-separated list of forbidden key names")
	prefixes := fs.String("prefixes", "", "comma-separated list of forbidden key prefixes")
	values := fs.String("values", "", "comma-separated list of forbidden values")
	format := fs.String("format", "text", "output format: text or json")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: envdiff block [flags] <file>")
	}

	path := fs.Arg(0)
	f, err := parser.Parse(path)
	if err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}

	opts := blocker.Options{
		ForbiddenKeys:     splitCSV(*keys),
		ForbiddenPrefixes: splitCSV(*prefixes),
		ForbiddenValues:   splitCSV(*values),
	}

	if len(opts.ForbiddenKeys)+len(opts.ForbiddenPrefixes)+len(opts.ForbiddenValues) == 0 {
		return fmt.Errorf("at least one of --keys, --prefixes, or --values must be set")
	}

	violations, err := blocker.Apply(f, opts)
	if err != nil {
		return err
	}

	switch strings.ToLower(*format) {
	case "json":
		return json.NewEncoder(os.Stdout).Encode(violations)
	default:
		if len(violations) == 0 {
			fmt.Println("No blocked entries found.")
			return nil
		}
		for _, v := range violations {
			fmt.Printf("BLOCKED  %s  %s\n", v.Key, v.Reason)
		}
		return fmt.Errorf("%d blocked entry/entries detected", len(violations))
	}
}
