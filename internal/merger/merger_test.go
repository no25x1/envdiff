package merger_test

import (
	"testing"

	"github.com/user/envdiff/internal/merger"
	"github.com/user/envdiff/internal/parser"
)

func makeFile(source string, pairs ...string) *parser.EnvFile {
	f := &parser.EnvFile{
		Entries: make(map[string]parser.Entry),
	}
	for i := 0; i+1 < len(pairs); i += 2 {
		key := pairs[i]
		f.Entries[key] = parser.Entry{Key: key, Value: pairs[i+1], Source: source}
		f.Order = append(f.Order, key)
	}
	return f
}

func TestMerge_NoFiles(t *testing.T) {
	res, err := merger.Merge(nil, merger.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.File.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(res.File.Entries))
	}
}

func TestMerge_StrategyLast_OverridesValues(t *testing.T) {
	a := makeFile("a.env", "FOO", "from_a", "BAR", "bar_a")
	b := makeFile("b.env", "FOO", "from_b")

	res, err := merger.Merge([]*parser.EnvFile{a, b}, merger.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := res.File.Entries["FOO"].Value; got != "from_b" {
		t.Errorf("StrategyLast: FOO = %q, want %q", got, "from_b")
	}
	if got := res.File.Entries["BAR"].Value; got != "bar_a" {
		t.Errorf("BAR = %q, want %q", got, "bar_a")
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0].Key != "FOO" {
		t.Errorf("expected 1 conflict on FOO, got %+v", res.Conflicts)
	}
}

func TestMerge_StrategyFirst_KeepsOriginal(t *testing.T) {
	a := makeFile("a.env", "FOO", "from_a")
	b := makeFile("b.env", "FOO", "from_b")

	res, err := merger.Merge([]*parser.EnvFile{a, b}, merger.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := res.File.Entries["FOO"].Value; got != "from_a" {
		t.Errorf("StrategyFirst: FOO = %q, want %q", got, "from_a")
	}
}

func TestMerge_StrategyError_ReturnsError(t *testing.T) {
	a := makeFile("a.env", "SECRET", "abc")
	b := makeFile("b.env", "SECRET", "xyz")

	_, err := merger.Merge([]*parser.EnvFile{a, b}, merger.StrategyError)
	if err == nil {
		t.Fatal("expected error for conflicting key, got nil")
	}
}

func TestMerge_NoConflicts_UniqueKeys(t *testing.T) {
	a := makeFile("a.env", "FOO", "foo")
	b := makeFile("b.env", "BAR", "bar")

	res, err := merger.Merge([]*parser.EnvFile{a, b}, merger.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(res.Conflicts))
	}
	if len(res.File.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.File.Entries))
	}
}

func TestMerge_NilFileSkipped(t *testing.T) {
	a := makeFile("a.env", "FOO", "foo")
	res, err := merger.Merge([]*parser.EnvFile{nil, a, nil}, merger.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.File.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(res.File.Entries))
	}
}
