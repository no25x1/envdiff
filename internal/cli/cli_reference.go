package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"envdiff/internal/parser"
	"envdiff/internal/referencer"
)

func runReference(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff reference <file> [--json]")
	}

	filePath := args[0]
	useJSON := len(args) >= 2 && args[1] == "--json"

	f, err := parser.Parse(filePath)
	if err != nil {
		return fmt.Errorf("reference: %w", err)
	}

	result, err := referencer.Analyse(f)
	if err != nil {
		return fmt.Errorf("reference: %w", err)
	}

	if useJSON {
		return writeReferenceJSON(result)
	}
	return writeReferenceText(result)
}

func writeReferenceText(r *referencer.Result) error {
	w := os.Stdout

	if len(r.Referrers) == 0 {
		fmt.Fprintln(w, "No cross-key references found.")
	} else {
		fmt.Fprintln(w, "References:")
		keys := make([]string, 0, len(r.Referrers))
		for k := range r.Referrers {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(w, "  %s  <-  %v\n", k, r.Referrers[k])
		}
	}

	if len(r.Undefined) > 0 {
		fmt.Fprintln(w, "\nUndefined references:")
		for _, u := range r.Undefined {
			fmt.Fprintf(w, "  %s\n", u)
		}
	}

	if len(r.Unused) > 0 {
		fmt.Fprintln(w, "\nUnreferenced keys:")
		for _, u := range r.Unused {
			fmt.Fprintf(w, "  %s\n", u)
		}
	}
	return nil
}

func writeReferenceJSON(r *referencer.Result) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
