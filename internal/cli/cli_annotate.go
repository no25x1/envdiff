package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/annotator"
	"github.com/user/envdiff/internal/exporter"
	"github.com/user/envdiff/internal/parser"
)

// runAnnotate implements the `annotate` sub-command.
// Usage: envdiff annotate <file> [--overwrite] [--out <file>] [--format dotenv|json]
func runAnnotate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("annotate requires a file argument")
	}

	filePath := args[0]
	overwrite := false
	outPath := ""
	format := "dotenv"

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--overwrite":
			overwrite = true
		case "--out":
			if i+1 >= len(args) {
				return fmt.Errorf("--out requires a value")
			}
			i++
			outPath = args[i]
		case "--format":
			if i+1 >= len(args) {
				return fmt.Errorf("--format requires a value")
			}
			i++
			format = strings.ToLower(args[i])
		}
	}

	env, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parse %s: %w", filePath, err)
	}

	opts := annotator.DefaultOptions()
	opts.OverwriteExisting = overwrite
	annotated := annotator.Apply(env, opts)

	exp, err := exporter.New(format)
	if err != nil {
		return fmt.Errorf("exporter: %w", err)
	}

	w := os.Stdout
	if outPath != "" {
		f, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("create %s: %w", outPath, err)
		}
		defer f.Close()
		w = f
	}

	return exp.Write(w, annotated)
}
