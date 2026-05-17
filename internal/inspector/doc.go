// Package inspector analyses individual entries within an env file,
// producing per-key metadata including value length, emptiness,
// sensitivity classification, and a best-effort type hint.
//
// Example usage:
//
//	res, err := inspector.Inspect(envFile)
//	for _, r := range res {
//		fmt.Printf("%s: type=%s sensitive=%v\n", r.Key, r.TypeHint, r.IsSensitive)
//	}
package inspector
