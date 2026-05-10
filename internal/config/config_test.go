package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/config"
)

func writeTempConfig(t *testing.T, cfg *config.Config) string {
	t.Helper()
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}
	path := filepath.Join(t.TempDir(), "envdiff.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func TestDefault_Values(t *testing.T) {
	cfg := config.Default()
	if cfg.OutputFormat != "text" {
		t.Errorf("expected output_format=text, got %q", cfg.OutputFormat)
	}
	if !cfg.MaskSecrets {
		t.Error("expected mask_secrets=true by default")
	}
	if cfg.ReconcileMode {
		t.Error("expected reconcile_mode=false by default")
	}
}

func TestLoadFile_Valid(t *testing.T) {
	src := &config.Config{
		BaseFile:     ".env.example",
		TargetFile:   ".env",
		OutputFormat: "json",
		MaskSecrets:  false,
	}
	path := writeTempConfig(t, src)

	cfg, err := config.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BaseFile != src.BaseFile {
		t.Errorf("base_file: want %q, got %q", src.BaseFile, cfg.BaseFile)
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("output_format: want json, got %q", cfg.OutputFormat)
	}
}

func TestLoadFile_Missing(t *testing.T) {
	_, err := config.LoadFile("/nonexistent/path/envdiff.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestValidate_MissingBaseFile(t *testing.T) {
	cfg := config.Default()
	cfg.TargetFile = ".env"
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for missing base_file")
	}
}

func TestValidate_InvalidOutputFormat(t *testing.T) {
	cfg := config.Default()
	cfg.BaseFile = ".env.example"
	cfg.TargetFile = ".env"
	cfg.OutputFormat = "yaml"
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for unsupported output_format")
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := config.Default()
	cfg.BaseFile = ".env.example"
	cfg.TargetFile = ".env"
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}
