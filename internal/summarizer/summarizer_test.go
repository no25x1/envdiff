package summarizer_test

import (
	"testing"

	"envdiff/internal/parser"
	"envdiff/internal/summarizer"
)

func makeFile(entries ...parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(key, value string) parser.Entry {
	return parser.Entry{Key: key, Value: value}
}

func TestAnalyse_NilFile(t *testing.T) {
	s := summarizer.Analyse(nil)
	if s.TotalKeys != 0 {
		t.Fatalf("expected 0 total keys, got %d", s.TotalKeys)
	}
}

func TestAnalyse_TotalKeys(t *testing.T) {
	f := makeFile(entry("FOO", "bar"), entry("BAZ", "qux"))
	s := summarizer.Analyse(f)
	if s.TotalKeys != 2 {
		t.Fatalf("expected 2, got %d", s.TotalKeys)
	}
}

func TestAnalyse_EmptyValues(t *testing.T) {
	f := makeFile(entry("A", ""), entry("B", "val"), entry("C", ""))
	s := summarizer.Analyse(f)
	if s.EmptyValues != 2 {
		t.Fatalf("expected 2 empty, got %d", s.EmptyValues)
	}
}

func TestAnalyse_SensitiveKeys(t *testing.T) {
	f := makeFile(
		entry("DB_PASSWORD", "secret"),
		entry("API_TOKEN", "tok"),
		entry("HOST", "localhost"),
	)
	s := summarizer.Analyse(f)
	if s.SensitiveKeys != 2 {
		t.Fatalf("expected 2 sensitive, got %d", s.SensitiveKeys)
	}
}

func TestAnalyse_LongestAndShortestKey(t *testing.T) {
	f := makeFile(entry("AB", "1"), entry("ABCDEFGH", "2"), entry("ABCD", "3"))
	s := summarizer.Analyse(f)
	if s.LongestKey != "ABCDEFGH" {
		t.Fatalf("expected ABCDEFGH, got %s", s.LongestKey)
	}
	if s.ShortestKey != "AB" {
		t.Fatalf("expected AB, got %s", s.ShortestKey)
	}
}

func TestAnalyse_TopPrefixes(t *testing.T) {
	f := makeFile(
		entry("DB_HOST", "h"),
		entry("DB_PORT", "p"),
		entry("DB_USER", "u"),
		entry("APP_ENV", "prod"),
		entry("APP_NAME", "envdiff"),
		entry("CACHE_TTL", "60"),
	)
	s := summarizer.Analyse(f)
	if len(s.TopPrefixes) == 0 {
		t.Fatal("expected top prefixes, got none")
	}
	if s.TopPrefixes[0].Prefix != "DB" {
		t.Fatalf("expected DB as top prefix, got %s", s.TopPrefixes[0].Prefix)
	}
	if s.TopPrefixes[0].Count != 3 {
		t.Fatalf("expected count 3, got %d", s.TopPrefixes[0].Count)
	}
}

func TestAnalyse_UniquePrefixes(t *testing.T) {
	f := makeFile(
		entry("DB_HOST", "h"),
		entry("DB_PORT", "p"),
		entry("APP_ENV", "prod"),
		entry("NOPREFIXKEY", "x"),
	)
	s := summarizer.Analyse(f)
	if s.UniquePrefix != 2 {
		t.Fatalf("expected 2 unique prefixes, got %d", s.UniquePrefix)
	}
}
