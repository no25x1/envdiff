package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envdiff/internal/cli"
)

func writeTempEnvScope(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnvScope: %v", err)
	}
	return p
}

func TestRun_ScopeMissingArg(t *testing.T) {
	err := cli.Run([]string{"scope", "-scope=PROD"})
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Fatalf("expected usage error, got %v", err)
	}
}

func TestRun_ScopeMissingFlag(t *testing.T) {
	p := writeTempEnvScope(t, "PROD_HOST=localhost\n")
	err := cli.Run([]string{"scope", p})
	if err == nil || !strings.Contains(err.Error(), "-scope") {
		t.Fatalf("expected -scope error, got %v", err)
	}
}

func TestRun_ScopeMissingFile(t *testing.T) {
	err := cli.Run([]string{"scope", "-scope=PROD", "/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRun_ScopeFiltersKeys(t *testing.T) {
	p := writeTempEnvScope(t, "PROD_HOST=prod.example.com\nSTAGING_HOST=staging.example.com\nPROD_PORT=443\n")
	out := captureStdout(t, func() {
		if err := cli.Run([]string{"scope", "-scope=PROD", p}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(out, "PROD_HOST") {
		t.Errorf("expected PROD_HOST in output")
	}
	if strings.Contains(out, "STAGING_HOST") {
		t.Errorf("unexpected STAGING_HOST in output")
	}
}

func TestRun_ScopeStripsPrefix(t *testing.T) {
	p := writeTempEnvScope(t, "PROD_HOST=prod.example.com\nPROD_PORT=443\n")
	out := captureStdout(t, func() {
		if err := cli.Run([]string{"scope", "-scope=PROD", "-strip", p}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	if strings.Contains(out, "PROD_") {
		t.Errorf("expected prefix to be stripped, got: %s", out)
	}
	if !strings.Contains(out, "HOST") || !strings.Contains(out, "PORT") {
		t.Errorf("expected stripped keys HOST and PORT, got: %s", out)
	}
}
