package renamer_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/parser"
	"github.com/your-org/envdiff/internal/renamer"
)

func makeFile(pairs ...string) *parser.EnvFile {
	f := &parser.EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		f.Entries = append(f.Entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return f
}

func TestApply_NoOptions_ReturnsCopy(t *testing.T) {
	f := makeFile("FOO", "1", "BAR", "2")
	res, err := renamer.Apply(f, renamer.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Changes) != 0 {
		t.Errorf("expected no changes, got %d", len(res.Changes))
	}
	if len(res.File.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.File.Entries))
	}
}

func TestApply_ExplicitMapping(t *testing.T) {
	f := makeFile("OLD_KEY", "val")
	res, err := renamer.Apply(f, renamer.Options{
		Mappings: map[string]string{"OLD_KEY": "NEW_KEY"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Changes) != 1 || res.Changes[0].NewKey != "NEW_KEY" {
		t.Errorf("expected rename to NEW_KEY, got %+v", res.Changes)
	}
	if res.File.Entries[0].Key != "NEW_KEY" {
		t.Errorf("expected entry key NEW_KEY, got %s", res.File.Entries[0].Key)
	}
}

func TestApply_AddPrefix(t *testing.T) {
	f := makeFile("HOST", "localhost", "PORT", "5432")
	res, err := renamer.Apply(f, renamer.Options{AddPrefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(res.Changes))
	}
	if res.File.Entries[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", res.File.Entries[0].Key)
	}
}

func TestApply_StripPrefix(t *testing.T) {
	f := makeFile("APP_FOO", "1", "APP_BAR", "2")
	res, err := renamer.Apply(f, renamer.Options{StripPrefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.File.Entries[0].Key != "FOO" {
		t.Errorf("expected FOO, got %s", res.File.Entries[0].Key)
	}
}

func TestApply_ToUpper(t *testing.T) {
	f := makeFile("db_host", "localhost")
	res, err := renamer.Apply(f, renamer.Options{ToUpper: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.File.Entries[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", res.File.Entries[0].Key)
	}
}

func TestApply_DuplicateAfterRename_ReturnsError(t *testing.T) {
	f := makeFile("APP_FOO", "1", "FOO", "2")
	_, err := renamer.Apply(f, renamer.Options{StripPrefix: "APP_"})
	if err == nil {
		t.Error("expected error for duplicate key after rename")
	}
}
