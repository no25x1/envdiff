package masked_test

import (
	"strings"
	"testing"

	"envdiff/internal/diff"
	"envdiff/internal/masked"
	"envdiff/internal/parser"
)

func makeEnvFile(entries []parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(k, v string) parser.Entry { return parser.Entry{Key: k, Value: v} }

func TestCompare_NilBase_ReturnsError(t *testing.T) {
	_, err := masked.Compare(nil, makeEnvFile(nil), masked.Options{})
	if err == nil {
		t.Fatal("expected error for nil base")
	}
}

func TestCompare_NilNext_ReturnsError(t *testing.T) {
	_, err := masked.Compare(makeEnvFile(nil), nil, masked.Options{})
	if err == nil {
		t.Fatal("expected error for nil next")
	}
}

func TestCompare_NonSensitiveKey_NotMasked(t *testing.T) {
	base := makeEnvFile([]parser.Entry{entry("APP_NAME", "myapp")})
	next := makeEnvFile([]parser.Entry{entry("APP_NAME", "myapp")})
	res, err := masked.Compare(base, next, masked.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if res[0].Masked {
		t.Error("expected non-sensitive key to not be masked")
	}
	if res[0].BaseVal != "myapp" {
		t.Errorf("unexpected base value: %s", res[0].BaseVal)
	}
}

func TestCompare_SensitiveKey_ValueMasked(t *testing.T) {
	base := makeEnvFile([]parser.Entry{entry("DB_PASSWORD", "s3cr3t")})
	next := makeEnvFile([]parser.Entry{entry("DB_PASSWORD", "n3ws3cr3t")})
	res, err := masked.Compare(base, next, masked.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !res[0].Masked {
		t.Error("expected sensitive key to be masked")
	}
	if res[0].BaseVal == "s3cr3t" || res[0].NextVal == "n3ws3cr3t" {
		t.Error("sensitive values should be masked")
	}
	if res[0].Status != diff.StatusModified {
		t.Errorf("expected modified status, got %v", res[0].Status)
	}
}

func TestCompare_ExtraPattern_Masked(t *testing.T) {
	base := makeEnvFile([]parser.Entry{entry("MY_CUSTOM_TOKEN", "abc123")})
	next := makeEnvFile([]parser.Entry{entry("MY_CUSTOM_TOKEN", "abc123")})
	res, err := masked.Compare(base, next, masked.Options{ExtraPatterns: []string{"token"}})
	if err != nil {
		t.Fatal(err)
	}
	if !res[0].Masked {
		t.Error("expected custom pattern key to be masked")
	}
}

func TestFormat_Added(t *testing.T) {
	r := masked.Result{Key: "FOO", Status: diff.StatusAdded, NextVal: "bar"}
	line := masked.Format(r)
	if !strings.HasPrefix(line, "+ FOO=") {
		t.Errorf("unexpected format: %s", line)
	}
}

func TestFormat_Masked_AppendsSuffix(t *testing.T) {
	r := masked.Result{Key: "API_KEY", Status: diff.StatusUnchanged, BaseVal: "***", Masked: true}
	line := masked.Format(r)
	if !strings.Contains(line, "[masked]") {
		t.Errorf("expected [masked] suffix, got: %s", line)
	}
}
