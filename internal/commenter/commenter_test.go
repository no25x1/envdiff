package commenter_test

import (
	"testing"

	"github.com/envdiff/envdiff/internal/commenter"
	"github.com/envdiff/envdiff/internal/parser"
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
	out, err := commenter.Apply(nil, commenter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 0 {
		t.Errorf("expected empty, got %d entries", len(out.Entries))
	}
}

func TestApply_SetComment_NoExisting(t *testing.T) {
	src := makeFile("DB_HOST", "localhost", "API_KEY", "secret")
	out, err := commenter.Apply(src, commenter.Options{
		Set: map[string]string{"API_KEY": "sensitive value"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out.Entries {
		if e.Key == "API_KEY" && e.Comment != "# sensitive value" {
			t.Errorf("expected '# sensitive value', got %q", e.Comment)
		}
		if e.Key == "DB_HOST" && e.Comment != "" {
			t.Errorf("expected empty comment for DB_HOST, got %q", e.Comment)
		}
	}
}

func TestApply_SetComment_OverwriteExisting(t *testing.T) {
	src := &parser.EnvFile{
		Entries: []parser.Entry{
			entry("API_KEY", "secret", "# old comment"),
		},
	}
	out, err := commenter.Apply(src, commenter.Options{
		Set:       map[string]string{"API_KEY": "new comment"},
		Overwrite: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[0].Comment != "# new comment" {
		t.Errorf("expected '# new comment', got %q", out.Entries[0].Comment)
	}
}

func TestApply_SetComment_NoOverwriteKeepsExisting(t *testing.T) {
	src := &parser.EnvFile{
		Entries: []parser.Entry{
			entry("API_KEY", "secret", "# keep me"),
		},
	}
	out, err := commenter.Apply(src, commenter.Options{
		Set:       map[string]string{"API_KEY": "new comment"},
		Overwrite: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[0].Comment != "# keep me" {
		t.Errorf("expected '# keep me', got %q", out.Entries[0].Comment)
	}
}

func TestApply_RemoveComment(t *testing.T) {
	src := &parser.EnvFile{
		Entries: []parser.Entry{
			entry("DB_PASS", "hunter2", "# remove this"),
			entry("DB_HOST", "localhost", "# keep this"),
		},
	}
	out, err := commenter.Apply(src, commenter.Options{
		Remove: []string{"DB_PASS"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out.Entries {
		if e.Key == "DB_PASS" && e.Comment != "" {
			t.Errorf("expected empty comment for DB_PASS, got %q", e.Comment)
		}
		if e.Key == "DB_HOST" && e.Comment != "# keep this" {
			t.Errorf("expected '# keep this' for DB_HOST, got %q", e.Comment)
		}
	}
}
