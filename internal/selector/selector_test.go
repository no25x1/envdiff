package selector_test

import (
	"testing"

	"envdiff/internal/parser"
	"envdiff/internal/selector"
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

func TestApply_NilFile_ReturnsEmpty(t *testing.T) {
	out, err := selector.Apply(nil, selector.Options{Keys: []string{"A"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 0 {
		t.Errorf("expected empty, got %d entries", len(out.Entries))
	}
}

func TestApply_NoOptions_ReturnsError(t *testing.T) {
	f := makeFile("A", "1")
	_, err := selector.Apply(f, selector.Options{})
	if err == nil {
		t.Fatal("expected error for empty options")
	}
}

func TestApply_ExplicitKeys(t *testing.T) {
	f := makeFile("DB_HOST", "localhost", "DB_PORT", "5432", "APP_NAME", "myapp")
	out, err := selector.Apply(f, selector.Options{Keys: []string{"DB_HOST", "APP_NAME"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := keys(out)
	if len(got) != 2 || got[0] != "DB_HOST" || got[1] != "APP_NAME" {
		t.Errorf("unexpected keys: %v", got)
	}
}

func TestApply_Prefix(t *testing.T) {
	f := makeFile("DB_HOST", "h", "DB_PORT", "5432", "APP_NAME", "x")
	out, err := selector.Apply(f, selector.Options{Prefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := keys(out)
	if len(got) != 2 || got[0] != "DB_HOST" || got[1] != "DB_PORT" {
		t.Errorf("unexpected keys: %v", got)
	}
}

func TestApply_Regex(t *testing.T) {
	f := makeFile("SECRET_KEY", "s", "API_TOKEN", "t", "APP_NAME", "n")
	out, err := selector.Apply(f, selector.Options{Regex: "(SECRET|TOKEN)"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := keys(out)
	if len(got) != 2 {
		t.Errorf("expected 2 keys, got %v", got)
	}
}

func TestApply_InvalidRegex_ReturnsError(t *testing.T) {
	f := makeFile("A", "1")
	_, err := selector.Apply(f, selector.Options{Regex: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestApply_Invert(t *testing.T) {
	f := makeFile("DB_HOST", "h", "DB_PORT", "p", "APP_NAME", "n")
	out, err := selector.Apply(f, selector.Options{Prefix: "DB_", Invert: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := keys(out)
	if len(got) != 1 || got[0] != "APP_NAME" {
		t.Errorf("unexpected keys: %v", got)
	}
}
