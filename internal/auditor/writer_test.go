package auditor_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/envdiff/envdiff/internal/auditor"
	"github.com/envdiff/envdiff/internal/diff"
)

func fixedLog() *auditor.Log {
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	return auditor.FromDiff([]diff.Result{
		{Key: "API_KEY", Status: diff.Added, BaseValue: "", CompareValue: "xyz"},
		{Key: "OLD_VAR", Status: diff.Removed, BaseValue: "old", CompareValue: ""},
	}, auditor.Options{Actor: "dev", File: ".env", Timestamp: ts})
}

func TestWrite_TextFormat_ContainsKey(t *testing.T) {
	var buf bytes.Buffer
	if err := fixedLog().Write(&buf, auditor.FormatText); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in text output, got:\n%s", out)
	}
	if !strings.Contains(out, "OLD_VAR") {
		t.Errorf("expected OLD_VAR in text output")
	}
}

func TestWrite_TextFormat_EmptyValuePlaceholder(t *testing.T) {
	var buf bytes.Buffer
	_ = fixedLog().Write(&buf, auditor.FormatText)
	if !strings.Contains(buf.String(), "(empty)") {
		t.Errorf("expected (empty) placeholder in text output")
	}
}

func TestWrite_JSONFormat_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := fixedLog().Write(&buf, auditor.FormatJSON); err != nil {
		t.Fatal(err)
	}
	var entries []auditor.Entry
	if err := json.Unmarshal(buf.Bytes(), &entries); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 JSON entries, got %d", len(entries))
	}
}

func TestWrite_JSONFormat_FieldsPresent(t *testing.T) {
	var buf bytes.Buffer
	_ = fixedLog().Write(&buf, auditor.FormatJSON)
	out := buf.String()
	for _, field := range []string{"timestamp", "actor", "file", "key", "action"} {
		if !strings.Contains(out, field) {
			t.Errorf("expected field %q in JSON output", field)
		}
	}
}

func TestWrite_UnknownFormat_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	err := fixedLog().Write(&buf, auditor.Format("xml"))
	if err == nil {
		t.Error("expected error for unknown format")
	}
}
