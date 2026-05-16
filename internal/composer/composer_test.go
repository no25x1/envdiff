package composer_test

import (
	"os"
	"path/filepath"
	"testing"

	"envdiff/internal/composer"
	"envdiff/internal/parser"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func keys(f *parser.EnvFile) []string {
	out := make([]string, len(f.Entries))
	for i, e := range f.Entries {
		out[i] = e.Key
	}
	return out
}

func TestApply_NoSources_ReturnsError(t *testing.T) {
	_, err := composer.Apply(composer.Options{})
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}

func TestApply_SingleSource_ReturnsCopy(t *testing.T) {
	p := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	f, err := composer.Apply(composer.Options{Sources: []composer.Source{{Path: p}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(f.Entries))
	}
}

func TestApply_Namespace_PrefixesKeys(t *testing.T) {
	p := writeTempEnv(t, "HOST=localhost\nPORT=5432\n")
	f, err := composer.Apply(composer.Options{
		Sources: []composer.Source{{Path: p, Namespace: "db"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range f.Entries {
		if len(e.Key) < 3 || e.Key[:3] != "DB_" {
			t.Errorf("expected DB_ prefix, got %q", e.Key)
		}
	}
}

func TestApply_FirstWriterWins_OnConflict(t *testing.T) {
	p1 := writeTempEnv(t, "FOO=first\n")
	p2 := writeTempEnv(t, "FOO=second\n")
	f, err := composer.Apply(composer.Options{
		Sources:             []composer.Source{{Path: p1}, {Path: p2}},
		OverwriteOnConflict: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Entries[0].Value != "first" {
		t.Errorf("expected first writer to win, got %q", f.Entries[0].Value)
	}
}

func TestApply_LastWriterWins_OnConflict(t *testing.T) {
	p1 := writeTempEnv(t, "FOO=first\n")
	p2 := writeTempEnv(t, "FOO=second\n")
	f, err := composer.Apply(composer.Options{
		Sources:             []composer.Source{{Path: p1}, {Path: p2}},
		OverwriteOnConflict: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Entries[0].Value != "second" {
		t.Errorf("expected last writer to win, got %q", f.Entries[0].Value)
	}
}

func TestApply_MissingFile_ReturnsError(t *testing.T) {
	p := filepath.Join(t.TempDir(), "nonexistent.env")
	_, err := composer.Apply(composer.Options{Sources: []composer.Source{{Path: p}}})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
