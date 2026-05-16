package cli

import (
	"flag"
	"fmt"
	"strings"

	"envdiff/internal/exporter"
	"envdiff/internal/parser"
	"envdiff/internal/tagger"
)

func runTag(args []string) error {
	fs := flag.NewFlagSet("tag", flag.ContinueOnError)
	add := fs.String("add", "", "comma-separated tags to add")
	remove := fs.String("remove", "", "comma-separated tags to remove")
	filter := fs.String("filter", "", "only output entries carrying this tag")
	keys := fs.String("keys", "", "comma-separated keys to restrict operations to")
	format := fs.String("format", "dotenv", "output format: dotenv or json")
	output := fs.String("output", "", "output file path (default: stdout)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: envdiff tag [flags] <file>")
	}

	srcPath := fs.Arg(0)
	src, err := parser.Parse(srcPath)
	if err != nil {
		return fmt.Errorf("parse %s: %w", srcPath, err)
	}

	opts := tagger.Options{
		FilterTag: *filter,
	}
	if *add != "" {
		opts.Add = splitCSV(*add)
	}
	if *remove != "" {
		opts.Remove = splitCSV(*remove)
	}
	if *keys != "" {
		opts.Keys = splitCSV(*keys)
	}

	result, err := tagger.Apply(src, opts)
	if err != nil {
		return fmt.Errorf("tag: %w", err)
	}

	exp, err := exporter.New(*format)
	if err != nil {
		return err
	}
	return exp.Write(result, *output)
}

func splitCSV(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
