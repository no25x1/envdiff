package masker_test

import (
	"testing"

	"github.com/user/envdiff/internal/masker"
	"github.com/user/envdiff/internal/parser"
)

func TestIsSensitive_DefaultPatterns(t *testing.T) {
	m := masker.New()
	sensitive := []string{"DB_PASSWORD", "API_KEY", "AUTH_TOKEN", "SECRET", "PRIVATE_KEY", "CREDENTIALS"}
	for _, k := range sensitive {
		if !m.IsSensitive(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}
}

func TestIsSensitive_CustomPatterns(t *testing.T) {
	import_regexp := func(s string) interface{} { return s }
	_ = import_regexp

	var customOpts masker.Options
	import (
		"regexp"
	)
	customOpts.Patterns = []*regexp.Regexp{regexp.MustCompile(`(?i)mysecret`)}
	customOpts.Placeholder = "REDACTED"
	m := masker.NewWithOptions(customOpts)

	if !m.IsSensitive("MYSECRET_VALUE") {
		t.Error("expected custom pattern to match")
	}
	if m.IsSensitive("PASSWORD") {
		t.Error("default patterns should not apply with custom options")
	}
}

func TestIsSensitive_CaseInsensitive(t *testing.T) {
	m := masker.New()
	if !m.IsSensitive("db_password") {
		t.Error("expected lowercase key to match")
	}
}

func TestMaskValue(t *testing.T) {
	m := masker.New()
	if got := m.MaskValue("DB_PASSWORD", "hunter2"); got != masker.DefaultPlaceholder {
		t.Errorf("expected placeholder, got %q", got)
	}
	if got := m.MaskValue("HOST", "localhost"); got != "localhost" {
		t.Errorf("expected value unchanged, got %q", got)
	}
}

func TestMaskEnv(t *testing.T) {
	m := masker.New()
	f := &parser.EnvFile{
		Path: ".env",
		Entries: []parser.Entry{
			{Key: "HOST", Value: "localhost"},
			{Key: "DB_PASSWORD", Value: "secret123"},
			{Key: "PORT", Value: "5432"},
		},
	}
	out := m.MaskEnv(f)
	if out == nil {
		t.Fatal("expected non-nil result")
	}
	if len(out.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out.Entries))
	}
	for _, e := range out.Entries {
		if e.Key == "DB_PASSWORD" && e.Value != masker.DefaultPlaceholder {
			t.Errorf("expected password masked, got %q", e.Value)
		}
		if e.Key == "HOST" && e.Value != "localhost" {
			t.Errorf("expected host unchanged, got %q", e.Value)
		}
	}
}

func TestMaskEnv_NilFile(t *testing.T) {
	m := masker.New()
	if got := m.MaskEnv(nil); got != nil {
		t.Errorf("expected nil for nil input, got %v", got)
	}
}

func TestNewWithOptions_DefaultPlaceholder(t *testing.T) {
	m := masker.NewWithOptions(masker.Options{})
	if got := m.MaskValue("API_KEY", "abc"); got != masker.DefaultPlaceholder {
		t.Errorf("expected default placeholder, got %q", got)
	}
}
