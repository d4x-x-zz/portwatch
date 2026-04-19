// Package history provides persistent storage of port scan snapshots.
package history

import (
	"encoding/json"
	"os"
	"time"
)

// Snapshot represents a recorded port scan result at a point in time.
type Snapshot struct {
	Timestamp time.Time `json:"timestamp"`
	Ports     []uint16  `json:"ports"`
}

// Store manages reading and writing snapshots to disk.
type Store struct {
	path string
}

// NewStore creates a Store backed by the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Save writes a snapshot to disk, overwriting any previous snapshot.
func (s *Store) Save(ports []uint16) error {
	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		Ports:     ports,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Load reads the most recent snapshot from disk.
// Returns nil, nil if no snapshot exists yet.
func (s *Store) Load() (*Snapshot, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, err
	}
	return &snap, nil
}
