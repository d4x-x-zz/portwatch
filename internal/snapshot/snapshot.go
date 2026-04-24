// Package snapshot provides functionality for saving and loading
// point-in-time captures of the observed open port set.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Entry represents a single snapshot on disk.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Ports     []int     `json:"ports"`
}

// Save writes the given port set as a JSON snapshot file under dir.
// Files are named by RFC3339 timestamp so they sort chronologically.
func Save(dir string, ports []int, now time.Time) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("snapshot: mkdir %s: %w", dir, err)
	}

	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)

	e := Entry{Timestamp: now.UTC(), Ports: sorted}
	data, err := json.Marshal(e)
	if err != nil {
		return "", fmt.Errorf("snapshot: marshal: %w", err)
	}

	name := now.UTC().Format("20060102T150405Z") + ".json"
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", fmt.Errorf("snapshot: write %s: %w", path, err)
	}
	return path, nil
}

// List returns all snapshot entries in dir sorted oldest-first.
func List(dir string) ([]Entry, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("snapshot: glob: %w", err)
	}
	sort.Strings(matches)

	out := make([]Entry, 0, len(matches))
	for _, p := range matches {
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, fmt.Errorf("snapshot: read %s: %w", p, err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("snapshot: decode %s: %w", p, err)
		}
		out = append(out, e)
	}
	return out, nil
}

// Prune removes the oldest snapshots so that at most keepLast remain.
func Prune(dir string, keepLast int) error {
	matches, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return fmt.Errorf("snapshot: glob: %w", err)
	}
	sort.Strings(matches)

	for len(matches) > keepLast {
		if err := os.Remove(matches[0]); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("snapshot: remove %s: %w", matches[0], err)
		}
		matches = matches[1:]
	}
	return nil
}
