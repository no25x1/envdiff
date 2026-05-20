package cli

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/exporter"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/renamer"
)

func runRename(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envdiff rename <file> [--add-prefix PREFIX] [--strip-prefix PREFIX] [--map OLD=NEW,...] [--output dotenv|json]")
	}

	filePath := args[0]

	var addPrefix, stripPrefix, mapFlag, outputFmt string
	for i := 1; i < len(args)-1; i++ {
		switch args[i] {
		case "--add-prefix":
			addPrefix = args[i+1]
			i++
		case "--strip-prefix":
			stripPrefix = args[i+1]
			i++
		case "--map":
			mapFlag = args[i+1]
			i++
		case "--output":
			outputFmt = args[i+1]
			i++
		}
	}

	if addPrefix == "" && stripPrefix == "" && mapFlag == "" {
		return fmt.Errorf("rename: at least one of --add-prefix, --strip-prefix, or --map must be specified")
	}

	file, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("rename: %w", err)
	}

	opts := renamer.Options{
		AddPrefix:   addPrefix,
		StripPrefix: stripPrefix,
	}

	if mapFlag != "" {
		opts.Mapping = parseRenameMapping(mapFlag)
	}

	result, err := renamer.Apply(file, opts)
	if err != nil {
		return fmt.Errorf("rename: %w", err)
	}

	if outputFmt == "" {
		outputFmt = "dotenv"
	}

	exp, err := exporter.New(outputFmt)
	if err != nil {
		return fmt.Errorf("rename: %w", err)
	}

	return exp.Write(result, nil)
}

func parseRenameMapping(raw string) map[string]string {
	m := make(map[string]string)
	for _, pair := range strings.Split(raw, ",") {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) == 2 {
			m[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return m
}
