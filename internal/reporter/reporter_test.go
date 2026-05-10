package reporter_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/reporter"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_NAME", Status: diff.Unchanged, OldValue: "myapp", NewValue: "myapp"},
		{Key: "DB_HOST", Status: diff.Modified, OldValue: "localhost", NewValue: "prod-db"},
		{Key: "NEW_KEY", Status: diff.Added, OldValue: "", NewValue: "new_val"},
		{Key: "OLD_KEY", Status: diff.Removed, OldValue: "old_val", NewValue: ""},
	}
}

func TestWriteText_ContainsAllEntries(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.Write(makeResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+ NEW_KEY=new_val") {
		t.Errorf("expected added key line, got:\n%s", out)
	}
	if !strings.Contains(out, "- OLD_KEY=old_val") {
		t.Errorf("expected removed key line, got:\n%s", out)
	}
	if !strings.Contains(out, "~ DB_HOST: localhost -> prod-db") {
		t.Errorf("expected modified key line, got:\n%s", out)
	}
	if !strings.Contains(out, "  APP_NAME=myapp") {
		t.Errorf("expected unchanged key line, got:\n%s", out)
	}
}

func TestWriteText_NoDifferences(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.Write([]diff.Result{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences found.") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}

func TestWriteJSON_ValidStructure(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf, reporter.FormatJSON)
	if err := r.Write(makeResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.HasPrefix(out, "[") || !strings.Contains(out, "]") {
		t.Errorf("expected JSON array, got: %s", out)
	}
	if !strings.Contains(out, `"key": "NEW_KEY"`) {
		t.Errorf("expected NEW_KEY in JSON output, got: %s", out)
	}
	if !strings.Contains(out, `"status": "added"`) {
		t.Errorf("expected added status in JSON output, got: %s", out)
	}
}

func TestSummary_Output(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.Summary(makeResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+1 added") {
		t.Errorf("expected +1 added in summary, got: %s", out)
	}
	if !strings.Contains(out, "-1 removed") {
		t.Errorf("expected -1 removed in summary, got: %s", out)
	}
	if !strings.Contains(out, "~1 modified") {
		t.Errorf("expected ~1 modified in summary, got: %s", out)
	}
	if !strings.Contains(out, "1 unchanged") {
		t.Errorf("expected 1 unchanged in summary, got: %s", out)
	}
}
