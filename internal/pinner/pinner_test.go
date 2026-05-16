package pinner_test

import (
	"testing"

	"envdiff/internal/parser"
	"envdiff/internal/pinner"
)

func makeFile(pairs ...string) *parser.EnvFile {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return &parser.EnvFile{Entries: entries}
}

func TestApply_NilFile_ReturnsError(t *testing.T) {
	_, _, err := pinner.Apply(nil, pinner.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for nil file")
	}
}

func TestApply_PinsAllKeys(t *testing.T) {
	src := makeFile("HOST", "localhost", "PORT", "5432")
	out, results, err := pinner.Apply(src, pinner.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out.Entries))
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Pinned {
			t.Errorf("expected entry %q to be pinned", r.Key)
		}
	}
}

func TestApply_EmptyValue_UsesPlaceholder(t *testing.T) {
	src := makeFile("SECRET", "")
	opts := pinner.DefaultOptions()
	out, results, err := pinner.Apply(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out.Entries))
	}
	if out.Entries[0].Value != opts.Placeholder {
		t.Errorf("expected placeholder %q, got %q", opts.Placeholder, out.Entries[0].Value)
	}
	if !results[0].Pinned {
		t.Error("expected result to be pinned")
	}
}

func TestApply_EmptyValue_SkippedWhenNoPinEmpty(t *testing.T) {
	src := makeFile("EMPTY", "")
	opts := pinner.DefaultOptions()
	opts.Placeholder = ""
	out, results, err := pinner.Apply(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(out.Entries))
	}
	if !results[0].Skipped {
		t.Error("expected result to be skipped")
	}
}

func TestApply_KeyFilter_RestrictsOutput(t *testing.T) {
	src := makeFile("A", "1", "B", "2", "C", "3")
	opts := pinner.DefaultOptions()
	opts.Keys = []string{"A", "C"}
	out, _, err := pinner.Apply(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out.Entries))
	}
	keys := map[string]bool{}
	for _, e := range out.Entries {
		keys[e.Key] = true
	}
	if !keys["A"] || !keys["C"] {
		t.Error("expected keys A and C in output")
	}
}

func TestApply_OutputIsSorted(t *testing.T) {
	src := makeFile("Z", "last", "A", "first", "M", "mid")
	out, _, err := pinner.Apply(src, pinner.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"A", "M", "Z"}
	for i, e := range out.Entries {
		if e.Key != expected[i] {
			t.Errorf("position %d: expected %q, got %q", i, expected[i], e.Key)
		}
	}
}
