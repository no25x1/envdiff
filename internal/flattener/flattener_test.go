package flattener_test

import (
	"testing"

	"github.com/user/envdiff/internal/flattener"
	"github.com/user/envdiff/internal/parser"
)

func makeFile(path string, pairs ...string) *parser.EnvFile {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return &parser.EnvFile{Path: path, Entries: entries}
}

func TestApply_NoFiles_ReturnsEmpty(t *testing.T) {
	res, err := flattener.Apply(nil, flattener.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.File.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(res.File.Entries))
	}
}

func TestApply_SingleFile_ReturnsCopy(t *testing.T) {
	f := makeFile("a.env", "FOO", "bar", "BAZ", "qux")
	res, err := flattener.Apply([]*parser.EnvFile{f}, flattener.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.File.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.File.Entries))
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(res.Conflicts))
	}
}

func TestApply_StrategyLast_OverridesValue(t *testing.T) {
	a := makeFile("a.env", "FOO", "from-a", "ONLY_A", "1")
	b := makeFile("b.env", "FOO", "from-b", "ONLY_B", "2")

	res, err := flattener.Apply([]*parser.EnvFile{a, b}, flattener.Options{Strategy: flattener.StrategyLast})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	vals := entryMap(res.File)
	if vals["FOO"] != "from-b" {
		t.Errorf("expected FOO=from-b, got %q", vals["FOO"])
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0].Key != "FOO" {
		t.Errorf("expected 1 conflict for FOO, got %+v", res.Conflicts)
	}
}

func TestApply_StrategyFirst_KeepsOriginal(t *testing.T) {
	a := makeFile("a.env", "FOO", "from-a")
	b := makeFile("b.env", "FOO", "from-b")

	res, err := flattener.Apply([]*parser.EnvFile{a, b}, flattener.Options{Strategy: flattener.StrategyFirst})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entryMap(res.File)["FOO"] != "from-a" {
		t.Errorf("expected FOO=from-a (first wins)")
	}
	if res.Conflicts[0].Chosen != "from-a" {
		t.Errorf("conflict Chosen should reflect first value")
	}
}

func TestApply_StrategyError_ReturnsError(t *testing.T) {
	a := makeFile("a.env", "FOO", "1")
	b := makeFile("b.env", "FOO", "2")

	_, err := flattener.Apply([]*parser.EnvFile{a, b}, flattener.Options{Strategy: flattener.StrategyError})
	if err == nil {
		t.Fatal("expected error for duplicate key with StrategyError")
	}
}

func TestApply_NoConflicts_UniqueKeys(t *testing.T) {
	a := makeFile("a.env", "A", "1")
	b := makeFile("b.env", "B", "2")

	res, err := flattener.Apply([]*parser.EnvFile{a, b}, flattener.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.File.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.File.Entries))
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts")
	}
}

func entryMap(f *parser.EnvFile) map[string]string {
	m := make(map[string]string, len(f.Entries))
	for _, e := range f.Entries {
		m[e.Key] = e.Value
	}
	return m
}
