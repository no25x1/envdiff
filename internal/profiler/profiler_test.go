package profiler_test

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/profiler"
)

func makeFile(entries []parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(key, value string) parser.Entry {
	return parser.Entry{Key: key, Value: value}
}

func TestAnalyse_NilFile(t *testing.T) {
	p := profiler.Analyse(nil)
	if p.TotalKeys != 0 {
		t.Fatalf("expected 0 keys, got %d", p.TotalKeys)
	}
}

func TestAnalyse_TotalKeys(t *testing.T) {
	f := makeFile([]parser.Entry{
		entry("FOO", "bar"),
		entry("BAZ", "qux"),
	})
	p := profiler.Analyse(f)
	if p.TotalKeys != 2 {
		t.Fatalf("expected 2, got %d", p.TotalKeys)
	}
}

func TestAnalyse_EmptyValues(t *testing.T) {
	f := makeFile([]parser.Entry{
		entry("FOO", ""),
		entry("BAR", "val"),
		entry("BAZ", ""),
	})
	p := profiler.Analyse(f)
	if p.EmptyValues != 2 {
		t.Fatalf("expected 2 empty, got %d", p.EmptyValues)
	}
}

func TestAnalyse_SensitiveKeys(t *testing.T) {
	f := makeFile([]parser.Entry{
		entry("DB_PASSWORD", "secret"),
		entry("API_TOKEN", "tok"),
		entry("APP_NAME", "myapp"),
	})
	p := profiler.Analyse(f)
	if p.SensitiveKeys != 2 {
		t.Fatalf("expected 2 sensitive, got %d", p.SensitiveKeys)
	}
}

func TestAnalyse_PrefixCounts(t *testing.T) {
	f := makeFile([]parser.Entry{
		entry("DB_HOST", "localhost"),
		entry("DB_PORT", "5432"),
		entry("APP_NAME", "envdiff"),
	})
	p := profiler.Analyse(f)
	if p.PrefixCounts["DB"] != 2 {
		t.Fatalf("expected DB=2, got %d", p.PrefixCounts["DB"])
	}
	if p.PrefixCounts["APP"] != 1 {
		t.Fatalf("expected APP=1, got %d", p.PrefixCounts["APP"])
	}
}

func TestTopPrefixes_Order(t *testing.T) {
	f := makeFile([]parser.Entry{
		entry("DB_HOST", "h"),
		entry("DB_PORT", "p"),
		entry("DB_NAME", "n"),
		entry("APP_ENV", "prod"),
		entry("APP_PORT", "8080"),
		entry("SVC_URL", "http"),
	})
	p := profiler.Analyse(f)
	top := p.TopPrefixes(2)
	if len(top) != 2 {
		t.Fatalf("expected 2 prefixes, got %d", len(top))
	}
	if top[0] != "DB" {
		t.Fatalf("expected DB first, got %s", top[0])
	}
}
