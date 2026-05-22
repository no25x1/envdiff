// Package masker provides utilities for detecting and masking sensitive
// values in .env files.
//
// Keys matching patterns such as PASSWORD, SECRET, TOKEN, API_KEY, AUTH, and
// PRIVATE_KEY are considered sensitive. Their values are replaced with a
// configurable placeholder (default "***") when masking is applied.
//
// Usage:
//
//	m := masker.New()
//	masked := m.MaskEnv(envFile)
//
// Custom patterns and placeholders can be supplied via NewWithOptions.
package masker
