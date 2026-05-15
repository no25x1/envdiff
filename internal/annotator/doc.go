// Package annotator attaches inline comments to .env entries based on
// configurable rules that match key prefixes or suffixes.
//
// Built-in rules flag common sensitive suffixes (_SECRET, _PASSWORD,
// _TOKEN, _KEY) and debug prefixes (DEBUG_). Custom rules can be
// supplied via Options.Rules.
//
// Example usage:
//
//	file, _ := parser.ParseFile(".env")
//	annotated := annotator.Apply(file, annotator.DefaultOptions())
package annotator
