package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/cli"
)

func writeTempEnvDefault(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestRun_DefaultMissingArgs(t *testing.T) {
	err := cli.Run([]string{"default"})
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Fatalf("expected usage error, got %v", err)
	}
}

func TestRun_DefaultMissingDefaultsFile(t *testing.T) {
	tgt := writeTempEnvDefault(t, "PORT=9999\n")
	err := cli.Run([]string{"default", "/nonexistent/defaults.env", tgt})
	if err == nil {
		t.Fatal("expected error for missing defaults file")
	}
}

func TestRun_DefaultFillsMissingKeys(t *testing.T) {
	def := writeTempEnvDefault(t, "HOST=localhost\nPORT=5432\n")
	tgt := writeTempEnvDefault(t, "PORT=9999\n")
	out := filepath.Join(t.TempDir(), "out.env")

	err := cli.Run([]string{"default", def, tgt, "-out", out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	body := string(data)
	if !strings.Contains(body, "HOST=localhost") {
		t.Errorf("expected HOST=localhost in output, got:\n%s", body)
	}
	if !strings.Contains(body, "PORT=9999") {
		t.Errorf("expected PORT=9999 (original) in output, got:\n%s", body)
	}
}

func TestRun_DefaultOverwrite(t *testing.T) {
	def := writeTempEnvDefault(t, "HOST=localhost\n")
	tgt := writeTempEnvDefault(t, "HOST=prod.example.com\n")
	out := filepath.Join(t.TempDir(), "out.env")

	err := cli.Run([]string{"default", def, tgt, "-overwrite", "-out", out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "HOST=localhost") {
		t.Errorf("expected HOST overwritten to localhost, got:\n%s", string(data))
	}
}

func TestRun_DefaultKeyFilter(t *testing.T) {
	def := writeTempEnvDefault(t, "HOST=localhost\nPORT=5432\nDB=mydb\n")
	tgt := writeTempEnvDefault(t, "")
	out := filepath.Join(t.TempDir(), "out.env")

	err := cli.Run([]string{"default", def, tgt, "-keys", "HOST,DB", "-out", out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	body := string(data)
	if strings.Contains(body, "PORT") {
		t.Errorf("PORT should have been filtered out, got:\n%s", body)
	}
	if !strings.Contains(body, "HOST") || !strings.Contains(body, "DB") {
		t.Errorf("expected HOST and DB in output, got:\n%s", body)
	}
}
