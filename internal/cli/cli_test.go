package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestRun_NoArgs(t *testing.T) {
	err := Run([]string{})
	if err == nil {
		t.Fatal("expected error for no args")
	}
}

func TestRun_UnknownCommand(t *testing.T) {
	err := Run([]string{"unknown"})
	if err == nil || !strings.Contains(err.Error(), "unknown command") {
		t.Fatalf("expected unknown command error, got %v", err)
	}
}

func TestRun_DiffMissingFiles(t *testing.T) {
	err := Run([]string{"diff"})
	if err == nil {
		t.Fatal("expected error when no files provided")
	}
}

func TestRun_DiffTextFormat(t *testing.T) {
	base := writeTempEnv(t, "FOO=bar\nSECRET_KEY=abc\n")
	target := writeTempEnv(t, "FOO=baz\nNEW_VAR=hello\n")

	err := Run([]string{"diff", base, target})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_DiffJSONFormat(t *testing.T) {
	base := writeTempEnv(t, "FOO=bar\n")
	target := writeTempEnv(t, "FOO=bar\n")

	err := Run([]string{"diff", "-format", "json", base, target})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_ReconcileMissingFiles(t *testing.T) {
	err := Run([]string{"reconcile"})
	if err == nil {
		t.Fatal("expected error when no files provided")
	}
}

func TestRun_ReconcileToStdout(t *testing.T) {
	base := writeTempEnv(t, "FOO=bar\nOLD=val\n")
	target := writeTempEnv(t, "FOO=bar\nNEW=hello\n")

	err := Run([]string{"reconcile", base, target})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_ReconcileToFile(t *testing.T) {
	base := writeTempEnv(t, "FOO=bar\n")
	target := writeTempEnv(t, "FOO=baz\n")
	out := filepath.Join(t.TempDir(), "out.env")

	err := Run([]string{"reconcile", "-out", out, base, target})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}
	if !strings.Contains(string(data), "FOO") {
		t.Errorf("expected output to contain FOO, got: %s", data)
	}
}
