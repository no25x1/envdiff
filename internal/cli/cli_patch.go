package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envdiff/internal/exporter"
	"github.com/yourorg/envdiff/internal/parser"
	"github.com/yourorg/envdiff/internal/patcher"
)

// runPatch applies a JSON patch file to an .env file and writes the result.
func runPatch(args []string) error {
	fs := flag.NewFlagSet("patch", flag.ContinueOnError)
	output := fs.String("output", "", "output file (default: stdout)")
	format := fs.String("format", "dotenv", "output format: dotenv|json")

	if err := fs.Parse(args); err != nil {
		return err
	}

	remaining := fs.Args()
	if len(remaining) < 2 {
		return fmt.Errorf("usage: envdiff patch <env-file> <patch-file> [flags]")
	}

	envPath := remaining[0]
	patchPath := remaining[1]

	envFile, err := parser.Parse(envPath)
	if err != nil {
		return fmt.Errorf("patch: parse env file: %w", err)
	}

	ops, err := patcher.LoadFile(patchPath)
	if err != nil {
		return fmt.Errorf("patch: load patch: %w", err)
	}

	result, err := patcher.Apply(envFile, ops)
	if err != nil {
		return fmt.Errorf("patch: apply: %w", err)
	}

	exp, err := exporter.New(*format)
	if err != nil {
		return fmt.Errorf("patch: exporter: %w", err)
	}

	w := os.Stdout
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			return fmt.Errorf("patch: create output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	return exp.Write(w, result)
}
