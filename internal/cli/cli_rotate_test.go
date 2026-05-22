package cli_test

import (
	"os"
	"strings"
	"testing"
)

func writeTempEnvRotate(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestRun_RotateMissingArg(t *testing.T) {
	err := Run([]string{"rotate"})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRun_RotateMissingMapFlag(t *testing.T) {
	p := writeTempEnvRotate(t, "OLD_KEY=value\n")
	err := Run([]string{"rotate", p})
	if err == nil {
		t.Fatal("expected error for missing --map flag")
	}
}

func TestRun_RotateMissingFile(t *testing.T) {
	err := Run([]string{"rotate", "/nonexistent.env", "--map", "A=B"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRun_RotateRenamesKey(t *testing.T) {
	p := writeTempEnvRotate(t, "OLD_KEY=hello\nKEEP=world\n")
	err := Run([]string{"rotate", p, "--map", "OLD_KEY=NEW_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_RotateInvalidMapping(t *testing.T) {
	p := writeTempEnvRotate(t, "A=1\n")
	err := Run([]string{"rotate", p, "--map", "BADFORMAT"})
	if err == nil {
		t.Fatal("expected error for bad mapping")
	}
	if !strings.Contains(err.Error(), "invalid mapping") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRun_RotateFailOnMissing(t *testing.T) {
	p := writeTempEnvRotate(t, "KEEP=value\n")
	err := Run([]string{"rotate", p, "--map", "GONE=NEW", "--fail-on-missing"})
	if err == nil {
		t.Fatal("expected error for missing key with --fail-on-missing")
	}
}
