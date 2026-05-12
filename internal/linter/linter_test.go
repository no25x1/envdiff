package linter_test

import (
	"testing"

	"github.com/envdiff/internal/linter"
	"github.com/envdiff/internal/parser"
)

func makeFile(entries []parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func TestLint_NoFindings(t *testing.T) {
	f := makeFile([]parser.Entry{
		{Key: "DATABASE_URL", Value: "postgres://localhost/db"},
		{Key: "PORT", Value: "8080"},
	})
	l := linter.New()
	findings := l.Lint(f)
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d: %v", len(findings), findings)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	f := makeFile([]parser.Entry{
		{Key: "SECRET", Value: ""},
	})
	l := linter.WithRules(linter.RuleNoEmptyValue)
	findings := l.Lint(f)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != linter.SeverityWarning {
		t.Errorf("expected warning severity, got %s", findings[0].Severity)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	f := makeFile([]parser.Entry{
		{Key: "db_host", Value: "localhost"},
	})
	l := linter.WithRules(linter.RuleUppercaseKey)
	findings := l.Lint(f)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != linter.SeverityError {
		t.Errorf("expected error severity, got %s", findings[0].Severity)
	}
}

func TestLint_WhitespaceInKey(t *testing.T) {
	f := makeFile([]parser.Entry{
		{Key: "MY KEY", Value: "value"},
	})
	l := linter.WithRules(linter.RuleNoWhitespaceInKey)
	findings := l.Lint(f)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

func TestLint_LeadingUnderscore(t *testing.T) {
	f := makeFile([]parser.Entry{
		{Key: "_INTERNAL", Value: "yes"},
	})
	l := linter.WithRules(linter.RuleNoLeadingUnderscore)
	findings := l.Lint(f)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != linter.SeverityWarning {
		t.Errorf("expected warning severity")
	}
}

func TestFinding_String(t *testing.T) {
	f := linter.Finding{Key: "FOO", Message: "bar", Severity: linter.SeverityError}
	s := f.String()
	if s != "[error] FOO: bar" {
		t.Errorf("unexpected string: %s", s)
	}
}
