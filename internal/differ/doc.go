// Package differ implements line-level diffing between two .env files,
// producing structured results with change types (added, removed, modified,
// unchanged) and a unified-style text formatter.
//
// Usage:
//
//	result := differ.Compare(baseFile, targetFile)
//	if result.HasChanges() {
//		fmt.Println(result.Summary())
//		fmt.Print(differ.Format(result))
//	}
package differ
