package differ_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/parser"
)

func makeFile(entries [][2]string) *parser.EnvFile {
	f := &parser.EnvFile{}
	for _, kv := range entries {
		f.Entries = append(f.Entries, parser.Entry{Key: kv[0], Value: kv[1]})
	}
	return f
}

func TestCompare_Unchanged(t *testing.T) {
	base := makeFile([][2]string{{"FOO", "bar"}, {"BAZ", "qux"}})
	target := makeFile([][2]string{{"FOO", "bar"}, {"BAZ", "qux"}})
	r := differ.Compare(base, target)
	if r.HasChanges() {
		t.Error("expected no changes")
	}
	if len(r.Lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(r.Lines))
	}
}

func TestCompare_Added(t *testing.T) {
	base := makeFile([][2]string{{"FOO", "bar"}})
	target := makeFile([][2]string{{"FOO", "bar"}, {"NEW", "val"}})
	r := differ.Compare(base, target)
	if !r.HasChanges() {
		t.Error("expected changes")
	}
	var found bool
	for _, l := range r.Lines {
		if l.Key == "NEW" && l.Change == differ.Added {
			found = true
		}
	}
	if !found {
		t.Error("expected NEW to be Added")
	}
}

func TestCompare_Removed(t *testing.T) {
	base := makeFile([][2]string{{"FOO", "bar"}, {"OLD", "val"}})
	target := makeFile([][2]string{{"FOO", "bar"}})
	r := differ.Compare(base, target)
	var found bool
	for _, l := range r.Lines {
		if l.Key == "OLD" && l.Change == differ.Removed {
			found = true
		}
	}
	if !found {
		t.Error("expected OLD to be Removed")
	}
}

func TestCompare_Modified(t *testing.T) {
	base := makeFile([][2]string{{"FOO", "old"}})
	target := makeFile([][2]string{{"FOO", "new"}})
	r := differ.Compare(base, target)
	if len(r.Lines) != 1 || r.Lines[0].Change != differ.Modified {
		t.Errorf("expected Modified, got %v", r.Lines)
	}
	if r.Lines[0].Before != "old" || r.Lines[0].After != "new" {
		t.Errorf("unexpected before/after: %+v", r.Lines[0])
	}
}

func TestSummary_Format(t *testing.T) {
	base := makeFile([][2]string{{"A", "1"}, {"B", "old"}})
	target := makeFile([][2]string{{"B", "new"}, {"C", "3"}})
	r := differ.Compare(base, target)
	s := r.Summary()
	if !strings.HasPrefix(s, "+") {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestFormat_ContainsSymbols(t *testing.T) {
	base := makeFile([][2]string{{"A", "1"}, {"B", "old"}})
	target := makeFile([][2]string{{"B", "new"}, {"C", "3"}})
	r := differ.Compare(base, target)
	out := differ.Format(r)
	if !strings.Contains(out, "- A=") {
		t.Errorf("expected removed line, got:\n%s", out)
	}
	if !strings.Contains(out, "~ B:") {
		t.Errorf("expected modified line, got:\n%s", out)
	}
	if !strings.Contains(out, "+ C=") {
		t.Errorf("expected added line, got:\n%s", out)
	}
}
