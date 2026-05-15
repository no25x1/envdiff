package cli

import (
	"flag"
	"fmt"
	"strings"

	"github.com/your-org/envdiff/internal/encoder"
	"github.com/your-org/envdiff/internal/exporter"
	"github.com/your-org/envdiff/internal/parser"
)

func runEncode(args []string) error {
	fs := flag.NewFlagSet("encode", flag.ContinueOnError)
	strategy := fs.String("strategy", "base64", "encoding strategy: base64 | url")
	decode := fs.Bool("decode", false, "decode instead of encode")
	keys := fs.String("keys", "", "comma-separated list of keys to process (default: all)")
	output := fs.String("output", "", "output file (default: stdout)")
	fmt_ := fs.String("format", "dotenv", "output format: dotenv | json")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("encode: input file required")
	}

	f, err := parser.Parse(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("encode: parse %s: %w", fs.Arg(0), err)
	}

	var keysOnly []string
	if *keys != "" {
		for _, k := range strings.Split(*keys, ",") {
			if k = strings.TrimSpace(k); k != "" {
				keysOnly = append(keysOnly, k)
			}
		}
	}

	out, err := encoder.Apply(f, encoder.Options{
		Strategy: encoder.Strategy(*strategy),
		Decode:   *decode,
		KeysOnly: keysOnly,
	})
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	exp, err := exporter.New(*fmt_)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	return exp.Write(out, *output)
}
