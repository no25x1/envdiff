// Package grapher builds a dependency graph from variable references
// within an env file, exposing topological ordering and cycle detection.
package grapher

import (
	"errors"
	"fmt"
	"regexp"
	"sort"

	"github.com/envdiff/envdiff/internal/parser"
)

var refPattern = regexp.MustCompile(`\$\{?([A-Za-z_][A-Za-z0-9_]*)\}?`)

// Edge represents a directed dependency: Key depends on Dep.
type Edge struct {
	Key string
	Dep string
}

// Graph holds the adjacency list and the source entries.
type Graph struct {
	Edges []Edge
	Order []string // topological order, populated by Build
}

// Build constructs a Graph from an EnvFile, returning an error if a cycle
// is detected.
func Build(f *parser.EnvFile) (*Graph, error) {
	if f == nil {
		return &Graph{}, nil
	}

	deps := make(map[string][]string, len(f.Entries))
	for _, e := range f.Entries {
		matches := refPattern.FindAllStringSubmatch(e.Value, -1)
		for _, m := range matches {
			ref := m[1]
			if ref != e.Key {
				deps[e.Key] = append(deps[e.Key], ref)
			}
		}
		if _, ok := deps[e.Key]; !ok {
			deps[e.Key] = nil
		}
	}

	var edges []Edge
	for key, ds := range deps {
		for _, d := range ds {
			edges = append(edges, Edge{Key: key, Dep: d})
		}
	}

	order, err := topoSort(deps)
	if err != nil {
		return nil, err
	}

	return &Graph{Edges: edges, Order: order}, nil
}

// topoSort performs Kahn's algorithm and returns an error on cycle.
func topoSort(deps map[string][]string) ([]string, error) {
	inDegree := make(map[string]int)
	adj := make(map[string][]string)

	for key, ds := range deps {
		if _, ok := inDegree[key]; !ok {
			inDegree[key] = 0
		}
		for _, d := range ds {
			adj[d] = append(adj[d], key)
			inDegree[key]++
			if _, ok := inDegree[d]; !ok {
				inDegree[d] = 0
			}
		}
	}

	var queue []string
	for k, v := range inDegree {
		if v == 0 {
			queue = append(queue, k)
		}
	}
	sort.Strings(queue)

	var order []string
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		order = append(order, n)
		neighbors := adj[n]
		sort.Strings(neighbors)
		for _, nb := range neighbors {
			inDegree[nb]--
			if inDegree[nb] == 0 {
				queue = append(queue, nb)
			}
		}
	}

	if len(order) != len(inDegree) {
		return nil, errors.New(fmt.Sprintf("cycle detected among %d variables", len(inDegree)-len(order)))
	}
	return order, nil
}
