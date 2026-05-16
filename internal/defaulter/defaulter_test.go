package defaulter_test

import (
	"testing"

	"github.com/user/envdiff/internal/defaulter"
	"github.com/user/envdiff/internal/parser"
)

func makeFile(pairs ...string) *parser.EnvFile {
	f := &parser.EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		f.Entries = append(f.Entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return f
}

func keys(results []defaulter.Result) []string {
	out := make([]string, len(results))
	for i, r := range results {
		out[i] = r.Key
	}
	return out
}

func TestApply_NilDefaults_ReturnsError(t *testing.T) {
	_, _, err := defaulter.Apply(nil, makeFile(), defaulter.Options{})
	if err == nil {
		t.Fatal("expected error for nil defaults")
	}
}

func TestApply_NilDst_TreatedAsEmpty(t *testing.T) {
	def := makeFile("HOST", "localhost")
	out, results, err := defaulter.Apply(def, nil, defaulter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 1 || out.Entries[0].Key != "HOST" {
		t.Fatalf("expected HOST entry, got %v", out.Entries)
	}
	if len(results) != 1 || results[0].Skipped {
		t.Fatalf("expected one non-skipped result, got %v", results)
	}
}

func TestApply_FillsMissingKeys(t *testing.T) {
	def := makeFile("HOST", "localhost", "PORT", "5432")
	dst := makeFile("PORT", "9999")
	out, results, err := defaulter.Apply(def, dst, defaulter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out.Entries))
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	// PORT already existed — should be skipped
	for _, r := range results {
		if r.Key == "PORT" && !r.Skipped {
			t.Fatal("PORT should be skipped")
		}
		if r.Key == "HOST" && r.Skipped {
			t.Fatal("HOST should not be skipped")
		}
	}
}

func TestApply_Overwrite_ReplacesExisting(t *testing.T) {
	def := makeFile("HOST", "localhost")
	dst := makeFile("HOST", "prod.example.com")
	out, results, err := defaulter.Apply(def, dst, defaulter.Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[0].Value != "localhost" {
		t.Fatalf("expected localhost, got %s", out.Entries[0].Value)
	}
	if results[0].Skipped {
		t.Fatal("result should not be skipped when overwrite=true")
	}
}

func TestApply_KeyFilter_OnlyProcessesAllowed(t *testing.T) {
	def := makeFile("HOST", "localhost", "PORT", "5432", "DB", "mydb")
	dst := makeFile()
	out, results, err := defaulter.Apply(def, dst, defaulter.Options{Keys: []string{"HOST", "DB"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out.Entries))
	}
	for _, r := range results {
		if r.Key == "PORT" {
			t.Fatal("PORT should have been filtered out")
		}
	}
	_ = keys(results)
}
