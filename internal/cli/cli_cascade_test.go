package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envdiff/internal/cli"
)

func writeTempEnvCascade(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestRun_CascadeMissingArgs(t *testing.T) {
	err := cli.Run([]string{"cascade"})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRun_CascadeMissingBaseFile(t *testing.T) {
	err := cli.Run([]string{"cascade", "/nonexistent/base.env", "/nonexistent/overlay.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRun_CascadeBasicOverride(t *testing.T) {
	base := writeTempEnvCascade(t, "A=base\nB=base\n")
	overlay := writeTempEnvCascade(t, "A=overlay\n")
	err := cli.Run([]string{"cascade", base, overlay})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_CascadeOutputToFile(t *testing.T) {
	base := writeTempEnvCascade(t, "X=1\nY=2\n")
	overlay := writeTempEnvCascade(t, "X=99\n")
	out := filepath.Join(t.TempDir(), "result.env")
	err := cli.Run([]string{"cascade", base, overlay, "--output", out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "X=99") {
		t.Errorf("expected X=99 in output, got: %s", data)
	}
	if !strings.Contains(string(data), "Y=2") {
		t.Errorf("expected Y=2 in output, got: %s", data)
	}
}

func TestRun_CascadeStopOnMissing(t *testing.T) {
	base := writeTempEnvCascade(t, "A=1\nB=2\n")
	overlay := writeTempEnvCascade(t, "A=x\n")
	err := cli.Run([]string{"cascade", base, overlay, "--stop-on-missing"})
	if err == nil {
		t.Fatal("expected error when B missing from overlay with --stop-on-missing")
	}
}

func TestRun_CascadeJSONFormat(t *testing.T) {
	base := writeTempEnvCascade(t, "K=v\n")
	overlay := writeTempEnvCascade(t, "K=w\n")
	out := filepath.Join(t.TempDir(), "result.json")
	err := cli.Run([]string{"cascade", base, overlay, "--format", "json", "--output", out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "\"K\"") {
		t.Errorf("expected JSON with key K, got: %s", data)
	}
}
