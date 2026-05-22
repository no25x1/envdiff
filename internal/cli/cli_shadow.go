package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/your-org/envdiff/internal/parser"
	"github.com/your-org/envdiff/internal/shadower"
)

func runShadow(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff shadow <base> <next> [--prefix PREFIX] [--format text|json]")
	}

	baseFile, err := parser.ParseFile(args[0])
	if err != nil {
		return fmt.Errorf("shadow: reading base file: %w", err)
	}
	nextFile, err := parser.ParseFile(args[1])
	if err != nil {
		return fmt.Errorf("shadow: reading next file: %w", err)
	}

	prefix := flagStringDefault(args[2:], "--prefix", "")
	format := flagStringDefault(args[2:], "--format", "text")

	opts := shadower.Options{
		Prefix:      prefix,
		IgnoreEqual: true,
	}

	results, err := shadower.Apply(baseFile, nextFile, opts)
	if err != nil {
		return fmt.Errorf("shadow: %w", err)
	}

	switch format {
	case "json":
		return writeShadowJSON(results)
	default:
		return writeShadowText(results)
	}
}

func writeShadowText(results []shadower.Shadow) error {
	if len(results) == 0 {
		fmt.Fprintln(os.Stdout, "no shadowed keys detected")
		return nil
	}
	for _, s := range results {
		fmt.Fprintf(os.Stdout, "~ %s: %q -> %q\n", s.Key, s.BaseValue, s.NextValue)
	}
	fmt.Fprintf(os.Stdout, "\n%d shadowed key(s)\n", len(results))
	return nil
}

func writeShadowJSON(results []shadower.Shadow) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
