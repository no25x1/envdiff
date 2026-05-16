// Package masker provides utilities for detecting and masking sensitive
// values in .env files.
//
// A Masker is configured with a set of key-name patterns (regular expressions)
// that identify sensitive fields such as passwords, tokens, and API keys.
// When a key matches a pattern its value is replaced with a configurable
// placeholder string (default: "***").
//
// Usage:
//
//	m := masker.New()
//	masked := m.MaskEnv(envFile)
//
// Custom patterns and placeholders can be provided via NewWithOptions.
package masker
