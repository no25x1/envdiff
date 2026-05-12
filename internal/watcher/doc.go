// Package watcher provides file-change detection for envdiff.
//
// It polls one or more files at a configurable interval and emits
// FileEvent values on the Events channel whenever a file's SHA-256
// digest changes.  The polling approach avoids OS-specific filesystem
// notification APIs and works reliably across Linux, macOS, and
// Windows.
//
// Typical usage:
//
//	w, err := watcher.New([]string{".env", ".env.production"}, 2*time.Second)
//	if err != nil { ... }
//	w.Start()
//	for ev := range w.Events {
//	    fmt.Printf("changed: %s\n", ev.Path)
//	}
package watcher
