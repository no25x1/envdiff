package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/profiler"
)

func runProfile(args []string, stdout io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff profile <file> [--json]")
	}

	filePath := args[0]
	useJSON := len(args) >= 2 && args[1] == "--json"

	f, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("profile: cannot parse %s: %w", filePath, err)
	}

	p := profiler.Analyse(f)

	if useJSON {
		return writeProfileJSON(stdout, p)
	}
	return writeProfileText(stdout, p)
}

func writeProfileText(w io.Writer, p profiler.Profile) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Total keys:\t%d\n", p.TotalKeys)
	fmt.Fprintf(tw, "Empty values:\t%d\n", p.EmptyValues)
	fmt.Fprintf(tw, "Sensitive keys:\t%d\n", p.SensitiveKeys)
	fmt.Fprintf(tw, "Commented lines:\t%d\n", p.CommentedLines)
	if len(p.PrefixCounts) > 0 {
		fmt.Fprintf(tw, "\nTop prefixes:\n")
		for _, prefix := range p.TopPrefixes(5) {
			fmt.Fprintf(tw, "  %s:\t%d\n", prefix, p.PrefixCounts[prefix])
		}
	}
	return tw.Flush()
}

func writeProfileJSON(w io.Writer, p profiler.Profile) error {
	type output struct {
		TotalKeys      int            `json:"total_keys"`
		EmptyValues    int            `json:"empty_values"`
		SensitiveKeys  int            `json:"sensitive_keys"`
		CommentedLines int            `json:"commented_lines"`
		PrefixCounts   map[string]int `json:"prefix_counts"`
	}
	out := output{
		TotalKeys:      p.TotalKeys,
		EmptyValues:    p.EmptyValues,
		SensitiveKeys:  p.SensitiveKeys,
		CommentedLines: p.CommentedLines,
		PrefixCounts:   p.PrefixCounts,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func init() {
	_ = os.Stderr // ensure os import used
}
