package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envdiff/internal/cli"
)

func writeTempEnvReference(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnvReference: %v", err)
	}
	return p
}

func TestRun_ReferenceMissingArg(t *testing.T) {
	err := cli.Run([]string{"envdiff", "reference"})
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Errorf("expected usage error, got %v", err)
	}
}

func TestRun_ReferenceMissingFile(t *testing.T) {
	err := cli.Run([]string{"envdiff", "reference", "/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRun_ReferenceNoRefs(t *testing.T) {
	p := writeTempEnvReference(t, "HOST=localhost\nPORT=5432\n")
	err := cli.Run([]string{"envdiff", "reference", p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_ReferenceWithRefs(t *testing.T) {
	p := writeTempEnvReference(t, "BASE=http://localhost\nURL=${BASE}/api\n")
	err := cli.Run([]string{"envdiff", "reference", p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_ReferenceJSONOutput(t *testing.T) {
	p := writeTempEnvReference(t, "HOST=db\nDSN=postgres://$HOST:5432/app\n")
	err := cli.Run([]string{"envdiff", "reference", p, "--json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
