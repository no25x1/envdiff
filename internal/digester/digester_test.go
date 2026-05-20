package digester_test

import (
	"testing"

	"github.com/user/envdiff/internal/digester"
	"github.com/user/envdiff/internal/parser"
)

func makeFile(entries ...parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(key, value string) parser.Entry {
	return parser.Entry{Key: key, Value: value}
}

func TestDigest_NilFile(t *testing.T) {
	_, err := digester.Digest(nil, digester.AlgorithmSHA256)
	if err == nil {
		t.Fatal("expected error for nil file")
	}
}

func TestDigest_UnsupportedAlgorithm(t *testing.T) {
	f := makeFile(entry("KEY", "val"))
	_, err := digester.Digest(f, digester.Algorithm("blake2"))
	if err == nil {
		t.Fatal("expected error for unsupported algorithm")
	}
}

func TestDigest_SHA256_Deterministic(t *testing.T) {
	f := makeFile(entry("A", "1"), entry("B", "2"))
	r1, err := digester.Digest(f, digester.AlgorithmSHA256)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r2, err := digester.Digest(f, digester.AlgorithmSHA256)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r1.Digest != r2.Digest {
		t.Errorf("expected same digest, got %q and %q", r1.Digest, r2.Digest)
	}
}

func TestDigest_OrderIndependent(t *testing.T) {
	f1 := makeFile(entry("A", "1"), entry("B", "2"))
	f2 := makeFile(entry("B", "2"), entry("A", "1"))
	r1, _ := digester.Digest(f1, digester.AlgorithmSHA256)
	r2, _ := digester.Digest(f2, digester.AlgorithmSHA256)
	if r1.Digest != r2.Digest {
		t.Errorf("expected order-independent digest, got %q and %q", r1.Digest, r2.Digest)
	}
}

func TestDigest_DifferentValues_DifferentDigest(t *testing.T) {
	f1 := makeFile(entry("KEY", "foo"))
	f2 := makeFile(entry("KEY", "bar"))
	r1, _ := digester.Digest(f1, digester.AlgorithmSHA256)
	r2, _ := digester.Digest(f2, digester.AlgorithmSHA256)
	if r1.Digest == r2.Digest {
		t.Error("expected different digests for different values")
	}
}

func TestDigest_KeyCount(t *testing.T) {
	f := makeFile(entry("X", "1"), entry("Y", "2"), entry("Z", "3"))
	r, err := digester.Digest(f, digester.AlgorithmMD5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.KeyCount != 3 {
		t.Errorf("expected KeyCount=3, got %d", r.KeyCount)
	}
}

func TestEqual_SameContent(t *testing.T) {
	f1 := makeFile(entry("A", "1"), entry("B", "2"))
	f2 := makeFile(entry("B", "2"), entry("A", "1"))
	ok, err := digester.Equal(f1, f2, digester.AlgorithmSHA1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected files to be equal")
	}
}

func TestEqual_DifferentContent(t *testing.T) {
	f1 := makeFile(entry("A", "1"))
	f2 := makeFile(entry("A", "2"))
	ok, err := digester.Equal(f1, f2, digester.AlgorithmSHA256)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected files to be unequal")
	}
}
