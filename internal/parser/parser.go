package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair from a .env file.
type Entry struct {
	Key     string
	Value   string
	Comment string
	LineNum int
}

// EnvFile holds all parsed entries from a .env file.
type EnvFile struct {
	Path    string
	Entries []Entry
	Index   map[string]*Entry
}

// Parse reads and parses a .env file at the given path.
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()

	env := &EnvFile{
		Path:  path,
		Index: make(map[string]*Entry),
	}

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		entry, err := parseLine(line, lineNum)
		if err != nil {
			return nil, fmt.Errorf("%s:%d: %w", path, lineNum, err)
		}

		env.Entries = append(env.Entries, entry)
		env.Index[entry.Key] = &env.Entries[len(env.Entries)-1]
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning %s: %w", path, err)
	}

	return env, nil
}

// Keys returns a slice of all keys in the order they appear in the file.
func (e *EnvFile) Keys() []string {
	keys := make([]string, len(e.Entries))
	for i, entry := range e.Entries {
		keys[i] = entry.Key
	}
	return keys
}

func parseLine(line string, lineNum int) (Entry, error) {
	comment := ""
	if idx := strings.Index(line, " #"); idx != -1 {
		comment = strings.TrimSpace(line[idx+2:])
		line = strings.TrimSpace(line[:idx])
	}

	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return Entry{}, fmt.Errorf("invalid line: %q", line)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.Trim(strings.TrimSpace(parts[1]), `"`)

	if key == "" {
		return Entry{}, fmt.Errorf("empty key in line: %q", line)
	}

	return Entry{Key: key, Value: value, Comment: comment, LineNum: lineNum}, nil
}
