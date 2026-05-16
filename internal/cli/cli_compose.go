package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"envdiff/internal/composer"
	"envdiff/internal/exporter"
)

// runCompose implements: envdiff compose [flags] file1 [file2 ...]
//
// Flags:
//
//	-ns     comma-separated NAMESPACE:file mappings, e.g. db:db.env,app:app.env
//	-overwrite  last writer wins on key conflict (default: first wins)
//	-out    output path (default: stdout)
//	-format dotenv|json (default: dotenv)
func runCompose(args []string) error {
	fs := flag.NewFlagSet("compose", flag.ContinueOnError)
	nsFlag := fs.String("ns", "", "namespace:file mappings (comma-separated)")
	overwrite := fs.Bool("overwrite", false, "last writer wins on conflict")
	outPath := fs.String("out", "", "output file path (default: stdout)")
	format := fs.String("format", "dotenv", "output format: dotenv or json")

	if err := fs.Parse(args); err != nil {
		return err
	}

	var sources []composer.Source

	// Parse -ns flag
	if *nsFlag != "" {
		for _, pair := range strings.Split(*nsFlag, ",") {
			parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
			if len(parts) != 2 {
				return fmt.Errorf("compose: invalid -ns entry %q, expected namespace:file", pair)
			}
			sources = append(sources, composer.Source{Namespace: parts[0], Path: parts[1]})
		}
	}

	// Remaining positional args are plain files (no namespace)
	for _, p := range fs.Args() {
		sources = append(sources, composer.Source{Path: p})
	}

	if len(sources) == 0 {
		return fmt.Errorf("compose: at least one source file is required")
	}

	result, err := composer.Apply(composer.Options{
		Sources:             sources,
		OverwriteOnConflict: *overwrite,
	})
	if err != nil {
		return err
	}

	w := os.Stdout
	if *outPath != "" {
		f, err := os.Create(*outPath)
		if err != nil {
			return fmt.Errorf("compose: creating output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	exp, err := exporter.New(*format)
	if err != nil {
		return fmt.Errorf("compose: %w", err)
	}
	return exp.Write(w, result)
}
