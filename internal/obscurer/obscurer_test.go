package obscurer_test

import (
	"testing"

	"github.com/envdiff/envdiff/internal/obscurer"
	"github.com/envdiff/envdiff/internal/parser"
)

func makeFile(pairs ...string) *parser.EnvFile {
	f := &parser.EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		f.Entries = append(f.Entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return f
}

func TestApply_NilFile_ReturnsEmpty(t *testing.T) {
	out := obscurer.Apply(nil, obscurer.DefaultOptions())
	if out == nil || len(out.Entries) != 0 {
		t.Fatal("expected empty EnvFile for nil input")
	}
}

func TestApply_EmptyValue_Unchanged(t *testing.T) {
	src := makeFile("EMPTY_KEY", "")
	out := obscurer.Apply(src, obscurer.DefaultOptions())
	if out.Entries[0].Value != "" {
		t.Errorf("expected empty value unchanged, got %q", out.Entries[0].Value)
	}
}

func TestApply_ShortValue_FullyMasked(t *testing.T) {
	src := makeFile("KEY", "ab")
	opts := obscurer.DefaultOptions() // ShowPrefix=2, ShowSuffix=2 => pre+suf >= len
	out := obscurer.Apply(src, opts)
	if out.Entries[0].Value != opts.Mask {
		t.Errorf("expected full mask, got %q", out.Entries[0].Value)
	}
}

func TestApply_LongValue_ShowsPrefixAndSuffix(t *testing.T) {
	src := makeFile("SECRET", "abcdefghij")
	opts := obscurer.DefaultOptions()
	out := obscurer.Apply(src, opts)
	got := out.Entries[0].Value
	if got[:2] != "ab" {
		t.Errorf("expected prefix 'ab', got %q", got[:2])
	}
	if got[len(got)-2:] != "ij" {
		t.Errorf("expected suffix 'ij', got %q", got[len(got)-2:])
	}
	if got[2:len(got)-2] != opts.Mask {
		t.Errorf("expected mask in middle, got %q", got[2:len(got)-2])
	}
}

func TestApply_KeyFilter_OnlyObscuresMatchingKeys(t *testing.T) {
	src := makeFile("SECRET", "abcdefghij", "PLAIN", "abcdefghij")
	opts := obscurer.DefaultOptions()
	opts.Keys = []string{"SECRET"}
	out := obscurer.Apply(src, opts)

	if out.Entries[0].Key != "SECRET" {
		t.Fatal("unexpected key order")
	}
	if out.Entries[0].Value == "abcdefghij" {
		t.Error("SECRET should have been obscured")
	}
	if out.Entries[1].Value != "abcdefghij" {
		t.Errorf("PLAIN should be unchanged, got %q", out.Entries[1].Value)
	}
}

func TestApply_CustomMaskAndWindow(t *testing.T) {
	src := makeFile("TOKEN", "ABCDE12345")
	opts := obscurer.Options{ShowPrefix: 3, ShowSuffix: 1, Mask: "--"}
	out := obscurer.Apply(src, opts)
	got := out.Entries[0].Value
	expected := "ABC--5"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
