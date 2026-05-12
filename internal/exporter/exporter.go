package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Format represents the output format for exporting env files.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatJSON   Format = "json"
	FormatShell  Format = "shell"
)

// Exporter writes env entries to a writer in a given format.
type Exporter struct {
	format Format
}

// New creates a new Exporter for the given format.
func New(format Format) (*Exporter, error) {
	switch format {
	case FormatDotenv, FormatJSON, FormatShell:
		return &Exporter{format: format}, nil
	default:
		return nil, fmt.Errorf("unsupported export format: %q", format)
	}
}

// Write exports the given env entries to w.
func (e *Exporter) Write(w io.Writer, entries []parser.Entry) error {
	switch e.format {
	case FormatDotenv:
		return writeDotenv(w, entries)
	case FormatJSON:
		return writeJSON(w, entries)
	case FormatShell:
		return writeShell(w, entries)
	}
	return nil
}

func sortedEntries(entries []parser.Entry) []parser.Entry {
	sorted := make([]parser.Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})
	return sorted
}

func writeDotenv(w io.Writer, entries []parser.Entry) error {
	for _, e := range sortedEntries(entries) {
		val := e.Value
		if strings.ContainsAny(val, " \t\n#") {
			val = fmt.Sprintf("%q", val)
		}
		if _, err := fmt.Fprintf(w, "%s=%s\n", e.Key, val); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, entries []parser.Entry) error {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(m)
}

func writeShell(w io.Writer, entries []parser.Entry) error {
	for _, e := range sortedEntries(entries) {
		escaped := strings.ReplaceAll(e.Value, "'", "'\\''")
		if _, err := fmt.Fprintf(w, "export %s='%s'\n", e.Key, escaped); err != nil {
			return err
		}
	}
	return nil
}
