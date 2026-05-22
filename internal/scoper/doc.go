// Package scoper narrows an EnvFile to a single named scope (prefix group).
//
// A scope is an uppercase prefix followed by an underscore, e.g. "PROD" matches
// keys such as PROD_HOST, PROD_PORT, PROD_DB_URL.
//
// Typical usage:
//
//	out, err := scoper.Apply(envFile, scoper.Options{
//		Scope:       "PROD",
//		StripPrefix: true,
//	})
//
// With StripPrefix enabled the returned file contains HOST, PORT, DB_URL — ready
// to be used as a clean, scope-agnostic configuration.
package scoper
