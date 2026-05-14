package redactor_test

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/redactor"
)

func makeFile(entries ...parser.Entry) parser.EnvFile {
	return parser.EnvFile{Path: "test.env", Entries: entries}
}

func entry(k, v string) parser.Entry { return parser.Entry{Key: k, Value: v} }

func TestApply_NonSensitiveUnchanged(t *testing.T) {
	f := makeFile(entry("APP_NAME", "myapp"), entry("PORT", "8080"))
	out := redactor.Apply(f, redactor.Options{})
	for _, e := range out.Entries {
		if e.Masked {
			t.Errorf("expected %q to be unmasked", e.Key)
		}
		if e.Value == "[REDACTED]" {
			t.Errorf("expected %q value to be unchanged", e.Key)
		}
	}
}

func TestApply_SensitiveKeysRedacted(t *testing.T) {
	sensitive := []string{"DB_PASSWORD", "API_KEY", "AUTH_TOKEN", "PRIVATE_KEY"}
	var entries []parser.Entry
	for _, k := range sensitive {
		entries = append(entries, entry(k, "supersecret"))
	}
	f := makeFile(entries...)
	out := redactor.Apply(f, redactor.Options{})
	for _, e := range out.Entries {
		if e.Value != "[REDACTED]" {
			t.Errorf("key %q: expected [REDACTED], got %q", e.Key, e.Value)
		}
		if !e.Masked {
			t.Errorf("key %q: expected Masked=true", e.Key)
		}
	}
}

func TestApply_CustomPlaceholder(t *testing.T) {
	f := makeFile(entry("SECRET_KEY", "abc123"))
	out := redactor.Apply(f, redactor.Options{Placeholder: "***"})
	if out.Entries[0].Value != "***" {
		t.Errorf("expected ***, got %q", out.Entries[0].Value)
	}
}

func TestApply_CustomPatterns(t *testing.T) {
	f := makeFile(entry("STRIPE_LIVE", "sk_live_xxx"), entry("APP_ENV", "prod"))
	out := redactor.Apply(f, redactor.Options{Patterns: []string{"live"}})
	if out.Entries[0].Value != "[REDACTED]" {
		t.Errorf("expected STRIPE_LIVE to be redacted")
	}
	if out.Entries[1].Value != "prod" {
		t.Errorf("expected APP_ENV to be unchanged")
	}
}

func TestApply_OriginalUnmodified(t *testing.T) {
	f := makeFile(entry("DB_PASSWORD", "original"))
	redactor.Apply(f, redactor.Options{})
	if f.Entries[0].Value != "original" {
		t.Error("Apply must not mutate the original EnvFile")
	}
}
