package cli

import (
	"fmt"
	"os"

	"github.com/envdiff/envdiff/internal/exporter"
	"github.com/envdiff/envdiff/internal/parser"
	"github.com/envdiff/envdiff/internal/sanitizer"
)

// runSanitize implements the `envdiff sanitize` sub-command.
// Usage: envdiff sanitize [--format=dotenv|json] <file>
func runSanitize(args []string) error {
	format := "dotenv"
	var fileArgs []string

	for _, a := range args {
		switch {
		case len(a) > 9 && a[:9] == "--format=":
			format = a[9:]
		default:
			fileArgs = append(fileArgs, a)
		}
	}

	if len(fileArgs) < 1 {
		return fmt.Errorf("sanitize requires a file argument")
	}

	f, err := parser.Parse(fileArgs[0])
	if err != nil {
		return fmt.Errorf("parse %s: %w", fileArgs[0], err)
	}

	cleaned := sanitizer.Apply(f, sanitizer.DefaultOptions())

	exp, err := exporter.New(format)
	if err != nil {
		return fmt.Errorf("exporter: %w", err)
	}

	if err := exp.Write(os.Stdout, cleaned); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}
