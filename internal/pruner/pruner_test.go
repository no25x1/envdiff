package pruner_test

import (
	"testing"

	"github.com/yourusername/envdiff/internal/parser"
	"github.com/yourusername/envdiff/internal/pruner"
)

func makeFile(entries ...parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(key, value string) parser.Entry {
	return parser.Entry{Key: key, Value: value}
}

func keys(f *parser.EnvFile) []string {
	out := make([]string, 0, len(f.Entries))
	for _, e := range f.Entries {
		out = append(out, e.Key)
	}
	return out
}

func TestApply_NilFile_ReturnsEmpty(t *testing.T) {
	opts := pruner.DefaultOptions()
	out, err := pruner.Apply(nil, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 0 {
		t.Errorf("expected empty file, got %d entries", len(out.Entries))
	}
}

func TestApply_NoOptionsEnabled_ReturnsError(t *testing.T) {
	f := makeFile(entry("FOO", "bar"))
	_, err := pruner.Apply(f, pruner.Options{})
	if err == nil {
		t.Fatal("expected error when no options enabled")
	}
}

func TestApply_RemoveEmpty_DropsEmptyValues(t *testing.T) {
	f := makeFile(
		entry("FOO", "bar"),
		entry("EMPTY", ""),
		entry("BAZ", "qux"),
	)
	out, err := pruner.Apply(f, pruner.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := keys(out)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d: %v", len(got), got)
	}
	for _, k := range got {
		if k == "EMPTY" {
			t.Error("EMPTY key should have been pruned")
		}
	}
}

func TestApply_Allowlist_PreservesEmptyKey(t *testing.T) {
	f := makeFile(
		entry("REQUIRED", ""),
		entry("OPTIONAL", ""),
	)
	opts := pruner.Options{
		RemoveEmpty: true,
		Allowlist:   []string{"REQUIRED"},
	}
	out, err := pruner.Apply(f, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := keys(out)
	if len(got) != 1 || got[0] != "REQUIRED" {
		t.Errorf("expected only REQUIRED, got %v", got)
	}
}

func TestApply_RemoveCommented_DropsPureComments(t *testing.T) {
	f := makeFile(
		parser.Entry{Key: "", Value: "", Comment: "# section header"},
		entry("FOO", "bar"),
	)
	opts := pruner.Options{RemoveCommented: true}
	out, err := pruner.Apply(f, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 1 || out.Entries[0].Key != "FOO" {
		t.Errorf("expected only FOO entry, got %v", out.Entries)
	}
}

func TestApply_OriginalUnmodified(t *testing.T) {
	f := makeFile(entry("A", ""), entry("B", "val"))
	orig := len(f.Entries)
	_, _ = pruner.Apply(f, pruner.DefaultOptions())
	if len(f.Entries) != orig {
		t.Error("Apply must not modify the original file")
	}
}
