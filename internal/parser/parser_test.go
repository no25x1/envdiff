package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestParse_BasicEntries(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDB_HOST=localhost\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(env.Entries))
	}
	if env.Index["APP_ENV"].Value != "production" {
		t.Errorf("expected APP_ENV=production, got %s", env.Index["APP_ENV"].Value)
	}
}

func TestParse_SkipsComments(t *testing.T) {
	path := writeTempEnv(t, "# comment\nKEY=value\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(env.Entries))
	}
}

func TestParse_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `SECRET="my secret value"` + "\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.Index["SECRET"].Value != "my secret value" {
		t.Errorf("unexpected value: %q", env.Index["SECRET"].Value)
	}
}

func TestParse_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "BADLINE\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for invalid line")
	}
}

func TestParse_InlineComment(t *testing.T) {
	path := writeTempEnv(t, "PORT=8080 # http port\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := env.Index["PORT"]
	if e.Value != "8080" {
		t.Errorf("expected value 8080, got %s", e.Value)
	}
	if e.Comment != "http port" {
		t.Errorf("expected comment 'http port', got %s", e.Comment)
	}
}
