package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"envdiff/internal/diff"
	"envdiff/internal/masker"
	"envdiff/internal/parser"
	"envdiff/internal/reporter"
	"envdiff/internal/watcher"
)

const defaultWatchInterval = 2 * time.Second

// runWatch starts a polling loop that re-diffs env files on change.
func runWatch(args []string, cfg *appConfig) error {
	if len(args) < 2 {
		return fmt.Errorf("watch requires <base> <compare> arguments")
	}
	base, compare := args[0], args[1]

	w, err := watcher.New([]string{base, compare}, defaultWatchInterval)
	if err != nil {
		return fmt.Errorf("watch: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Watching %s and %s for changes (Ctrl-C to stop)\n", base, compare)
	w.Start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case ev, ok := <-w.Events:
			if !ok {
				return nil
			}
			fmt.Fprintf(os.Stderr, "\n[%s] change detected in %s\n", ev.At.Format("15:04:05"), ev.Path)
			if err := diffAndReport(base, compare, cfg); err != nil {
				fmt.Fprintf(os.Stderr, "diff error: %v\n", err)
			}
		case <-sigs:
			w.Stop()
			fmt.Fprintln(os.Stderr, "\nwatcher stopped")
			return nil
		}
	}
}

func diffAndReport(base, compare string, cfg *appConfig) error {
	bf, err := parser.Parse(base)
	if err != nil {
		return err
	}
	cf, err := parser.Parse(compare)
	if err != nil {
		return err
	}

	msk := masker.New(cfg.SensitivePatterns)
	bf = msk.MaskEnv(bf)
	cf = msk.MaskEnv(cf)

	results := diff.Diff(bf, cf)
	rep := reporter.New(os.Stdout)
	return rep.WriteText(results)
}

// appConfig is a minimal config holder used by watch/diff helpers.
type appConfig struct {
	SensitivePatterns []string
}
