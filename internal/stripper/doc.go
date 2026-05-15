// Package stripper provides functionality to remove entries from an EnvFile
// based on explicit key lists, key prefix/suffix patterns, or value predicates
// such as empty-value removal.
//
// It is non-destructive: the source EnvFile is never modified; a new EnvFile
// containing only the retained entries is returned.
//
// Example usage:
//
//	out := stripper.Apply(src, stripper.Options{
//		Prefixes:    []string{"STAGING_"},
//		EmptyValues: true,
//	})
package stripper
