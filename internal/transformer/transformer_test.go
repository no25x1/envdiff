package transformer_test

import (
	"testing"

	"github.com/envdiff/envdiff/internal/parser"
	"github.com/envdiff/envdiff/internal/transformer"
)

func makeFile(pairs ...string) parser.EnvFile {
	f := parser.EnvFile{Path: "test.env"}
	for i := 0; i+1 < len(pairs); i += 2 {
		f.Entries = append(f.Entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return f
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	f := makeFile("FOO", "bar", "BAZ", "qux")
	out := transformer.Apply(f, transformer.Options{})
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out.Entries))
	}
}

func TestApply_PrefixAdd(t *testing.T) {
	f := makeFile("HOST", "localhost", "PORT", "5432")
	out := transformer.Apply(f, transformer.Options{PrefixAdd: "APP_"})
	if out.Entries[0].Key != "APP_HOST" {
		t.Errorf("expected APP_HOST, got %s", out.Entries[0].Key)
	}
	if out.Entries[1].Key != "APP_PORT" {
		t.Errorf("expected APP_PORT, got %s", out.Entries[1].Key)
	}
}

func TestApply_PrefixStrip_DropsNonMatching(t *testing.T) {
	f := makeFile("APP_HOST", "localhost", "DB_PORT", "5432", "APP_NAME", "myapp")
	out := transformer.Apply(f, transformer.Options{PrefixStrip: "APP_"})
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries after strip, got %d", len(out.Entries))
	}
	if out.Entries[0].Key != "HOST" {
		t.Errorf("expected HOST, got %s", out.Entries[0].Key)
	}
}

func TestApply_KeyToUpper(t *testing.T) {
	f := makeFile("db_host", "localhost")
	out := transformer.Apply(f, transformer.Options{KeyToUpper: true})
	if out.Entries[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", out.Entries[0].Key)
	}
}

func TestApply_KeyToLower(t *testing.T) {
	f := makeFile("DB_HOST", "localhost")
	out := transformer.Apply(f, transformer.Options{KeyToLower: true})
	if out.Entries[0].Key != "db_host" {
		t.Errorf("expected db_host, got %s", out.Entries[0].Key)
	}
}

func TestApply_CustomTransform_DropsEntry(t *testing.T) {
	f := makeFile("SECRET", "abc", "HOST", "localhost")
	dropSecret := func(k, v string) (string, string, bool) {
		if k == "SECRET" {
			return "", "", false
		}
		return k, v, true
	}
	out := transformer.Apply(f, transformer.Options{Custom: []transformer.TransformFunc{dropSecret}})
	if len(out.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out.Entries))
	}
	if out.Entries[0].Key != "HOST" {
		t.Errorf("expected HOST, got %s", out.Entries[0].Key)
	}
}

func TestApply_PrefixAddAndStrip_Combined(t *testing.T) {
	f := makeFile("APP_KEY", "val")
	out := transformer.Apply(f, transformer.Options{PrefixStrip: "APP_", PrefixAdd: "NEW_"})
	if out.Entries[0].Key != "NEW_KEY" {
		t.Errorf("expected NEW_KEY, got %s", out.Entries[0].Key)
	}
}
