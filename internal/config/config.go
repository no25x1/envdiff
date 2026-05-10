package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the runtime configuration for envdiff.
type Config struct {
	// BaseFile is the reference .env file (e.g. .env.example)
	BaseFile string `json:"base_file"`
	// TargetFile is the environment-specific .env file to compare
	TargetFile string `json:"target_file"`
	// OutputFormat controls reporter output: "text" or "json"
	OutputFormat string `json:"output_format"`
	// MaskSecrets enables secret masking in output
	MaskSecrets bool `json:"mask_secrets"`
	// SensitivePatterns are additional regex patterns to treat as secrets
	SensitivePatterns []string `json:"sensitive_patterns,omitempty"`
	// ReconcileMode writes missing keys from base into target when true
	ReconcileMode bool `json:"reconcile_mode"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		OutputFormat: "text",
		MaskSecrets:  true,
		ReconcileMode: false,
	}
}

// LoadFile reads a JSON config file from the given path.
func LoadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: reading file %q: %w", path, err)
	}
	cfg := Default()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("config: parsing JSON: %w", err)
	}
	return cfg, nil
}

// Validate checks that required fields are set and values are valid.
func (c *Config) Validate() error {
	if c.BaseFile == "" {
		return fmt.Errorf("config: base_file must not be empty")
	}
	if c.TargetFile == "" {
		return fmt.Errorf("config: target_file must not be empty")
	}
	switch c.OutputFormat {
	case "text", "json":
		// valid
	default:
		return fmt.Errorf("config: unsupported output_format %q (want \"text\" or \"json\")", c.OutputFormat)
	}
	return nil
}
