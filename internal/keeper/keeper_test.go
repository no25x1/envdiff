package keeper_test

import (
	"testing"

	"envdiff/internal/keeper"
	"envdiff/internal/parser"
)

func makeFile(pairs ...string) *parser.EnvFile {
	f := &parser.EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		f.Entries = append(f.Entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return f
}

func keys(f *parser.EnvFile) []string {
	out := make([]string, len(f.Entries))
	for i, e := range f.Entries {
		out[i] = e.Key
	}
	return out
}

func TestApply_NilFile_ReturnsEmpty(t *testing.T) {
	result, err := keeper.Apply(nil, keeper.Options{Keys: []string{"A"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Entries) != 0 {
		t.Fatalf("expected empty, got %d entries", len(result.Entries))
	}
}

func TestApply_NoOptions_ReturnsError(t *testing.T) {
	src := makeFile("A", "1")
	_, err := keeper.Apply(src, keeper.Options{})
	if err == nil {
		t.Fatal("expected error for empty options")
	}
}

func TestApply_ExplicitKeys_KeepsOnly(t *testing.T) {
	src := makeFile("KEEP", "yes", "DROP", "no", "ALSO_KEEP", "maybe")
	result, err := keeper.Apply(src, keeper.Options{Keys: []string{"KEEP", "ALSO_KEEP"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := keys(result)
	if len(got) != 2 {
		t.Fatalf("expected 2 keys, got %v", got)
	}
	for _, k := range got {
		if k != "KEEP" && k != "ALSO_KEEP" {
			t.Errorf("unexpected key %q in result", k)
		}
	}
}

func TestApply_Prefix_KeepsMatchingKeys(t *testing.T) {
	src := makeFile("APP_HOST", "localhost", "APP_PORT", "8080", "DB_HOST", "db")
	result, err := keeper.Apply(src, keeper.Options{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := keys(result)
	if len(got) != 2 {
		t.Fatalf("expected 2 keys, got %v", got)
	}
}

func TestApply_KeysAndPrefix_UnionRetained(t *testing.T) {
	src := makeFile("APP_HOST", "h", "DB_HOST", "d", "EXTRA", "e")
	result, err := keeper.Apply(src, keeper.Options{Keys: []string{"EXTRA"}, Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result.Entries))
	}
}

func TestCount_ReturnsCorrectNumber(t *testing.T) {
	src := makeFile("A", "1", "B", "2", "C", "3")
	n, err := keeper.Count(src, keeper.Options{Keys: []string{"A", "C"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2, got %d", n)
	}
}
