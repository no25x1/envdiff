package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/envdiff/internal/comparator"
	"github.com/envdiff/internal/parser"
)

func runCompare(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("compare", flag.ContinueOnError)
	conflictsOnly := fs.Bool("conflicts", false, "show only keys with conflicting values")
	fs.SetOutput(out)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 2 {
		return fmt.Errorf("compare requires at least two files: envdiff compare <file1> <file2> [fileN...]")
	}

	files := make(map[string]*parser.EnvFile, fs.NArg())
	for _, path := range fs.Args() {
		ef, err := parser.Parse(path)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		files[path] = ef
	}

	result := comparator.Compare(files)

	keys := result.Keys
	if *conflictsOnly {
		keys = result.Conflicts()
		if len(keys) == 0 {
			fmt.Fprintln(out, "no conflicts found")
			return nil
		}
	}

	// header row
	fmt.Fprintf(out, "%-30s", "KEY")
	for _, label := range result.Files {
		fmt.Fprintf(out, "  %-20s", label)
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, strings.Repeat("-", 30+len(result.Files)*22))

	for _, key := range keys {
		fmt.Fprintf(out, "%-30s", key)
		for _, label := range result.Files {
			v := result.Matrix[key][label]
			if v == "" {
				v = "(absent)"
			}
			if len(v) > 20 {
				v = v[:17] + "..."
			}
			fmt.Fprintf(out, "  %-20s", v)
		}
		fmt.Fprintln(out)
	}
	return nil
}

func init() {
	_ = os.Stderr // ensure os import used
}
