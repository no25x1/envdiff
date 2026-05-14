// Package duplicates detects duplicate keys within or across .env files.
package duplicates

import (
	"fmt"

	"github.com/user/envdiff/internal/parser"
)

// Finding describes a duplicate key occurrence.
type Finding struct {
	Key    string
	File   string
	Line   int
	Source string // "intra" for within-file, "inter" for cross-file
}

// Report holds all duplicate findings from an analysis.
type Report struct {
	Findings []Finding
}

// HasDuplicates returns true if any findings were recorded.
func (r *Report) HasDuplicates() bool {
	return len(r.Findings) > 0
}

// Summary returns a human-readable summary line.
func (r *Report) Summary() string {
	if !r.HasDuplicates() {
		return "no duplicate keys found"
	}
	return fmt.Sprintf("%d duplicate key(s) found", len(r.Findings))
}

// CheckIntra detects duplicate keys within a single EnvFile.
func CheckIntra(file *parser.EnvFile) *Report {
	seen := make(map[string]int) // key -> first line number
	var findings []Finding

	for _, entry := range file.Entries {
		if firstLine, ok := seen[entry.Key]; ok {
			_ = firstLine
			findings = append(findings, Finding{
				Key:    entry.Key,
				File:   file.Path,
				Line:   entry.Line,
				Source: "intra",
			})
		} else {
			seen[entry.Key] = entry.Line
		}
	}

	return &Report{Findings: findings}
}

// CheckInter detects keys that appear in more than one of the provided files.
func CheckInter(files []*parser.EnvFile) *Report {
	// map key -> list of (file, line) where it appears
	type occurrence struct {
		file string
		line int
	}
	index := make(map[string][]occurrence)

	for _, f := range files {
		for _, entry := range f.Entries {
			index[entry.Key] = append(index[entry.Key], occurrence{file: f.Path, line: entry.Line})
		}
	}

	var findings []Finding
	for key, occurrences := range index {
		if len(occurrences) > 1 {
			for _, occ := range occurrences {
				findings = append(findings, Finding{
					Key:    key,
					File:   occ.file,
					Line:   occ.line,
					Source: "inter",
				})
			}
		}
	}

	return &Report{Findings: findings}
}
