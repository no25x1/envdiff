package cli_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/envdiff/envdiff/internal/cli"
)

func writeTempEnvSchema(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func writeTempSchemaFile(t *testing.T, rules []map[string]interface{}) string {
	t.Helper()
	data, _ := json.Marshal(map[string]interface{}{"rules": rules})
	path := filepath.Join(t.TempDir(), "schema.json")
	_ = os.WriteFile(path, data, 0644)
	return path
}

func TestRun_SchemaMissingArgs(t *testing.T) {
	err := cli.Run([]string{"schema"})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRun_SchemaMissingSchemaFile(t *testing.T) {
	env := writeTempEnvSchema(t, "APP_ENV=test\n")
	err := cli.Run([]string{"schema", "/no/such/schema.json", env})
	if err == nil {
		t.Fatal("expected error for missing schema file")
	}
}

func TestRun_SchemaValidFile(t *testing.T) {
	schemaPath := writeTempSchemaFile(t, []map[string]interface{}{
		{"key": "APP_ENV", "required": true},
	})
	env := writeTempEnvSchema(t, "APP_ENV=production\n")
	if err := cli.Run([]string{"schema", schemaPath, env}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_SchemaInvalidFile(t *testing.T) {
	schemaPath := writeTempSchemaFile(t, []map[string]interface{}{
		{"key": "DB_URL", "required": true},
	})
	env := writeTempEnvSchema(t, "APP_ENV=production\n")
	err := cli.Run([]string{"schema", schemaPath, env})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestRun_SchemaPatternFailure(t *testing.T) {
	schemaPath := writeTempSchemaFile(t, []map[string]interface{}{
		{"key": "PORT", "required": true, "pattern": `^\d+$`},
	})
	env := writeTempEnvSchema(t, "PORT=abc\n")
	if err := cli.Run([]string{"schema", schemaPath, env}); err == nil {
		t.Fatal("expected pattern violation error")
	}
}
