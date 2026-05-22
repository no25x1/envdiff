package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/envdiff/envdiff/internal/grapher"
	"github.com/envdiff/envdiff/internal/parser"
)

func runGraph(args []string, stdout io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff graph <file> [--format text|json]")
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

	envFile, err := parser.Parse(f)
	if err != nil {
		return fmt.Errorf("parse %s: %w", filePath, err)
	}

	g, err := grapher.Build(envFile)
	if err != nil {
		return fmt.Errorf("graph: %w", err)
	}

	switch format {
	case "json":
		return writeGraphJSON(g, stdout)
	default:
		return writeGraphText(g, stdout)
	}
}

func writeGraphText(g *grapher.Graph, w io.Writer) error {
	if len(g.Edges) == 0 {
		fmt.Fprintln(w, "No variable references found.")
		return nil
	}
	fmt.Fprintln(w, "Dependencies:")
	for _, e := range g.Edges {
		fmt.Fprintf(w, "  %s -> %s\n", e.Key, e.Dep)
	}
	fmt.Fprintln(w, "\nEvaluation order:")
	for i, k := range g.Order {
		fmt.Fprintf(w, "  %d. %s\n", i+1, k)
	}
	return nil
}

func writeGraphJSON(g *grapher.Graph, w io.Writer) error {
	out := struct {
		Edges []grapher.Edge `json:"edges"`
		Order []string       `json:"order"`
	}{
		Edges: g.Edges,
		Order: g.Order,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
