package auditor

import (
	"encoding/json"
	"fmt"
	"io"
)

// Format controls the output format of the audit log.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Write serialises the audit Log to w in the requested format.
func (l *Log) Write(w io.Writer, format Format) error {
	switch format {
	case FormatJSON:
		return l.writeJSON(w)
	case FormatText:
		return l.writeText(w)
	default:
		return fmt.Errorf("auditor: unsupported format %q", format)
	}
}

func (l *Log) writeJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(l.Entries)
}

func (l *Log) writeText(w io.Writer) error {
	for _, e := range l.Entries {
		line := fmt.Sprintf("%s\t%s\t%s\t%s\t%s -> %s\n",
			e.Timestamp.Format("2006-01-02T15:04:05Z"),
			e.Actor,
			e.File,
			e.Key,
			quoteVal(e.OldValue),
			quoteVal(e.NewValue),
		)
		if _, err := fmt.Fprint(w, line); err != nil {
			return err
		}
	}
	return nil
}

func quoteVal(v string) string {
	if v == "" {
		return "(empty)"
	}
	return v
}
