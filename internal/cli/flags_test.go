package cli

import (
	"strings"
	"testing"
)

func TestUsageText_ContainsCommands(t *testing.T) {
	usage := UsageText()
	for _, keyword := range []string{"diff", "reconcile", "format", "config", "out"} {
		if !strings.Contains(usage, keyword) {
			t.Errorf("usage text missing keyword %q", keyword)
		}
	}
}

func TestVersionString_Format(t *testing.T) {
	v := versionInfo{Version: "1.2.3", Commit: "abc123"}
	s := VersionString(v)
	if !strings.Contains(s, "1.2.3") {
		t.Errorf("expected version in string, got: %s", s)
	}
	if !strings.Contains(s, "abc123") {
		t.Errorf("expected commit in string, got: %s", s)
	}
}

func TestVersionString_Dev(t *testing.T) {
	s := VersionString(defaultVersion)
	if !strings.Contains(s, "dev") {
		t.Errorf("expected dev version, got: %s", s)
	}
}
