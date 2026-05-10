package parser

// EnvFile is a map of environment variable names to their string values.
// It is the canonical in-memory representation used throughout envdiff.
type EnvFile map[string]string

// Keys returns a sorted slice of all variable names in the file.
func (e EnvFile) Keys() []string {
	keys := make([]string, 0, len(e))
	for k := range e {
		keys = append(keys, k)
	}
	return keys
}

// Merge returns a new EnvFile that contains all entries from e, with any
// entries from other overwriting matching keys.
func (e EnvFile) Merge(other EnvFile) EnvFile {
	result := make(EnvFile, len(e)+len(other))
	for k, v := range e {
		result[k] = v
	}
	for k, v := range other {
		result[k] = v
	}
	return result
}

// Clone returns a shallow copy of the EnvFile.
func (e EnvFile) Clone() EnvFile {
	c := make(EnvFile, len(e))
	for k, v := range e {
		c[k] = v
	}
	return c
}
