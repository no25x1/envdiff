// Package pinpointer identifies entries whose values look like they contain
// environment-specific tokens (URLs, IPs, hostnames) that may need updating
// when promoting across environments.
package pinpointer

import (
	"net"
	"net/url"
	"regexp"
	"strings"

	"envdiff/internal/parser"
)

var (
	hostPattern = regexp.MustCompile(`(?i)(localhost|127\.0\.0\.1|0\.0\.0\.0)`)
	envTagPattern = regexp.MustCompile(`(?i)(dev|staging|stage|prod|production|local|test)`)
)

// Finding describes a single entry that contains an environment-pinned value.
type Finding struct {
	Key    string
	Value  string
	Reason string
}

// Apply scans the given EnvFile and returns all findings where values appear
// to be pinned to a specific environment.
func Apply(f *parser.EnvFile) ([]Finding, error) {
	if f == nil {
		return nil, nil
	}

	var findings []Finding
	for _, e := range f.Entries {
		if reason, ok := inspect(e.Value); ok {
			findings = append(findings, Finding{
				Key:    e.Key,
				Value:  e.Value,
				Reason: reason,
			})
		}
	}
	return findings, nil
}

func inspect(val string) (string, bool) {
	if val == "" {
		return "", false
	}

	// Check for localhost / loopback addresses.
	if hostPattern.MatchString(val) {
		return "contains loopback/localhost reference", true
	}

	// Check for raw IP addresses.
	if ip := net.ParseIP(strings.TrimSpace(val)); ip != nil {
		return "value is a raw IP address", true
	}

	// Check for URLs with environment-tagged hostnames.
	if u, err := url.Parse(val); err == nil && u.Host != "" {
		host := strings.ToLower(u.Hostname())
		if envTagPattern.MatchString(host) {
			return "URL hostname contains environment tag", true
		}
	}

	// Check for plain env-tagged strings (e.g. "staging-api", "prod_db").
	if envTagPattern.MatchString(val) {
		return "value contains environment tag", true
	}

	return "", false
}
