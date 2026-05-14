package interpolator_test

import (
	"testing"

	"github.com/user/envdiff/internal/interpolator"
	"github.com/user/envdiff/internal/parser"
)

func makeFile(entries ...parser.Entry) parser.EnvFile {
	return parser.EnvFile{Path: "test.env", Entries: entries}
}

func entry(k, v string) parser.Entry { return parser.Entry{Key: k, Value: v} }

func TestApply_NoReferences(t *testing.T) {
	f := makeFile(entry("FOO", "bar"), entry("BAZ", "qux"))
	out, err := interpolator.Apply(f, interpolator.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[0].Value != "bar" || out.Entries[1].Value != "qux" {
		t.Errorf("values changed unexpectedly: %+v", out.Entries)
	}
}

func TestApply_CurlyBraceRef(t *testing.T) {
	f := makeFile(entry("BASE", "/app"), entry("DATA", "${BASE}/data"))
	out, err := interpolator.Apply(f, interpolator.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out.Entries[1].Value; got != "/app/data" {
		t.Errorf("expected /app/data, got %q", got)
	}
}

func TestApply_BareDollarRef(t *testing.T) {
	f := makeFile(entry("HOST", "localhost"), entry("URL", "http://$HOST:8080"))
	out, err := interpolator.Apply(f, interpolator.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out.Entries[1].Value; got != "http://localhost:8080" {
		t.Errorf("expected http://localhost:8080, got %q", got)
	}
}

func TestApply_OverrideWins(t *testing.T) {
	f := makeFile(entry("HOST", "localhost"), entry("URL", "http://${HOST}"))
	opts := interpolator.Options{Overrides: map[string]string{"HOST": "prod.example.com"}}
	out, err := interpolator.Apply(f, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out.Entries[1].Value; got != "http://prod.example.com" {
		t.Errorf("expected prod.example.com URL, got %q", got)
	}
}

func TestApply_MissingRef_SilentByDefault(t *testing.T) {
	f := makeFile(entry("URL", "http://${MISSING_HOST}"))
	out, err := interpolator.Apply(f, interpolator.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out.Entries[0].Value; got != "http://" {
		t.Errorf("expected empty substitution, got %q", got)
	}
}

func TestApply_MissingRef_FailOnMissing(t *testing.T) {
	f := makeFile(entry("URL", "http://${MISSING_HOST}"))
	_, err := interpolator.Apply(f, interpolator.Options{FailOnMissing: true})
	if err == nil {
		t.Fatal("expected error for missing reference, got nil")
	}
}

func TestApply_PreservesComments(t *testing.T) {
	f := makeFile(parser.Entry{Key: "FOO", Value: "bar", Comment: "# important"})
	out, err := interpolator.Apply(f, interpolator.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[0].Comment != "# important" {
		t.Errorf("comment was lost: %q", out.Entries[0].Comment)
	}
}
