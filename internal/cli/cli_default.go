package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/defaulter"
	"github.com/user/envdiff/internal/parser"
)

func runDefault(args []string) error {
	fs := flag.NewFlagSet("default", flag.ContinueOnError)
	overwrite := fs.Bool("overwrite", false, "overwrite existing keys with default values")
	keysFlag := fs.String("keys", "", "comma-separated list of keys to process (default: all)")
	outFlag := fs.String("out", "", "output file path (default: stdout)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 2 {
		return fmt.Errorf("usage: envdiff default <defaults.env> <target.env> [flags]")
	}

	defaultsPath := fs.Arg(0)
	targetPath := fs.Arg(1)

	defFile, err := parser.Parse(defaultsPath)
	if err != nil {
		return fmt.Errorf("default: reading defaults file: %w", err)
	}

	tgtFile, err := parser.Parse(targetPath)
	if err != nil {
		return fmt.Errorf("default: reading target file: %w", err)
	}

	var keyList []string
	if *keysFlag != "" {
		for _, k := range strings.Split(*keysFlag, ",") {
			if k = strings.TrimSpace(k); k != "" {
				keyList = append(keyList, k)
			}
		}
	}

	opts := defaulter.Options{
		Overwrite: *overwrite,
		Keys:      keyList,
	}

	out, results, err := defaulter.Apply(defFile, tgtFile, opts)
	if err != nil {
		return err
	}

	// Print summary to stderr.
	applied, skipped := 0, 0
	for _, r := range results {
		if r.Skipped {
			skipped++
		} else {
			applied++
		}
	}
	fmt.Fprintf(os.Stderr, "defaulter: applied=%d skipped=%d\n", applied, skipped)

	// Write output.
	w := os.Stdout
	if *outFlag != "" {
		f, err := os.Create(*outFlag)
		if err != nil {
			return fmt.Errorf("default: creating output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	for _, e := range out.Entries {
		fmt.Fprintf(w, "%s=%s\n", e.Key, e.Value)
	}
	return nil
}
