package validator_test

import (
	"testing"

	"github.com/envdiff/internal/parser"
	"github.com/envdiff/internal/validator"
)

func makeFile(entries []parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func TestValidate_NoViolations(t *testing.T) {
	f := makeFile([]parser.Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_HOST", Value: "localhost"},
	})
	v := validator.New()
	results := v.Validate(f)
	if len(results) != 0 {
		t.Fatalf("expected no violations, got %d: %+v", len(results), results)
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	f := makeFile([]parser.Entry{
		{Key: "API_KEY", Value: ""},
	})
	v := validator.New()
	results := v.Validate(f)
	if !containsRule(results, "no-empty-value") {
		t.Errorf("expected no-empty-value violation, got %+v", results)
	}
}

func TestValidate_LowercaseKey(t *testing.T) {
	f := makeFile([]parser.Entry{
		{Key: "db_host", Value: "localhost"},
	})
	v := validator.New()
	results := v.Validate(f)
	if !containsRule(results, "uppercase-key") {
		t.Errorf("expected uppercase-key violation, got %+v", results)
	}
}

func TestValidate_KeyWithSpaces(t *testing.T) {
	f := makeFile([]parser.Entry{
		{Key: "MY KEY", Value: "val"},
	})
	v := validator.New()
	results := v.Validate(f)
	if !containsRule(results, "no-spaces-in-key") {
		t.Errorf("expected no-spaces-in-key violation, got %+v", results)
	}
}

func TestValidate_CustomRule(t *testing.T) {
	customRule := validator.Rule{
		Name: "no-localhost",
		Check: func(key, value string) error {
			if value == "localhost" {
				return fmt.Errorf("key %q must not use localhost in production", key)
			}
			return nil
		},
	}
	f := makeFile([]parser.Entry{
		{Key: "DB_HOST", Value: "localhost"},
	})
	v := validator.WithRules(customRule)
	results := v.Validate(f)
	if !containsRule(results, "no-localhost") {
		t.Errorf("expected no-localhost violation, got %+v", results)
	}
}

func TestValidate_MultipleViolations(t *testing.T) {
	f := makeFile([]parser.Entry{
		{Key: "bad key", Value: ""},
	})
	v := validator.New()
	results := v.Validate(f)
	if len(results) < 2 {
		t.Errorf("expected at least 2 violations, got %d", len(results))
	}
}

func containsRule(results []validator.Result, rule string) bool {
	for _, r := range results {
		if r.Rule == rule {
			return true
		}
	}
	return false
}
