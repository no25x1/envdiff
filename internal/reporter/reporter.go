package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter writes diff results to an output writer.
type Reporter struct {
	w      io.Writer
	format Format
}

// New creates a new Reporter writing to w in the given format.
func New(w io.Writer, format Format) *Reporter {
	return &Reporter{w: w, format: format}
}

// Write outputs the diff results according to the reporter's format.
func (r *Reporter) Write(results []diff.Result) error {
	switch r.format {
	case FormatJSON:
		return r.writeJSON(results)
	default:
		return r.writeText(results)
	}
}

func (r *Reporter) writeText(results []diff.Result) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(r.w, "No differences found.")
		return err
	}
	for _, res := range results {
		var line string
		switch res.Status {
		case diff.Added:
			line = fmt.Sprintf("+ %s=%s", res.Key, res.NewValue)
		case diff.Removed:
			line = fmt.Sprintf("- %s=%s", res.Key, res.OldValue)
		case diff.Modified:
			line = fmt.Sprintf("~ %s: %s -> %s", res.Key, res.OldValue, res.NewValue)
		case diff.Unchanged:
			line = fmt.Sprintf("  %s=%s", res.Key, res.NewValue)
		}
		if _, err := fmt.Fprintln(r.w, line); err != nil {
			return err
		}
	}
	return nil
}

func (r *Reporter) writeJSON(results []diff.Result) error {
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, res := range results {
		sb.WriteString(fmt.Sprintf(
			"  {\"key\": %q, \"status\": %q, \"old_value\": %q, \"new_value\": %q}",
			res.Key, res.Status, res.OldValue, res.NewValue,
		))
		if i < len(results)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("]\n")
	_, err := fmt.Fprint(r.w, sb.String())
	return err
}

// Summary prints a short summary line of the diff results.
func (r *Reporter) Summary(results []diff.Result) error {
	added, removed, modified, unchanged := 0, 0, 0, 0
	for _, res := range results {
		switch res.Status {
		case diff.Added:
			added++
		case diff.Removed:
			removed++
		case diff.Modified:
			modified++
		case diff.Unchanged:
			unchanged++
		}
	}
	_, err := fmt.Fprintf(r.w, "Summary: +%d added, -%d removed, ~%d modified, %d unchanged\n",
		added, removed, modified, unchanged)
	return err
}
