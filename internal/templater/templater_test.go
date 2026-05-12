package templater_test

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/templater"
)

func makeFile(entries ...parser.Entry) parser.EnvFile {
	return parser.EnvFile{Entries: entries}
}

func entry(k, v string) parser.Entry {
	return parser.Entry{Key: k, Value: v}
}

func TestRender_NoPlaceholders(t *testing.T) {
	tmpl := makeFile(entry("FOO", "bar"), entry("BAZ", "qux"))
	vals := makeFile()
	res := templater.Render(tmpl, vals)
	if len(res.Missing) != 0 {
		t.Fatalf("expected no missing, got %v", res.Missing)
	}
	if res.Entries[0].Value != "bar" || res.Entries[1].Value != "qux" {
		t.Errorf("unexpected entries: %v", res.Entries)
	}
}

func TestRender_SubstitutesCurlyBrace(t *testing.T) {
	tmpl := makeFile(entry("DB_URL", "postgres://${DB_HOST}:${DB_PORT}/db"))
	vals := makeFile(entry("DB_HOST", "localhost"), entry("DB_PORT", "5432"))
	res := templater.Render(tmpl, vals)
	want := "postgres://localhost:5432/db"
	if res.Entries[0].Value != want {
		t.Errorf("got %q, want %q", res.Entries[0].Value, want)
	}
	if len(res.Missing) != 0 {
		t.Errorf("unexpected missing: %v", res.Missing)
	}
}

func TestRender_SubstitutesBareDollar(t *testing.T) {
	tmpl := makeFile(entry("GREETING", "hello $NAME"))
	vals := makeFile(entry("NAME", "world"))
	res := templater.Render(tmpl, vals)
	if res.Entries[0].Value != "hello world" {
		t.Errorf("got %q", res.Entries[0].Value)
	}
}

func TestRender_TracksMissingKeys(t *testing.T) {
	tmpl := makeFile(entry("URL", "https://${HOST}:${PORT}"))
	vals := makeFile(entry("HOST", "example.com"))
	res := templater.Render(tmpl, vals)
	if len(res.Missing) != 1 || res.Missing[0] != "PORT" {
		t.Errorf("expected [PORT] missing, got %v", res.Missing)
	}
	// value should still contain unresolved placeholder
	if res.Entries[0].Value != "https://example.com:${PORT}" {
		t.Errorf("unexpected value: %q", res.Entries[0].Value)
	}
}

func TestRender_DeduplicatesMissingKeys(t *testing.T) {
	tmpl := makeFile(
		entry("A", "${MISSING}"),
		entry("B", "${MISSING}"),
	)
	res := templater.Render(tmpl, makeFile())
	if len(res.Missing) != 1 {
		t.Errorf("expected 1 unique missing key, got %d: %v", len(res.Missing), res.Missing)
	}
}

func TestValidate_NoErrors(t *testing.T) {
	tmpl := makeFile(entry("FOO", "${BAR}"))
	vals := makeFile(entry("BAR", "baz"))
	errs := templater.Validate(tmpl, vals)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidate_ReturnsErrors(t *testing.T) {
	tmpl := makeFile(entry("FOO", "${MISSING_KEY}"))
	errs := templater.Validate(tmpl, makeFile())
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0] != "missing value for placeholder: MISSING_KEY" {
		t.Errorf("unexpected error message: %q", errs[0])
	}
}
