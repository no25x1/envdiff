package grapher_test

import (
	"testing"

	"github.com/envdiff/envdiff/internal/grapher"
	"github.com/envdiff/envdiff/internal/parser"
)

func makeFile(entries ...parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func entry(key, value string) parser.Entry {
	return parser.Entry{Key: key, Value: value}
}

func TestBuild_NilFile(t *testing.T) {
	g, err := grapher.Build(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.Edges) != 0 || len(g.Order) != 0 {
		t.Error("expected empty graph for nil file")
	}
}

func TestBuild_NoReferences(t *testing.T) {
	f := makeFile(entry("A", "hello"), entry("B", "world"))
	g, err := grapher.Build(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.Edges) != 0 {
		t.Errorf("expected 0 edges, got %d", len(g.Edges))
	}
	if len(g.Order) != 2 {
		t.Errorf("expected 2 nodes in order, got %d", len(g.Order))
	}
}

func TestBuild_SingleReference(t *testing.T) {
	f := makeFile(entry("BASE", "/app"), entry("DIR", "${BASE}/data"))
	g, err := grapher.Build(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(g.Edges))
	}
	e := g.Edges[0]
	if e.Key != "DIR" || e.Dep != "BASE" {
		t.Errorf("unexpected edge: %+v", e)
	}
}

func TestBuild_TopologicalOrder(t *testing.T) {
	f := makeFile(
		entry("C", "${B}-suffix"),
		entry("B", "${A}-mid"),
		entry("A", "root"),
	)
	g, err := grapher.Build(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	pos := make(map[string]int)
	for i, k := range g.Order {
		pos[k] = i
	}
	if pos["A"] >= pos["B"] {
		t.Errorf("expected A before B in order")
	}
	if pos["B"] >= pos["C"] {
		t.Errorf("expected B before C in order")
	}
}

func TestBuild_CycleDetected(t *testing.T) {
	f := makeFile(
		entry("X", "${Y}"),
		entry("Y", "${X}"),
	)
	_, err := grapher.Build(f)
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
}

func TestBuild_SelfReference_NotEdge(t *testing.T) {
	f := makeFile(entry("A", "${A}"))
	g, err := grapher.Build(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(g.Edges) != 0 {
		t.Errorf("self-reference should not create an edge")
	}
}
