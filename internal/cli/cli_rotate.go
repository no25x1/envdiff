package cli

import (
	"fmt"
	"strings"

	"github.com/envdiff/envdiff/internal/exporter"
	"github.com/envdiff/envdiff/internal/parser"
	"github.com/envdiff/envdiff/internal/rotator"
)

// runRotate implements the `envdiff rotate` sub-command.
//
// Usage: envdiff rotate <file> --map OLD=NEW[,OLD=NEW...] [--fail-on-missing] [--format dotenv|json]
func runRotate(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("rotate: usage: envdiff rotate <file> --map OLD=NEW")
	}

	filePath := args[0]
	rest := args[1:]

	var rawMap string
	failOnMissing := false
	format := "dotenv"

	for i := 0; i < len(rest); i++ {
		switch rest[i] {
		case "--map":
			if i+1 >= len(rest) {
				return fmt.Errorf("rotate: --map requires a value")
			}
			i++
			rawMap = rest[i]
		case "--fail-on-missing":
			failOnMissing = true
		case "--format":
			if i+1 >= len(rest) {
				return fmt.Errorf("rotate: --format requires a value")
			}
			i++
			format = rest[i]
		}
	}

	if rawMap == "" {
		return fmt.Errorf("rotate: --map flag is required")
	}

	mappings, err := parseRotateMappings(rawMap)
	if err != nil {
		return fmt.Errorf("rotate: %w", err)
	}

	src, err := parser.Parse(filePath)
	if err != nil {
		return fmt.Errorf("rotate: %w", err)
	}

	res, err := rotator.Apply(src, rotator.Options{
		Mappings:      mappings,
		FailOnMissing: failOnMissing,
	})
	if err != nil {
		return fmt.Errorf("rotate: %w", err)
	}

	if len(res.Missing) > 0 {
		fmt.Printf("warning: keys not found: %s\n", strings.Join(res.Missing, ", "))
	}

	exp, err := exporter.New(format)
	if err != nil {
		return fmt.Errorf("rotate: %w", err)
	}
	return exp.Write(res.File, defaultWriter)
}

func parseRotateMappings(raw string) ([]rotator.Mapping, error) {
	var out []rotator.Mapping
	for _, pair := range strings.Split(raw, ",") {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid mapping %q: expected OLD=NEW", pair)
		}
		out = append(out, rotator.Mapping{OldKey: parts[0], NewKey: parts[1]})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no valid mappings provided")
	}
	return out, nil
}
