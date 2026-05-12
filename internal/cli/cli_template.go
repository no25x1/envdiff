package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/templater"
)

// runTemplate handles the `template` sub-command.
// Usage: envdiff template --template <tmpl.env> --values <values.env> [--validate]
func runTemplate(args []string) error {
	fs := flag.NewFlagSet("template", flag.ContinueOnError)
	tmplPath := fs.String("template", "", "path to the template .env file (required)")
	valsPath := fs.String("values", "", "path to the values .env file (required)")
	validateOnly := fs.Bool("validate", false, "only validate placeholders, do not render output")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *tmplPath == "" || *valsPath == "" {
		return fmt.Errorf("template: --template and --values are required")
	}

	tmplFile, err := parser.Parse(*tmplPath)
	if err != nil {
		return fmt.Errorf("template: cannot parse template file: %w", err)
	}

	valsFile, err := parser.Parse(*valsPath)
	if err != nil {
		return fmt.Errorf("template: cannot parse values file: %w", err)
	}

	if *validateOnly {
		return runTemplateValidate(tmplFile, valsFile)
	}

	res := templater.Render(tmplFile, valsFile)

	if len(res.Missing) > 0 {
		_, _ = fmt.Fprintf(os.Stderr, "warning: unresolved placeholders: %s\n",
			strings.Join(res.Missing, ", "))
	}

	for _, e := range res.Entries {
		fmt.Printf("%s=%s\n", e.Key, e.Value)
	}
	return nil
}

func runTemplateValidate(tmplFile, valsFile parser.EnvFile) error {
	errs := templater.Validate(tmplFile, valsFile)
	if len(errs) == 0 {
		fmt.Println("template: all placeholders resolved successfully")
		return nil
	}
	for _, e := range errs {
		_, _ = fmt.Fprintln(os.Stderr, "template error:", e)
	}
	return fmt.Errorf("template: validation failed with %d error(s)", len(errs))
}
