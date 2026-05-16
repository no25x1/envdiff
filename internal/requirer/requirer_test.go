package requirer_test

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/requirer"
)

func makeFile(pairs ...string) *parser.EnvFile {
	f := &parser.EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		f.Entries = append(f.Entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return f
}

func TestApply_NilFile_ReturnsNil(t *testing.T) {
	v, err := requirer.Apply(nil, requirer.Options{Keys: []string{"FOO"}})
	if err != nil || v != nil {
		t.Fatalf("expected nil, nil; got %v, %v", v, err)
	}
}

func TestApply_NoRequiredKeys_ReturnsNil(t *testing.T) {
	f := makeFile("FOO", "bar")
	v, err := requirer.Apply(f, requirer.Options{})
	if err != nil || v != nil {
		t.Fatalf("expected nil, nil; got %v, %v", v, err)
	}
}

func TestApply_AllPresent_NoViolations(t *testing.T) {
	f := makeFile("FOO", "bar", "BAZ", "qux")
	v, err := requirer.Apply(f, requirer.Options{Keys: []string{"FOO", "BAZ"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %v", v)
	}
}

func TestApply_MissingKey_ReturnsViolation(t *testing.T) {
	f := makeFile("FOO", "bar")
	v, err := requirer.Apply(f, requirer.Options{Keys: []string{"FOO", "MISSING"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Key != "MISSING" {
		t.Errorf("expected key MISSING, got %s", v[0].Key)
	}
}

func TestApply_EmptyValue_ViolationByDefault(t *testing.T) {
	f := makeFile("FOO", "")
	v, err := requirer.Apply(f, requirer.Options{Keys: []string{"FOO"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 1 {
		t.Fatalf("expected 1 violation for empty value, got %d", len(v))
	}
	if v[0].Key != "FOO" {
		t.Errorf("unexpected key %s", v[0].Key)
	}
}

func TestApply_EmptyValue_AllowedWhenFlagSet(t *testing.T) {
	f := makeFile("FOO", "")
	v, err := requirer.Apply(f, requirer.Options{Keys: []string{"FOO"}, AllowEmpty: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 0 {
		t.Fatalf("expected no violations when AllowEmpty=true, got %v", v)
	}
}

func TestApply_MultipleViolations(t *testing.T) {
	f := makeFile("PRESENT", "val")
	opts := requirer.Options{Keys: []string{"PRESENT", "GONE1", "GONE2"}}
	v, err := requirer.Apply(f, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(v))
	}
}
