package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/your-org/envdiff/internal/exporter"
	"github.com/your-org/envdiff/internal/importer"
)

func runImport(args []string) error {
	fs := flag.NewFlagSet("import", flag.ContinueOnError)
	format := fs.String("format", "", "input format: json, dotenv (default: inferred from extension)")
	prefix := fs.String("prefix", "", "key prefix to add to all imported keys")
	output := fs.String("output", "", "output file (default: stdout)")
	outFmt := fs.String("out-format", "dotenv", "output format: dotenv, json")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("import: requires <file> argument")
	}

	srcPath := fs.Arg(0)

	opts := importer.Options{
		Format: importer.Format(*format),
		Prefix: *prefix,
	}

	file, err := importer.Apply(srcPath, opts)
	if err != nil {
		return err
	}

	exp, err := exporter.New(*outFmt)
	if err != nil {
		return fmt.Errorf("import: %w", err)
	}

	w := os.Stdout
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			return fmt.Errorf("import: create output: %w", err)
		}
		defer f.Close()
		w = f
	}

	return exp.Write(w, file)
}
