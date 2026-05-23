package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"envdiff/internal/masked"
	"envdiff/internal/parser"
)

func runMaskedDiff(args []string, stdout io.Writer) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff maskeddiff <base> <next> [--format=text|json] [--extra-patterns=p1,p2]")
	}

	baseFile, err := parser.Parse(args[0])
	if err != nil {
		return fmt.Errorf("maskeddiff: loading base: %w", err)
	}
	nextFile, err := parser.Parse(args[1])
	if err != nil {
		return fmt.Errorf("maskeddiff: loading next: %w", err)
	}

	format := "text"
	var extraPatterns []string
	for _, a := range args[2:] {
		if strings.HasPrefix(a, "--format=") {
			format = strings.TrimPrefix(a, "--format=")
		} else if strings.HasPrefix(a, "--extra-patterns=") {
			raw := strings.TrimPrefix(a, "--extra-patterns=")
			for _, p := range strings.Split(raw, ",") {
				if p != "" {
					extraPatterns = append(extraPatterns, p)
				}
			}
		}
	}

	results, err := masked.Compare(baseFile, nextFile, masked.Options{ExtraPatterns: extraPatterns})
	if err != nil {
		return fmt.Errorf("maskeddiff: %w", err)
	}

	switch format {
	case "json":
		return writeMaskedDiffJSON(stdout, results)
	default:
		return writeMaskedDiffText(stdout, results)
	}
}

func writeMaskedDiffText(w io.Writer, results []masked.Result) error {
	if len(results) == 0 {
		fmt.Fprintln(w, "No differences.")
		return nil
	}
	for _, r := range results {
		fmt.Fprintln(w, masked.Format(r))
	}
	return nil
}

func writeMaskedDiffJSON(w io.Writer, results []masked.Result) error {
	type row struct {
		Key    string `json:"key"`
		Status string `json:"status"`
		Base   string `json:"base_value"`
		Next   string `json:"next_value"`
		Masked bool   `json:"masked"`
	}
	rows := make([]row, 0, len(results))
	for _, r := range results {
		rows = append(rows, row{
			Key:    r.Key,
			Status: r.Status.String(),
			Base:   r.BaseVal,
			Next:   r.NextVal,
			Masked: r.Masked,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}

func init() {
	_ = os.Stderr // ensure os imported
}
