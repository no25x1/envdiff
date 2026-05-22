// Package shadower identifies keys in a target .env file that override
// (shadow) keys defined in a base .env file. This is useful when layering
// environment configurations — e.g. a shared base with per-environment
// overrides — to make implicit overrides explicit and auditable.
//
// Usage:
//
//	results, err := shadower.Apply(baseFile, nextFile, shadower.DefaultOptions())
//	for _, s := range results {
//		fmt.Printf("%s: %q -> %q\n", s.Key, s.BaseValue, s.NextValue)
//	}
package shadower
