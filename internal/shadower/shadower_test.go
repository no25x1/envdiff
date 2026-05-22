package shadower_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/parser"
	"github.com/your-org/envdiff/internal/shadower"
)

func makeFile(entries ...parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(key, value string) parser.Entry {
	return parser.Entry{Key: key, Value: value}
}

func TestApply_NilBase_ReturnsError(t *testing.T) {
	_, err := shadower.Apply(nil, makeFile(), shadower.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for nil base")
	}
}

func TestApply_NilNext_ReturnsError(t *testing.T) {
	_, err := shadower.Apply(makeFile(), nil, shadower.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for nil next")
	}
}

func TestApply_NoShadows(t *testing.T) {
	base := makeFile(entry("FOO", "bar"))
	next := makeFile(entry("BAZ", "qux"))
	results, err := shadower.Apply(base, next, shadower.DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 shadows, got %d", len(results))
	}
}

func TestApply_DetectsShadow(t *testing.T) {
	base := makeFile(entry("DB_HOST", "localhost"), entry("PORT", "5432"))
	next := makeFile(entry("DB_HOST", "prod.db.example.com"), entry("PORT", "5432"))
	results, err := shadower.Apply(base, next, shadower.DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 shadow, got %d", len(results))
	}
	if results[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", results[0].Key)
	}
}

func TestApply_IgnoreEqual_False_IncludesIdentical(t *testing.T) {
	base := makeFile(entry("PORT", "5432"))
	next := makeFile(entry("PORT", "5432"))
	opts := shadower.Options{IgnoreEqual: false}
	results, err := shadower.Apply(base, next, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 shadow, got %d", len(results))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	base := makeFile(entry("DB_HOST", "localhost"), entry("APP_NAME", "old"))
	next := makeFile(entry("DB_HOST", "prod"), entry("APP_NAME", "new"))
	opts := shadower.Options{Prefix: "APP_", IgnoreEqual: true}
	results, err := shadower.Apply(base, next, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Key != "APP_NAME" {
		t.Fatalf("expected only APP_NAME shadow, got %+v", results)
	}
}
