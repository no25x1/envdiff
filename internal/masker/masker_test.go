package masker_test

import (
	"regexp"
	"testing"

	"github.com/user/envdiff/internal/masker"
	"github.com/user/envdiff/internal/parser"
)

func TestIsSensitive_DefaultPatterns(t *testing.T) {
	m := masker.New()
	sensitive := []string{"DB_PASSWORD", "API_SECRET", "AUTH_TOKEN", "API_KEY", "PRIVATE_KEY", "AUTH_HEADER", "CREDENTIALS"}
	for _, k := range sensitive {
		if !m.IsSensitive(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}
}

func TestIsSensitive_CustomPatterns(t *testing.T) {
	m := masker.NewWithOptions([]*regexp.Regexp{
		regexp.MustCompile(`(?i)custom`),
	}, "")
	if !m.IsSensitive("MY_CUSTOM_VAR") {
		t.Error("expected MY_CUSTOM_VAR to be sensitive")
	}
	if m.IsSensitive("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD not to be sensitive with custom patterns")
	}
}

func TestIsSensitive_CaseInsensitive(t *testing.T) {
	m := masker.New()
	if !m.IsSensitive("db_password") {
		t.Error("expected lowercase db_password to be sensitive")
	}
}

func TestMaskValue(t *testing.T) {
	m := masker.New()
	if got := m.MaskValue("DB_PASSWORD", "secret123"); got != "***" {
		t.Errorf("expected ***, got %q", got)
	}
	if got := m.MaskValue("APP_ENV", "production"); got != "production" {
		t.Errorf("expected production, got %q", got)
	}
}

func TestMaskEnv(t *testing.T) {
	m := masker.New()
	f := &parser.EnvFile{
		Path: ".env",
		Entries: []parser.Entry{
			{Key: "APP_NAME", Value: "myapp"},
			{Key: "DB_PASSWORD", Value: "supersecret"},
			{Key: "API_TOKEN", Value: "tok_abc123"},
		},
	}
	out := m.MaskEnv(f)
	if len(out.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out.Entries))
	}
	if out.Entries[0].Value != "myapp" {
		t.Errorf("expected myapp, got %q", out.Entries[0].Value)
	}
	if out.Entries[1].Value != "***" {
		t.Errorf("expected ***, got %q", out.Entries[1].Value)
	}
	if out.Entries[2].Value != "***" {
		t.Errorf("expected ***, got %q", out.Entries[2].Value)
	}
}

func TestMaskEnv_NilFile(t *testing.T) {
	m := masker.New()
	out := m.MaskEnv(nil)
	if out == nil {
		t.Fatal("expected non-nil result for nil input")
	}
	if len(out.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(out.Entries))
	}
}

func TestNewWithOptions_CustomPlaceholder(t *testing.T) {
	m := masker.NewWithOptions(nil, "[REDACTED]")
	if got := m.MaskValue("DB_PASSWORD", "secret"); got != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", got)
	}
}
