package cloner_test

import (
	"testing"

	"envdiff/internal/cloner"
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

func TestApply_NilSource_ReturnsError(t *testing.T) {
	_, err := cloner.Apply(nil, makeFile(), cloner.Options{})
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestApply_NilDest_TreatedAsEmpty(t *testing.T) {
	src := makeFile("FOO", "bar")
	out, err := cloner.Apply(src, nil, cloner.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 1 || out.Entries[0].Key != "FOO" {
		t.Fatalf("expected FOO, got %v", keys(out))
	}
}

func TestApply_ClonesAllKeys(t *testing.T) {
	src := makeFile("A", "1", "B", "2")
	dst := makeFile("C", "3")
	out, err := cloner.Apply(src, dst, cloner.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out.Entries))
	}
}

func TestApply_KeyFilter(t *testing.T) {
	src := makeFile("A", "1", "B", "2", "C", "3")
	out, err := cloner.Apply(src, nil, cloner.Options{Keys: []string{"A", "C"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2, got %d: %v", len(out.Entries), keys(out))
	}
}

func TestApply_StripAndAddPrefix(t *testing.T) {
	src := makeFile("PROD_HOST", "example.com", "PROD_PORT", "443")
	out, err := cloner.Apply(src, nil, cloner.Options{
		StripPrefix: "PROD_",
		AddPrefix:   "STAGING_",
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.Entries[0].Key != "STAGING_HOST" {
		t.Errorf("expected STAGING_HOST, got %s", out.Entries[0].Key)
	}
	if out.Entries[1].Key != "STAGING_PORT" {
		t.Errorf("expected STAGING_PORT, got %s", out.Entries[1].Key)
	}
}

func TestApply_NoOverwriteByDefault(t *testing.T) {
	src := makeFile("FOO", "new")
	dst := makeFile("FOO", "old")
	out, err := cloner.Apply(src, dst, cloner.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if out.Entries[0].Value != "old" {
		t.Errorf("expected old value preserved, got %s", out.Entries[0].Value)
	}
}

func TestApply_OverwriteEnabled(t *testing.T) {
	src := makeFile("FOO", "new")
	dst := makeFile("FOO", "old")
	out, err := cloner.Apply(src, dst, cloner.Options{Overwrite: true})
	if err != nil {
		t.Fatal(err)
	}
	if out.Entries[0].Value != "new" {
		t.Errorf("expected new value, got %s", out.Entries[0].Value)
	}
}
