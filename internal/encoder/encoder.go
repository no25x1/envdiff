// Package encoder provides value encoding/decoding utilities for .env entries,
// supporting base64 and URL encoding strategies.
package encoder

import (
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/your-org/envdiff/internal/parser"
)

// Strategy defines the encoding algorithm to apply.
type Strategy string

const (
	Base64  Strategy = "base64"
	URLEnc  Strategy = "url"
)

// Options controls how encoding/decoding is performed.
type Options struct {
	Strategy Strategy
	Decode   bool // if true, decode instead of encode
	KeysOnly []string // if non-empty, only process these keys
}

// Apply encodes or decodes values in f according to opts.
func Apply(f parser.EnvFile, opts Options) (parser.EnvFile, error) {
	out := parser.EnvFile{Path: f.Path}
	for _, e := range f.Entries {
		if len(opts.KeysOnly) > 0 && !contains(opts.KeysOnly, e.Key) {
			out.Entries = append(out.Entries, e)
			continue
		}
		v, err := transform(e.Value, opts)
		if err != nil {
			return parser.EnvFile{}, fmt.Errorf("key %s: %w", e.Key, err)
		}
		out.Entries = append(out.Entries, parser.Entry{Key: e.Key, Value: v, Raw: e.Raw})
	}
	return out, nil
}

func transform(value string, opts Options) (string, error) {
	switch opts.Strategy {
	case Base64:
		if opts.Decode {
			b, err := base64.StdEncoding.DecodeString(value)
			if err != nil {
				return "", fmt.Errorf("base64 decode: %w", err)
			}
			return string(b), nil
		}
		return base64.StdEncoding.EncodeToString([]byte(value)), nil
	case URLEnc:
		if opts.Decode {
			v, err := url.QueryUnescape(value)
			if err != nil {
				return "", fmt.Errorf("url decode: %w", err)
			}
			return v, nil
		}
		return url.QueryEscape(value), nil
	default:
		return "", fmt.Errorf("unknown strategy %q", opts.Strategy)
	}
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
