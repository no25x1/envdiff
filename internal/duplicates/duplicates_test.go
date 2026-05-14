package duplicates_test

import (
	"testing"

	"github.com/user/envdiff/internal/duplicates"
	"github.com/user/envdiff/internal/parser"
)

func makeFile(path string, pairs ...string) *parser.EnvFile {
	if len(pairs)%2 != 0 {
		panic("pairs must be even")
	}
	var entries []parser.Entry
	for i := 0; i < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{
			Key:   pairs[i],
			Value: pairs[i+1],
			Line:  i/2 + 1,
		})
	}
	return &parser.EnvFile{Path: path, Entries: entries}
}

func TestCheckIntra_NoDuplicates(t *testing.T) {
	f := makeFile("a.env", "FOO", "1", "BAR", "2")
	report := duplicates.CheckIntra(f)
	if report.HasDuplicates() {
		t.Fatalf("expected no duplicates, got %d", len(report.Findings))
	}
	if report.Summary() != "no duplicate keys found" {
		t.Errorf("unexpected summary: %s", report.Summary())
	}
}

func TestCheckIntra_FindsDuplicates(t *testing.T) {
	f := makeFile("a.env", "FOO", "1", "BAR", "2", "FOO", "3")
	report := duplicates.CheckIntra(f)
	if !report.HasDuplicates() {
		t.Fatal("expected duplicates, got none")
	}
	if len(report.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(report.Findings))
	}
	if report.Findings[0].Key != "FOO" {
		t.Errorf("expected key FOO, got %s", report.Findings[0].Key)
	}
	if report.Findings[0].Source != "intra" {
		t.Errorf("expected source intra, got %s", report.Findings[0].Source)
	}
}

func TestCheckIntra_MultiDuplicates(t *testing.T) {
	f := makeFile("a.env", "FOO", "1", "FOO", "2", "FOO", "3")
	report := duplicates.CheckIntra(f)
	if len(report.Findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(report.Findings))
	}
}

func TestCheckInter_NoDuplicates(t *testing.T) {
	a := makeFile("a.env", "FOO", "1")
	b := makeFile("b.env", "BAR", "2")
	report := duplicates.CheckInter([]*parser.EnvFile{a, b})
	if report.HasDuplicates() {
		t.Fatalf("expected no duplicates, got %d", len(report.Findings))
	}
}

func TestCheckInter_FindsDuplicates(t *testing.T) {
	a := makeFile("a.env", "FOO", "1", "BAR", "2")
	b := makeFile("b.env", "FOO", "99", "BAZ", "3")
	report := duplicates.CheckInter([]*parser.EnvFile{a, b})
	if !report.HasDuplicates() {
		t.Fatal("expected duplicates, got none")
	}
	// FOO appears in both files -> 2 findings
	if len(report.Findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(report.Findings))
	}
	for _, f := range report.Findings {
		if f.Key != "FOO" {
			t.Errorf("expected key FOO, got %s", f.Key)
		}
		if f.Source != "inter" {
			t.Errorf("expected source inter, got %s", f.Source)
		}
	}
}

func TestSummary_WithFindings(t *testing.T) {
	f := makeFile("a.env", "X", "1", "X", "2")
	report := duplicates.CheckIntra(f)
	if report.Summary() != "1 duplicate key(s) found" {
		t.Errorf("unexpected summary: %s", report.Summary())
	}
}
