package truncator

import (
	"strings"
	"testing"

	"github.com/your-org/envdiff/internal/parser"
)

func makeFile(pairs ...string) *parser.EnvFile {
	f := &parser.EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		f.Entries = append(f.Entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return f
}

func TestApply_NilFile_ReturnsEmpty(t *testing.T) {
	out := Apply(nil, DefaultOptions())
	if out == nil || len(out.Entries) != 0 {
		t.Fatal("expected empty EnvFile for nil input")
	}
}

func TestApply_ShortValues_Unchanged(t *testing.T) {
	src := makeFile("KEY", "short")
	out := Apply(src, DefaultOptions())
	if out.Entries[0].Value != "short" {
		t.Fatalf("expected 'short', got %q", out.Entries[0].Value)
	}
}

func TestApply_LongValue_Truncated(t *testing.T) {
	long := strings.Repeat("a", 100)
	src := makeFile("KEY", long)
	opts := DefaultOptions() // MaxLen=64, Suffix="..."
	out := Apply(src, opts)
	got := out.Entries[0].Value
	if len([]rune(got)) != 64 {
		t.Fatalf("expected length 64, got %d", len([]rune(got)))
	}
	if !strings.HasSuffix(got, "...") {
		t.Fatalf("expected suffix '...', got %q", got)
	}
}

func TestApply_ExactLength_Unchanged(t *testing.T) {
	val := strings.Repeat("b", 64)
	src := makeFile("KEY", val)
	out := Apply(src, DefaultOptions())
	if out.Entries[0].Value != val {
		t.Fatal("value at exact MaxLen should not be truncated")
	}
}

func TestApply_KeysOnly_LimitsScope(t *testing.T) {
	long := strings.Repeat("x", 100)
	src := makeFile("KEEP", long, "SKIP", long)
	opts := Options{MaxLen: 10, Suffix: "…", KeysOnly: []string{"KEEP"}}
	out := Apply(src, opts)

	for _, e := range out.Entries {
		switch e.Key {
		case "KEEP":
			if len([]rune(e.Value)) != 10 {
				t.Fatalf("KEEP: expected len 10, got %d", len([]rune(e.Value)))
			}
		case "SKIP":
			if e.Value != long {
				t.Fatal("SKIP should be unchanged")
			}
		}
	}
}

func TestApply_NoSuffix_PlainCut(t *testing.T) {
	src := makeFile("K", "abcdefghij")
	opts := Options{MaxLen: 5, Suffix: ""}
	out := Apply(src, opts)
	if out.Entries[0].Value != "abcde" {
		t.Fatalf("expected 'abcde', got %q", out.Entries[0].Value)
	}
}
