package aliaser_test

import (
	"testing"

	"envdiff/internal/aliaser"
	"envdiff/internal/parser"
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
	out, errs := aliaser.Apply(nil, aliaser.Options{})
	if out == nil {
		t.Fatal("expected non-nil result")
	}
	if len(out.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(out.Entries))
	}
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestApply_NoMappings_ReturnsCopy(t *testing.T) {
	f := makeFile("DB_HOST", "localhost", "DB_PORT", "5432")
	out, errs := aliaser.Apply(f, aliaser.Options{})
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(out.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out.Entries))
	}
}

func TestApply_AddsAlias(t *testing.T) {
	f := makeFile("DB_HOST", "localhost")
	opts := aliaser.Options{
		Mappings: []aliaser.Mapping{{Alias: "DATABASE_HOST", Source: "DB_HOST"}},
	}
	out, errs := aliaser.Apply(f, opts)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out.Entries))
	}
	if out.Entries[1].Key != "DATABASE_HOST" || out.Entries[1].Value != "localhost" {
		t.Errorf("unexpected alias entry: %+v", out.Entries[1])
	}
}

func TestApply_SourceNotFound_ReturnsError(t *testing.T) {
	f := makeFile("FOO", "bar")
	opts := aliaser.Options{
		Mappings: []aliaser.Mapping{{Alias: "BAZ", Source: "MISSING"}},
	}
	out, errs := aliaser.Apply(f, opts)
	if len(errs) == 0 {
		t.Fatal("expected an error for missing source")
	}
	if len(out.Entries) != 1 {
		t.Errorf("expected original 1 entry, got %d", len(out.Entries))
	}
}

func TestApply_ExistingAlias_NoOverwrite(t *testing.T) {
	f := makeFile("SRC", "hello", "ALIAS_KEY", "original")
	opts := aliaser.Options{
		Mappings: []aliaser.Mapping{{Alias: "ALIAS_KEY", Source: "SRC"}},
		Overwrite: false,
	}
	out, errs := aliaser.Apply(f, opts)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if out.Entries[1].Value != "original" {
		t.Errorf("expected value to remain 'original', got %q", out.Entries[1].Value)
	}
}

func TestApply_ExistingAlias_WithOverwrite(t *testing.T) {
	f := makeFile("SRC", "hello", "ALIAS_KEY", "original")
	opts := aliaser.Options{
		Mappings: []aliaser.Mapping{{Alias: "ALIAS_KEY", Source: "SRC"}},
		Overwrite: true,
	}
	out, errs := aliaser.Apply(f, opts)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if out.Entries[1].Value != "hello" {
		t.Errorf("expected overwritten value 'hello', got %q", out.Entries[1].Value)
	}
}
