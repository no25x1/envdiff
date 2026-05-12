package cli

import (
	"fmt"
	"os"

	"github.com/envdiff/internal/linter"
	"github.com/envdiff/internal/parser"
)

// runLint parses the given .env file and runs the linter against it.
// It prints findings to stdout and exits with code 1 if any errors are found.
func runLint(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("lint requires a file argument: envdiff lint <file>")
	}

	filePath := args[0]

	envFile, err := parser.Parse(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	l := linter.New()
	findings := l.Lint(envFile)

	if len(findings) == 0 {
		fmt.Printf("✔  %s: no lint findings\n", filePath)
		return nil
	}

	errorCount := 0
	for _, f := range findings {
		fmt.Fprintln(os.Stdout, f.String())
		if f.Severity == linter.SeverityError {
			errorCount++
		}
	}

	fmt.Printf("\n%d finding(s) in %s\n", len(findings), filePath)

	if errorCount > 0 {
		return fmt.Errorf("%d error(s) found", errorCount)
	}
	return nil
}
