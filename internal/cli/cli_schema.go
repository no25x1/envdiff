package cli

import (
	"fmt"
	"os"

	"github.com/envdiff/envdiff/internal/parser"
	"github.com/envdiff/envdiff/internal/schema"
)

// runSchema validates one or more env files against a JSON schema.
// Usage: envdiff schema <schema.json> <file.env> [file2.env ...]
func runSchema(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff schema <schema.json> <file.env> [...]")
	}
	schemaPath := args[0]
	envPaths := args[1:]

	s, err := schema.LoadFile(schemaPath)
	if err != nil {
		return err
	}

	exitCode := 0
	for _, path := range envPaths {
		f, err := parser.Parse(path)
		if err != nil {
			return fmt.Errorf("schema: parse %q: %w", path, err)
		}
		violations := s.Validate(f)
		if len(violations) == 0 {
			fmt.Fprintf(os.Stdout, "%s: OK\n", path)
			continue
		}
		exitCode = 1
		for _, v := range violations {
			fmt.Fprintf(os.Stdout, "%s: [%s] %s\n", path, v.Key, v.Message)
		}
	}
	if exitCode != 0 {
		return fmt.Errorf("schema validation failed")
	}
	return nil
}
