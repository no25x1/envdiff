// Package profiler provides statistical analysis of EnvFile contents.
//
// It computes key counts, empty value ratios, sensitive key detection,
// and prefix-based grouping to give operators a quick overview of an
// environment file without exposing secret values.
//
// Usage:
//
//	f, _ := parser.ParseFile("production.env")
//	profile := profiler.Analyse(f)
//	fmt.Println(profile.TotalKeys, profile.SensitiveKeys)
package profiler
