// Package redactor scans an EnvFile and replaces the values of sensitive
// entries with a configurable placeholder string.
//
// Sensitivity is determined by matching the entry key (case-insensitively)
// against a list of patterns. A built-in set of patterns covers common
// secrets such as passwords, tokens, API keys and certificates. Callers
// can supply their own pattern list to override the defaults.
//
// The original EnvFile is never mutated; Apply always returns a new copy.
package redactor
