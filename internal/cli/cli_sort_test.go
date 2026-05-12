package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/cli"
)

func writeTempEnvSort(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRun_SortMissingArg(t *testing.T) {
	err := cli.Run([]string{"sort"})
	if err == nil || !strings.Contains(err.Error(), "missing required argument") {
		t.Errorf("expected missing-arg error, got %v", err)
	}
}

func TestRun_SortMissingFile(t *testing.T) {
	err := cli.Run([]string{"sort", "/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRun_SortAlphabetical(t *testing.T) {
	p := writeTempEnvSort(t, "ZEBRA=1\nAPPLE=2\nMANGO=3\n")
	err := cli.Run([]string{"sort", "-alpha", p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_SortGroupByPrefix(t *testing.T) {
	p := writeTempEnvSort(t, "DB_HOST=localhost\nAPP_NAME=myapp\nDB_PORT=5432\nAPP_ENV=prod\n")
	err := cli.Run([]string{"sort", "-group", p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_SortOutputFile(t *testing.T) {
	p := writeTempEnvSort(t, "Z=1\nA=2\n")
	out := filepath.Join(t.TempDir(), "sorted.env")
	err := cli.Run([]string{"sort", "-alpha", "-output", out, p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if !strings.Contains(string(data), "A=") {
		t.Errorf("output missing expected key, got: %s", data)
	}
}
