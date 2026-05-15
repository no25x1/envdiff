package patcher_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/parser"
	"github.com/yourorg/envdiff/internal/patcher"
)

func makeFile(pairs ...string) parser.EnvFile {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return parser.EnvFile{Entries: entries}
}

func keys(f parser.EnvFile) []string {
	out := make([]string, len(f.Entries))
	for i, e := range f.Entries {
		out[i] = e.Key
	}
	return out
}

func TestApply_SetNewKey(t *testing.T) {
	src := makeFile("HOST", "localhost")
	out, err := patcher.Apply(src, []patcher.Op{{Kind: patcher.OpSet, Key: "PORT", Value: "5432"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out.Entries))
	}
	if out.Entries[1].Value != "5432" {
		t.Errorf("expected PORT=5432, got %s", out.Entries[1].Value)
	}
}

func TestApply_SetExistingKey(t *testing.T) {
	src := makeFile("HOST", "localhost", "PORT", "3000")
	out, err := patcher.Apply(src, []patcher.Op{{Kind: patcher.OpSet, Key: "PORT", Value: "8080"}})
	if err != nil {
		t.Fatal(err)
	}
	if out.Entries[1].Value != "8080" {
		t.Errorf("expected PORT=8080, got %s", out.Entries[1].Value)
	}
}

func TestApply_DeleteKey(t *testing.T) {
	src := makeFile("HOST", "localhost", "PORT", "3000")
	out, err := patcher.Apply(src, []patcher.Op{{Kind: patcher.OpDelete, Key: "PORT"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Entries) != 1 || out.Entries[0].Key != "HOST" {
		t.Errorf("unexpected entries: %v", keys(out))
	}
}

func TestApply_DeleteMissingKey_ReturnsError(t *testing.T) {
	src := makeFile("HOST", "localhost")
	_, err := patcher.Apply(src, []patcher.Op{{Kind: patcher.OpDelete, Key: "MISSING"}})
	if err == nil {
		t.Fatal("expected error for missing key delete")
	}
}

func TestApply_RenameKey(t *testing.T) {
	src := makeFile("HOST", "localhost", "PORT", "3000")
	out, err := patcher.Apply(src, []patcher.Op{{Kind: patcher.OpRename, Key: "HOST", NewKey: "DB_HOST"}})
	if err != nil {
		t.Fatal(err)
	}
	if out.Entries[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", out.Entries[0].Key)
	}
}

func TestApply_RenameToExistingKey_ReturnsError(t *testing.T) {
	src := makeFile("HOST", "localhost", "PORT", "3000")
	_, err := patcher.Apply(src, []patcher.Op{{Kind: patcher.OpRename, Key: "HOST", NewKey: "PORT"}})
	if err == nil {
		t.Fatal("expected error for rename collision")
	}
}

func TestApply_UnknownOp_ReturnsError(t *testing.T) {
	src := makeFile("HOST", "localhost")
	_, err := patcher.Apply(src, []patcher.Op{{Kind: "upsert", Key: "HOST"}})
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}
