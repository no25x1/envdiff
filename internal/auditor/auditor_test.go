package auditor_test

import (
	"testing"
	"time"

	"github.com/envdiff/envdiff/internal/auditor"
	"github.com/envdiff/envdiff/internal/diff"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.Added, BaseValue: "", CompareValue: "localhost"},
		{Key: "DB_PASS", Status: diff.Removed, BaseValue: "secret", CompareValue: ""},
		{Key: "APP_ENV", Status: diff.Modified, BaseValue: "dev", CompareValue: "prod"},
		{Key: "APP_NAME", Status: diff.Unchanged, BaseValue: "myapp", CompareValue: "myapp"},
	}
}

func TestFromDiff_EntryCount(t *testing.T) {
	results := makeResults()
	log := auditor.FromDiff(results, auditor.Options{Actor: "ci", File: ".env"})
	if len(log.Entries) != len(results) {
		t.Fatalf("expected %d entries, got %d", len(results), len(log.Entries))
	}
}

func TestFromDiff_FieldsPopulated(t *testing.T) {
	ts := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	results := makeResults()
	log := auditor.FromDiff(results, auditor.Options{Actor: "alice", File: ".env.prod", Timestamp: ts})

	e := log.Entries[0]
	if e.Actor != "alice" {
		t.Errorf("expected actor alice, got %s", e.Actor)
	}
	if e.File != ".env.prod" {
		t.Errorf("expected file .env.prod, got %s", e.File)
	}
	if !e.Timestamp.Equal(ts) {
		t.Errorf("unexpected timestamp %v", e.Timestamp)
	}
	if e.Action != "added" {
		t.Errorf("expected action added, got %s", e.Action)
	}
	if e.NewValue != "localhost" {
		t.Errorf("expected new value localhost, got %s", e.NewValue)
	}
}

func TestFromDiff_DefaultTimestamp(t *testing.T) {
	before := time.Now().UTC()
	log := auditor.FromDiff(makeResults(), auditor.Options{})
	after := time.Now().UTC()

	for _, e := range log.Entries {
		if e.Timestamp.Before(before) || e.Timestamp.After(after) {
			t.Errorf("timestamp %v out of expected range", e.Timestamp)
		}
	}
}

func TestSummary_Counts(t *testing.T) {
	log := auditor.FromDiff(makeResults(), auditor.Options{})
	summary := log.Summary()
	expected := "audit: added=1 removed=1 modified=1 unchanged=1"
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}

func TestFromDiff_EmptyResults(t *testing.T) {
	log := auditor.FromDiff(nil, auditor.Options{Actor: "bot"})
	if len(log.Entries) != 0 {
		t.Errorf("expected 0 entries for nil results")
	}
	if s := log.Summary(); s != "audit: added=0 removed=0 modified=0 unchanged=0" {
		t.Errorf("unexpected summary: %s", s)
	}
}
