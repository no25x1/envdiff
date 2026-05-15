package cli

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/exporter"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/promote"
)

// runPromote implements the `envdiff promote <src> <dst>` sub-command.
// Flags:
//
//	--keys=A,B,C   only promote listed keys
//	--overwrite    replace existing keys in dst
//	--format=dotenv|json  output format (default: dotenv)
//	--out=<file>   write result to file instead of stdout
func runPromote(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("promote requires <source> and <destination> arguments")
	}

	srcPath := args[0]
	dstPath := args[1]

	var keyList []string
	overwrite := false
	format := "dotenv"
	outPath := ""

	for _, a := range args[2:] {
		switch {
		case strings.HasPrefix(a, "--keys="):
			raw := strings.TrimPrefix(a, "--keys=")
			for _, k := range strings.Split(raw, ",") {
				if k != "" {
					keyList = append(keyList, k)
				}
			}
		case a == "--overwrite":
			overwrite = true
		case strings.HasPrefix(a, "--format="):
			format = strings.TrimPrefix(a, "--format=")
		case strings.HasPrefix(a, "--out="):
			outPath = strings.TrimPrefix(a, "--out=")
		}
	}

	src, err := parser.Parse(srcPath)
	if err != nil {
		return fmt.Errorf("promote: reading source: %w", err)
	}

	dst, err := parser.Parse(dstPath)
	if err != nil {
		// dst may not exist yet — start with empty file.
		dst = parser.EnvFile{Path: dstPath}
	}

	out, results, err := promote.Apply(src, dst, promote.Options{
		Keys:      keyList,
		Overwrite: overwrite,
	})
	if err != nil {
		return err
	}

	for _, r := range results {
		fmt.Printf("  [%s] %s\n", r.Action, r.Key)
	}

	exp, err := exporter.New(format)
	if err != nil {
		return err
	}

	dest := outPath
	if dest == "" {
		dest = "-" // stdout sentinel handled by exporter
	}
	return exp.Write(out, dest)
}
