package cli

import "fmt"

// UsageText returns the full CLI usage string shown to users.
func UsageText() string {
	return `envdiff — diff and reconcile .env files across environments

Usage:
  envdiff diff   [options] <base.env> <target.env>
  envdiff reconcile [options] <base.env> <target.env>

Commands:
  diff        Show differences between two .env files
  reconcile   Merge target keys into base, producing a reconciled .env

Diff options:
  -format string   Output format: text (default) or json
  -config string   Path to envdiff config file

Reconcile options:
  -out string      Write output to file instead of stdout
  -config string   Path to envdiff config file

Examples:
  envdiff diff .env .env.production
  envdiff diff -format json .env .env.staging
  envdiff reconcile -out merged.env .env .env.production
`
}

// versionInfo holds build-time version metadata.
type versionInfo struct {
	Version string
	Commit  string
}

// defaultVersion is the fallback when no build tags are set.
var defaultVersion = versionInfo{Version: "dev", Commit: "none"}

// VersionString formats the version info for display.
func VersionString(v versionInfo) string {
	return fmt.Sprintf("envdiff %s (commit: %s)", v.Version, v.Commit)
}

// DefaultVersionString returns the version string using the default build-time version.
func DefaultVersionString() string {
	return VersionString(defaultVersion)
}
