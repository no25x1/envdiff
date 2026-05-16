package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envdiff/internal/cli"
)

func writeTempEnvClone(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnvClone: %v", err)
	}
	return p
}

func TestRun_CloneMissingArg(t *testing.T) {
	var buf bytes.Buffer
	err := cli.Run([]string{"clone"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing source arg")
	}
}

func TestRun_CloneMissingSourceFile(t *testing.T) {
	var buf bytes.Buffer
	err := cli.Run([]string{"clone", "/no/such/file.env"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing source file")
	}
}

func TestRun_CloneAllKeys(t *testing.T) {
	src := writeTempEnvClone(t, "FOO=bar\nBAZ=qux\n")
	var buf bytes.Buffer
	if err := cli.Run([]string{"clone", src}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "FOO") || !strings.Contains(out, "BAZ") {
		t.Errorf("expected FOO and BAZ in output, got:\n%s", out)
	}
}

func TestRun_CloneWithPrefixSwap(t *testing.T) {
	src := writeTempEnvClone(t, "PROD_HOST=example.com\nPROD_PORT=443\n")
	var buf bytes.Buffer
	err := cli.Run([]string{"clone", "--strip-prefix=PROD_", "--add-prefix=STAGING_", src}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "STAGING_HOST") {
		t.Errorf("expected STAGING_HOST in output, got:\n%s", out)
	}
	if strings.Contains(out, "PROD_HOST") {
		t.Errorf("did not expect PROD_HOST in output, got:\n%s", out)
	}
}

func TestRun_CloneKeyFilter(t *testing.T) {
	src := writeTempEnvClone(t, "A=1\nB=2\nC=3\n")
	var buf bytes.Buffer
	err := cli.Run([]string{"clone", "--keys=A,C", src}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "B=") {
		t.Errorf("B should have been filtered out, got:\n%s", out)
	}
	if !strings.Contains(out, "A=") || !strings.Contains(out, "C=") {
		t.Errorf("expected A and C in output, got:\n%s", out)
	}
}
