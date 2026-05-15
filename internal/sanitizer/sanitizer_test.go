package sanitizer_test

import (
	"testing"

	"github.com/envdiff/envdiff/internal/parser"
	"github.com/envdiff/envdiff/internal/sanitizer"
)

func makeFile(entries ...parser.Entry) parser.EnvFile {
	return parser.EnvFile{Path: "test.env", Entries: entries}
}

func entry(k, v string) parser.Entry {
	return parser.Entry{Key: k, Value: v, Raw: k + "=" + v}
}

func TestApply_NoOptions_ReturnsCopy(t *testing.T) {
	f := makeFile(entry("FOO", "bar"), entry("BAZ", "qux"))
	out := sanitizer.Apply(f, sanitizer.Options{})
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out.Entries))
	}
	if out.Entries[0].Value != "bar" {
		t.Errorf("unexpected value: %q", out.Entries[0].Value)
	}
}

func TestApply_TrimWhitespace(t *testing.T) {
	f := makeFile(entry("KEY", "  hello  "))
	out := sanitizer.Apply(f, sanitizer.Options{TrimWhitespace: true})
	if out.Entries[0].Value != "hello" {
		t.Errorf("expected trimmed value, got %q", out.Entries[0].Value)
	}
}

func TestApply_RemoveNullBytes(t *testing.T) {
	f := makeFile(entry("KEY", "hel\x00lo"))
	out := sanitizer.Apply(f, sanitizer.Options{RemoveNullBytes: true})
	if out.Entries[0].Value != "hello" {
		t.Errorf("expected null bytes removed, got %q", out.Entries[0].Value)
	}
}

func TestApply_NormaliseLineEndings(t *testing.T) {
	f := makeFile(entry("KEY", "line1\r\nline2\rline3"))
	out := sanitizer.Apply(f, sanitizer.Options{NormaliseLineEndings: true})
	want := "line1\nline2\nline3"
	if out.Entries[0].Value != want {
		t.Errorf("expected %q, got %q", want, out.Entries[0].Value)
	}
}

func TestApply_CollapseBlankValues(t *testing.T) {
	f := makeFile(entry("KEY", "   "))
	out := sanitizer.Apply(f, sanitizer.Options{CollapseBlankValues: true})
	if out.Entries[0].Value != "" {
		t.Errorf("expected empty string, got %q", out.Entries[0].Value)
	}
}

func TestApply_DefaultOptions_AllPasses(t *testing.T) {
	f := makeFile(
		entry("A", "  val\r\n  "),
		entry("B", "\x00data\x00"),
		entry("C", "  \t  "),
	)
	out := sanitizer.Apply(f, sanitizer.DefaultOptions())
	if out.Entries[0].Value != "val" {
		t.Errorf("A: expected 'val', got %q", out.Entries[0].Value)
	}
	if out.Entries[1].Value != "data" {
		t.Errorf("B: expected 'data', got %q", out.Entries[1].Value)
	}
	if out.Entries[2].Value != "" {
		t.Errorf("C: expected empty string, got %q", out.Entries[2].Value)
	}
}

func TestApply_PreservesPath(t *testing.T) {
	f := makeFile(entry("K", "v"))
	f.Path = "/etc/app/.env"
	out := sanitizer.Apply(f, sanitizer.DefaultOptions())
	if out.Path != f.Path {
		t.Errorf("path not preserved: got %q", out.Path)
	}
}
