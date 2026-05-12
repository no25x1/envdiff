// Package sorter provides sorting and grouping operations for parsed .env
// files.
//
// Entries can be sorted alphabetically by key, grouped by shared key prefix
// (the segment before the first underscore), and optionally reversed. All
// operations return a new EnvFile, leaving the original unchanged.
//
// Example usage:
//
//	out := sorter.Apply(envFile, sorter.Options{
//		Alphabetical:  true,
//		GroupByPrefix: true,
//	})
package sorter
