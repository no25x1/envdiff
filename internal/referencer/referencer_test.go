package referencer_test

import (
	"testing"

	"envdiff/internal/parser"
	"envdiff/internal/referencer"
)

func makeFile(entries ...parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(key, value string) parser.Entry {
	return parser.Entry{Key: key, Value: value}
}

func TestAnalyse_NilFile(t *testing.T) {
	_, err := referencer.Analyse(nil)
	if err == nil {
		t.Fatal("expected error for nil file")
	}
}

func TestAnalyse_NoReferences(t *testing.T) {
	f := makeFile(entry("HOST", "localhost"), entry("PORT", "5432"))
	r, err := referencer.Analyse(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Referrers) != 0 {
		t.Errorf("expected no referrers, got %v", r.Referrers)
	}
	if len(r.Undefined) != 0 {
		t.Errorf("expected no undefined, got %v", r.Undefined)
	}
	if len(r.Unused) != 2 {
		t.Errorf("expected 2 unused keys, got %v", r.Unused)
	}
}

func TestAnalyse_CurlyBraceReference(t *testing.T) {
	f := makeFile(
		entry("BASE", "http://localhost"),
		entry("URL", "${BASE}/api"),
	)
	r, err := referencer.Analyse(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if refs, ok := r.Referrers["BASE"]; !ok || len(refs) != 1 || refs[0] != "URL" {
		t.Errorf("expected BASE referenced by URL, got %v", r.Referrers)
	}
	if len(r.Undefined) != 0 {
		t.Errorf("expected no undefined keys, got %v", r.Undefined)
	}
}

func TestAnalyse_BareDollarReference(t *testing.T) {
	f := makeFile(
		entry("HOST", "db"),
		entry("DSN", "postgres://$HOST:5432/mydb"),
	)
	r, err := referencer.Analyse(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if refs := r.Referrers["HOST"]; len(refs) != 1 || refs[0] != "DSN" {
		t.Errorf("expected HOST referenced by DSN, got %v", refs)
	}
}

func TestAnalyse_UndefinedReference(t *testing.T) {
	f := makeFile(entry("URL", "${MISSING_VAR}/path"))
	r, err := referencer.Analyse(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Undefined) != 1 || r.Undefined[0] != "MISSING_VAR" {
		t.Errorf("expected MISSING_VAR in undefined, got %v", r.Undefined)
	}
}

func TestAnalyse_SelfReferenceIgnored(t *testing.T) {
	f := makeFile(entry("PATH", "$PATH:/usr/local/bin"))
	r, err := referencer.Analyse(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if refs, ok := r.Referrers["PATH"]; ok {
		t.Errorf("self-reference should be ignored, got %v", refs)
	}
}
