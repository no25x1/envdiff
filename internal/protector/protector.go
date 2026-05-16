// Package protector provides key-level write protection for .env files.
// Keys marked as protected cannot be modified or deleted during reconcile
// or patch operations without an explicit override.
package protector

import (
	"fmt"
	"strings"

	"envdiff/internal/parser"
)

// Options controls protection behaviour.
type Options struct {
	// Keys is an explicit list of keys to protect.
	Keys []string
	// Prefixes protects any key whose name starts with one of these prefixes.
	Prefixes []string
	// AllowOverride disables enforcement and only emits warnings.
	AllowOverride bool
}

// Violation describes a single protection breach.
type Violation struct {
	Key    string
	Reason string
}

func (v Violation) Error() string {
	return fmt.Sprintf("protected key %q: %s", v.Key, v.Reason)
}

// Apply checks whether any entries in next would overwrite or remove
// protected keys defined in base. It returns the list of violations found.
// When opts.AllowOverride is false, the first violation also causes an error.
func Apply(base, next *parser.EnvFile, opts Options) ([]Violation, error) {
	if base == nil || next == nil {
		return nil, nil
	}

	protected := buildSet(opts)

	baseMap := make(map[string]string, len(base.Entries))
	for _, e := range base.Entries {
		baseMap[e.Key] = e.Value
	}

	nextMap := make(map[string]string, len(next.Entries))
	for _, e := range next.Entries {
		nextMap[e.Key] = e.Value
	}

	var violations []Violation
	for key, baseVal := range baseMap {
		if !isProtected(key, protected, opts.Prefixes) {
			continue
		}
		nextVal, exists := nextMap[key]
		if !exists {
			violations = append(violations, Violation{Key: key, Reason: "key removed"})
		} else if nextVal != baseVal {
			violations = append(violations, Violation{Key: key, Reason: fmt.Sprintf("value changed from %q to %q", baseVal, nextVal)})
		}
	}

	if len(violations) > 0 && !opts.AllowOverride {
		return violations, violations[0]
	}
	return violations, nil
}

func buildSet(opts Options) map[string]struct{} {
	s := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		s[k] = struct{}{}
	}
	return s
}

func isProtected(key string, set map[string]struct{}, prefixes []string) bool {
	if _, ok := set[key]; ok {
		return true
	}
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}
