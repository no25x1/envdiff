package placeholder_test

import (
	"testing"

	"envdiff/internal/parser"
	"envdiff/internal/placeholder"
)

func makeFile(entries ...parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(k, v string) parser.Entry { return parser.Entry{Key: k, Value: v} }

func TestApply_NilFile_ReturnsNil(t *testing.T) {
	findings, err := placeholder.Apply(nil, placeholder.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if findings != nil {
		t.Errorf("expected nil findings, got %v", findings)
	}
}

func TestApply_NoPlaceholders(t *testing.T) {
	f := makeFile(entry("DB_HOST", "localhost"), entry("PORT", "5432"))
	findings, err := placeholder.Apply(f, placeholder.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

func TestApply_DetectsChangeme(t *testing.T) {
	f := makeFile(entry("API_KEY", "changeme"), entry("HOST", "prod.example.com"))
	findings, err := placeholder.Apply(f, placeholder.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "API_KEY" {
		t.Errorf("expected key API_KEY, got %s", findings[0].Key)
	}
}

func TestApply_DetectsAngleBracket(t *testing.T) {
	f := makeFile(entry("SECRET", "<your-secret>"))
	findings, err := placeholder.Apply(f, placeholder.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 1 {
		t.Errorf("expected 1 finding, got %d", len(findings))
	}
}

func TestApply_CaseInsensitive(t *testing.T) {
	f := makeFile(entry("TOKEN", "CHANGEME"), entry("PW", "Placeholder"))
	findings, err := placeholder.Apply(f, placeholder.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 2 {
		t.Errorf("expected 2 findings, got %d", len(findings))
	}
}

func TestApply_CustomPatterns(t *testing.T) {
	f := makeFile(entry("KEY", "REDACTED"), entry("OTHER", "changeme"))
	findings, err := placeholder.Apply(f, placeholder.Options{
		Patterns: []string{`(?i)^redacted$`},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 1 || findings[0].Key != "KEY" {
		t.Errorf("expected 1 finding for KEY, got %v", findings)
	}
}

func TestApply_InvalidPattern_ReturnsError(t *testing.T) {
	f := makeFile(entry("K", "v"))
	_, err := placeholder.Apply(f, placeholder.Options{Patterns: []string{`[invalid`}})
	if err == nil {
		t.Error("expected error for invalid pattern, got nil")
	}
}
