// Package digester computes and compares cryptographic digests of env files.
// It supports multiple hash algorithms and can detect whether two env files
// are semantically equivalent regardless of key ordering.
package digester

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Algorithm represents a supported hash algorithm.
type Algorithm string

const (
	AlgorithmSHA256 Algorithm = "sha256"
	AlgorithmSHA1   Algorithm = "sha1"
	AlgorithmMD5    Algorithm = "md5"
)

// Result holds the digest output for an env file.
type Result struct {
	Algorithm Algorithm
	Digest    string
	KeyCount  int
}

// Digest computes a deterministic hash of the given env file using the
// specified algorithm. Keys are sorted before hashing so that ordering
// differences do not affect the output.
func Digest(f *parser.EnvFile, alg Algorithm) (Result, error) {
	if f == nil {
		return Result{}, fmt.Errorf("digester: nil env file")
	}

	h, err := newHash(alg)
	if err != nil {
		return Result{}, err
	}

	entries := make([]parser.Entry, len(f.Entries))
	copy(entries, f.Entries)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	var sb strings.Builder
	for _, e := range entries {
		sb.WriteString(e.Key)
		sb.WriteByte('=')
		sb.WriteString(e.Value)
		sb.WriteByte('\n')
	}

	_, _ = fmt.Fprint(h, sb.String())

	return Result{
		Algorithm: alg,
		Digest:    hex.EncodeToString(h.Sum(nil)),
		KeyCount:  len(entries),
	}, nil
}

// Equal returns true when both env files produce the same digest.
func Equal(a, b *parser.EnvFile, alg Algorithm) (bool, error) {
	ra, err := Digest(a, alg)
	if err != nil {
		return false, err
	}
	rb, err := Digest(b, alg)
	if err != nil {
		return false, err
	}
	return ra.Digest == rb.Digest, nil
}

func newHash(alg Algorithm) (hash.Hash, error) {
	switch alg {
	case AlgorithmSHA256:
		return sha256.New(), nil
	case AlgorithmSHA1:
		return sha1.New(), nil
	case AlgorithmMD5:
		return md5.New(), nil
	default:
		return nil, fmt.Errorf("digester: unsupported algorithm %q", alg)
	}
}
