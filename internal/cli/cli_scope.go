package cli

import (
	"flag"
	"fmt"
	"os"

	"envdiff/internal/parser"
	"envdiff/internal/scoper"
)

func runScope(args []string) error {
	fs := flag.NewFlagSet("scope", flag.ContinueOnError)
	scope := fs.String("scope", "", "scope prefix to extract (required)")
	strip := fs.Bool("strip", false, "strip the scope prefix from keys")
	keepUnscoped := fs.Bool("keep-unscoped", false, "include keys with no recognised scope")
	output := fs.String("output", "", "write result to file instead of stdout")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("scope: usage: envdiff scope -scope=<NAME> [flags] <file>")
	}
	if *scope == "" {
		return fmt.Errorf("scope: -scope flag is required")
	}

	src, err := parser.ParseFile(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("scope: %w", err)
	}

	out, err := scoper.Apply(src, scoper.Options{
		Scope:        *scope,
		StripPrefix:  *strip,
		KeepUnscoped: *keepUnscoped,
	})
	if err != nil {
		return fmt.Errorf("scope: %w", err)
	}

	w := os.Stdout
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			return fmt.Errorf("scope: %w", err)
		}
		defer f.Close()
		w = f
	}

	for _, e := range out.Entries {
		if e.Comment != "" {
			fmt.Fprintf(w, "# %s\n", e.Comment)
		}
		fmt.Fprintf(w, "%s=%s\n", e.Key, e.Value)
	}
	return nil
}
