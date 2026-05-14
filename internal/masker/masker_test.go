package masker_test

import (
	"testing"

	"envdiff/internal/masker"
)

func TestIsSensitive_DefaultPatterns(t *testing.T) {
	m := masker.New()

	sensitiveKeys := []string{
		"DB_PASSWORD",
		"API_KEY",
		"AWS_SECRET_ACCESS_KEY",
		"AUTH_TOKEN",
		"PRIVATE_KEY",
		"SIGNING_KEY",
	}
	for _, key := range sensitiveKeys {
		if !m.IsSensitive(key) {
			t.Errorf("expected key %q to be sensitive", key)
		}
	}

	safeKeys := []string{
		"APP_ENV",
		"PORT",
		"DATABASE_HOST",
		"LOG_LEVEL",
	}
	for _, key := range safeKeys {
		if m.IsSensitive(key) {
			t.Errorf("expected key %q to NOT be sensitive", key)
		}
	}
}

func TestIsSensitive_CustomPatterns(t *testing.T) {
	m := masker.New("CUSTOM_SECRET", "INTERNAL")

	if !m.IsSensitive("MY_CUSTOM_SECRET_VALUE") {
		t.Error("expected CUSTOM_SECRET pattern to match")
	}
	if !m.IsSensitive("INTERNAL_KEY") {
		t.Error("expected INTERNAL pattern to match")
	}
	if m.IsSensitive("API_KEY") {
		t.Error("default patterns should not apply with custom patterns")
	}
}

func TestIsSensitive_CaseInsensitive(t *testing.T) {
	m := masker.New()

	// Keys should be matched regardless of case
	cases := []struct {
		key       string
		wantMatch bool
	}{
		{"db_password", true},
		{"Api_Key", true},
		{"aws_secret_access_key", true},
		{"app_env", false},
		{"port", false},
	}
	for _, tc := range cases {
		got := m.IsSensitive(tc.key)
		if got != tc.wantMatch {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.wantMatch)
		}
	}
}

func TestMaskEnv(t *testing.T) {
	m := masker.New()

	env := map[string]string{
		"APP_ENV":     "production",
		"DB_PASSWORD": "supersecret",
		"PORT":        "8080",
		"API_KEY":     "abc123",
	}

	masked := m.MaskEnv(env)

	if masked["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV to be unchanged, got %q", masked["APP_ENV"])
	}
	if masked["PORT"] != "8080" {
		t.Errorf("expected PORT to be unchanged, got %q", masked["PORT"])
	}
	if masked["DB_PASSWORD"] != masker.MaskedValue {
		t.Errorf("expected DB_PASSWORD to be masked, got %q", masked["DB_PASSWORD"])
	}
	if masked["API_KEY"] != masker.MaskedValue {
		t.Errorf("expected API_KEY to be masked, got %q", masked["API_KEY"])
	}

	// Ensure original is not mutated
	if env["DB_PASSWORD"] != "supersecret" {
		t.Error("original env map should not be mutated")
	}
}

func TestMaskValue(t *testing.T) {
	m := masker.New()

	if got := m.MaskValue("API_KEY", "abc123"); got != masker.MaskedValue {
		t.Errorf("expected masked value, got %q", got)
	}
	if got := m.MaskValue("APP_ENV", "staging"); got != "staging" {
		t.Errorf("expected original value, got %q", got)
	}
}
