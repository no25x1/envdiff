package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/parser"
)

func runLineDiff(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("linediff requires <base> <target> [--format=text|json]")
	}

	baseFile := args[0]
	targetFile := args[1]
	format := "text"
	for _, a := range args[2:] {
		if a == "--format=json" {
			format = "json"
		} else if a == "--format=text" {
			format = "text"
		}
	}

	base, err := parser.Parse(baseFile)
	if err != nil {
		return fmt.Errorf("reading base file: %w", err)
	}
	target, err := parser.Parse(targetFile)
	if err != nil {
		return fmt.Errorf("reading target file: %w", err)
	}

	result := differ.Compare(base, target)

	switch format {
	case "json":
		return writeLineDiffJSON(result)
	default:
		return writeLineDiffText(result)
	}
}

func writeLineDiffText(r *differ.Result) error {
	if !r.HasChanges() {
		fmt.Println("No differences found.")
		return nil
	}
	fmt.Printf("Summary: %s\n", r.Summary())
	fmt.Print(differ.Format(r))
	return nil
}

func writeLineDiffJSON(r *differ.Result) error {
	type jsonLine struct {
		Key    string `json:"key"`
		Before string `json:"before"`
		After  string `json:"after"`
		Change string `json:"change"`
	}
	var lines []jsonLine
	for _, l := range r.Lines {
		lines = append(lines, jsonLine{
			Key:    l.Key,
			Before: l.Before,
			After:  l.After,
			Change: string(l.Change),
		})
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(map[string]interface{}{
		"summary": r.Summary(),
		"lines":   lines,
	})
}
