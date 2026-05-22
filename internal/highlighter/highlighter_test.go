package highlighter_test

import (
	"strings"
	"testing"

	"github.com/envdiff/envdiff/internal/diff"
	"github.com/envdiff/envdiff/internal/highlighter"
)

func makeResult(status diff.Status, key, oldVal, newVal string) diff.Result {
	return diff.Result{Key: key, Status: status, OldValue: oldVal, NewValue: newVal}
}

func TestWrite_NoResults_PrintsNoDifferences(t *testing.T) {
	var buf strings.Builder
	err := highlighter.Write(&buf, nil, highlighter.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected 'No differences' message, got: %q", buf.String())
	}
}

func TestWrite_Added_ContainsPlusPrefix(t *testing.T) {
	var buf strings.Builder
	results := []diff.Result{makeResult(diff.Added, "FOO", "", "bar")}
	_ = highlighter.Write(&buf, results, highlighter.Options{NoColor: true})
	if !strings.Contains(buf.String(), "+ FOO=bar") {
		t.Errorf("expected '+ FOO=bar', got: %q", buf.String())
	}
}

func TestWrite_Removed_ContainsMinusPrefix(t *testing.T) {
	var buf strings.Builder
	results := []diff.Result{makeResult(diff.Removed, "BAR", "old", "")}
	_ = highlighter.Write(&buf, results, highlighter.Options{NoColor: true})
	if !strings.Contains(buf.String(), "- BAR=old") {
		t.Errorf("expected '- BAR=old', got: %q", buf.String())
	}
}

func TestWrite_Modified_ContainsTilde(t *testing.T) {
	var buf strings.Builder
	results := []diff.Result{makeResult(diff.Modified, "BAZ", "v1", "v2")}
	_ = highlighter.Write(&buf, results, highlighter.Options{NoColor: true})
	out := buf.String()
	if !strings.Contains(out, "~ BAZ") {
		t.Errorf("expected '~ BAZ', got: %q", out)
	}
	if !strings.Contains(out, "v1") || !strings.Contains(out, "v2") {
		t.Errorf("expected both old and new values in output, got: %q", out)
	}
}

func TestWrite_Unchanged_ContainsKey(t *testing.T) {
	var buf strings.Builder
	results := []diff.Result{makeResult(diff.Unchanged, "QUX", "same", "same")}
	_ = highlighter.Write(&buf, results, highlighter.Options{NoColor: true})
	if !strings.Contains(buf.String(), "QUX=same") {
		t.Errorf("expected 'QUX=same', got: %q", buf.String())
	}
}

func TestWrite_WithColor_ContainsAnsiCodes(t *testing.T) {
	var buf strings.Builder
	results := []diff.Result{makeResult(diff.Added, "KEY", "", "val")}
	_ = highlighter.Write(&buf, results, highlighter.DefaultOptions())
	if !strings.Contains(buf.String(), "\033[") {
		t.Errorf("expected ANSI escape codes in coloured output, got: %q", buf.String())
	}
}

func TestWrite_NoColor_NoAnsiCodes(t *testing.T) {
	var buf strings.Builder
	results := []diff.Result{makeResult(diff.Removed, "KEY", "val", "")}
	_ = highlighter.Write(&buf, results, highlighter.Options{NoColor: true})
	if strings.Contains(buf.String(), "\033[") {
		t.Errorf("expected no ANSI codes when NoColor=true, got: %q", buf.String())
	}
}
