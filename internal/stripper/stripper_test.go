package stripper_test

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/stripper"
)

func makeFile(pairs ...string) *parser.EnvFile {
	f := &parser.EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		f.Entries = append(f.Entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return f
}

func keys(f *parser.EnvFile) []string {
	out := make([]string, len(f.Entries))
	for i, e := range f.Entries {
		out[i] = e.Key
	}
	return out
}

func TestApply_NoOptions_ReturnsCopy(t *testing.T) {
	src := makeFile("A", "1", "B", "2")
	out := stripper.Apply(src, stripper.Options{})
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out.Entries))
	}
}

func TestApply_ExplicitKeys(t *testing.T) {
	src := makeFile("A", "1", "B", "2", "C", "3")
	out := stripper.Apply(src, stripper.Options{Keys: []string{"B"}})
	for _, e := range out.Entries {
		if e.Key == "B" {
			t.Fatal("expected B to be stripped")
		}
	}
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out.Entries))
	}
}

func TestApply_Prefix(t *testing.T) {
	src := makeFile("STAGING_HOST", "h", "STAGING_PORT", "80", "PROD_HOST", "p")
	out := stripper.Apply(src, stripper.Options{Prefixes: []string{"STAGING_"}})
	for _, e := range out.Entries {
		if e.Key == "STAGING_HOST" || e.Key == "STAGING_PORT" {
			t.Fatalf("expected %s to be stripped", e.Key)
		}
	}
	if len(out.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out.Entries))
	}
}

func TestApply_Suffix(t *testing.T) {
	src := makeFile("DB_URL_OLD", "x", "API_KEY_OLD", "y", "DB_HOST", "z")
	out := stripper.Apply(src, stripper.Options{Suffixes: []string{"_OLD"}})
	if len(out.Entries) != 1 || out.Entries[0].Key != "DB_HOST" {
		t.Fatalf("unexpected entries: %v", keys(out))
	}
}

func TestApply_EmptyValues(t *testing.T) {
	src := makeFile("A", "", "B", "hello", "C", "")
	out := stripper.Apply(src, stripper.Options{EmptyValues: true})
	if len(out.Entries) != 1 || out.Entries[0].Key != "B" {
		t.Fatalf("unexpected entries: %v", keys(out))
	}
}

func TestApply_OriginalUnmodified(t *testing.T) {
	src := makeFile("X", "1", "Y", "2")
	stripper.Apply(src, stripper.Options{Keys: []string{"X"}})
	if len(src.Entries) != 2 {
		t.Fatal("original EnvFile should not be modified")
	}
}
