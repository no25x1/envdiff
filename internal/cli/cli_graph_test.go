package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/envdiff/envdiff/internal/cli"
)

func writeTempEnvGraph(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRun_GraphMissingArg(t *testing.T) {
	var buf bytes.Buffer
	err := cli.Run([]string{"graph"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing arg")
	}
}

func TestRun_GraphMissingFile(t *testing.T) {
	var buf bytes.Buffer
	err := cli.Run([]string{"graph", "/nonexistent/.env"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRun_GraphNoRefs(t *testing.T) {
	p := writeTempEnvGraph(t, "A=hello\nB=world\n")
	var buf bytes.Buffer
	if err := cli.Run([]string{"graph", p}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No variable references") {
		t.Errorf("expected no-references message, got: %s", buf.String())
	}
}

func TestRun_GraphWithRefs(t *testing.T) {
	p := writeTempEnvGraph(t, "BASE=/app\nDIR=${BASE}/data\n")
	var buf bytes.Buffer
	if err := cli.Run([]string{"graph", p}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DIR -> BASE") {
		t.Errorf("expected edge DIR -> BASE in output, got: %s", out)
	}
	if !strings.Contains(out, "Evaluation order") {
		t.Errorf("expected evaluation order section, got: %s", out)
	}
}

func TestRun_GraphJSONFormat(t *testing.T) {
	p := writeTempEnvGraph(t, "BASE=/app\nDIR=${BASE}/data\n")
	var buf bytes.Buffer
	if err := cli.Run([]string{"graph", p, "--format", "json"}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "\"edges\"") {
		t.Errorf("expected JSON with edges field, got: %s", out)
	}
	if !strings.Contains(out, "\"order\"") {
		t.Errorf("expected JSON with order field, got: %s", out)
	}
}

func TestRun_GraphCycleError(t *testing.T) {
	p := writeTempEnvGraph(t, "X=${Y}\nY=${X}\n")
	var buf bytes.Buffer
	err := cli.Run([]string{"graph", p}, &buf)
	if err == nil {
		t.Fatal("expected cycle error")
	}
	if !strings.Contains(err.Error(), "cycle") {
		t.Errorf("expected cycle in error message, got: %v", err)
	}
}
