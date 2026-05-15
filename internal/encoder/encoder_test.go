package encoder_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/encoder"
	"github.com/your-org/envdiff/internal/parser"
)

func makeFile(entries ...parser.Entry) parser.EnvFile {
	return parser.EnvFile{Path: "test.env", Entries: entries}
}

func entry(k, v string) parser.Entry { return parser.Entry{Key: k, Value: v} }

func TestApply_Base64Encode(t *testing.T) {
	f := makeFile(entry("SECRET", "hello"), entry("OTHER", "world"))
	out, err := encoder.Apply(f, encoder.Options{Strategy: encoder.Base64})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[0].Value != "aGVsbG8=" {
		t.Errorf("expected base64 of 'hello', got %q", out.Entries[0].Value)
	}
	if out.Entries[1].Value != "d29ybGQ=" {
		t.Errorf("expected base64 of 'world', got %q", out.Entries[1].Value)
	}
}

func TestApply_Base64Decode(t *testing.T) {
	f := makeFile(entry("SECRET", "aGVsbG8="))
	out, err := encoder.Apply(f, encoder.Options{Strategy: encoder.Base64, Decode: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[0].Value != "hello" {
		t.Errorf("expected 'hello', got %q", out.Entries[0].Value)
	}
}

func TestApply_URLEncode(t *testing.T) {
	f := makeFile(entry("QUERY", "foo bar&baz"))
	out, err := encoder.Apply(f, encoder.Options{Strategy: encoder.URLEnc})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[0].Value != "foo+bar%26baz" {
		t.Errorf("unexpected url-encoded value: %q", out.Entries[0].Value)
	}
}

func TestApply_URLDecode(t *testing.T) {
	f := makeFile(entry("QUERY", "foo+bar%26baz"))
	out, err := encoder.Apply(f, encoder.Options{Strategy: encoder.URLEnc, Decode: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[0].Value != "foo bar&baz" {
		t.Errorf("unexpected decoded value: %q", out.Entries[0].Value)
	}
}

func TestApply_KeysOnly_SkipsOthers(t *testing.T) {
	f := makeFile(entry("A", "hello"), entry("B", "world"))
	out, err := encoder.Apply(f, encoder.Options{Strategy: encoder.Base64, KeysOnly: []string{"A"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[0].Value != "aGVsbG8=" {
		t.Errorf("expected A to be encoded")
	}
	if out.Entries[1].Value != "world" {
		t.Errorf("expected B to be unchanged, got %q", out.Entries[1].Value)
	}
}

func TestApply_UnknownStrategy(t *testing.T) {
	f := makeFile(entry("X", "val"))
	_, err := encoder.Apply(f, encoder.Options{Strategy: "rot13"})
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestApply_InvalidBase64Decode(t *testing.T) {
	f := makeFile(entry("BAD", "not-valid-base64!!!"))
	_, err := encoder.Apply(f, encoder.Options{Strategy: encoder.Base64, Decode: true})
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}
