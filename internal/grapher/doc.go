// Package grapher analyses variable reference relationships within an
// EnvFile and produces a directed acyclic graph (DAG) of dependencies.
//
// It supports:
//   - Detecting which variables reference other variables via $VAR or ${VAR}
//   - Topological ordering so variables can be resolved in dependency order
//   - Cycle detection, returning an error when circular references exist
//
// Example usage:
//
//	g, err := grapher.Build(envFile)
//	if err != nil {
//	    log.Fatalf("cycle: %v", err)
//	}
//	fmt.Println(g.Order) // safe evaluation order
package grapher
