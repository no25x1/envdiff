// Package masker provides utilities for detecting and masking sensitive
// environment variable values before display or export.
//
// A Masker instance can be created with default sensitive key patterns (such as
// PASSWORD, SECRET, TOKEN, KEY, etc.) or with custom patterns supplied by the
// caller.
//
// Example usage:
//
//	m := masker.New()
//	masked := m.MaskEnv(envFile)
//
// Sensitive detection is case-insensitive and pattern-based. Values are replaced
// with a configurable placeholder (default: "****").
package masker
