package cli_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envdiff/internal/cli"
)

func writeTempEnvMasked(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRun_MaskedDiffMissingArgs(t *testing.T) {
	var buf bytes.Buffer
	err := cli.Run([]string{"maskeddiff"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRun_MaskedDiffMissingFile(t *testing.T) {
	var buf bytes.Buffer
	err := cli.Run([]string{"maskeddiff", "/no/such/base.env", "/no/such/next.env"}, &buf)
	if err == nil {
		t.Fatal("expected error for missing files")
	}
}

func TestRun_MaskedDiff_TextOutput(t *testing.T) {
	base := writeTempEnvMasked(t, "APP_NAME=myapp\nDB_PASSWORD=secret\n")
	next := writeTempEnvMasked(t, "APP_NAME=myapp\nDB_PASSWORD=newsecret\nNEW_KEY=hello\n")
	var buf bytes.Buffer
	err := cli.Run([]string{"maskeddiff", base, next}, &buf)
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Error("expected DB_PASSWORD in output")
	}
	if strings.Contains(out, "newsecret") {
		t.Error("sensitive value should be masked in output")
	}
	if !strings.Contains(out, "NEW_KEY") {
		t.Error("expected NEW_KEY in output")
	}
}

func TestRun_MaskedDiff_JSONOutput(t *testing.T) {
	base := writeTempEnvMasked(t, "API_KEY=abc\n")
	next := writeTempEnvMasked(t, "API_KEY=xyz\n")
	var buf bytes.Buffer
	err := cli.Run([]string{"maskeddiff", base, next, "--format=json"}, &buf)
	if err != nil {
		t.Fatal(err)
	}
	var rows []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &rows); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(rows) == 0 {
		t.Fatal("expected at least one result")
	}
	if rows[0]["masked"] != true {
		t.Error("expected API_KEY to be masked")
	}
}

func TestRun_MaskedDiff_ExtraPattern(t *testing.T) {
	base := writeTempEnvMasked(t, "MY_CUSTOM_CRED=val1\n")
	next := writeTempEnvMasked(t, "MY_CUSTOM_CRED=val2\n")
	var buf bytes.Buffer
	err := cli.Run([]string{"maskeddiff", base, next, "--extra-patterns=cred"}, &buf)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(buf.String(), "val1") || strings.Contains(buf.String(), "val2") {
		t.Error("custom pattern values should be masked")
	}
}
