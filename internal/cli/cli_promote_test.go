package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/cli"
)

func writeTempEnvPromote(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnvPromote: %v", err)
	}
	return p
}

func TestRun_PromoteMissingArgs(t *testing.T) {
	err := cli.Run([]string{"promote"})
	if err == nil || !strings.Contains(err.Error(), "requires") {
		t.Errorf("expected requires error, got %v", err)
	}
}

func TestRun_PromoteMissingSourceFile(t *testing.T) {
	err := cli.Run([]string{"promote", "/nonexistent/src.env", "/nonexistent/dst.env"})
	if err == nil {
		t.Fatal("expected error for missing source file")
	}
}

func TestRun_PromoteAllKeys(t *testing.T) {
	src := writeTempEnvPromote(t, "FOO=bar\nBAZ=qux\n")
	dst := writeTempEnvPromote(t, "")
	err := cli.Run([]string{"promote", src, dst, "--out=-"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_PromoteWithKeyFilter(t *testing.T) {
	src := writeTempEnvPromote(t, "FOO=1\nBAR=2\nBAZ=3\n")
	dst := writeTempEnvPromote(t, "")
	err := cli.Run([]string{"promote", src, dst, "--keys=FOO,BAZ", "--out=-"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_PromoteOverwrite(t *testing.T) {
	src := writeTempEnvPromote(t, "FOO=new\n")
	dst := writeTempEnvPromote(t, "FOO=old\n")
	err := cli.Run([]string{"promote", src, dst, "--overwrite", "--out=-"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_PromoteJSONFormat(t *testing.T) {
	src := writeTempEnvPromote(t, "KEY=value\n")
	dst := writeTempEnvPromote(t, "")
	err := cli.Run([]string{"promote", src, dst, "--format=json", "--out=-"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
