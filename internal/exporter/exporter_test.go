package exporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/exporter"
	"github.com/user/envdiff/internal/parser"
)

func makeEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "SECRET_KEY", Value: "s3cr3t value"},
	}
}

func TestNew_ValidFormats(t *testing.T) {
	for _, f := range []exporter.Format{"dotenv", "json", "shell"} {
		_, err := exporter.New(f)
		if err != nil {
			t.Errorf("expected no error for format %q, got %v", f, err)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := exporter.New("yaml")
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}

func TestWrite_DotenvFormat(t *testing.T) {
	e, _ := exporter.New(exporter.FormatDotenv)
	var buf bytes.Buffer
	if err := e.Write(&buf, makeEntries()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production in output, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output, got:\n%s", out)
	}
	// value with space should be quoted
	if !strings.Contains(out, `"s3cr3t value"`) {
		t.Errorf("expected quoted value for SECRET_KEY, got:\n%s", out)
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	e, _ := exporter.New(exporter.FormatJSON)
	var buf bytes.Buffer
	if err := e.Write(&buf, makeEntries()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]string
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if m["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", m["DB_HOST"])
	}
}

func TestWrite_ShellFormat(t *testing.T) {
	e, _ := exporter.New(exporter.FormatShell)
	var buf bytes.Buffer
	if err := e.Write(&buf, makeEntries()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export APP_ENV='production'") {
		t.Errorf("expected export APP_ENV='production' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "export DB_HOST='localhost'") {
		t.Errorf("expected export DB_HOST='localhost' in output, got:\n%s", out)
	}
}

func TestWrite_SortedOutput(t *testing.T) {
	e, _ := exporter.New(exporter.FormatDotenv)
	var buf bytes.Buffer
	_ = e.Write(&buf, makeEntries())
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if !strings.HasPrefix(lines[0], "APP_ENV") {
		t.Errorf("expected first line to be APP_ENV (sorted), got %q", lines[0])
	}
}
