package filter_test

import (
	"testing"

	"github.com/user/envdiff/internal/filter"
	"github.com/user/envdiff/internal/parser"
)

func makeFile(entries ...parser.Entry) parser.EnvFile {
	return parser.EnvFile{Path: "test.env", Entries: entries}
}

func entry(key, value string) parser.Entry {
	return parser.Entry{Key: key, Value: value}
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	file := makeFile(entry("DB_HOST", "localhost"), entry("APP_PORT", "8080"))
	out, err := filter.Apply(file, filter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out.Entries))
	}
}

func TestApply_Prefix(t *testing.T) {
	file := makeFile(entry("DB_HOST", "localhost"), entry("DB_PORT", "5432"), entry("APP_PORT", "8080"))
	out, err := filter.Apply(file, filter.Options{Prefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out.Entries))
	}
}

func TestApply_Suffix(t *testing.T) {
	file := makeFile(entry("DB_HOST", "localhost"), entry("APP_HOST", "example.com"), entry("APP_PORT", "8080"))
	out, err := filter.Apply(file, filter.Options{Suffix: "_HOST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out.Entries))
	}
}

func TestApply_Pattern(t *testing.T) {
	file := makeFile(entry("DB_HOST", "localhost"), entry("APP_SECRET", "abc"), entry("APP_PORT", "8080"))
	out, err := filter.Apply(file, filter.Options{Pattern: "^APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out.Entries))
	}
}

func TestApply_InvalidPattern(t *testing.T) {
	file := makeFile(entry("KEY", "val"))
	_, err := filter.Apply(file, filter.Options{Pattern: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}

func TestApply_Exclude(t *testing.T) {
	file := makeFile(entry("DB_HOST", "localhost"), entry("APP_PORT", "8080"), entry("APP_SECRET", "xyz"))
	out, err := filter.Apply(file, filter.Options{Prefix: "DB_", Exclude: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Errorf("expected 2 entries after exclusion, got %d", len(out.Entries))
	}
	for _, e := range out.Entries {
		if e.Key == "DB_HOST" {
			t.Error("excluded key DB_HOST should not be present")
		}
	}
}

func TestApply_PreservesPath(t *testing.T) {
	file := parser.EnvFile{Path: "prod.env", Entries: []parser.Entry{entry("KEY", "val")}}
	out, err := filter.Apply(file, filter.Options{Prefix: "KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Path != "prod.env" {
		t.Errorf("expected path prod.env, got %s", out.Path)
	}
}
