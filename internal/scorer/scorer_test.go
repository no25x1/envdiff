package scorer_test

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/scorer"
)

func makeFile(entries ...parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(key, value string) parser.Entry {
	return parser.Entry{Key: key, Value: value}
}

func TestScore_NilFile(t *testing.T) {
	r := scorer.Score(nil)
	if r.Score != 0 {
		t.Fatalf("expected 0 for nil file, got %d", r.Score)
	}
}

func TestScore_EmptyFile(t *testing.T) {
	r := scorer.Score(makeFile())
	if r.Score != 0 {
		t.Fatalf("expected 0 for empty file, got %d", r.Score)
	}
}

func TestScore_PerfectFile(t *testing.T) {
	f := makeFile(
		entry("APP_NAME", "myapp"),
		entry("APP_PORT", "8080"),
		entry("DB_HOST", "localhost"),
	)
	r := scorer.Score(f)
	if r.Score != 100 {
		t.Fatalf("expected 100 for perfect file, got %d", r.Score)
	}
	if r.Max != 100 {
		t.Fatalf("expected max 100, got %d", r.Max)
	}
}

func TestScore_EmptyValuesPenalised(t *testing.T) {
	f := makeFile(
		entry("APP_NAME", ""),
		entry("APP_PORT", ""),
	)
	r := scorer.Score(f)
	if r.Score >= 100 {
		t.Fatalf("expected score < 100 when values are empty, got %d", r.Score)
	}
}

func TestScore_LowercaseKeysPenalised(t *testing.T) {
	f := makeFile(
		entry("app_name", "myapp"),
		entry("app_port", "8080"),
	)
	r := scorer.Score(f)
	if r.Score >= 100 {
		t.Fatalf("expected score < 100 for lowercase keys, got %d", r.Score)
	}
}

func TestScore_InlineSecretPenalised(t *testing.T) {
	f := makeFile(
		entry("APP_SECRET", "supersecret123"),
		entry("APP_PORT", "8080"),
	)
	r := scorer.Score(f)
	if r.Score >= 100 {
		t.Fatalf("expected score < 100 for inline secret, got %d", r.Score)
	}
}

func TestScore_SecretWithPlaceholderNotPenalised(t *testing.T) {
	f := makeFile(
		entry("APP_SECRET", "${SECRET_FROM_VAULT}"),
		entry("APP_PORT", "8080"),
	)
	r := scorer.Score(f)
	if r.Score < 80 {
		t.Fatalf("expected high score for placeholder secret, got %d", r.Score)
	}
}

func TestScore_BreakdownLength(t *testing.T) {
	f := makeFile(entry("KEY", "val"))
	r := scorer.Score(f)
	if len(r.Breakdown) != 4 {
		t.Fatalf("expected 4 breakdown rules, got %d", len(r.Breakdown))
	}
}

func TestScore_WhitespaceKeyPenalised(t *testing.T) {
	f := makeFile(
		entry("BAD KEY", "value"),
	)
	r := scorer.Score(f)
	if r.Score >= 100 {
		t.Fatalf("expected score < 100 for key with whitespace, got %d", r.Score)
	}
}
