package cli_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/cli"
)

func writeTempEnvProfile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestRun_ProfileMissingArg(t *testing.T) {
	var buf bytes.Buffer
	err := cli.Run([]string{"profile"}, &buf)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_ProfileMissingFile(t *testing.T) {
	var buf bytes.Buffer
	err := cli.Run([]string{"profile", "/nonexistent/.env"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRun_ProfileTextOutput(t *testing.T) {
	content := "DB_HOST=localhost\nDB_PASSWORD=secret\nAPP_NAME=envdiff\nEMPTY=\n"
	p := writeTempEnvProfile(t, content)
	var buf bytes.Buffer
	if err := cli.Run([]string{"profile", p}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Total keys") {
		t.Error("expected 'Total keys' in output")
	}
	if !strings.Contains(out, "Sensitive keys") {
		t.Error("expected 'Sensitive keys' in output")
	}
}

func TestRun_ProfileJSONOutput(t *testing.T) {
	content := "API_TOKEN=tok\nAPI_KEY=key\nSVC_URL=http://example.com\n"
	p := writeTempEnvProfile(t, content)
	var buf bytes.Buffer
	if err := cli.Run([]string{"profile", p, "--json"}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := result["total_keys"]; !ok {
		t.Error("expected total_keys in JSON")
	}
	if _, ok := result["sensitive_keys"]; !ok {
		t.Error("expected sensitive_keys in JSON")
	}
}

func TestRun_ProfilePrefixCounts(t *testing.T) {
	content := "DB_HOST=h\nDB_PORT=5432\nDB_NAME=mydb\nAPP_ENV=prod\n"
	p := writeTempEnvProfile(t, content)
	var buf bytes.Buffer
	if err := cli.Run([]string{"profile", p, "--json"}, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result struct {
		PrefixCounts map[string]int `json:"prefix_counts"`
	}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("json decode: %v", err)
	}
	if result.PrefixCounts["DB"] != 3 {
		t.Errorf("expected DB=3, got %d", result.PrefixCounts["DB"])
	}
}
