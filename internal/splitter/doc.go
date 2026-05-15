// Package splitter provides functionality to split a single EnvFile into
// multiple EnvFile instances, each containing keys that share a common prefix.
//
// This is useful when a large monolithic .env file needs to be broken apart
// into service-specific files (e.g. DB_, AWS_, APP_) for deployment or
// secret-management purposes.
//
// Example usage:
//
//	result, err := splitter.Apply(src, splitter.Options{
//		Prefixes: map[string]string{
//			"db":  "DB_",
//			"aws": "AWS_",
//		},
//		IncludeUnmatched: true,
//	})
package splitter
