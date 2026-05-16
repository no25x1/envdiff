package cli

import (
	"flag"
	"fmt"
	"strings"

	"envdiff/internal/parser"
	"envdiff/internal/pinner"
)

func runPin(args []string) error {
	fs := flag.NewFlagSet("pin", flag.ContinueOnError)
	keys := fs.String("keys", "", "comma-separated list of keys to pin (default: all)")
	placeholder := fs.String("placeholder", "CHANGE_ME", "placeholder for empty values")
	pinEmpty := fs.Bool("pin-empty", false, "include entries with empty values as-is")
	output := fs.String("output", "", "write pinned file to path (default: stdout)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("pin: usage: envdiff pin [flags] <file>")
	}

	src, err := parser.ParseFile(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("pin: %w", err)
	}

	opts := pinner.Options{
		Placeholder: *placeholder,
		PinEmpty:    *pinEmpty,
	}
	if *keys != "" {
		for _, k := range strings.Split(*keys, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				opts.Keys = append(opts.Keys, k)
			}
		}
	}

	pinned, results, err := pinner.Apply(src, opts)
	if err != nil {
		return fmt.Errorf("pin: %w", err)
	}

	skipped := 0
	for _, r := range results {
		if r.Skipped {
			skipped++
		}
	}

	var sb strings.Builder
	for _, e := range pinned.Entries {
		if e.Comment != "" {
			sb.WriteString("# " + e.Comment + "\n")
		}
		sb.WriteString(e.Key + "=" + e.Value + "\n")
	}

	if *output != "" {
		if err := writeFile(*output, sb.String()); err != nil {
			return fmt.Errorf("pin: write output: %w", err)
		}
		fmt.Printf("pinned %d entries (%d skipped) → %s\n", len(pinned.Entries), skipped, *output)
	} else {
		fmt.Print(sb.String())
	}

	return nil
}
