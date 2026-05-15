// Package typecheck validates that env file values conform to declared types.
package typecheck

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/envdiff/envdiff/internal/parser"
)

// Type represents a declared value type.
type Type string

const (
	TypeString  Type = "string"
	TypeInt     Type = "int"
	TypeFloat   Type = "float"
	TypeBool    Type = "bool"
	TypeURL     Type = "url"
	TypeIP      Type = "ip"
	TypeEmail   Type = "email"
)

// Rule maps a key pattern to an expected type.
type Rule struct {
	Pattern string
	Type    Type
}

// Violation describes a type mismatch found in an env file.
type Violation struct {
	Key      string
	Value    string
	Expected Type
	Reason   string
}

func (v Violation) Error() string {
	return fmt.Sprintf("key %q value %q: expected %s — %s", v.Key, v.Value, v.Expected, v.Reason)
}

var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

// Check validates entries in file against the provided rules.
// Keys not matched by any rule are skipped.
func Check(file *parser.EnvFile, rules []Rule) []Violation {
	if file == nil {
		return nil
	}
	var violations []Violation
	for _, entry := range file.Entries {
		for _, rule := range rules {
			matched, _ := regexp.MatchString("(?i)"+rule.Pattern, entry.Key)
			if !matched {
				continue
			}
			if reason, ok := validate(entry.Value, rule.Type); !ok {
				violations = append(violations, Violation{
					Key:      entry.Key,
					Value:    entry.Value,
					Expected: rule.Type,
					Reason:   reason,
				})
			}
			break
		}
	}
	return violations
}

func validate(value string, t Type) (string, bool) {
	switch t {
	case TypeInt:
		if _, err := strconv.ParseInt(value, 10, 64); err != nil {
			return "not a valid integer", false
		}
	case TypeFloat:
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return "not a valid float", false
		}
	case TypeBool:
		v := strings.ToLower(value)
		if v != "true" && v != "false" && v != "1" && v != "0" {
			return "not a valid bool (true/false/1/0)", false
		}
	case TypeURL:
		u, err := url.ParseRequestURI(value)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return "not a valid URL", false
		}
	case TypeIP:
		if net.ParseIP(value) == nil {
			return "not a valid IP address", false
		}
	case TypeEmail:
		if !emailRe.MatchString(value) {
			return "not a valid email address", false
		}
	}
	return "", true
}
