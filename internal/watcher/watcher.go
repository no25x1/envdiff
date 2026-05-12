package watcher

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

// FileEvent represents a change detected in a watched file.
type FileEvent struct {
	Path    string
	OldHash string
	NewHash string
	At      time.Time
}

// Watcher polls files for changes and emits events.
type Watcher struct {
	paths    []string
	hashes   map[string]string
	interval time.Duration
	Events   chan FileEvent
	stop     chan struct{}
}

// New creates a Watcher that polls the given paths every interval.
func New(paths []string, interval time.Duration) (*Watcher, error) {
	w := &Watcher{
		paths:    paths,
		hashes:   make(map[string]string),
		interval: interval,
		Events:   make(chan FileEvent, 8),
		stop:     make(chan struct{}),
	}
	for _, p := range paths {
		h, err := hashFile(p)
		if err != nil {
			return nil, fmt.Errorf("watcher: initial hash %s: %w", p, err)
		}
		w.hashes[p] = h
	}
	return w, nil
}

// Start begins polling in a background goroutine.
func (w *Watcher) Start() {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.poll()
			case <-w.stop:
				return
			}
		}
	}()
}

// Stop halts the watcher and closes the Events channel.
func (w *Watcher) Stop() {
	close(w.stop)
	close(w.Events)
}

func (w *Watcher) poll() {
	for _, p := range w.paths {
		h, err := hashFile(p)
		if err != nil {
			continue
		}
		if old, ok := w.hashes[p]; ok && old != h {
			w.Events <- FileEvent{Path: p, OldHash: old, NewHash: h, At: time.Now()}
		}
		w.hashes[p] = h
	}
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
