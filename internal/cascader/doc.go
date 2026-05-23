// Package cascader implements layered env-file cascading.
//
// It accepts an ordered slice of parsed env files and produces a single
// merged file where the first file defines the authoritative set of keys.
// Each subsequent "overlay" layer may override values for keys that already
// exist in the base; keys that are new in an overlay are silently dropped,
// preserving the schema of the base file.
//
// Typical use-case: apply environment-specific overrides (.env.production)
// on top of a shared base (.env) without accidentally introducing new keys.
package cascader
