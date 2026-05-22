package counter_test

import (
	"testing"

	"github.com/envdiff/envdiff/internal/counter"
	"github.com/envdiff/envdiff/internal/parser"
)

func makeFile(entries []parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(key, value string) parser.Entry {
	return parser.Entry{Key: key, Value: value}
}

func TestApply_NilFile(t *testing.T) {
	_, err := counter.Apply(nil)
	if err == nil {
		t.Fatal("expected error for nil file")
	}
}

func TestApply_EmptyFile(t *testing.T) {
	r, err := counter.Apply(makeFile(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Total != 0 {
		t.Errorf("expected Total=0, got %d", r.Total)
	}
}

func TestApply_TotalCount(t *testing.T) {
	f := makeFile([]parser.Entry{
		entry("FOO", "bar"),
		entry("BAZ", "qux"),
	})
	r, err := counter.Apply(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Total != 2 {
		t.Errorf("expected Total=2, got %d", r.Total)
	}
}

func TestApply_EmptyAndNonEmpty(t *testing.T) {
	f := makeFile([]parser.Entry{
		entry("A", ""),
		entry("B", "hello"),
		entry("C", ""),
	})
	r, _ := counter.Apply(f)
	if r.Empty != 2 {
		t.Errorf("expected Empty=2, got %d", r.Empty)
	}
	if r.NonEmpty != 1 {
		t.Errorf("expected NonEmpty=1, got %d", r.NonEmpty)
	}
}

func TestApply_UniqueAndDuplicated(t *testing.T) {
	f := makeFile([]parser.Entry{
		entry("A", "same"),
		entry("B", "same"),
		entry("C", "unique"),
	})
	r, _ := counter.Apply(f)
	// "same" appears twice → duplicated bucket; "unique" once → unique bucket
	if r.Unique != 1 {
		t.Errorf("expected Unique=1, got %d", r.Unique)
	}
	if r.Duplicated != 1 {
		t.Errorf("expected Duplicated=1, got %d", r.Duplicated)
	}
}

func TestApply_TopPrefixes(t *testing.T) {
	f := makeFile([]parser.Entry{
		entry("DB_HOST", "localhost"),
		entry("DB_PORT", "5432"),
		entry("DB_NAME", "mydb"),
		entry("APP_ENV", "prod"),
		entry("NOPREFIXKEY", "val"),
	})
	r, err := counter.Apply(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.TopPrefixes["DB"] != 3 {
		t.Errorf("expected DB prefix count=3, got %d", r.TopPrefixes["DB"])
	}
	if r.TopPrefixes["APP"] != 1 {
		t.Errorf("expected APP prefix count=1, got %d", r.TopPrefixes["APP"])
	}
}
