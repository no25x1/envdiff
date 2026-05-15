package patcher_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envdiff/internal/patcher"
)

func writeTempPatch(t *testing.T, ops []map[string]string) string {
	t.Helper()
	data, err := json.Marshal(map[string]interface{}{"ops": ops})
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(t.TempDir(), "patch.json")
	if err := os.WriteFile(p, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadFile_SetOp(t *testing.T) {
	p := writeTempPatch(t, []map[string]string{
		{"kind": "set", "key": "PORT", "value": "9090"},
	})
	ops, err := patcher.LoadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(ops) != 1 {
		t.Fatalf("expected 1 op, got %d", len(ops))
	}
	if ops[0].Kind != patcher.OpSet || ops[0].Key != "PORT" || ops[0].Value != "9090" {
		t.Errorf("unexpected op: %+v", ops[0])
	}
}

func TestLoadFile_MultipleOps(t *testing.T) {
	p := writeTempPatch(t, []map[string]string{
		{"kind": "set", "key": "A", "value": "1"},
		{"kind": "delete", "key": "B"},
		{"kind": "rename", "key": "C", "new_key": "D"},
	})
	ops, err := patcher.LoadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(ops) != 3 {
		t.Fatalf("expected 3 ops, got %d", len(ops))
	}
	if ops[1].Kind != patcher.OpDelete {
		t.Errorf("expected delete op, got %s", ops[1].Kind)
	}
	if ops[2].NewKey != "D" {
		t.Errorf("expected new_key D, got %s", ops[2].NewKey)
	}
}

func TestLoadFile_MissingFile(t *testing.T) {
	_, err := patcher.LoadFile("/nonexistent/patch.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFile_InvalidJSON(t *testing.T) {
	p := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(p, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := patcher.LoadFile(p)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
