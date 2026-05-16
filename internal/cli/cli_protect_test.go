package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envdiff/internal/cli"
)

func writeTempEnvProtect(t *testing.T, content string) string {
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

func TestRun_ProtectMissingArgs(t *testing.T) {
	err := cli.Run([]string{"protect"})
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Fatalf("expected usage error, got %v", err)
	}
}

func TestRun_ProtectMissingBaseFile(t *testing.T) {
	err := cli.Run([]string{"protect", filepath.Join(t.TempDir(), "missing.env"), "/dev/null"})
	if err == nil {
		t.Fatal("expected error for missing base file")
	}
}

func TestRun_ProtectNoViolations(t *testing.T) {
	base := writeTempEnvProtect(t, "SECRET=same\n")
	next := writeTempEnvProtect(t, "SECRET=same\n")
	err := cli.Run([]string{"protect", base, next, "--keys=SECRET"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRun_ProtectViolationReturnsError(t *testing.T) {
	base := writeTempEnvProtect(t, "SECRET=old\n")
	next := writeTempEnvProtect(t, "SECRET=new\n")
	err := cli.Run([]string{"protect", base, next, "--keys=SECRET"})
	if err == nil {
		t.Fatal("expected violation error")
	}
}

func TestRun_ProtectAllowOverrideNoError(t *testing.T) {
	base := writeTempEnvProtect(t, "SECRET=old\n")
	next := writeTempEnvProtect(t, "SECRET=new\n")
	err := cli.Run([]string{"protect", base, next, "--keys=SECRET", "--allow-override"})
	if err != nil {
		t.Fatalf("expected no error with allow-override, got %v", err)
	}
}

func TestRun_ProtectJSONOutput(t *testing.T) {
	base := writeTempEnvProtect(t, "SECRET=old\n")
	next := writeTempEnvProtect(t, "SECRET=new\n")
	err := cli.Run([]string{"protect", base, next, "--keys=SECRET", "--allow-override", "--format=json"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
