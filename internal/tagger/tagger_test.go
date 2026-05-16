package tagger_test

import (
	"testing"

	"envdiff/internal/parser"
	"envdiff/internal/tagger"
)

func makeFile(entries []parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(key, value, comment string) parser.Entry {
	return parser.Entry{Key: key, Value: value, Comment: comment}
}

func TestApply_NilFile_ReturnsEmpty(t *testing.T) {
	out, err := tagger.Apply(nil, tagger.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(out.Entries))
	}
}

func TestApply_AddTag_AllKeys(t *testing.T) {
	src := makeFile([]parser.Entry{
		entry("DB_HOST", "localhost", ""),
		entry("APP_ENV", "prod", ""),
	})
	out, err := tagger.Apply(src, tagger.Options{Add: []string{"production"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out.Entries {
		tags := tagger.Tags(e.Comment)
		if len(tags) != 1 || tags[0] != "production" {
			t.Errorf("key %s: expected tag production, got %v", e.Key, tags)
		}
	}
}

func TestApply_AddTag_SpecificKeys(t *testing.T) {
	src := makeFile([]parser.Entry{
		entry("SECRET", "abc", ""),
		entry("PORT", "8080", ""),
	})
	out, err := tagger.Apply(src, tagger.Options{Add: []string{"sensitive"}, Keys: []string{"SECRET"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out.Entries {
		tags := tagger.Tags(e.Comment)
		if e.Key == "SECRET" && (len(tags) != 1 || tags[0] != "sensitive") {
			t.Errorf("SECRET should have tag sensitive, got %v", tags)
		}
		if e.Key == "PORT" && len(tags) != 0 {
			t.Errorf("PORT should have no tags, got %v", tags)
		}
	}
}

func TestApply_RemoveTag(t *testing.T) {
	src := makeFile([]parser.Entry{
		entry("DB_HOST", "localhost", "# @tag:staging @tag:production"),
	})
	out, err := tagger.Apply(src, tagger.Options{Remove: []string{"staging"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags := tagger.Tags(out.Entries[0].Comment)
	if len(tags) != 1 || tags[0] != "production" {
		t.Errorf("expected only production tag, got %v", tags)
	}
}

func TestApply_FilterTag_ReturnsOnlyMatching(t *testing.T) {
	src := makeFile([]parser.Entry{
		entry("A", "1", "# @tag:prod"),
		entry("B", "2", "# @tag:dev"),
		entry("C", "3", "# @tag:prod @tag:dev"),
	})
	out, err := tagger.Apply(src, tagger.Options{FilterTag: "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out.Entries))
	}
}

func TestApply_AddTag_NoDuplicate(t *testing.T) {
	src := makeFile([]parser.Entry{
		entry("KEY", "val", "# @tag:prod"),
	})
	out, err := tagger.Apply(src, tagger.Options{Add: []string{"prod"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags := tagger.Tags(out.Entries[0].Comment)
	if len(tags) != 1 {
		t.Errorf("expected exactly 1 tag, got %v", tags)
	}
}
