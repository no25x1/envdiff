// Package pinpointer scans an EnvFile for values that appear to be pinned to
// a specific environment — such as loopback addresses, raw IPs, or strings
// containing environment tags like "dev", "staging", or "prod".
//
// It is useful as a pre-promotion check to surface entries that likely need
// updating before an env file is promoted to a new environment.
//
// Usage:
//
//	findings, err := pinpointer.Apply(envFile)
//	for _, f := range findings {
//	    fmt.Printf("%s = %q  (%s)\n", f.Key, f.Value, f.Reason)
//	}
package pinpointer
