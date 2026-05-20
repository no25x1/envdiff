// Package renamer provides utilities for renaming keys in .env files.
//
// Keys can be renamed via an explicit mapping, by adding a prefix,
// or by stripping an existing prefix. Entries whose keys do not
// match any rule are passed through unchanged.
//
// Example usage:
//
//	result, err := renamer.Apply(file, renamer.Options{
//		Mapping:     map[string]string{"OLD_KEY": "NEW_KEY"},
//		AddPrefix:   "APP_",
//		StripPrefix: "LEGACY_",
//	})
package renamer
