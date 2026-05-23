package cli

import (
	"fmt"
	"os"

	"envdiff/internal/cascader"
	"envdiff/internal/exporter"
	"envdiff/internal/parser"
)

func runCascade(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff cascade <base> <overlay> [overlay...] [--stop-on-missing] [--format dotenv|json]")
	}

	stopOnMissing := false
	format := "dotenv"
	output := ""
	var filePaths []string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--stop-on-missing":
			stopOnMissing = true
		case "--format":
			if i+1 >= len(args) {
				return fmt.Errorf("--format requires a value")
			}
			i++
			format = args[i]
		case "--output":
			if i+1 >= len(args) {
				return fmt.Errorf("--output requires a value")
			}
			i++
			output = args[i]
		default:
			filePaths = append(filePaths, args[i])
		}
	}

	if len(filePaths) < 2 {
		return fmt.Errorf("cascade requires at least a base file and one overlay")
	}

	files := make([]*parser.EnvFile, 0, len(filePaths))
	for _, p := range filePaths {
		f, err := parser.Parse(p)
		if err != nil {
			return fmt.Errorf("cascade: cannot parse %s: %w", p, err)
		}
		files = append(files, f)
	}

	opts := cascader.DefaultOptions()
	opts.StopOnMissing = stopOnMissing

	result, err := cascader.Apply(files, opts)
	if err != nil {
		return fmt.Errorf("cascade: %w", err)
	}

	w := os.Stdout
	if output != "" {
		f, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("cascade: cannot create output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	expr, err := exporter.New(format)
	if err != nil {
		return fmt.Errorf("cascade: %w", err)
	}
	return expr.Write(w, result)
}
