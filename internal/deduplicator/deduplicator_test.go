package deduplicator_test

import (
	"testing"

	"github.com/user/envdiff/internal/deduplicator"
	"github.com/user/envdiff/internal/parser"
)

func makeFile(entries ...parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Path: "test.env", Entries: entries}
}

func entry(k, v string) parser.Entry { return parser.Entry{Key: k, Value: v} }

func TestApply_NoDuplicates_ReturnsCopy(t *testing.T) {
	src := makeFile(entry("A", "1"), entry("B", "2"))
	res := deduplicator.Apply(src, deduplicator.DefaultOptions())
	if len(res.File.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.File.Entries))
	}
	if len(res.Removed) != 0 {
		t.Fatalf("expected no duplicates reported, got %d", len(res.Removed))
	}
}

func TestApply_KeepLast_RemovesEarlier(t *testing.T) {
	src := makeFile(entry("A", "first"), entry("B", "only"), entry("A", "last"))
	res := deduplicator.Apply(src, deduplicator.Options{Strategy: deduplicator.KeepLast})
	if len(res.File.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.File.Entries))
	}
	for _, e := range res.File.Entries {
		if e.Key == "A" && e.Value != "last" {
			t.Errorf("expected last value for A, got %q", e.Value)
		}
	}
	if len(res.Removed) != 1 || res.Removed[0].Key != "A" {
		t.Errorf("expected duplicate report for A")
	}
}

func TestApply_KeepFirst_RemovesLater(t *testing.T) {
	src := makeFile(entry("X", "first"), entry("X", "second"), entry("Y", "only"))
	res := deduplicator.Apply(src, deduplicator.Options{Strategy: deduplicator.KeepFirst})
	if len(res.File.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(res.File.Entries))
	}
	for _, e := range res.File.Entries {
		if e.Key == "X" && e.Value != "first" {
			t.Errorf("expected first value for X, got %q", e.Value)
		}
	}
}

func TestApply_OrderPreserved_KeepFirst(t *testing.T) {
	src := makeFile(entry("C", "1"), entry("A", "2"), entry("C", "3"), entry("B", "4"))
	res := deduplicator.Apply(src, deduplicator.Options{Strategy: deduplicator.KeepFirst})
	keys := []string{}
	for _, e := range res.File.Entries {
		keys = append(keys, e.Key)
	}
	expected := []string{"C", "A", "B"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("position %d: expected %q, got %q", i, k, keys[i])
		}
	}
}

func TestApply_DuplicateCount(t *testing.T) {
	src := makeFile(entry("K", "1"), entry("K", "2"), entry("K", "3"))
	res := deduplicator.Apply(src, deduplicator.DefaultOptions())
	if len(res.Removed) != 1 {
		t.Fatalf("expected 1 duplicate entry, got %d", len(res.Removed))
	}
	if res.Removed[0].Count != 3 {
		t.Errorf("expected count 3, got %d", res.Removed[0].Count)
	}
}
