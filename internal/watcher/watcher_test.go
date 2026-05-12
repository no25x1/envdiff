package watcher_test

import (
	"os"
	"testing"
	"time"

	"envdiff/internal/watcher"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "watcher-*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestNew_HashesFiles(t *testing.T) {
	p := writeTempFile(t, "KEY=value\n")
	w, err := watcher.New([]string{p}, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w.Stop()
}

func TestNew_MissingFile(t *testing.T) {
	_, err := watcher.New([]string{"/nonexistent/.env"}, 50*time.Millisecond)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestWatcher_DetectsChange(t *testing.T) {
	p := writeTempFile(t, "KEY=original\n")
	w, err := watcher.New([]string{p}, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w.Start()
	defer w.Stop()

	time.Sleep(30 * time.Millisecond)
	if err := os.WriteFile(p, []byte("KEY=changed\n"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case ev := <-w.Events:
		if ev.Path != p {
			t.Errorf("expected path %s, got %s", p, ev.Path)
		}
		if ev.OldHash == ev.NewHash {
			t.Error("expected hashes to differ")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatcher_NoEventWhenUnchanged(t *testing.T) {
	p := writeTempFile(t, "KEY=stable\n")
	w, err := watcher.New([]string{p}, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w.Start()
	defer w.Stop()

	select {
	case ev := <-w.Events:
		t.Errorf("unexpected event: %+v", ev)
	case <-time.After(120 * time.Millisecond):
		// expected: no event
	}
}
