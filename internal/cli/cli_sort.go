package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/exporter"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/sorter"
)

// runSort implements the `sort` sub-command.
// Usage: envdiff sort [flags] <file>
func runSort(args []string) error {
	fs := flag.NewFlagSet("sort", flag.ContinueOnError)
	alpha := fs.Bool("alpha", true, "sort keys alphabetically")
	group := fs.Bool("group", false, "group keys by prefix (segment before first '_')")
	reverse := fs.Bool("reverse", false, "reverse the final order")
	output := fs.String("output", "", "write result to file instead of stdout")
	format := fs.String("format", "dotenv", "output format: dotenv or json")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("sort: missing required argument <file>")
	}

	path := fs.Arg(0)
	f, err := parser.Parse(path)
	if err != nil {
		return fmt.Errorf("sort: parse %q: %w", path, err)
	}

	sorted := sorter.Apply(f, sorter.Options{
		Alphabetical:  *alpha,
		GroupByPrefix: *group,
		Reverse:       *reverse,
	})

	exp, err := exporter.New(*format)
	if err != nil {
		return fmt.Errorf("sort: %w", err)
	}

	w := os.Stdout
	if *output != "" {
		f2, err := os.Create(*output)
		if err != nil {
			return fmt.Errorf("sort: create output file: %w", err)
		}
		defer f2.Close()
		w = f2
	}

	return exp.Write(w, sorted)
}
