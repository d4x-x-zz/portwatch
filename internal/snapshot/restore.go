package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a saved snapshot entry on disk.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Ports     []int     `json:"ports"`
}

// Load reads the snapshot at the given path and returns its entry.
func Load(path string) (*Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open %s: %w", path, err)
	}
	defer f.Close()

	var e Entry
	if err := json.NewDecoder(f).Decode(&e); err != nil {
		return nil, fmt.Errorf("snapshot: decode %s: %w", path, err)
	}
	return &e, nil
}

// Latest returns the most-recent snapshot entry from dir, or nil if none exist.
func Latest(dir string) (*Entry, error) {
	entries, err := List(dir)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, nil
	}
	// List returns entries sorted oldest-first; last element is newest.
	newest := entries[len(entries)-1]
	return Load(filepath.Join(dir, newest.Name()))
}
