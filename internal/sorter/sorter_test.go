package sorter_test

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/sorter"
)

func makeFile(keys ...string) parser.EnvFile {
	entries := make([]parser.Entry, len(keys))
	for i, k := range keys {
		entries[i] = parser.Entry{Key: k, Value: "v"}
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

func TestApply_Alphabetical(t *testing.T) {
	f := makeFile("ZEBRA", "APPLE", "MANGO")
	out := sorter.Apply(f, sorter.Options{Alphabetical: true})
	got := keys(out)
	want := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range want {
		if got[i] != k {
			t.Errorf("pos %d: got %q want %q", i, got[i], k)
		}
	}
}

func TestApply_Reverse(t *testing.T) {
	f := makeFile("A", "B", "C")
	out := sorter.Apply(f, sorter.Options{Alphabetical: true, Reverse: true})
	got := keys(out)
	want := []string{"C", "B", "A"}
	for i, k := range want {
		if got[i] != k {
			t.Errorf("pos %d: got %q want %q", i, got[i], k)
		}
	}
}

func TestApply_GroupByPrefix(t *testing.T) {
	f := makeFile("DB_HOST", "APP_NAME", "DB_PORT", "APP_ENV")
	out := sorter.Apply(f, sorter.Options{GroupByPrefix: true})
	got := keys(out)
	// DB_ group comes first (insertion order), then APP_
	if got[0] != "DB_HOST" || got[1] != "DB_PORT" {
		t.Errorf("expected DB group first, got %v", got)
	}
	if got[2] != "APP_NAME" || got[3] != "APP_ENV" {
		t.Errorf("expected APP group second, got %v", got)
	}
}

func TestApply_NoOptions_ReturnsCopy(t *testing.T) {
	f := makeFile("Z", "A", "M")
	out := sorter.Apply(f, sorter.Options{})
	if len(out.Entries) != len(f.Entries) {
		t.Fatalf("length mismatch")
	}
	for i, e := range f.Entries {
		if out.Entries[i].Key != e.Key {
			t.Errorf("pos %d changed without options", i)
		}
	}
}

func TestApply_AlphabeticalAndGroup(t *testing.T) {
	f := makeFile("DB_PORT", "APP_NAME", "DB_HOST", "APP_ENV")
	out := sorter.Apply(f, sorter.Options{Alphabetical: true, GroupByPrefix: true})
	got := keys(out)
	// After alpha sort: APP_ENV, APP_NAME, DB_HOST, DB_PORT
	// Group by prefix preserves APP then DB order from alpha-sorted input.
	if got[0] != "APP_ENV" || got[1] != "APP_NAME" {
		t.Errorf("expected APP group first after alpha+group, got %v", got)
	}
	if got[2] != "DB_HOST" || got[3] != "DB_PORT" {
		t.Errorf("expected DB group second after alpha+group, got %v", got)
	}
}
