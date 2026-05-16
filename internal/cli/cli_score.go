package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/scorer"
)

func runScore(args []string, stdout io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff score <file> [--format text|json]")
	}

	filePath := args[0]
	format := "text"
	for i, a := range args {
		if a == "--format" && i+1 < len(args) {
			format = args[i+1]
		}
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open %s: %w", filePath, err)
	}
	defer f.Close()

	env, err := parser.Parse(f)
	if err != nil {
		return fmt.Errorf("parse %s: %w", filePath, err)
	}

	result := scorer.Score(env)

	switch format {
	case "json":
		return writeScoreJSON(stdout, result)
	default:
		return writeScoreText(stdout, result)
	}
}

func writeScoreText(w io.Writer, r scorer.Result) error {
	fmt.Fprintf(w, "Score: %d / %d\n", r.Score, r.Max)
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Breakdown:")
	for _, br := range r.Breakdown {
		status := "✓"
		if br.Failed > 0 {
			status = "✗"
		}
		fmt.Fprintf(w, "  %s %-30s passed=%d failed=%d weight=%d\n",
			status, br.Rule, br.Passed, br.Failed, br.Weight)
		if br.Comment != "" {
			fmt.Fprintf(w, "    %s\n", br.Comment)
		}
	}
	return nil
}

func writeScoreJSON(w io.Writer, r scorer.Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
