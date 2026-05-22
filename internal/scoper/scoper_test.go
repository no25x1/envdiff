package scoper_test

import (
	"testing"

	"envdiff/internal/parser"
	"envdiff/internal/scoper"
)

func makeFile(entries ...parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(k, v string) parser.Entry { return parser.Entry{Key: k, Value: v} }

func keys(f *parser.EnvFile) []string {
	out := make([]string, len(f.Entries))
	for i, e := range f.Entries {
		out[i] = e.Key
	}
	return out
}

func TestApply_NilFile_ReturnsError(t *testing.T) {
	_, err := scoper.Apply(nil, scoper.Options{Scope: "PROD"})
	if err == nil {
		t.Fatal("expected error for nil file")
	}
}

func TestApply_EmptyScope_ReturnsError(t *testing.T) {
	f := makeFile(entry("PROD_HOST", "localhost"))
	_, err := scoper.Apply(f, scoper.Options{Scope: ""})
	if err == nil {
		t.Fatal("expected error for empty scope")
	}
}

func TestApply_FiltersToScope(t *testing.T) {
	f := makeFile(
		entry("PROD_HOST", "prod.example.com"),
		entry("STAGING_HOST", "staging.example.com"),
		entry("PROD_PORT", "443"),
	)
	out, err := scoper.Apply(f, scoper.Options{Scope: "PROD"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out.Entries))
	}
}

func TestApply_StripPrefix(t *testing.T) {
	f := makeFile(
		entry("PROD_HOST", "prod.example.com"),
		entry("PROD_PORT", "443"),
	)
	out, err := scoper.Apply(f, scoper.Options{Scope: "PROD", StripPrefix: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := keys(out)
	for _, k := range got {
		if k == "PROD_HOST" || k == "PROD_PORT" {
			t.Errorf("prefix not stripped: %s", k)
		}
	}
	if got[0] != "HOST" || got[1] != "PORT" {
		t.Errorf("unexpected keys: %v", got)
	}
}

func TestApply_KeepUnscoped(t *testing.T) {
	f := makeFile(
		entry("PROD_HOST", "prod.example.com"),
		entry("STAGING_HOST", "staging.example.com"),
		entry("APP_NAME", "myapp"), // no matching scope peer
	)
	out, err := scoper.Apply(f, scoper.Options{Scope: "PROD", KeepUnscoped: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// PROD_HOST + APP_NAME (unscoped)
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d: %v", len(out.Entries), keys(out))
	}
}

func TestApply_CaseInsensitiveScope(t *testing.T) {
	f := makeFile(
		entry("prod_host", "localhost"),
		entry("staging_host", "remote"),
	)
	out, err := scoper.Apply(f, scoper.Options{Scope: "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out.Entries))
	}
}
