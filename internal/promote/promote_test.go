package promote_test

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/promote"
)

func makeFile(path string, pairs ...string) parser.EnvFile {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return parser.EnvFile{Path: path, Entries: entries}
}

func TestApply_PromotesAllKeys(t *testing.T) {
	src := makeFile("staging.env", "FOO", "bar", "BAZ", "qux")
	dst := makeFile("prod.env")
	out, results, err := promote.Apply(src, dst, promote.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out.Entries))
	}
	for _, r := range results {
		if r.Action != "promoted" {
			t.Errorf("expected promoted, got %q for key %s", r.Action, r.Key)
		}
	}
}

func TestApply_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := makeFile("staging.env", "FOO", "new")
	dst := makeFile("prod.env", "FOO", "old")
	out, results, err := promote.Apply(src, dst, promote.Options{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[0].Value != "old" {
		t.Errorf("expected old value to be kept, got %q", out.Entries[0].Value)
	}
	if results[0].Action != "skipped_exists" {
		t.Errorf("expected skipped_exists, got %q", results[0].Action)
	}
}

func TestApply_OverwritesExistingWhenEnabled(t *testing.T) {
	src := makeFile("staging.env", "FOO", "new")
	dst := makeFile("prod.env", "FOO", "old")
	out, _, err := promote.Apply(src, dst, promote.Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[0].Value != "new" {
		t.Errorf("expected new value, got %q", out.Entries[0].Value)
	}
}

func TestApply_FiltersByAllowlist(t *testing.T) {
	src := makeFile("staging.env", "FOO", "1", "BAR", "2", "BAZ", "3")
	dst := makeFile("prod.env")
	out, results, err := promote.Apply(src, dst, promote.Options{Keys: []string{"FOO", "BAZ"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out.Entries))
	}
	skipped := 0
	for _, r := range results {
		if r.Action == "skipped_filter" {
			skipped++
		}
	}
	if skipped != 1 {
		t.Errorf("expected 1 skipped_filter, got %d", skipped)
	}
}

func TestApply_EmptySourcePath_ReturnsError(t *testing.T) {
	src := makeFile("", "FOO", "bar")
	dst := makeFile("prod.env")
	_, _, err := promote.Apply(src, dst, promote.Options{})
	if err == nil {
		t.Fatal("expected error for empty source path")
	}
}
