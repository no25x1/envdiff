package shadower

import (
	"fmt"
	"strings"
)

// parseImportLine extracts the alias and path from a single import line.
// It handles both aliased imports (e.g., `alias "path/to/pkg"`) and
// plain imports (e.g., `"path/to/pkg"`).
func parseImportLine(line string) (alias, path string, err error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return "", "", fmt.Errorf("empty import line")
	}

	parts := strings.Fields(line)
	switch len(parts) {
	case 1:
		path = strings.Trim(parts[0], `"`)
		return "", path, nil
	case 2:
		alias = parts[0]
		path = strings.Trim(parts[1], `"`)
		return alias, path, nil
	default:
		return "", "", fmt.Errorf("invalid import line: %q", line)
	}
}
