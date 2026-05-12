package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/envdiff/internal/parser"
	"github.com/envdiff/internal/snapshot"
)

func makeEnvFile(path string, entries []parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Path: path, Entries: entries}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "snap.json")

	env := makeEnvFile(".env", []parser.Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "PORT", Value: "8080"},
	})

	if err := snapshot.Save(env, dest); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Fatal("expected snapshot file to exist")
	}
}

func TestLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "snap.json")

	original := makeEnvFile(".env.prod", []parser.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "SECRET_KEY", Value: "abc123"},
	})

	if err := snapshot.Save(original, dest); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	snap, err := snapshot.Load(dest)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if snap.Source != ".env.prod" {
		t.Errorf("Source = %q, want .env.prod", snap.Source)
	}
	if len(snap.Entries) != 2 {
		t.Errorf("Entries len = %d, want 2", len(snap.Entries))
	}
	if snap.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
	if snap.Timestamp.Location() != time.UTC {
		t.Error("Timestamp should be UTC")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestToEnvFile_ReturnsEnvFile(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "snap.json")

	original := makeEnvFile(".env", []parser.Entry{
		{Key: "FOO", Value: "bar"},
	})

	_ = snapshot.Save(original, dest)
	snap, _ := snapshot.Load(dest)

	env := snap.ToEnvFile()
	if env.Path != ".env" {
		t.Errorf("Path = %q, want .env", env.Path)
	}
	if len(env.Entries) != 1 || env.Entries[0].Key != "FOO" {
		t.Errorf("unexpected entries: %+v", env.Entries)
	}
}
