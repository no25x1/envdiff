package reconcile_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/reconcile"
)

func makeEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "HOST", Status: diff.Unchanged, SourceVal: "localhost", TargetVal: "localhost"},
		{Key: "PORT", Status: diff.Modified, SourceVal: "3000", TargetVal: "8080"},
		{Key: "OLD_KEY", Status: diff.Removed, SourceVal: "gone"},
		{Key: "NEW_KEY", Status: diff.Added, TargetVal: "fresh"},
	}
}

func TestPlan_Actions(t *testing.T) {
	steps := reconcile.Plan(makeEntries())
	actions := map[string]reconcile.Action{}
	for _, s := range steps {
		actions[s.Key] = s.Action
	}
	if actions["HOST"] != reconcile.ActionKeep {
		t.Errorf("HOST should be KEEP, got %s", actions["HOST"])
	}
	if actions["PORT"] != reconcile.ActionUpdate {
		t.Errorf("PORT should be UPDATE, got %s", actions["PORT"])
	}
	if actions["OLD_KEY"] != reconcile.ActionRemove {
		t.Errorf("OLD_KEY should be REMOVE, got %s", actions["OLD_KEY"])
	}
	if actions["NEW_KEY"] != reconcile.ActionAdd {
		t.Errorf("NEW_KEY should be ADD, got %s", actions["NEW_KEY"])
	}
}

func TestApply_ReconcilesMaps(t *testing.T) {
	source := parser.EnvFile{"HOST": "localhost", "PORT": "3000", "OLD_KEY": "gone"}
	steps := reconcile.Plan(makeEntries())
	result := reconcile.Apply(source, steps)

	if result["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", result["PORT"])
	}
	if _, ok := result["OLD_KEY"]; ok {
		t.Error("OLD_KEY should have been removed")
	}
	if result["NEW_KEY"] != "fresh" {
		t.Errorf("expected NEW_KEY=fresh, got %s", result["NEW_KEY"])
	}
	if result["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %s", result["HOST"])
	}
}

func TestRender_SortedOutput(t *testing.T) {
	env := parser.EnvFile{"Z_KEY": "last", "A_KEY": "first", "M_KEY": "mid"}
	out := reconcile.Render(env)
	expected := "A_KEY=first\nM_KEY=mid\nZ_KEY=last\n"
	if out != expected {
		t.Errorf("unexpected render output:\ngot:  %q\nwant: %q", out, expected)
	}
}
