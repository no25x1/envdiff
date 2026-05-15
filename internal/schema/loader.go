package schema

import (
	"encoding/json"
	"fmt"
	"os"
)

// schemaFile mirrors the JSON structure on disk.
type schemaFile struct {
	Rules []struct {
		Key      string `json:"key"`
		Required bool   `json:"required"`
		Pattern  string `json:"pattern,omitempty"`
	} `json:"rules"`
}

// LoadFile reads a JSON schema definition from path and returns a Schema.
func LoadFile(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("schema: read %q: %w", path, err)
	}
	var sf schemaFile
	if err := json.Unmarshal(data, &sf); err != nil {
		return nil, fmt.Errorf("schema: parse %q: %w", path, err)
	}
	s := &Schema{}
	for _, r := range sf.Rules {
		s.Rules = append(s.Rules, FieldRule{
			Key:      r.Key,
			Required: r.Required,
			Pattern:  r.Pattern,
		})
	}
	return s, nil
}
