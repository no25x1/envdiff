package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/deduplicator"
	"github.com/user/envdiff/internal/exporter"
	"github.com/user/envdiff/internal/parser"
)

func runDedup(args []string) error {
	fs := flag.NewFlagSet("dedup", flag.ContinueOnError)
	strategy := fs.String("strategy", "last", "which occurrence to keep: first|last")
	output := fs.String("output", "", "output file (default: stdout)")
	format := fs.String("format", "dotenv", "output format: dotenv|json")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: envdiff dedup [flags] <file>")
	}

	filePath := fs.Arg(0)
	env, err := parser.Parse(filePath)
	if err != nil {
		return fmt.Errorf("parse %q: %w", filePath, err)
	}

	var strat deduplicator.Strategy
	switch *strategy {
	case "first":
		strat = deduplicator.KeepFirst
	case "last":
		strat = deduplicator.KeepLast
	default:
		return fmt.Errorf("unknown strategy %q: choose first or last", *strategy)
	}

	res := deduplicator.Apply(env, deduplicator.Options{Strategy: strat})

	if len(res.Removed) > 0 {
		fmt.Fprintf(os.Stderr, "dedup: removed duplicates for %d key(s):\n", len(res.Removed))
		for _, d := range res.Removed {
			fmt.Fprintf(os.Stderr, "  %s (%d occurrences)\n", d.Key, d.Count)
		}
	}

	exp, err := exporter.New(*format)
	if err != nil {
		return err
	}

	var out *os.File
	if *output == "" {
		out = os.Stdout
	} else {
		out, err = os.Create(*output)
		if err != nil {
			return fmt.Errorf("create output file: %w", err)
		}
		defer out.Close()
	}

	return exp.Write(res.File, out)
}
