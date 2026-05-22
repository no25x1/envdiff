// Package merger combines multiple .env files into a single EnvFile.
//
// Merge accepts a slice of EnvFile pointers and a Strategy that controls
// how duplicate keys are resolved:
//
//   - StrategyLast  – the last file's value wins (default)
//   - StrategyFirst – the first file's value is kept
//   - StrategyError – any duplicate key causes an error
//
// The order of entries in the result follows the order in which keys are
// first encountered across all source files.
package merger
