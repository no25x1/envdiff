package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"envdiff/internal/parser"
	"envdiff/internal/protector"
)

func runProtect(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envdiff protect <base> <next> [--keys=A,B] [--prefixes=P] [--allow-override] [--format=text|json]")
	}

	basePath := args[0]
	nextPath := args[1]

	var keys, prefixes []string
	allowOverride := false
	format := "text"

	for _, a := range args[2:] {
		switch {
		case strings.HasPrefix(a, "--keys="):
			keys = splitCSV(strings.TrimPrefix(a, "--keys="))
		case strings.HasPrefix(a, "--prefixes="):
			prefixes = splitCSV(strings.TrimPrefix(a, "--prefixes="))
		case a == "--allow-override":
			allowOverride = true
		case strings.HasPrefix(a, "--format="):
			format = strings.TrimPrefix(a, "--format=")
		}
	}

	base, err := parser.Parse(basePath)
	if err != nil {
		return fmt.Errorf("reading base file: %w", err)
	}
	next, err := parser.Parse(nextPath)
	if err != nil {
		return fmt.Errorf("reading next file: %w", err)
	}

	violations, applyErr := protector.Apply(base, next, protector.Options{
		Keys:          keys,
		Prefixes:      prefixes,
		AllowOverride: allowOverride,
	})

	if format == "json" {
		type jsonViolation struct {
			Key    string `json:"key"`
			Reason string `json:"reason"`
		}
		out := make([]jsonViolation, len(violations))
		for i, v := range violations {
			out[i] = jsonViolation{Key: v.Key, Reason: v.Reason}
		}
		return json.NewEncoder(os.Stdout).Encode(out)
	}

	if len(violations) == 0 {
		fmt.Println("No protection violations found.")
	} else {
		for _, v := range violations {
			fmt.Printf("VIOLATION  %-30s  %s\n", v.Key, v.Reason)
		}
	}
	return applyErr
}
