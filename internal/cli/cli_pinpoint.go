package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"envdiff/internal/parser"
	"envdiff/internal/pinpointer"
)

func runPinpoint(args []string, stdout io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff pinpoint <file> [--format=text|json]")
	}

	filePath := args[0]
	format := "text"
	for _, a := range args[1:] {
		if len(a) > 9 && a[:9] == "--format=" {
			format = a[9:]
		}
	}

	f, err := parser.Parse(filePath)
	if err != nil {
		return fmt.Errorf("parse %s: %w", filePath, err)
	}

	findings, err := pinpointer.Apply(f)
	if err != nil {
		return fmt.Errorf("pinpoint: %w", err)
	}

	switch format {
	case "json":
		return writePinpointJSON(stdout, findings)
	default:
		return writePinpointText(stdout, findings)
	}
}

func writePinpointText(w io.Writer, findings []pinpointer.Finding) error {
	if len(findings) == 0 {
		fmt.Fprintln(w, "no environment-pinned values detected")
		return nil
	}
	for _, f := range findings {
		fmt.Fprintf(w, "%-30s  %s\n", f.Key, f.Reason)
	}
	fmt.Fprintf(w, "\n%d finding(s)\n", len(findings))
	return nil
}

func writePinpointJSON(w io.Writer, findings []pinpointer.Finding) error {
	type row struct {
		Key    string `json:"key"`
		Value  string `json:"value"`
		Reason string `json:"reason"`
	}
	rows := make([]row, len(findings))
	for i, f := range findings {
		rows[i] = row{Key: f.Key, Value: f.Value, Reason: f.Reason}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}

func init() {
	_ = os.Stderr // ensure os import used
}
