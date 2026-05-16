package masker_test

import (
	"testing"

	"github.com/user/envdiff/internal/masker"
	"github.com/user/envdiff/internal/parser"
)

func TestIsSensitive_DefaultPatterns(t *testing.T) {
	m := masker.New()
	sensitive := []string{"DB_PASSWORD", "API_SECRET", "AUTH_TOKEN", "PRIVATE_KEY", "API_KEY", "AWS_CREDENTIAL"}
	for _, key := range sensitive {
		if !m.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_CustomPatterns(t *testing.T) {
	m := masker.NewWithOptions([]string{`(?i)mysecret`}, "REDACTED")
	if !m.IsSensitive("MY_MYSECRET_KEY") {
		t.Error("expected custom pattern to match")
	}
	if m.IsSensitive("DB_PASSWORD") {
		t.Error("default pattern should not match with custom-only masker")
	}
}

func TestIsSensitive_CaseInsensitive(t *testing.T) {
	m := masker.New()
	if !m.IsSensitive("db_password") {
		t.Error("expected lowercase key to match")
	}
	if !m.IsSensitive("DB_PASSWORD") {
		t.Error("expected uppercase key to match")
	}
}

func TestMaskValue(t *testing.T) {
	m := masker.New()
	if got := m.MaskValue("DB_PASSWORD", "hunter2"); got != masker.DefaultPlaceholder {
		t.Errorf("expected placeholder, got %q", got)
	}
	if got := m.MaskValue("APP_NAME", "myapp"); got != "myapp" {
		t.Errorf("expected original value, got %q", got)
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
	masked := m.MaskEnv(f)
	if masked == nil {
		t.Fatal("expected non-nil result")
	}
	if len(masked.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(masked.Entries))
	}
	if masked.Entries[0].Value != "myapp" {
		t.Errorf("expected APP_NAME unchanged, got %q", masked.Entries[0].Value)
	}
	if masked.Entries[1].Value != masker.DefaultPlaceholder {
		t.Errorf("expected DB_PASSWORD masked, got %q", masked.Entries[1].Value)
	}
	if masked.Entries[2].Value != masker.DefaultPlaceholder {
		t.Errorf("expected API_TOKEN masked, got %q", masked.Entries[2].Value)
	}
	// original should be unchanged
	if f.Entries[1].Value != "supersecret" {
		t.Error("original file should not be modified")
	}
}

func TestMaskEnv_NilFile(t *testing.T) {
	m := masker.New()
	if got := m.MaskEnv(nil); got != nil {
		t.Errorf("expected nil for nil input, got %v", got)
	}
}
