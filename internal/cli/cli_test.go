package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envdiff/internal/cli"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRun_NoArgs(t *testing.T) {
	code := cli.Run([]string{})
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}

func TestRun_UnknownCommand(t *testing.T) {
	code := cli.Run([]string{"unknown"})
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}

func TestRun_Version(t *testing.T) {
	code := cli.Run([]string{"version"})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRun_DiffMissingFiles(t *testing.T) {
	code := cli.Run([]string{"diff", "nonexistent.env", "also-missing.env"})
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}

func TestRun_DiffTextFormat(t *testing.T) {
	base := writeTempEnv(t, "APP_ENV=production\nDB_HOST=db.prod\n")
	compare := writeTempEnv(t, "APP_ENV=staging\nDB_HOST=db.prod\nNEW_KEY=value\n")
	code := cli.Run([]string{"diff", base, compare, "text"})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRun_ValidatePass(t *testing.T) {
	f := writeTempEnv(t, "APP_NAME=myapp\nDB_HOST=localhost\n")
	code := cli.Run([]string{"validate", f})
	// DB_HOST=localhost triggers no-empty-value passes, but uppercase passes
	// no-empty-value passes since localhost is non-empty
	if code != 0 {
		t.Errorf("expected exit code 0 for valid file, got %d", code)
	}
}

func TestRun_ValidateFail(t *testing.T) {
	f := writeTempEnv(t, "app_name=myapp\nDB_HOST=\n")
	code := cli.Run([]string{"validate", f})
	if code == 0 {
		t.Errorf("expected non-zero exit code for invalid file, got 0")
	}
}

func TestRun_ValidateMissingFile(t *testing.T) {
	code := cli.Run([]string{"validate", "no-such-file.env"})
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}

func TestRun_ReconcileMissingArgs(t *testing.T) {
	code := cli.Run([]string{"reconcile"})
	if code != 1 {
		t.Errorf("expected exit code 1, got %d", code)
	}
}
