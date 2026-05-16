// Package importer reads env entries from external formats (JSON, YAML, TOML-style)
// and converts them into an EnvFile for use with other envdiff tools.
package importer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/your-org/envdiff/internal/parser"
)

// Format represents a supported import format.
type Format string

const (
	FormatJSON   Format = "json"
	FormatDotenv Format = "dotenv"
)

// Options controls import behaviour.
type Options struct {
	Format    Format
	Prefix    string // optional key prefix to add on import
	Overwrite bool   // if true, duplicate keys use the new value
}

// Apply reads the file at path and returns a populated EnvFile.
// Format is inferred from the file extension when opts.Format is empty.
func Apply(path string, opts Options) (*parser.EnvFile, error) {
	if opts.Format == "" {
		opts.Format = inferFormat(path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("importer: read %s: %w", path, err)
	}

	switch opts.Format {
	case FormatJSON:
		return fromJSON(data, opts)
	case FormatDotenv:
		return parser.Parse(path)
	default:
		return nil, fmt.Errorf("importer: unsupported format %q", opts.Format)
	}
}

func inferFormat(path string) Format {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		return FormatJSON
	default:
		return FormatDotenv
	}
}

func fromJSON(data []byte, opts Options) (*parser.EnvFile, error) {
	var raw map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("importer: parse JSON: %w", err)
	}

	keys := make([]string, 0, len(raw))
	for k := range raw {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	file := &parser.EnvFile{}
	for _, k := range keys {
		key := opts.Prefix + k
		file.Entries = append(file.Entries, parser.Entry{
			Key:   key,
			Value: raw[k],
		})
	}
	return file, nil
}
