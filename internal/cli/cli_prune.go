package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/envdiff/internal/exporter"
	"github.com/yourusername/envdiff/internal/parser"
	"github.com/yourusername/envdiff/internal/pruner"
)

func runPrune(args []string) error {
	fs := flag.NewFlagSet("prune", flag.ContinueOnError)
	removeEmpty := fs.Bool("empty", true, "remove keys with empty values")
	removeComments := fs.Bool("comments", false, "remove pure-comment / blank lines")
	allowlist := fs.String("allow", "", "comma-separated keys to never prune")
	output := fs.String("out", "", "write result to file instead of stdout")
	format := fs.String("format", "dotenv", "output format: dotenv or json")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("prune: usage: envdiff prune <file>")
	}

	filePath := fs.Arg(0)
	f, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("prune: cannot read %q: %w", filePath, err)
	}

	var allow []string
	if *allowlist != "" {
		for _, k := range strings.Split(*allowlist, ",") {
			if k = strings.TrimSpace(k); k != "" {
				allow = append(allow, k)
			}
		}
	}

	opts := pruner.Options{
		RemoveEmpty:     *removeEmpty,
		RemoveCommented: *removeComments,
		Allowlist:       allow,
	}

	result, err := pruner.Apply(f, opts)
	if err != nil {
		return fmt.Errorf("prune: %w", err)
	}

	w := os.Stdout
	if *output != "" {
		w, err = os.Create(*output)
		if err != nil {
			return fmt.Errorf("prune: cannot create output file: %w", err)
		}
		defer w.Close()
	}

	exp, err := exporter.New(*format)
	if err != nil {
		return fmt.Errorf("prune: %w", err)
	}
	return exp.Write(w, result)
}
