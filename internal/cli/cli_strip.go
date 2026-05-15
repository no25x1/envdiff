package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/exporter"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/stripper"
)

func runStrip(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff strip <file> [--keys K1,K2] [--prefixes P1,P2] [--suffixes S1,S2] [--empty] [--out <file>]")
	}

	filePath := args[0]
	rest := args[1:]

	var keyList, prefixList, suffixList []string
	var stripEmpty bool
	outPath := ""
	format := "dotenv"

	for i := 0; i < len(rest); i++ {
		switch rest[i] {
		case "--keys":
			if i+1 >= len(rest) {
				return fmt.Errorf("--keys requires a value")
			}
			i++
			keyList = strings.Split(rest[i], ",")
		case "--prefixes":
			if i+1 >= len(rest) {
				return fmt.Errorf("--prefixes requires a value")
			}
			i++
			prefixList = strings.Split(rest[i], ",")
		case "--suffixes":
			if i+1 >= len(rest) {
				return fmt.Errorf("--suffixes requires a value")
			}
			i++
			suffixList = strings.Split(rest[i], ",")
		case "--empty":
			stripEmpty = true
		case "--out":
			if i+1 >= len(rest) {
				return fmt.Errorf("--out requires a value")
			}
			i++
			outPath = rest[i]
		case "--format":
			if i+1 >= len(rest) {
				return fmt.Errorf("--format requires a value")
			}
			i++
			format = rest[i]
		}
	}

	src, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parse %s: %w", filePath, err)
	}

	result := stripper.Apply(src, stripper.Options{
		Keys:        keyList,
		Prefixes:    prefixList,
		Suffixes:    suffixList,
		EmptyValues: stripEmpty,
	})

	w := os.Stdout
	if outPath != "" {
		f, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("create output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	exp, err := exporter.New(format, w)
	if err != nil {
		return err
	}
	return exp.Write(result)
}
