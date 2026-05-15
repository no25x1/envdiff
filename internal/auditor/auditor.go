// Package auditor provides change-audit logging for .env file operations.
// It records who changed what key, when, and from which value to which value.
package auditor

import (
	"fmt"
	"time"

	"github.com/envdiff/envdiff/internal/diff"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Actor     string    `json:"actor"`
	File      string    `json:"file"`
	Key       string    `json:"key"`
	Action    string    `json:"action"` // added, removed, modified, unchanged
	OldValue  string    `json:"old_value,omitempty"`
	NewValue  string    `json:"new_value,omitempty"`
}

// Log holds a collection of audit entries.
type Log struct {
	Entries []Entry
}

// Options configures audit log generation.
type Options struct {
	Actor     string
	File      string
	Timestamp time.Time // zero value uses time.Now()
}

// FromDiff builds an audit Log from a slice of diff.Result records.
func FromDiff(results []diff.Result, opts Options) *Log {
	ts := opts.Timestamp
	if ts.IsZero() {
		ts = time.Now().UTC()
	}

	log := &Log{}
	for _, r := range results {
		e := Entry{
			Timestamp: ts,
			Actor:     opts.Actor,
			File:      opts.File,
			Key:       r.Key,
			Action:    string(r.Status),
			OldValue:  r.BaseValue,
			NewValue:  r.CompareValue,
		}
		log.Entries = append(log.Entries, e)
	}
	return log
}

// Summary returns a human-readable summary line for the audit log.
func (l *Log) Summary() string {
	counts := map[string]int{}
	for _, e := range l.Entries {
		counts[e.Action]++
	}
	return fmt.Sprintf("audit: added=%d removed=%d modified=%d unchanged=%d",
		counts["added"], counts["removed"], counts["modified"], counts["unchanged"])
}
