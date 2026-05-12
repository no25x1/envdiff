package cli

import (
	"fmt"
	"os"

	"github.com/envdiff/internal/config"
	"github.com/envdiff/internal/diff"
	"github.com/envdiff/internal/masker"
	"github.com/envdiff/internal/parser"
	"github.com/envdiff/internal/reconcile"
	"github.com/envdiff/internal/reporter"
	"github.com/envdiff/internal/validator"
)

// Run is the entry point for the CLI.
func Run(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, UsageText)
		return 1
	}

	cmd := args[0]
	remaining := args[1:]

	switch cmd {
	case "diff":
		return runDiff(remaining)
	case "reconcile":
		return runReconcile(remaining)
	case "validate":
		return runValidate(remaining)
	case "version":
		fmt.Println(VersionString)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		return 1
	}
}

func runDiff(args []string) int {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "diff requires <base> <compare> arguments")
		return 1
	}
	cfg := loadConfig()
	base, err := parser.Parse(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading base file: %v\n", err)
		return 1
	}
	compare, err := parser.Parse(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading compare file: %v\n", err)
		return 1
	}
	m := masker.New(cfg.SensitivePatterns...)
	results := diff.Diff(base, compare)
	format := "text"
	if len(args) > 2 {
		format = args[2]
	}
	r := reporter.New(os.Stdout, format, m)
	if err := r.Write(results); err != nil {
		fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
		return 1
	}
	return 0
}

func runReconcile(args []string) int {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "reconcile requires <base> <target> arguments")
		return 1
	}
	base, err := parser.Parse(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading base file: %v\n", err)
		return 1
	}
	target, err := parser.Parse(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading target file: %v\n", err)
		return 1
	}
	plan := reconcile.Plan(base, target)
	result := reconcile.Apply(target.ToMap(), plan)
	fmt.Print(reconcile.Render(result))
	return 0
}

func runValidate(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "validate requires <file> argument")
		return 1
	}
	f, err := parser.Parse(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
		return 1
	}
	v := validator.New()
	results := v.Validate(f)
	if len(results) == 0 {
		fmt.Println("validation passed: no issues found")
		return 0
	}
	for _, r := range results {
		fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", r.Rule, r.Key, r.Message)
	}
	fmt.Fprintf(os.Stderr, "%d issue(s) found\n", len(results))
	return 2
}

func loadConfig() *config.Config {
	cfg, err := config.LoadFile(".envdiff.yaml")
	if err != nil {
		return config.Default()
	}
	return cfg
}
