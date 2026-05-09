package diff

import (
	"testing"

	"github.com/user/envdiff/internal/parser"
)

func makeEnvFile(path string, pairs map[string]string) *parser.EnvFile {
	env := &parser.EnvFile{
		Path:  path,
		Index: make(map[string]*parser.Entry),
	}
	for k, v := range pairs {
		e := parser.Entry{Key: k, Value: v}
		env.Entries = append(env.Entries, e)
		env.Index[k] = &env.Entries[len(env.Entries)-1]
	}
	return env
}

func TestDiff_Added(t *testing.T) {
	base := makeEnvFile(".env.base", map[string]string{})
	target := makeEnvFile(".env.target", map[string]string{"NEW_KEY": "value"})
	result := Diff(base, target)
	if len(result.Changes) != 1 || result.Changes[0].Type != Added {
		t.Errorf("expected one Added change, got %+v", result.Changes)
	}
}

func TestDiff_Removed(t *testing.T) {
	base := makeEnvFile(".env.base", map[string]string{"OLD_KEY": "value"})
	target := makeEnvFile(".env.target", map[string]string{})
	result := Diff(base, target)
	if len(result.Changes) != 1 || result.Changes[0].Type != Removed {
		t.Errorf("expected one Removed change, got %+v", result.Changes)
	}
}

func TestDiff_Modified(t *testing.T) {
	base := makeEnvFile(".env.base", map[string]string{"KEY": "old"})
	target := makeEnvFile(".env.target", map[string]string{"KEY": "new"})
	result := Diff(base, target)
	if len(result.Changes) != 1 || result.Changes[0].Type != Modified {
		t.Errorf("expected one Modified change, got %+v", result.Changes)
	}
	if result.Changes[0].OldValue != "old" || result.Changes[0].NewValue != "new" {
		t.Errorf("unexpected old/new values: %+v", result.Changes[0])
	}
}

func TestDiff_Unchanged(t *testing.T) {
	base := makeEnvFile(".env.base", map[string]string{"KEY": "same"})
	target := makeEnvFile(".env.target", map[string]string{"KEY": "same"})
	result := Diff(base, target)
	if len(result.Changes) != 1 || result.Changes[0].Type != Unchanged {
		t.Errorf("expected one Unchanged change, got %+v", result.Changes)
	}
}

func TestDiff_Summary(t *testing.T) {
	base := makeEnvFile("base", map[string]string{"A": "1", "B": "2"})
	target := makeEnvFile("target", map[string]string{"A": "changed", "C": "new"})
	result := Diff(base, target)
	added, removed, modified, _ := result.Summary()
	if added != 1 || removed != 1 || modified != 1 {
		t.Errorf("unexpected summary: added=%d removed=%d modified=%d", added, removed, modified)
	}
}
