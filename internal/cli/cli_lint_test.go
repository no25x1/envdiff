package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envdiff/internal/cli"
)

func writeTempEnvLint(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRun_LintMissingArg(t *testing.T) {
	err := cli.Run([]string{"lint"})
	if err == nil {
		t.Fatal("expected error for missing file argument")
	}
}

func TestRun_LintMissingFile(t *testing.T) {
	err := cli.Run([]string{"lint", "/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRun_LintCleanFile(t *testing.T) {
	p := writeTempEnvLint(t, "DATABASE_URL=postgres://localhost/db\nPORT=8080\n")
	err := cli.Run([]string{"lint", p})
	if err != nil {
		t.Fatalf("expected no error for clean file, got: %v", err)
	}
}

func TestRun_LintFileWithErrors(t *testing.T) {
	p := writeTempEnvLint(t, "db_host=localhost\nSECRET=\n")
	err := cli.Run([]string{"lint", p})
	if err == nil {
		t.Fatal("expected error due to lowercase key")
	}
}

func TestRun_LintFileWithWarningsOnly(t *testing.T) {
	// Leading underscore and empty value are warnings, not errors.
	p := writeTempEnvLint(t, "_INTERNAL=\n")
	err := cli.Run([]string{"lint", p})
	// warnings alone should not cause a non-nil error return
	if err != nil {
		t.Fatalf("expected no error for warnings-only file, got: %v", err)
	}
}

func TestRun_LintEmptyFile(t *testing.T) {
	// An empty .env file should be considered valid with no errors.
	p := writeTempEnvLint(t, "")
	err := cli.Run([]string{"lint", p})
	if err != nil {
		t.Fatalf("expected no error for empty file, got: %v", err)
	}
}
