// Package patcher applies a set of patch operations (set, delete, rename) to an EnvFile.
package patcher

import (
	"fmt"

	"github.com/yourorg/envdiff/internal/parser"
)

// OpKind describes the type of patch operation.
type OpKind string

const (
	OpSet    OpKind = "set"
	OpDelete OpKind = "delete"
	OpRename OpKind = "rename"
)

// Op is a single patch operation.
type Op struct {
	Kind    OpKind
	Key     string
	Value   string // used by OpSet
	NewKey  string // used by OpRename
}

// Apply applies the given ops to src and returns a new EnvFile.
// Operations are applied in order; an error is returned on the first
// invalid operation (e.g. deleting a key that does not exist).
func Apply(src parser.EnvFile, ops []Op) (parser.EnvFile, error) {
	// Work on a mutable map for fast lookup; preserve insertion order via slice.
	index := make(map[string]int, len(src.Entries))
	entries := make([]parser.Entry, len(src.Entries))
	copy(entries, src.Entries)
	for i, e := range entries {
		index[e.Key] = i
	}

	for _, op := range ops {
		switch op.Kind {
		case OpSet:
			if i, ok := index[op.Key]; ok {
				entries[i].Value = op.Value
			} else {
				entries = append(entries, parser.Entry{Key: op.Key, Value: op.Value})
				index[op.Key] = len(entries) - 1
			}
		case OpDelete:
			i, ok := index[op.Key]
			if !ok {
				return parser.EnvFile{}, fmt.Errorf("patcher: delete: key %q not found", op.Key)
			}
			entries = append(entries[:i], entries[i+1:]...)
			// Rebuild index after deletion.
			index = make(map[string]int, len(entries))
			for j, e := range entries {
				index[e.Key] = j
			}
		case OpRename:
			i, ok := index[op.Key]
			if !ok {
				return parser.EnvFile{}, fmt.Errorf("patcher: rename: key %q not found", op.Key)
			}
			if _, exists := index[op.NewKey]; exists {
				return parser.EnvFile{}, fmt.Errorf("patcher: rename: target key %q already exists", op.NewKey)
			}
			entries[i].Key = op.NewKey
			delete(index, op.Key)
			index[op.NewKey] = i
		default:
			return parser.EnvFile{}, fmt.Errorf("patcher: unknown op kind %q", op.Kind)
		}
	}

	return parser.EnvFile{Path: src.Path, Entries: entries}, nil
}
