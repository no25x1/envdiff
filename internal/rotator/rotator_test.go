package rotator_test

import (
	"testing"

	"github.com/envdiff/envdiff/internal/parser"
	"github.com/envdiff/envdiff/internal/rotator"
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
	res, err := rotator.Apply(nil, rotator.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.File.Entries) != 0 {
		t.Errorf("expected empty file, got %d entries", len(res.File.Entries))
	}
}

func TestApply_NoMappings_ReturnsCopy(t *testing.T) {
	src := makeFile("DB_HOST", "localhost", "DB_PORT", "5432")
	res, err := rotator.Apply(src, rotator.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.File.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.File.Entries))
	}
	if len(res.Rotated) != 0 {
		t.Errorf("expected no rotated keys, got %v", res.Rotated)
	}
}

func TestApply_RotatesKey(t *testing.T) {
	src := makeFile("OLD_KEY", "value1", "KEEP", "value2")
	opts := rotator.Options{
		Mappings: []rotator.Mapping{{OldKey: "OLD_KEY", NewKey: "NEW_KEY"}},
	}
	res, err := rotator.Apply(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ks := keys(res.File)
	if ks[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY, got %s", ks[0])
	}
	if res.File.Entries[0].Value != "value1" {
		t.Errorf("expected value preserved, got %s", res.File.Entries[0].Value)
	}
	if len(res.Rotated) != 1 || res.Rotated[0] != "NEW_KEY" {
		t.Errorf("unexpected rotated list: %v", res.Rotated)
	}
}

func TestApply_MissingKey_ReportedNotError(t *testing.T) {
	src := makeFile("KEEP", "v")
	opts := rotator.Options{
		Mappings: []rotator.Mapping{{OldKey: "MISSING", NewKey: "NEW"}},
	}
	res, err := rotator.Apply(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "MISSING" {
		t.Errorf("expected missing list [MISSING], got %v", res.Missing)
	}
}

func TestApply_FailOnMissing_ReturnsError(t *testing.T) {
	src := makeFile("KEEP", "v")
	opts := rotator.Options{
		Mappings:      []rotator.Mapping{{OldKey: "GONE", NewKey: "NEW"}},
		FailOnMissing: true,
	}
	_, err := rotator.Apply(src, opts)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestApply_MultipleRotations(t *testing.T) {
	src := makeFile("A", "1", "B", "2", "C", "3")
	opts := rotator.Options{
		Mappings: []rotator.Mapping{
			{OldKey: "A", NewKey: "ALPHA"},
			{OldKey: "C", NewKey: "GAMMA"},
		},
	}
	res, err := rotator.Apply(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ks := keys(res.File)
	if ks[0] != "ALPHA" || ks[1] != "B" || ks[2] != "GAMMA" {
		t.Errorf("unexpected keys: %v", ks)
	}
	if len(res.Rotated) != 2 {
		t.Errorf("expected 2 rotated, got %d", len(res.Rotated))
	}
}
