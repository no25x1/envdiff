package typecheck_test

import (
	"testing"

	"github.com/envdiff/envdiff/internal/parser"
	"github.com/envdiff/envdiff/internal/typecheck"
)

func makeFile(pairs ...string) *parser.EnvFile {
	f := &parser.EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		f.Entries = append(f.Entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return f
}

func rules(key string, t typecheck.Type) []typecheck.Rule {
	return []typecheck.Rule{{Pattern: "^" + key + "$", Type: t}}
}

func TestCheck_NilFile(t *testing.T) {
	if v := typecheck.Check(nil, rules("X", typecheck.TypeInt)); len(v) != 0 {
		t.Fatalf("expected no violations, got %d", len(v))
	}
}

func TestCheck_ValidInt(t *testing.T) {
	v := typecheck.Check(makeFile("PORT", "8080"), rules("PORT", typecheck.TypeInt))
	if len(v) != 0 {
		t.Fatalf("unexpected violation: %v", v)
	}
}

func TestCheck_InvalidInt(t *testing.T) {
	v := typecheck.Check(makeFile("PORT", "abc"), rules("PORT", typecheck.TypeInt))
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Key != "PORT" {
		t.Errorf("expected key PORT, got %s", v[0].Key)
	}
}

func TestCheck_ValidBool(t *testing.T) {
	for _, val := range []string{"true", "false", "1", "0", "True", "FALSE"} {
		v := typecheck.Check(makeFile("ENABLED", val), rules("ENABLED", typecheck.TypeBool))
		if len(v) != 0 {
			t.Errorf("value %q should be valid bool, got violation", val)
		}
	}
}

func TestCheck_InvalidBool(t *testing.T) {
	v := typecheck.Check(makeFile("ENABLED", "yes"), rules("ENABLED", typecheck.TypeBool))
	if len(v) != 1 {
		t.Fatalf("expected violation for 'yes', got %d", len(v))
	}
}

func TestCheck_ValidURL(t *testing.T) {
	v := typecheck.Check(makeFile("API_URL", "https://example.com/path"), rules("API_URL", typecheck.TypeURL))
	if len(v) != 0 {
		t.Fatalf("unexpected violation: %v", v)
	}
}

func TestCheck_InvalidURL(t *testing.T) {
	v := typecheck.Check(makeFile("API_URL", "not-a-url"), rules("API_URL", typecheck.TypeURL))
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
}

func TestCheck_ValidIP(t *testing.T) {
	v := typecheck.Check(makeFile("HOST_IP", "192.168.1.1"), rules("HOST_IP", typecheck.TypeIP))
	if len(v) != 0 {
		t.Fatalf("unexpected violation: %v", v)
	}
}

func TestCheck_ValidEmail(t *testing.T) {
	v := typecheck.Check(makeFile("ADMIN_EMAIL", "admin@example.com"), rules("ADMIN_EMAIL", typecheck.TypeEmail))
	if len(v) != 0 {
		t.Fatalf("unexpected violation: %v", v)
	}
}

func TestCheck_InvalidEmail(t *testing.T) {
	v := typecheck.Check(makeFile("ADMIN_EMAIL", "not-an-email"), rules("ADMIN_EMAIL", typecheck.TypeEmail))
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
}

func TestCheck_UnmatchedKeySkipped(t *testing.T) {
	v := typecheck.Check(makeFile("OTHER", "abc"), rules("PORT", typecheck.TypeInt))
	if len(v) != 0 {
		t.Fatalf("unmatched key should be skipped, got %d violations", len(v))
	}
}

func TestCheck_MultipleViolations(t *testing.T) {
	file := makeFile("PORT", "bad", "ENABLED", "maybe", "HOST", "valid-string")
	rs := []typecheck.Rule{
		{Pattern: "^PORT$", Type: typecheck.TypeInt},
		{Pattern: "^ENABLED$", Type: typecheck.TypeBool},
		{Pattern: "^HOST$", Type: typecheck.TypeString},
	}
	v := typecheck.Check(file, rs)
	if len(v) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(v))
	}
}
