package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/your-org/envdiff/internal/cli"
)

func writeTempEnvEncode(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempEnvEncode: %v", err)
	}
	return p
}

func TestRun_EncodeMissingArg(t *testing.T) {
	err := cli.Run([]string{"encode"})
	if err == nil || !strings.Contains(err.Error(), "input file required") {
		t.Errorf("expected input file required error, got %v", err)
	}
}

func TestRun_EncodeMissingFile(t *testing.T) {
	err := cli.Run([]string{"encode", "/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRun_EncodeBase64(t *testing.T) {
	p := writeTempEnvEncode(t, "SECRET=hello\nOTHER=world\n")
	out := filepath.Join(t.TempDir(), "out.env")
	err := cli.Run([]string{"encode", "-strategy", "base64", "-output", out, p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, _ := os.ReadFile(out)
	if !strings.Contains(string(b), "aGVsbG8=") {
		t.Errorf("expected base64 of 'hello' in output, got:\n%s", b)
	}
}

func TestRun_EncodeBase64_KeysOnly(t *testing.T) {
	p := writeTempEnvEncode(t, "SECRET=hello\nOTHER=world\n")
	out := filepath.Join(t.TempDir(), "out.env")
	err := cli.Run([]string{"encode", "-strategy", "base64", "-keys", "SECRET", "-output", out, p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, _ := os.ReadFile(out)
	s := string(b)
	if !strings.Contains(s, "aGVsbG8=") {
		t.Errorf("expected SECRET to be encoded")
	}
	if !strings.Contains(s, "world") {
		t.Errorf("expected OTHER to remain as 'world'")
	}
}

func TestRun_EncodeUnknownStrategy(t *testing.T) {
	p := writeTempEnvEncode(t, "X=val\n")
	err := cli.Run([]string{"encode", "-strategy", "rot13", p})
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}
