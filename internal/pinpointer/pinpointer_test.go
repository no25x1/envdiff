package pinpointer_test

import (
	"testing"

	"envdiff/internal/parser"
	"envdiff/internal/pinpointer"
)

func makeFile(entries ...parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(k, v string) parser.Entry { return parser.Entry{Key: k, Value: v} }

func TestApply_NilFile_ReturnsNil(t *testing.T) {
	findings, err := pinpointer.Apply(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if findings != nil {
		t.Errorf("expected nil findings, got %v", findings)
	}
}

func TestApply_NoFindings_CleanValues(t *testing.T) {
	f := makeFile(
		entry("APP_NAME", "myapp"),
		entry("PORT", "8080"),
		entry("TIMEOUT", "30s"),
	)
	findings, err := pinpointer.Apply(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d", len(findings))
	}
}

func TestApply_DetectsLocalhost(t *testing.T) {
	f := makeFile(entry("DB_HOST", "localhost"))
	findings, _ := pinpointer.Apply(f)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "DB_HOST" {
		t.Errorf("unexpected key: %s", findings[0].Key)
	}
}

func TestApply_DetectsLoopbackIP(t *testing.T) {
	f := makeFile(entry("REDIS_HOST", "127.0.0.1"))
	findings, _ := pinpointer.Apply(f)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

func TestApply_DetectsRawIP(t *testing.T) {
	f := makeFile(entry("SERVICE_IP", "10.0.0.5"))
	findings, _ := pinpointer.Apply(f)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

func TestApply_DetectsEnvTaggedURL(t *testing.T) {
	f := makeFile(entry("API_URL", "https://staging-api.example.com/v1"))
	findings, _ := pinpointer.Apply(f)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Reason == "" {
		t.Error("expected non-empty reason")
	}
}

func TestApply_DetectsEnvTagInPlainValue(t *testing.T) {
	f := makeFile(entry("ENV_LABEL", "prod"))
	findings, _ := pinpointer.Apply(f)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

func TestApply_MultipleFindings(t *testing.T) {
	f := makeFile(
		entry("CLEAN", "value"),
		entry("DB", "localhost:5432"),
		entry("API", "https://dev.api.io"),
	)
	findings, _ := pinpointer.Apply(f)
	if len(findings) != 2 {
		t.Errorf("expected 2 findings, got %d", len(findings))
	}
}
