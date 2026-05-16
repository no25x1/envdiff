package protector_test

import (
	"testing"

	"envdiff/internal/parser"
	"envdiff/internal/protector"
)

func makeFile(entries ...parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(k, v string) parser.Entry { return parser.Entry{Key: k, Value: v} }

func TestApply_NilFiles_ReturnsNil(t *testing.T) {
	v, err := protector.Apply(nil, nil, protector.Options{})
	if err != nil || len(v) != 0 {
		t.Fatalf("expected no violations, got %v / %v", v, err)
	}
}

func TestApply_NoProtectedKeys_NoViolations(t *testing.T) {
	base := makeFile(entry("FOO", "bar"))
	next := makeFile(entry("FOO", "changed"))
	v, err := protector.Apply(base, next, protector.Options{})
	if err != nil || len(v) != 0 {
		t.Fatalf("expected no violations, got %v", v)
	}
}

func TestApply_ExplicitKey_ValueChanged(t *testing.T) {
	base := makeFile(entry("SECRET", "old"))
	next := makeFile(entry("SECRET", "new"))
	v, err := protector.Apply(base, next, protector.Options{Keys: []string{"SECRET"}})
	if err == nil {
		t.Fatal("expected error for changed protected key")
	}
	if len(v) != 1 || v[0].Key != "SECRET" {
		t.Fatalf("unexpected violations: %v", v)
	}
}

func TestApply_ExplicitKey_Removed(t *testing.T) {
	base := makeFile(entry("SECRET", "val"))
	next := makeFile(entry("OTHER", "val"))
	v, err := protector.Apply(base, next, protector.Options{Keys: []string{"SECRET"}})
	if err == nil {
		t.Fatal("expected error for removed protected key")
	}
	if len(v) != 1 || v[0].Reason != "key removed" {
		t.Fatalf("unexpected violations: %v", v)
	}
}

func TestApply_PrefixProtection_DetectsChange(t *testing.T) {
	base := makeFile(entry("PROD_DB_PASS", "secret"), entry("PROD_API_KEY", "key"))
	next := makeFile(entry("PROD_DB_PASS", "secret"), entry("PROD_API_KEY", "changed"))
	v, err := protector.Apply(base, next, protector.Options{Prefixes: []string{"PROD_"}})
	if err == nil {
		t.Fatal("expected violation for prefix-protected key")
	}
	if len(v) != 1 || v[0].Key != "PROD_API_KEY" {
		t.Fatalf("unexpected violations: %v", v)
	}
}

func TestApply_AllowOverride_NoError(t *testing.T) {
	base := makeFile(entry("SECRET", "old"))
	next := makeFile(entry("SECRET", "new"))
	v, err := protector.Apply(base, next, protector.Options{
		Keys:          []string{"SECRET"},
		AllowOverride: true,
	})
	if err != nil {
		t.Fatalf("expected no error with AllowOverride, got %v", err)
	}
	if len(v) != 1 {
		t.Fatalf("expected 1 violation (warning), got %d", len(v))
	}
}

func TestApply_UnchangedProtectedKey_NoViolation(t *testing.T) {
	base := makeFile(entry("SECRET", "same"))
	next := makeFile(entry("SECRET", "same"))
	v, err := protector.Apply(base, next, protector.Options{Keys: []string{"SECRET"}})
	if err != nil || len(v) != 0 {
		t.Fatalf("expected no violations for unchanged key, got %v / %v", v, err)
	}
}
