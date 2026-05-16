package freezer_test

import (
	"testing"

	"envdiff/internal/freezer"
	"envdiff/internal/parser"
)

func makeFile(entries ...parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(k, v string) parser.Entry { return parser.Entry{Key: k, Value: v} }

func TestFreeze_NilFile(t *testing.T) {
	if got := freezer.Freeze(nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestFreeze_ReturnsSortedEntries(t *testing.T) {
	f := makeFile(entry("ZEBRA", "1"), entry("ALPHA", "2"), entry("MIDDLE", "3"))
	got := freezer.Freeze(f)
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
	keys := []string{got[0].Key, got[1].Key, got[2].Key}
	want := []string{"ALPHA", "MIDDLE", "ZEBRA"}
	for i, k := range want {
		if keys[i] != k {
			t.Errorf("position %d: want %s, got %s", i, k, keys[i])
		}
	}
}

func TestFreeze_HashIsDeterministic(t *testing.T) {
	f := makeFile(entry("KEY", "secret"))
	a := freezer.Freeze(f)
	b := freezer.Freeze(f)
	if a[0].Hash != b[0].Hash {
		t.Errorf("hashes differ: %s vs %s", a[0].Hash, b[0].Hash)
	}
}

func TestCheck_NilInputs(t *testing.T) {
	if got := freezer.Check(nil, nil); got != nil {
		t.Fatalf("expected nil violations")
	}
}

func TestCheck_NoViolations(t *testing.T) {
	f := makeFile(entry("DB_PASS", "hunter2"), entry("API_KEY", "abc"))
	frozen := freezer.Freeze(f)
	violations := freezer.Check(f, frozen)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestCheck_DetectsChangedValue(t *testing.T) {
	original := makeFile(entry("DB_PASS", "hunter2"))
	frozen := freezer.Freeze(original)

	modified := makeFile(entry("DB_PASS", "newpassword"))
	violations := freezer.Check(modified, frozen)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "DB_PASS" {
		t.Errorf("wrong key: %s", violations[0].Key)
	}
}

func TestCheck_IgnoresNewKeys(t *testing.T) {
	original := makeFile(entry("EXISTING", "v1"))
	frozen := freezer.Freeze(original)

	extended := makeFile(entry("EXISTING", "v1"), entry("NEW_KEY", "v2"))
	violations := freezer.Check(extended, frozen)
	if len(violations) != 0 {
		t.Errorf("expected no violations for new keys, got %v", violations)
	}
}
