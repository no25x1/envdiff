package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/cli"
)

func writeTempEnvAnnotate(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRun_AnnotateMissingArg(t *testing.T) {
	err := cli.Run([]string{"annotate"})
	if err == nil || !strings.Contains(err.Error(), "requires a file") {
		t.Fatalf("expected missing-arg error, got %v", err)
	}
}

func TestRun_AnnotateMissingFile(t *testing.T) {
	err := cli.Run([]string{"annotate", "/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRun_AnnotateCleanFile(t *testing.T) {
	p := writeTempEnvAnnotate(t, "APP_NAME=envdiff\nPORT=8080\n")
	err := cli.Run([]string{"annotate", p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_AnnotateSensitiveKeys(t *testing.T) {
	p := writeTempEnvAnnotate(t, "DB_PASSWORD=secret\nAPI_TOKEN=tok\n")
	out := filepath.Join(t.TempDir(), "out.env")
	err := cli.Run([]string{"annotate", p, "--out", out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "sensitive") {
		t.Errorf("expected sensitive annotation in output:\n%s", data)
	}
}

func TestRun_AnnotateOverwrite(t *testing.T) {
	p := writeTempEnvAnnotate(t, "DB_SECRET=x\n")
	out := filepath.Join(t.TempDir(), "out.env")
	err := cli.Run([]string{"annotate", p, "--overwrite", "--out", out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
