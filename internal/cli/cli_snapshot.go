package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/envdiff/internal/diff"
	"github.com/envdiff/internal/parser"
	"github.com/envdiff/internal/reporter"
	"github.com/envdiff/internal/snapshot"
)

// runSnapshot saves a snapshot of the given env file.
// Usage: envdiff snapshot <file> [--out <dest>]
func runSnapshot(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("snapshot: requires <file> argument")
	}

	envPath := args[0]
	env, err := parser.Parse(envPath)
	if err != nil {
		return fmt.Errorf("snapshot: parse %s: %w", envPath, err)
	}

	dest := defaultSnapshotPath(envPath)
	for i, a := range args {
		if a == "--out" && i+1 < len(args) {
			dest = args[i+1]
		}
	}

	if err := snapshot.Save(env, dest); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "snapshot saved: %s\n", dest)
	return nil
}

// runSnapshotDiff diffs the current env file against a saved snapshot.
// Usage: envdiff snapshot-diff <file> <snapshot>
func runSnapshotDiff(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("snapshot-diff: requires <file> and <snapshot> arguments")
	}

	current, err := parser.Parse(args[0])
	if err != nil {
		return fmt.Errorf("snapshot-diff: parse current: %w", err)
	}

	snap, err := snapshot.Load(args[1])
	if err != nil {
		return fmt.Errorf("snapshot-diff: load snapshot: %w", err)
	}

	results := diff.Diff(snap.ToEnvFile(), current)

	rep := reporter.New(os.Stdout)
	if err := rep.WriteText(results); err != nil {
		return fmt.Errorf("snapshot-diff: report: %w", err)
	}

	rep.Summary(results)
	return nil
}

func defaultSnapshotPath(envPath string) string {
	base := filepath.Base(envPath)
	timestamp := time.Now().UTC().Format("20060102T150405Z")
	return fmt.Sprintf("%s.%s.snapshot.json", base, timestamp)
}
