package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/grouper"
	"github.com/user/envdiff/internal/parser"
)

// runGroup implements the `envdiff group <file>` sub-command.
// Flags:
//
//	--delimiter  string   key prefix delimiter (default "_")
//	--ungrouped          include keys with no delimiter in an unnamed group
func runGroup(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff group <file> [--delimiter=_] [--ungrouped]")
	}

	filePath := args[0]
	delimiter := "_"
	includeUngrouped := false

	for _, arg := range args[1:] {
		if strings.HasPrefix(arg, "--delimiter=") {
			delimiter = strings.TrimPrefix(arg, "--delimiter=")
		} else if arg == "--ungrouped" {
			includeUngrouped = true
		}
	}

	f, err := parser.Parse(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", filePath)
		}
		return fmt.Errorf("parse error: %w", err)
	}

	groups := grouper.Apply(f, grouper.Options{
		Delimiter:        delimiter,
		IncludeUngrouped: includeUngrouped,
	})

	if len(groups) == 0 {
		fmt.Println("no groups found")
		return nil
	}

	for _, g := range groups {
		name := g.Name
		if name == "" {
			name = "(ungrouped)"
		}
		fmt.Printf("[%s] (%d keys)\n", name, len(g.Entries))
		for _, e := range g.Entries {
			fmt.Printf("  %s=%s\n", e.Key, e.Value)
		}
	}
	return nil
}
