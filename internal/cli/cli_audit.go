package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/envdiff/envdiff/internal/auditor"
	"github.com/envdiff/envdiff/internal/diff"
	"github.com/envdiff/envdiff/internal/parser"
)

func runAudit(args []string) error {
	fs := flag.NewFlagSet("audit", flag.ContinueOnError)
	actor := fs.String("actor", "", "identity of the actor performing the change")
	format := fs.String("format", "text", "output format: text or json")

	if err := fs.Parse(args); err != nil {
		return err
	}

	remaining := fs.Args()
	if len(remaining) < 2 {
		return fmt.Errorf("audit requires <base> <compare> arguments")
	}

	baseFile, err := parser.Parse(remaining[0])
	if err != nil {
		return fmt.Errorf("audit: reading base file: %w", err)
	}

	cmpFile, err := parser.Parse(remaining[1])
	if err != nil {
		return fmt.Errorf("audit: reading compare file: %w", err)
	}

	results := diff.Diff(baseFile, cmpFile)

	log := auditor.FromDiff(results, auditor.Options{
		Actor: *actor,
		File:  remaining[1],
	})

	if err := log.Write(os.Stdout, auditor.Format(*format)); err != nil {
		return fmt.Errorf("audit: writing output: %w", err)
	}

	fmt.Fprintln(os.Stderr, log.Summary())
	return nil
}
