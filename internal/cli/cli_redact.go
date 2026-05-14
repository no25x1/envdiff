package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/exporter"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/redactor"
)

// runRedact handles the "redact" sub-command.
// Usage: envdiff redact [--placeholder=X] [--patterns=a,b] [--format=dotenv|json] <file>
func runRedact(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("redact: file argument required")
	}

	var placeholder string
	var rawPatterns string
	var format string
	var positional []string

	for _, a := range args {
		switch {
		case strings.HasPrefix(a, "--placeholder="):
			placeholder = strings.TrimPrefix(a, "--placeholder=")
		case strings.HasPrefix(a, "--patterns="):
			rawPatterns = strings.TrimPrefix(a, "--patterns=")
		case strings.HasPrefix(a, "--format="):
			format = strings.TrimPrefix(a, "--format=")
		default:
			positional = append(positional, a)
		}
	}

	if len(positional) == 0 {
		return fmt.Errorf("redact: file argument required")
	}
	filePath := positional[0]

	f, err := parser.Parse(filePath)
	if err != nil {
		return fmt.Errorf("redact: %w", err)
	}

	opts := redactor.Options{Placeholder: placeholder}
	if rawPatterns != "" {
		for _, p := range strings.Split(rawPatterns, ",") {
			if p = strings.TrimSpace(p); p != "" {
				opts.Patterns = append(opts.Patterns, p)
			}
		}
	}

	redacted := redactor.Apply(f, opts)

	if format == "" {
		format = "dotenv"
	}
	exp, err := exporter.New(format)
	if err != nil {
		return fmt.Errorf("redact: %w", err)
	}
	return exp.Write(os.Stdout, redacted)
}
