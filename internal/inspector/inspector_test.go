package inspector_test

import (
	"testing"

	"envdiff/internal/inspector"
	"envdiff/internal/parser"
)

func makeFile(entries ...parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(key, value string) parser.Entry {
	return parser.Entry{Key: key, Value: value}
}

func TestInspect_NilFile(t *testing.T) {
	res, err := inspector.Inspect(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != nil {
		t.Fatalf("expected nil result for nil file")
	}
}

func TestInspect_EmptyValue(t *testing.T) {
	f := makeFile(entry("FOO", ""))
	res, err := inspector.Inspect(f)
	if err != nil {
		t.Fatal(err)
	}
	if !res[0].IsEmpty {
		t.Error("expected IsEmpty=true")
	}
	if res[0].TypeHint != "empty" {
		t.Errorf("expected type hint 'empty', got %q", res[0].TypeHint)
	}
}

func TestInspect_SensitiveKey(t *testing.T) {
	f := makeFile(entry("DB_PASSWORD", "s3cr3t"))
	res, _ := inspector.Inspect(f)
	if !res[0].IsSensitive {
		t.Error("expected IsSensitive=true for DB_PASSWORD")
	}
}

func TestInspect_NonSensitiveKey(t *testing.T) {
	f := makeFile(entry("APP_NAME", "myapp"))
	res, _ := inspector.Inspect(f)
	if res[0].IsSensitive {
		t.Error("expected IsSensitive=false for APP_NAME")
	}
}

func TestInspect_TypeHint_Bool(t *testing.T) {
	f := makeFile(entry("FEATURE_FLAG", "true"))
	res, _ := inspector.Inspect(f)
	if res[0].TypeHint != "bool" {
		t.Errorf("expected 'bool', got %q", res[0].TypeHint)
	}
}

func TestInspect_TypeHint_Number(t *testing.T) {
	f := makeFile(entry("PORT", "8080"))
	res, _ := inspector.Inspect(f)
	if res[0].TypeHint != "number" {
		t.Errorf("expected 'number', got %q", res[0].TypeHint)
	}
}

func TestInspect_TypeHint_URL(t *testing.T) {
	f := makeFile(entry("BASE_URL", "https://example.com"))
	res, _ := inspector.Inspect(f)
	if res[0].TypeHint != "url" {
		t.Errorf("expected 'url', got %q", res[0].TypeHint)
	}
}

func TestInspect_LengthPopulated(t *testing.T) {
	f := makeFile(entry("VAR", "hello"))
	res, _ := inspector.Inspect(f)
	if res[0].Length != 5 {
		t.Errorf("expected length 5, got %d", res[0].Length)
	}
}
