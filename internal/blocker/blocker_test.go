package blocker_test

import (
	"testing"

	"envdiff/internal/blocker"
	"envdiff/internal/parser"
)

func makeFile(entries ...parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(key, value string) parser.Entry {
	return parser.Entry{Key: key, Value: value}
}

func TestApply_NilFile_ReturnsNil(t *testing.T) {
	v, err := blocker.Apply(nil, blocker.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != nil {
		t.Fatalf("expected nil violations, got %v", v)
	}
}

func TestApply_NoViolations(t *testing.T) {
	f := makeFile(entry("APP_NAME", "myapp"), entry("PORT", "8080"))
	v, err := blocker.Apply(f, blocker.Options{
		ForbiddenKeys: []string{"SECRET"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %v", v)
	}
}

func TestApply_ForbiddenKey(t *testing.T) {
	f := makeFile(entry("SECRET", "abc"), entry("PORT", "8080"))
	v, _ := blocker.Apply(f, blocker.Options{
		ForbiddenKeys: []string{"SECRET"},
	})
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Key != "SECRET" {
		t.Errorf("expected key SECRET, got %q", v[0].Key)
	}
}

func TestApply_ForbiddenPrefix(t *testing.T) {
	f := makeFile(entry("INTERNAL_HOST", "localhost"), entry("PORT", "9000"))
	v, _ := blocker.Apply(f, blocker.Options{
		ForbiddenPrefixes: []string{"INTERNAL_"},
	})
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Key != "INTERNAL_HOST" {
		t.Errorf("unexpected key: %q", v[0].Key)
	}
}

func TestApply_ForbiddenValue(t *testing.T) {
	f := makeFile(entry("DB_PASS", "changeme"), entry("API_KEY", "real-key"))
	v, _ := blocker.Apply(f, blocker.Options{
		ForbiddenValues: []string{"changeme"},
	})
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Key != "DB_PASS" {
		t.Errorf("unexpected key: %q", v[0].Key)
	}
}

func TestApply_MultipleViolations(t *testing.T) {
	f := makeFile(
		entry("SECRET", "abc"),
		entry("INTERNAL_X", "val"),
		entry("TOKEN", "changeme"),
	)
	v, _ := blocker.Apply(f, blocker.Options{
		ForbiddenKeys:     []string{"SECRET"},
		ForbiddenPrefixes: []string{"INTERNAL_"},
		ForbiddenValues:   []string{"changeme"},
	})
	if len(v) != 3 {
		t.Fatalf("expected 3 violations, got %d", len(v))
	}
}
