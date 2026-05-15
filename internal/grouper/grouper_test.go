package grouper_test

import (
	"testing"

	"github.com/user/envdiff/internal/grouper"
	"github.com/user/envdiff/internal/parser"
)

func makeFile(entries ...parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(key, value string) parser.Entry {
	return parser.Entry{Key: key, Value: value}
}

func TestApply_GroupsByPrefix(t *testing.T) {
	f := makeFile(
		entry("DB_HOST", "localhost"),
		entry("DB_PORT", "5432"),
		entry("APP_NAME", "envdiff"),
		entry("APP_ENV", "prod"),
	)
	groups := grouper.Apply(f, grouper.Options{})
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Name != "APP" {
		t.Errorf("expected first group APP, got %s", groups[0].Name)
	}
	if groups[1].Name != "DB" {
		t.Errorf("expected second group DB, got %s", groups[1].Name)
	}
}

func TestApply_UngroupedSkippedByDefault(t *testing.T) {
	f := makeFile(
		entry("DB_HOST", "localhost"),
		entry("STANDALONE", "value"),
	)
	groups := grouper.Apply(f, grouper.Options{})
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Name != "DB" {
		t.Errorf("expected group DB, got %s", groups[0].Name)
	}
}

func TestApply_IncludeUngrouped(t *testing.T) {
	f := makeFile(
		entry("DB_HOST", "localhost"),
		entry("STANDALONE", "value"),
	)
	groups := grouper.Apply(f, grouper.Options{IncludeUngrouped: true})
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	// ungrouped sorts before DB
	if groups[0].Name != "" {
		t.Errorf("expected first group to be ungrouped, got %q", groups[0].Name)
	}
}

func TestApply_CustomDelimiter(t *testing.T) {
	f := makeFile(
		entry("DB.HOST", "localhost"),
		entry("DB.PORT", "5432"),
		entry("APP.NAME", "envdiff"),
	)
	groups := grouper.Apply(f, grouper.Options{Delimiter: "."})
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
}

func TestApply_EmptyFile(t *testing.T) {
	f := makeFile()
	groups := grouper.Apply(f, grouper.Options{})
	if len(groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(groups))
	}
}
