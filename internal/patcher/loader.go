package patcher

import (
	"encoding/json"
	"fmt"
	"os"
)

// PatchFile is the JSON structure used to load a patch from disk.
type PatchFile struct {
	Ops []jsonOp `json:"ops"`
}

type jsonOp struct {
	Kind   string `json:"kind"`
	Key    string `json:"key"`
	Value  string `json:"value,omitempty"`
	NewKey string `json:"new_key,omitempty"`
}

// LoadFile reads a JSON patch file and returns a slice of Ops.
func LoadFile(path string) ([]Op, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("patcher: read patch file: %w", err)
	}

	var pf PatchFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("patcher: parse patch file: %w", err)
	}

	ops := make([]Op, 0, len(pf.Ops))
	for _, raw := range pf.Ops {
		ops = append(ops, Op{
			Kind:   OpKind(raw.Kind),
			Key:    raw.Key,
			Value:  raw.Value,
			NewKey: raw.NewKey,
		})
	}
	return ops, nil
}
