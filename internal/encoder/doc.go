// Package encoder provides value encoding and decoding for .env file entries.
//
// Supported strategies:
//
//   - base64: standard Base64 encoding (RFC 4648)
//   - url:    URL query-string encoding (percent-encoding)
//
// Use Apply to transform all or a subset of keys in an EnvFile.
// Pass Decode: true in Options to reverse the transformation.
//
// Example:
//
//	out, err := encoder.Apply(f, encoder.Options{
//		Strategy: encoder.Base64,
//		KeysOnly: []string{"DB_PASSWORD", "API_SECRET"},
//	})
package encoder
