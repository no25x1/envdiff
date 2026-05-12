package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/envdiff/envdiff/internal/exporter"
	"github.com/envdiff/envdiff/internal/parser"
	"github.com/envdiff/envdiff/internal/transformer"
)

// runTransform handles the "transform" sub-command.
// Usage: envdiff transform [flags] <file>
func runTransform(args []string) error {
	fs := flag.NewFlagSet("transform", flag.ContinueOnError)
	prefixAdd := fs.String("prefix-add", "", "prefix to prepend to every key")
	prefixStrip := fs.String("prefix-strip", "", "prefix to strip from every key (non-matching entries are dropped)")
	keyUpper := fs.Bool("key-upper", false, "convert all keys to upper-case")
	keyLower := fs.Bool("key-lower", false, "convert all keys to lower-case")
	outFormat := fs.String("format", "dotenv", "output format: dotenv or json")
	outFile := fs.String("out", "", "write output to file instead of stdout")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("transform: missing <file> argument")
	}

	if *keyUpper && *keyLower {
		return fmt.Errorf("transform: --key-upper and --key-lower are mutually exclusive")
	}

	filePath := fs.Arg(0)
	envFile, err := parser.Parse(filePath)
	if err != nil {
		return fmt.Errorf("transform: parse %s: %w", filePath, err)
	}

	opts := transformer.Options{
		PrefixAdd:   *prefixAdd,
		PrefixStrip: *prefixStrip,
		KeyToUpper:  *keyUpper,
		KeyToLower:  *keyLower,
	}

	result := transformer.Apply(envFile, opts)

	w := os.Stdout
	if *outFile != "" {
		f, err := os.Create(*outFile)
		if err != nil {
			return fmt.Errorf("transform: create output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	exp, err := exporter.New(strings.ToLower(*outFormat), w)
	if err != nil {
		return fmt.Errorf("transform: %w", err)
	}
	return exp.Write(result)
}
