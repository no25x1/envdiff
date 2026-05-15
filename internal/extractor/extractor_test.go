package extractor_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/extractor"
	"github.com/your-org/envdiff/internal/parser"
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

func TestApply_NilFile_ReturnsError(t *testing.T) {
	_, err := extractor.Apply(nil, extractor.Options{Keys: []string{"A"}})
	if err == nil {
		t.Fatal("expected error for nil file")
	}
}

func TestApply_NoOptions_ReturnsError(t *testing.T) {
	f := makeFile("A", "1")
	_, err := extractor.Apply(f, extractor.Options{})
	if err == nil {
		t.Fatal("expected error when no keys or prefix given")
	}
}

func TestApply_ExplicitKeys(t *testing.T) {
	f := makeFile("A", "1", "B", "2", "C", "3")
	out, err := extractor.Apply(f, extractor.Options{Keys: []string{"A", "C"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := keys(out)
	if len(got) != 2 || got[0] != "A" || got[1] != "C" {
		t.Errorf("expected [A C], got %v", got)
	}
}

func TestApply_Prefix(t *testing.T) {
	f := makeFile("DB_HOST", "localhost", "DB_PORT", "5432", "APP_NAME", "myapp")
	out, err := extractor.Apply(f, extractor.Options{Prefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := keys(out)
	if len(got) != 2 {
		t.Errorf("expected 2 entries, got %d: %v", len(got), got)
	}
}

func TestApply_StripPrefix(t *testing.T) {
	f := makeFile("DB_HOST", "localhost", "DB_PORT", "5432")
	out, err := extractor.Apply(f, extractor.Options{Prefix: "DB_", StripPrefix: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := keys(out)
	if got[0] != "HOST" || got[1] != "PORT" {
		t.Errorf("expected [HOST PORT], got %v", got)
	}
}

func TestApply_KeysAndPrefix_Combined(t *testing.T) {
	f := makeFile("DB_HOST", "h", "APP_NAME", "n", "SECRET", "s")
	out, err := extractor.Apply(f, extractor.Options{Keys: []string{"SECRET"}, Prefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out.Entries))
	}
}
