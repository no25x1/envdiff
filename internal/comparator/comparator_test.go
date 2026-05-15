package comparator_test

import (
	"testing"

	"github.com/envdiff/internal/comparator"
	"github.com/envdiff/internal/parser"
)

func makeFile(pairs ...string) *parser.EnvFile {
	ef := &parser.EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		ef.Entries = append(ef.Entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return ef
}

func TestCompare_UnionOfKeys(t *testing.T) {
	files := map[string]*parser.EnvFile{
		"dev":  makeFile("APP_ENV", "dev", "DB_HOST", "localhost"),
		"prod": makeFile("APP_ENV", "prod", "SECRET", "x"),
	}
	r := comparator.Compare(files)
	if len(r.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d: %v", len(r.Keys), r.Keys)
	}
}

func TestCompare_MatrixValues(t *testing.T) {
	files := map[string]*parser.EnvFile{
		"dev":  makeFile("APP_ENV", "dev"),
		"prod": makeFile("APP_ENV", "prod"),
	}
	r := comparator.Compare(files)
	if r.Matrix["APP_ENV"]["dev"] != "dev" {
		t.Errorf("expected dev, got %q", r.Matrix["APP_ENV"]["dev"])
	}
	if r.Matrix["APP_ENV"]["prod"] != "prod" {
		t.Errorf("expected prod, got %q", r.Matrix["APP_ENV"]["prod"])
	}
}

func TestCompare_AbsentKeys(t *testing.T) {
	files := map[string]*parser.EnvFile{
		"dev":  makeFile("ONLY_DEV", "1"),
		"prod": makeFile("ONLY_PROD", "2"),
	}
	r := comparator.Compare(files)
	if len(r.Absent["ONLY_DEV"]) != 1 || r.Absent["ONLY_DEV"][0] != "prod" {
		t.Errorf("expected ONLY_DEV absent in prod, got %v", r.Absent["ONLY_DEV"])
	}
}

func TestConflicts_DetectsValueMismatch(t *testing.T) {
	files := map[string]*parser.EnvFile{
		"dev":  makeFile("APP_ENV", "dev", "SHARED", "same"),
		"prod": makeFile("APP_ENV", "prod", "SHARED", "same"),
	}
	r := comparator.Compare(files)
	conflicts := r.Conflicts()
	if len(conflicts) != 1 || conflicts[0] != "APP_ENV" {
		t.Errorf("expected [APP_ENV] conflict, got %v", conflicts)
	}
}

func TestConflicts_NoFalsePositives(t *testing.T) {
	files := map[string]*parser.EnvFile{
		"dev":  makeFile("KEY", "val"),
		"prod": makeFile("KEY", "val"),
	}
	r := comparator.Compare(files)
	if len(r.Conflicts()) != 0 {
		t.Errorf("expected no conflicts, got %v", r.Conflicts())
	}
}

func TestCompare_SortedFileLabels(t *testing.T) {
	files := map[string]*parser.EnvFile{
		"z": makeFile("K", "1"),
		"a": makeFile("K", "2"),
	}
	r := comparator.Compare(files)
	if r.Files[0] != "a" || r.Files[1] != "z" {
		t.Errorf("expected sorted labels [a z], got %v", r.Files)
	}
}
