package cli_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/your-org/envdiff/internal/cli"
)

func writeTempJSONImport(t *testing.T, data map[string]string) string {
	t.Helper()
	b, _ := json.Marshal(data)
	p := filepath.Join(t.TempDir(), "input.json")
	_ = os.WriteFile(p, b, 0644)
	return p
}

func TestRun_ImportMissingArg(t *testing.T) {
	err := cli.Run([]string{"import"})
	if err == nil || !strings.Contains(err.Error(), "requires") {
		t.Fatalf("expected requires error, got %v", err)
	}
}

func TestRun_ImportMissingFile(t *testing.T) {
	err := cli.Run([]string{"import", "/no/such/file.json"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRun_ImportJSONToDotenv(t *testing.T) {
	path := writeTempJSONImport(t, map[string]string{"HOST": "db", "PORT": "5432"})
	err := cli.Run([]string{"import", path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_ImportWithPrefix(t *testing.T) {
	path := writeTempJSONImport(t, map[string]string{"NAME": "svc"})
	out := filepath.Join(t.TempDir(), "out.env")
	err := cli.Run([]string{"import", "-prefix", "APP_", "-output", out, path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "APP_NAME") {
		t.Errorf("expected APP_NAME in output, got: %s", string(data))
	}
}

func TestRun_ImportJSONOutput(t *testing.T) {
	path := writeTempJSONImport(t, map[string]string{"KEY": "value"})
	out := filepath.Join(t.TempDir(), "out.json")
	err := cli.Run([]string{"import", "-out-format", "json", "-output", out, path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "KEY") {
		t.Errorf("expected KEY in JSON output, got: %s", string(data))
	}
}
