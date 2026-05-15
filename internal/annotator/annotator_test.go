package annotator_test

import (
	"testing"

	"github.com/user/envdiff/internal/annotator"
	"github.com/user/envdiff/internal/parser"
)

func makeFile(pairs ...string) *parser.EnvFile {
	f := &parser.EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		f.Entries = append(f.Entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return f
}

func entry(key, value, comment string) parser.Entry {
	return parser.Entry{Key: key, Value: value, Comment: comment}
}

func TestApply_NilFile_ReturnsEmpty(t *testing.T) {
	out := annotator.Apply(nil, annotator.DefaultOptions())
	if len(out.Entries) != 0 {
		t.Fatalf("expected empty, got %d entries", len(out.Entries))
	}
}

func TestApply_NoMatchingRules_Unchanged(t *testing.T) {
	src := makeFile("APP_NAME", "envdiff", "PORT", "8080")
	out := annotator.Apply(src, annotator.DefaultOptions())
	for _, e := range out.Entries {
		if e.Comment != "" {
			t.Errorf("expected no comment for %s, got %q", e.Key, e.Comment)
		}
	}
}

func TestApply_SensitiveSuffix_GetsComment(t *testing.T) {
	src := makeFile("DB_PASSWORD", "hunter2", "API_TOKEN", "abc123")
	out := annotator.Apply(src, annotator.DefaultOptions())
	for _, e := range out.Entries {
		if e.Comment == "" {
			t.Errorf("expected annotation for %s", e.Key)
		}
	}
}

func TestApply_DebugPrefix_GetsComment(t *testing.T) {
	src := makeFile("DEBUG_VERBOSE", "true")
	out := annotator.Apply(src, annotator.DefaultOptions())
	if len(out.Entries) == 0 || out.Entries[0].Comment == "" {
		t.Fatal("expected annotation for DEBUG_ prefix key")
	}
}

func TestApply_ExistingComment_PreservedByDefault(t *testing.T) {
	src := &parser.EnvFile{
		Entries: []parser.Entry{
			{Key: "DB_SECRET", Value: "x", Comment: "# my custom note"},
		},
	}
	out := annotator.Apply(src, annotator.DefaultOptions())
	if out.Entries[0].Comment != "# my custom note" {
		t.Errorf("expected original comment preserved, got %q", out.Entries[0].Comment)
	}
}

func TestApply_OverwriteExisting_ReplacesComment(t *testing.T) {
	src := &parser.EnvFile{
		Entries: []parser.Entry{
			{Key: "DB_SECRET", Value: "x", Comment: "# my custom note"},
		},
	}
	opts := annotator.DefaultOptions()
	opts.OverwriteExisting = true
	out := annotator.Apply(src, opts)
	if out.Entries[0].Comment == "# my custom note" {
		t.Error("expected comment to be overwritten")
	}
}

func TestApply_CustomRule_Matches(t *testing.T) {
	src := makeFile("INTERNAL_FLAG", "1")
	opts := annotator.Options{
		Rules: []annotator.Rule{
			{KeyPrefix: "INTERNAL_", Comment: "internal use only"},
		},
	}
	out := annotator.Apply(src, opts)
	if out.Entries[0].Comment != "# internal use only" {
		t.Errorf("unexpected comment: %q", out.Entries[0].Comment)
	}
}
