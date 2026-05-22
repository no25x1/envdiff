package stacker_test

import (
	"testing"

	"envdiff/internal/parser"
	"envdiff/internal/stacker"
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

func TestApply_NoFiles_ReturnsError(t *testing.T) {
	_, err := stacker.Apply(nil, stacker.Options{})
	if err == nil {
		t.Fatal("expected error for empty slice")
	}
}

func TestApply_NilFile_ReturnsError(t *testing.T) {
	_, err := stacker.Apply([]*parser.EnvFile{nil}, stacker.Options{AllowEmpty: true})
	if err == nil {
		t.Fatal("expected error for nil file")
	}
}

func TestApply_EmptyFile_ErrorWithoutFlag(t *testing.T) {
	_, err := stacker.Apply([]*parser.EnvFile{{Entries: nil}}, stacker.Options{})
	if err == nil {
		t.Fatal("expected error for empty file when AllowEmpty is false")
	}
}

func TestApply_EmptyFile_AllowedWithFlag(t *testing.T) {
	f, err := stacker.Apply([]*parser.EnvFile{{Entries: nil}}, stacker.Options{AllowEmpty: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(f.Entries))
	}
}

func TestApply_LaterFileOverrides(t *testing.T) {
	a := makeFile("HOST", "localhost", "PORT", "5432")
	b := makeFile("PORT", "6543", "DEBUG", "true")

	f, err := stacker.Apply([]*parser.EnvFile{a, b}, stacker.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := map[string]string{}
	for _, e := range f.Entries {
		got[e.Key] = e.Value
	}

	if got["PORT"] != "6543" {
		t.Errorf("PORT: want 6543, got %s", got["PORT"])
	}
	if got["HOST"] != "localhost" {
		t.Errorf("HOST: want localhost, got %s", got["HOST"])
	}
	if got["DEBUG"] != "true" {
		t.Errorf("DEBUG: want true, got %s", got["DEBUG"])
	}
}

func TestApply_Prefix_AddedToAllKeys(t *testing.T) {
	a := makeFile("FOO", "bar", "BAZ", "qux")

	f, err := stacker.Apply([]*parser.EnvFile{a}, stacker.Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, e := range f.Entries {
		if len(e.Key) < 4 || e.Key[:4] != "APP_" {
			t.Errorf("key %q missing APP_ prefix", e.Key)
		}
	}
}

func TestApply_OrderPreserved(t *testing.T) {
	a := makeFile("A", "1", "B", "2")
	b := makeFile("C", "3", "A", "99")

	f, err := stacker.Apply([]*parser.EnvFile{a, b}, stacker.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"A", "B", "C"}
	got := keys(f)
	if len(got) != len(want) {
		t.Fatalf("key count: want %d, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("pos %d: want %s, got %s", i, want[i], got[i])
		}
	}
}
