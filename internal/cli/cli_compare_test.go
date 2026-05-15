package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envdiff/internal/cli"
)

func writeTempEnvCompare(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRun_CompareMissingArgs(t *testing.T) {
	var buf bytes.Buffer
	err := cli.Run([]string{"compare"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing args")
	}
	if !strings.Contains(err.Error(), "at least two files") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRun_CompareMissingFile(t *testing.T) {
	var buf bytes.Buffer
	err := cli.Run([]string{"compare", "/no/such/a.env", "/no/such/b.env"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRun_CompareOutputContainsKeys(t *testing.T) {
	a := writeTempEnvCompare(t, "APP_ENV=dev\nDB_HOST=localhost\n")
	b := writeTempEnvCompare(t, "APP_ENV=prod\nDB_HOST=localhost\n")
	var buf bytes.Buffer
	if err := cli.Run([]string{"compare", a, b}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_ENV") {
		t.Errorf("expected APP_ENV in output, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got:\n%s", out)
	}
}

func TestRun_CompareConflictsOnly(t *testing.T) {
	a := writeTempEnvCompare(t, "APP_ENV=dev\nSHARED=same\n")
	b := writeTempEnvCompare(t, "APP_ENV=prod\nSHARED=same\n")
	var buf bytes.Buffer
	if err := cli.Run([]string{"compare", "-conflicts", a, b}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_ENV") {
		t.Errorf("expected APP_ENV in conflicts output")
	}
	if strings.Contains(out, "SHARED") {
		t.Errorf("SHARED should not appear in conflicts-only output")
	}
}

func TestRun_CompareNoConflicts(t *testing.T) {
	a := writeTempEnvCompare(t, "KEY=val\n")
	b := writeTempEnvCompare(t, "KEY=val\n")
	var buf bytes.Buffer
	if err := cli.Run([]string{"compare", "-conflicts", a, b}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no conflicts") {
		t.Errorf("expected 'no conflicts' message, got: %s", buf.String())
	}
}
