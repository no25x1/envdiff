package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"envdiff/internal/inspector"
	"envdiff/internal/parser"
)

func runInspect(args []string, out io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff inspect <file> [--json]")
	}

	filePath := args[0]
	useJSON := len(args) > 1 && args[1] == "--json"

	f, err := parser.Parse(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", filePath)
		}
		return fmt.Errorf("parse error: %w", err)
	}

	results, err := inspector.Inspect(f)
	if err != nil {
		return fmt.Errorf("inspect error: %w", err)
	}

	if useJSON {
		return writeInspectJSON(results, out)
	}
	return writeInspectText(results, out)
}

func writeInspectText(results []inspector.Result, out io.Writer) error {
	for _, r := range results {
		sensitive := ""
		if r.IsSensitive {
			sensitive = " [sensitive]"
		}
		empty := ""
		if r.IsEmpty {
			empty = " [empty]"
		}
		fmt.Fprintf(out, "%-30s type=%-8s len=%-4d%s%s\n",
			r.Key, r.TypeHint, r.Length, sensitive, empty)
	}
	return nil
}

func writeInspectJSON(results []inspector.Result, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
