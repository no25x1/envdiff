package capitalizer_test

import (
	"testing"

	"github.com/envdiff/envdiff/internal/capitalizer"
	"github.com/envdiff/envdiff/internal/parser"
)

func makeFile(pairs ...string) *parser.EnvFile {
	f := &parser.EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		f.Entries = append(f.Entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return f
}

func keys(f *parser.EnvFile) []string {
	out := make([]string, len(f.Entries))
	for i, e := range f.Entries {
		out[i] = e.Key
	}
	return out
}

func TestApply_NilFile_ReturnsEmpty(t *testing.T) {
	out := capitalizer.Apply(nil, capitalizer.Options{KeysToUpper: true})
	if len(out.Entries) != 0 {
		t.Fatalf("expected empty, got %d entries", len(out.Entries))
	}
}

func TestApply_KeysToUpper(t *testing.T) {
	f := makeFile("db_host", "localhost", "db_port", "5432")
	out := capitalizer.Apply(f, capitalizer.Options{KeysToUpper: true})
	for _, e := range out.Entries {
		if e.Key != strings.ToUpper(e.Key) {
			t.Errorf("expected upper-case key, got %q", e.Key)
		}
	}
}

func TestApply_KeysToLower(t *testing.T) {
	f := makeFile("DB_HOST", "localhost", "APP_ENV", "prod")
	out := capitalizer.Apply(f, capitalizer.Options{KeysToLower: true})
	for _, e := range out.Entries {
		if e.Key != strings.ToLower(e.Key) {
			t.Errorf("expected lower-case key, got %q", e.Key)
		}
	}
}

func TestApply_ValuesToUpper(t *testing.T) {
	f := makeFile("ENV", "production")
	out := capitalizer.Apply(f, capitalizer.Options{ValuesToUpper: true})
	if out.Entries[0].Value != "PRODUCTION" {
		t.Errorf("expected PRODUCTION, got %q", out.Entries[0].Value)
	}
}

func TestApply_ValuesToLower(t *testing.T) {
	f := makeFile("ENV", "PRODUCTION")
	out := capitalizer.Apply(f, capitalizer.Options{ValuesToLower: true})
	if out.Entries[0].Value != "production" {
		t.Errorf("expected production, got %q", out.Entries[0].Value)
	}
}

func TestApply_OnlyKeys_LimitsTransformation(t *testing.T) {
	f := makeFile("db_host", "localhost", "app_env", "dev")
	out := capitalizer.Apply(f, capitalizer.Options{
		KeysToUpper: true,
		OnlyKeys:    []string{"db_host"},
	})
	if out.Entries[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %q", out.Entries[0].Key)
	}
	if out.Entries[1].Key != "app_env" {
		t.Errorf("expected app_env unchanged, got %q", out.Entries[1].Key)
	}
}

func TestApply_NoOptions_ReturnsCopy(t *testing.T) {
	f := makeFile("KEY", "value")
	out := capitalizer.Apply(f, capitalizer.Options{})
	if out.Entries[0].Key != "KEY" || out.Entries[0].Value != "value" {
		t.Error("expected unchanged copy")
	}
}
