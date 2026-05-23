package cascader_test

import (
	"testing"

	"envdiff/internal/cascader"
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

func TestApply_NoFiles_ReturnsError(t *testing.T) {
	_, err := cascader.Apply(nil, cascader.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for empty file list")
	}
}

func TestApply_NilFile_ReturnsError(t *testing.T) {
	_, err := cascader.Apply([]*parser.EnvFile{nil}, cascader.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for nil file")
	}
}

func TestApply_SingleFile_ReturnsCopy(t *testing.T) {
	base := makeFile("A", "1", "B", "2")
	out, err := cascader.Apply([]*parser.EnvFile{base}, cascader.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out.Entries))
	}
}

func TestApply_OverlayOverridesMatchingKey(t *testing.T) {
	base := makeFile("A", "base", "B", "base")
	overlay := makeFile("A", "overlay")
	out, err := cascader.Apply([]*parser.EnvFile{base, overlay}, cascader.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out.Entries {
		if e.Key == "A" && e.Value != "overlay" {
			t.Errorf("expected A=overlay, got %s", e.Value)
		}
		if e.Key == "B" && e.Value != "base" {
			t.Errorf("expected B=base, got %s", e.Value)
		}
	}
}

func TestApply_OverlayNewKeyIgnored(t *testing.T) {
	base := makeFile("A", "1")
	overlay := makeFile("A", "2", "Z", "99")
	out, err := cascader.Apply([]*parser.EnvFile{base, overlay}, cascader.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out.Entries {
		if e.Key == "Z" {
			t.Error("key Z from overlay should be ignored")
		}
	}
}

func TestApply_StopOnMissing_ReturnsError(t *testing.T) {
	base := makeFile("A", "1", "B", "2")
	overlay := makeFile("A", "x")
	opts := cascader.DefaultOptions()
	opts.StopOnMissing = true
	_, err := cascader.Apply([]*parser.EnvFile{base, overlay}, opts)
	if err == nil {
		t.Fatal("expected error when key B missing from overlay")
	}
}

func TestApply_MultipleOverlays(t *testing.T) {
	base := makeFile("A", "1", "B", "2", "C", "3")
	l1 := makeFile("A", "l1")
	l2 := makeFile("B", "l2", "C", "l2")
	out, err := cascader.Apply([]*parser.EnvFile{base, l1, l2}, cascader.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expect := map[string]string{"A": "l1", "B": "l2", "C": "l2"}
	for _, e := range out.Entries {
		if v, ok := expect[e.Key]; ok && e.Value != v {
			t.Errorf("key %s: expected %s, got %s", e.Key, v, e.Value)
		}
	}
}
