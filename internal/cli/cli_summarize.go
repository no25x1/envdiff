package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"envdiff/internal/parser"
	"envdiff/internal/summarizer"
)

func runSummarize(args []string, stdout io.Writer) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff summarize <file> [--json]")
	}

	filePath := args[0]
	useJSON := len(args) >= 2 && args[1] == "--json"

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open %s: %w", filePath, err)
	}
	defer f.Close()

	env, err := parser.Parse(f)
	if err != nil {
		return fmt.Errorf("parse %s: %w", filePath, err)
	}

	s := summarizer.Analyse(env)

	if useJSON {
		return writeSummarizeJSON(stdout, s)
	}
	return writeSummarizeText(stdout, s)
}

func writeSummarizeText(w io.Writer, s summarizer.Summary) error {
	fmt.Fprintf(w, "Total keys      : %d\n", s.TotalKeys)
	fmt.Fprintf(w, "Empty values    : %d\n", s.EmptyValues)
	fmt.Fprintf(w, "Sensitive keys  : %d\n", s.SensitiveKeys)
	fmt.Fprintf(w, "Unique prefixes : %d\n", s.UniquePrefix)
	fmt.Fprintf(w, "Longest key     : %s\n", s.LongestKey)
	fmt.Fprintf(w, "Shortest key    : %s\n", s.ShortestKey)
	if len(s.TopPrefixes) > 0 {
		fmt.Fprintln(w, "Top prefixes:")
		for _, p := range s.TopPrefixes {
			fmt.Fprintf(w, "  %-20s %d\n", p.Prefix, p.Count)
		}
	}
	return nil
}

func writeSummarizeJSON(w io.Writer, s summarizer.Summary) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}
