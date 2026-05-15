package splitter_test

import (
	"testing"

	"envdiff/internal/parser"
	"envdiff/internal/splitter"
)

func makeFile(pairs ...string) *parser.EnvFile {
	f := &parser.EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		f.Entries = append(f.Entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return f
}

func TestApply_NoFiles_ErrorOnNilSrc(t *testing.T) {
	_, err := splitter.Apply(nil, splitter.Options{Prefixes: map[string]string{"db": "DB_"}})
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestApply_NoPrefixes_ReturnsError(t *testing.T) {
	_, err := splitter.Apply(makeFile("DB_HOST", "localhost"), splitter.Options{})
	if err == nil {
		t.Fatal("expected error when no prefixes defined")
	}
}

func TestApply_SplitsByPrefix(t *testing.T) {
	src := makeFile("DB_HOST", "localhost", "AWS_KEY", "abc", "DB_PORT", "5432")
	opts := splitter.Options{
		Prefixes: map[string]string{"db": "DB_", "aws": "AWS_"},
	}
	res, err := splitter.Apply(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res["db"].Entries) != 2 {
		t.Errorf("expected 2 db entries, got %d", len(res["db"].Entries))
	}
	if len(res["aws"].Entries) != 1 {
		t.Errorf("expected 1 aws entry, got %d", len(res["aws"].Entries))
	}
}

func TestApply_UnmatchedDroppedByDefault(t *testing.T) {
	src := makeFile("DB_HOST", "localhost", "APP_NAME", "myapp")
	opts := splitter.Options{
		Prefixes: map[string]string{"db": "DB_"},
	}
	res, err := splitter.Apply(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res["other"]; ok {
		t.Error("did not expect 'other' group when IncludeUnmatched is false")
	}
}

func TestApply_IncludeUnmatched(t *testing.T) {
	src := makeFile("DB_HOST", "localhost", "APP_NAME", "myapp", "LOG_LEVEL", "info")
	opts := splitter.Options{
		Prefixes:         map[string]string{"db": "DB_"},
		IncludeUnmatched: true,
	}
	res, err := splitter.Apply(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res["other"].Entries) != 2 {
		t.Errorf("expected 2 unmatched entries, got %d", len(res["other"].Entries))
	}
}

func TestApply_CaseInsensitivePrefix(t *testing.T) {
	src := makeFile("db_host", "localhost")
	opts := splitter.Options{
		Prefixes: map[string]string{"db": "DB_"},
	}
	res, err := splitter.Apply(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res["db"].Entries) != 1 {
		t.Errorf("expected 1 db entry (case-insensitive), got %d", len(res["db"].Entries))
	}
}
