package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/config"
	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/masker"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/reconcile"
	"github.com/user/envdiff/internal/reporter"
)

// Run parses CLI arguments and dispatches to the appropriate subcommand.
func Run(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: envdiff <diff|reconcile> [options]")
	}

	switch args[0] {
	case "diff":
		return runDiff(args[1:])
	case "reconcile":
		return runReconcile(args[1:])
	default:
		return fmt.Errorf("unknown command %q; expected diff or reconcile", args[0])
	}
}

func runDiff(args []string) error {
	fs := flag.NewFlagSet("diff", flag.ContinueOnError)
	format := fs.String("format", "text", "output format: text or json")
	cfgFile := fs.String("config", "", "path to config file")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 2 {
		return errors.New("diff requires two .env file arguments")
	}

	cfg, err := loadConfig(*cfgFile)
	if err != nil {
		return err
	}

	base, err := parser.Parse(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("parsing base file: %w", err)
	}
	target, err := parser.Parse(fs.Arg(1))
	if err != nil {
		return fmt.Errorf("parsing target file: %w", err)
	}

	m := masker.New(cfg.SensitivePatterns, cfg.MaskValue)
	results := diff.Diff(base, target)

	rep := reporter.New(os.Stdout, m)
	switch *format {
	case "json":
		return rep.WriteJSON(results)
	default:
		rep.WriteText(results)
		rep.Summary(results)
		return nil
	}
}

func runReconcile(args []string) error {
	fs := flag.NewFlagSet("reconcile", flag.ContinueOnError)
	cfgFile := fs.String("config", "", "path to config file")
	outFile := fs.String("out", "", "write reconciled output to file (default: stdout)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 2 {
		return errors.New("reconcile requires two .env file arguments")
	}

	_, err := loadConfig(*cfgFile)
	if err != nil {
		return err
	}

	base, err := parser.Parse(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("parsing base file: %w", err)
	}
	target, err := parser.Parse(fs.Arg(1))
	if err != nil {
		return fmt.Errorf("parsing target file: %w", err)
	}

	plan := reconcile.Plan(base, target)
	result := reconcile.Apply(base.Entries, target.Entries, plan)
	output := reconcile.Render(result)

	w := os.Stdout
	if *outFile != "" {
		f, err := os.Create(*outFile)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		w = f
	}
	_, err = fmt.Fprint(w, output)
	return err
}

func loadConfig(path string) (*config.Config, error) {
	if path == "" {
		return config.Default(), nil
	}
	return config.LoadFile(path)
}
