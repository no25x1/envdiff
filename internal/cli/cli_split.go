package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"envdiff/internal/parser"
	"envdiff/internal/splitter"
)

func runSplit(args []string) error {
	fs := flag.NewFlagSet("split", flag.ContinueOnError)
	prefixFlag := fs.String("prefixes", "", "comma-separated name=PREFIX pairs, e.g. db=DB_,aws=AWS_")
	includeOther := fs.Bool("include-other", false, "include unmatched keys in an 'other' group")
	outDir := fs.String("out", ".", "directory to write split .env files")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("split: usage: envdiff split [flags] <file>")
	}
	if *prefixFlag == "" {
		return fmt.Errorf("split: --prefixes is required")
	}

	prefixes, err := parsePrefixFlag(*prefixFlag)
	if err != nil {
		return err
	}

	srcPath := fs.Arg(0)
	src, err := parser.Parse(srcPath)
	if err != nil {
		return fmt.Errorf("split: cannot parse %s: %w", srcPath, err)
	}

	result, err := splitter.Apply(src, splitter.Options{
		Prefixes:         prefixes,
		IncludeUnmatched: *includeOther,
	})
	if err != nil {
		return fmt.Errorf("split: %w", err)
	}

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		return fmt.Errorf("split: cannot create output directory: %w", err)
	}

	for group, file := range result {
		outPath := fmt.Sprintf("%s/%s.env", *outDir, group)
		if err := writeEnvFile(outPath, file); err != nil {
			return fmt.Errorf("split: cannot write %s: %w", outPath, err)
		}
		fmt.Printf("wrote %d keys to %s\n", len(file.Entries), outPath)
	}
	return nil
}

func parsePrefixFlag(raw string) (map[string]string, error) {
	prefixes := make(map[string]string)
	for _, pair := range strings.Split(raw, ",") {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("split: invalid prefix pair %q (want name=PREFIX)", pair)
		}
		prefixes[parts[0]] = parts[1]
	}
	return prefixes, nil
}

func writeEnvFile(path string, file *parser.EnvFile) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, e := range file.Entries {
		if _, err := fmt.Fprintf(f, "%s=%s\n", e.Key, e.Value); err != nil {
			return err
		}
	}
	return nil
}
