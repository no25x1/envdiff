package trimmer_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/parser"
	"github.com/your-org/envdiff/internal/trimmer"
)

func makeFile(entries []parser.Entry) parser.EnvFile {
	return parser.EnvFile{Path: "test.env", Entries: entries}
}

func entry(key, value string) parser.Entry {
	return parser.Entry{Key: key, Value: value}
}

func TestApply_NoOptions_ReturnsCopy(t *testing.T) {
	f := makeFile([]parser.Entry{entry("KEY", "value")})
	out := trimmer.Apply(f, trimmer.Options{})
	if len(out.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out.Entries))
	}
	if out.Entries[0].Key != "KEY" || out.Entries[0].Value != "value" {
		t.Errorf("unexpected entry: %+v", out.Entries[0])
	}
}

func TestApply_TrimValues(t *testing.T) {
	f := makeFile([]parser.Entry{entry("KEY", "  hello world  ")})
	opts := trimmer.DefaultOptions()
	out := trimmer.Apply(f, opts)
	if out.Entries[0].Value != "hello world" {
		t.Errorf("expected trimmed value, got %q", out.Entries[0].Value)
	}
}

func TestApply_TrimKeys(t *testing.T) {
	f := makeFile([]parser.Entry{entry("  KEY  ", "value")})
	opts := trimmer.DefaultOptions()
	out := trimmer.Apply(f, opts)
	if out.Entries[0].Key != "KEY" {
		t.Errorf("expected trimmed key, got %q", out.Entries[0].Key)
	}
}

func TestApply_RemoveEmptyValues_DropsEntries(t *testing.T) {
	f := makeFile([]parser.Entry{
		entry("PRESENT", "yes"),
		entry("EMPTY", ""),
		entry("WHITESPACE", "   "),
	})
	opts := trimmer.Options{
		TrimValues:        true,
		TrimKeys:          true,
		RemoveEmptyValues: true,
	}
	out := trimmer.Apply(f, opts)
	if len(out.Entries) != 1 {
		t.Fatalf("expected 1 entry after drop, got %d", len(out.Entries))
	}
	if out.Entries[0].Key != "PRESENT" {
		t.Errorf("unexpected surviving key: %q", out.Entries[0].Key)
	}
}

func TestApply_RemoveEmptyValues_False_KeepsEntries(t *testing.T) {
	f := makeFile([]parser.Entry{
		entry("EMPTY", ""),
	})
	opts := trimmer.DefaultOptions() // RemoveEmptyValues is false by default
	out := trimmer.Apply(f, opts)
	if len(out.Entries) != 1 {
		t.Fatalf("expected entry to be kept, got %d entries", len(out.Entries))
	}
}

func TestApply_PreservesPath(t *testing.T) {
	f := makeFile(nil)
	f.Path = "/some/path/.env"
	out := trimmer.Apply(f, trimmer.DefaultOptions())
	if out.Path != f.Path {
		t.Errorf("path not preserved: got %q", out.Path)
	}
}
