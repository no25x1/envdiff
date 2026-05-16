package freezer_test

import (
	"os"
	"path/filepath"
	"testing"

	"envdiff/internal/freezer"
	"envdiff/internal/parser"
)

func TestSave_CreatesFile(t *testing.T) {
	f := makeFile(entry("KEY", "value"))
	entries := freezer.Freeze(f)

	dir := t.TempDir()
	path := filepath.Join(dir, "freeze.json")

	if err := freezer.Save(path, entries); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestLoad_RoundTrip(t *testing.T) {
	f := makeFile(
		parser.Entry{Key: "A", Value: "1"},
		parser.Entry{Key: "B", Value: "2"},
	)
	original := freezer.Freeze(f)

	dir := t.TempDir()
	path := filepath.Join(dir, "freeze.json")

	if err := freezer.Save(path, original); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := freezer.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded) != len(original) {
		t.Fatalf("length mismatch: want %d, got %d", len(original), len(loaded))
	}
	for i := range original {
		if loaded[i].Key != original[i].Key || loaded[i].Hash != original[i].Hash {
			t.Errorf("entry %d mismatch: %+v vs %+v", i, original[i], loaded[i])
		}
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := freezer.Load("/nonexistent/freeze.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
