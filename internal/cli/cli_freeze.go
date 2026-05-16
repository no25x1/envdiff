package cli

import (
	"fmt"
	"os"

	"envdiff/internal/freezer"
	"envdiff/internal/parser"
)

const defaultFreezePath = ".env.freeze"

func runFreeze(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff freeze <file> [freeze-path]")
	}
	envPath := args[0]
	freezePath := defaultFreezePath
	if len(args) >= 2 {
		freezePath = args[1]
	}

	f, err := parser.Parse(envPath)
	if err != nil {
		return fmt.Errorf("freeze: parse %s: %w", envPath, err)
	}

	entries := freezer.Freeze(f)
	if err := freezer.Save(freezePath, entries); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Froze %d keys to %s\n", len(entries), freezePath)
	return nil
}

func runFreezeCheck(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff freeze-check <file> [freeze-path]")
	}
	envPath := args[0]
	freezePath := defaultFreezePath
	if len(args) >= 2 {
		freezePath = args[1]
	}

	f, err := parser.Parse(envPath)
	if err != nil {
		return fmt.Errorf("freeze-check: parse %s: %w", envPath, err)
	}
	frozen, err := freezer.Load(freezePath)
	if err != nil {
		return err
	}

	violations := freezer.Check(f, frozen)
	if len(violations) == 0 {
		fmt.Fprintln(os.Stdout, "No frozen keys have changed.")
		return nil
	}
	for _, v := range violations {
		fmt.Fprintf(os.Stdout, "CHANGED  %s\n  expected hash: %s\n  actual hash:   %s\n",
			v.Key, v.Expected, v.Actual)
	}
	return fmt.Errorf("%d frozen key(s) changed", len(violations))
}
