package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"envdiff/internal/exporter"
	"envdiff/internal/parser"
	"envdiff/internal/selector"
)

func runSelect(args []string) error {
	fs := flag.NewFlagSet("select", flag.ContinueOnError)
	keys := fs.String("keys", "", "comma-separated list of keys to select")
	prefix := fs.String("prefix", "", "select keys with this prefix")
	regex := fs.String("regex", "", "select keys matching this regex")
	invert := fs.Bool("invert", false, "invert selection")
	format := fs.String("format", "dotenv", "output format: dotenv or json")
	output := fs.String("output", "", "output file (default: stdout)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("select: usage: envdiff select [flags] <file>")
	}

	src, err := parser.Parse(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("select: failed to parse %s: %w", fs.Arg(0), err)
	}

	var keyList []string
	if *keys != "" {
		for _, k := range strings.Split(*keys, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				keyList = append(keyList, k)
			}
		}
	}

	opts := selector.Options{
		Keys:   keyList,
		Prefix: *prefix,
		Regex:  *regex,
		Invert: *invert,
	}

	result, err := selector.Apply(src, opts)
	if err != nil {
		return fmt.Errorf("select: %w", err)
	}

	exp, err := exporter.New(*format)
	if err != nil {
		return fmt.Errorf("select: %w", err)
	}

	w := os.Stdout
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			return fmt.Errorf("select: cannot create output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	return exp.Write(w, result)
}
