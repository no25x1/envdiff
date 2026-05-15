package schema_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/envdiff/envdiff/internal/parser"
	"github.com/envdiff/envdiff/internal/schema"
)

func makeFile(entries ...parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(k, v string) parser.Entry { return parser.Entry{Key: k, Value: v} }

func writeTempSchema(t *testing.T, rules []map[string]interface{}) string {
	t.Helper()
	data, _ := json.Marshal(map[string]interface{}{"rules": rules})
	f, err := os.CreateTemp(t.TempDir(), "schema-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.Write(data)
	_ = f.Close()
	return f.Name()
}

func TestValidate_NoViolations(t *testing.T) {
	s := &schema.Schema{
		Rules: []schema.FieldRule{
			{Key: "APP_ENV", Required: true},
		},
	}
	f := makeFile(entry("APP_ENV", "production"))
	if v := s.Validate(f); len(v) != 0 {
		t.Fatalf("expected no violations, got %v", v)
	}
}

func TestValidate_MissingRequiredKey(t *testing.T) {
	s := &schema.Schema{
		Rules: []schema.FieldRule{
			{Key: "DB_URL", Required: true},
		},
	}
	f := makeFile(entry("APP_ENV", "staging"))
	v := s.Validate(f)
	if len(v) != 1 || v[0].Key != "DB_URL" {
		t.Fatalf("expected DB_URL violation, got %v", v)
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	s := &schema.Schema{
		Rules: []schema.FieldRule{
			{Key: "PORT", Required: true, Pattern: `^\d+$`},
		},
	}
	f := makeFile(entry("PORT", "not-a-number"))
	v := s.Validate(f)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %v", v)
	}
}

func TestValidate_PatternMatch(t *testing.T) {
	s := &schema.Schema{
		Rules: []schema.FieldRule{
			{Key: "PORT", Required: true, Pattern: `^\d+$`},
		},
	}
	f := makeFile(entry("PORT", "8080"))
	if v := s.Validate(f); len(v) != 0 {
		t.Fatalf("expected no violations, got %v", v)
	}
}

func TestLoadFile_Valid(t *testing.T) {
	path := writeTempSchema(t, []map[string]interface{}{
		{"key": "APP_ENV", "required": true},
		{"key": "PORT", "required": false, "pattern": `^\d+$`},
	})
	s, err := schema.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(s.Rules))
	}
}

func TestLoadFile_Missing(t *testing.T) {
	_, err := schema.LoadFile("/nonexistent/schema.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
