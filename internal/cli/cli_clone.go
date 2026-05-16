package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"envdiff/internal/cloner"
	"envdiff/internal/exporter"
	"envdiff/internal/parser"
)

func runClone(args []string) error {
	fs := flag.NewFlagSet("clone", flag.ContinueOnError)
	strip := fs.String("strip-prefix", "", "prefix to strip from source keys")
	add := fs.String("add-prefix", "", "prefix to add to destination keys")
	overwrite := fs.Bool("overwrite", false, "overwrite existing keys in destination")
	keysFlag := fs.String("keys", "", "comma-separated list of source keys to clone")
	format := fs.String("format", "dotenv", "output format: dotenv or json")

	if err := fs.Parse(args); err != nil {
		return err
	}

	remaining := fs.Args()
	if len(remaining) < 1 {
		return fmt.Errorf("clone: usage: envdiff clone [flags] <source> [destination]")
	}

	srcPath := remaining[0]
	var dstPath string
	if len(remaining) >= 2 {
		dstPath = remaining[1]
	}

	src, err := parser.Parse(srcPath)
	if err != nil {
		return fmt.Errorf("clone: reading source: %w", err)
	}

	var dst *parser.EnvFile
	if dstPath != "" {
		if _, statErr := os.Stat(dstPath); statErr == nil {
			dst, err = parser.Parse(dstPath)
			if err != nil {
				return fmt.Errorf("clone: reading destination: %w", err)
			}
		}
	}

	var keyList []string
	if *keysFlag != "" {
		for _, k := range strings.Split(*keysFlag, ",") {
			if k = strings.TrimSpace(k); k != "" {
				keyList = append(keyList, k)
			}
		}
	}

	opts := cloner.Options{
		StripPrefix: *strip,
		AddPrefix:   *add,
		Overwrite:   *overwrite,
		Keys:        keyList,
	}

	out, err := cloner.Apply(src, dst, opts)
	if err != nil {
		return fmt.Errorf("clone: %w", err)
	}

	exp, err := exporter.New(*format)
	if err != nil {
		return fmt.Errorf("clone: %w", err)
	}
	return exp.Write(os.Stdout, out)
}
