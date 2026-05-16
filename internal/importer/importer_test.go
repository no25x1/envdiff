package importer_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/envdiff/internal/importer"
)

func writeTempJSON(t *testing.T, data map[string]string) string {
	t.Helper()
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.CreateTemp(t.TempDir(), "*.json")
	if err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(f.Name(), b, 0644)
	return f.Name()
}

func writeTempDotenv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestApply_JSONFormat(t *testing.T) {
	path := writeTempJSON(t, map[string]string{"HOST": "localhost", "PORT": "5432"})
	file, err := importer.Apply(path, importer.Options{Format: importer.FormatJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(file.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(file.Entries))
	}
}

func TestApply_InferredJSONFormat(t *testing.T) {
	path := writeTempJSON(t, map[string]string{"KEY": "val"})
	file, err := importer.Apply(path, importer.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(file.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(file.Entries))
	}
}

func TestApply_PrefixAdded(t *testing.T) {
	path := writeTempJSON(t, map[string]string{"NAME": "alice"})
	file, err := importer.Apply(path, importer.Options{Format: importer.FormatJSON, Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file.Entries[0].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME, got %s", file.Entries[0].Key)
	}
}

func TestApply_DotenvFormat(t *testing.T) {
	path := writeTempDotenv(t, "FOO=bar\nBAZ=qux\n")
	file, err := importer.Apply(path, importer.Options{Format: importer.FormatDotenv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(file.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(file.Entries))
	}
}

func TestApply_MissingFile(t *testing.T) {
	_, err := importer.Apply("/no/such/file.json", importer.Options{})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestApply_UnsupportedFormat(t *testing.T) {
	path := writeTempDotenv(t, "X=1")
	_, err := importer.Apply(path, importer.Options{Format: "yaml"})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}
