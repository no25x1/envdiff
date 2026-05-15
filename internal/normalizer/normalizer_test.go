package normalizer_test

import (
	"testing"

	"github.com/user/envdiff/internal/normalizer"
	"github.com/user/envdiff/internal/parser"
)

func makeFile(entries ...parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(key, val string) parser.Entry {
	return parser.Entry{Key: key, Value: val}
}

func TestApply_NilFile_ReturnsEmpty(t *testing.T) {
	out := normalizer.Apply(nil, normalizer.DefaultOptions())
	if out == nil {
		t.Fatal("expected non-nil EnvFile")
	}
	if len(out.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(out.Entries))
	}
}

func TestApply_KeyToUpper(t *testing.T) {
	src := makeFile(entry("db_host", "localhost"), entry("api_key", "abc123"))
	out := normalizer.Apply(src, normalizer.Options{KeyToUpper: true})
	for _, e := range out.Entries {
		for _, ch := range e.Key {
			if ch >= 'a' && ch <= 'z' {
				t.Errorf("key %q still has lowercase letters", e.Key)
			}
		}
	}
}

func TestApply_StripDoubleQuotes(t *testing.T) {
	src := makeFile(entry("HOST", `"localhost"`))
	out := normalizer.Apply(src, normalizer.Options{StripQuotes: true})
	if got := out.Entries[0].Value; got != "localhost" {
		t.Errorf("expected 'localhost', got %q", got)
	}
}

func TestApply_StripSingleQuotes(t *testing.T) {
	src := makeFile(entry("HOST", "'localhost'"))
	out := normalizer.Apply(src, normalizer.Options{StripQuotes: true})
	if got := out.Entries[0].Value; got != "localhost" {
		t.Errorf("expected 'localhost', got %q", got)
	}
}

func TestApply_TrimValues(t *testing.T) {
	src := makeFile(entry("HOST", "  localhost  "))
	out := normalizer.Apply(src, normalizer.Options{TrimValues: true})
	if got := out.Entries[0].Value; got != "localhost" {
		t.Errorf("expected 'localhost', got %q", got)
	}
}

func TestApply_CollapseEmptyValues(t *testing.T) {
	src := makeFile(entry("EMPTY", "   "))
	out := normalizer.Apply(src, normalizer.Options{CollapseEmptyValues: true})
	if got := out.Entries[0].Value; got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestApply_NoOptions_ReturnsCopy(t *testing.T) {
	src := makeFile(entry("db_host", "  'val'  "))
	out := normalizer.Apply(src, normalizer.Options{})
	if got := out.Entries[0].Key; got != "db_host" {
		t.Errorf("expected key unchanged, got %q", got)
	}
	if got := out.Entries[0].Value; got != "  'val'  " {
		t.Errorf("expected value unchanged, got %q", got)
	}
}

func TestApply_DefaultOptions_AllSteps(t *testing.T) {
	src := makeFile(
		entry("db_host", `  "localhost"  `),
		entry("empty_val", "   "),
	)
	out := normalizer.Apply(src, normalizer.DefaultOptions())

	if got := out.Entries[0].Key; got != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %q", got)
	}
	if got := out.Entries[0].Value; got != "localhost" {
		t.Errorf("expected 'localhost', got %q", got)
	}
	if got := out.Entries[1].Value; got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}
