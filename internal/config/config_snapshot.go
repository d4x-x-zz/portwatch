package config

import "errors"

// Snapshot controls periodic state snapshots written to disk so the
// daemon can resume cleanly after a restart.
type Snapshot struct {
	// Enabled turns snapshot persistence on or off.
	Enabled bool `toml:"enabled"`

	// Path is the file where the snapshot is written.
	Path string `toml:"path"`

	// KeepLast is the number of rolling backup snapshots to retain.
	// 0 means keep only the current snapshot.
	KeepLast int `toml:"keep_last"`
}

// DefaultSnapshot returns a Snapshot with sensible defaults.
func DefaultSnapshot() Snapshot {
	return Snapshot{
		Enabled:  true,
		Path:     "/var/lib/portwatch/snapshot.json",
		KeepLast: 3,
	}
}

// validateSnapshot checks that the Snapshot configuration is coherent.
func validateSnapshot(s Snapshot) error {
	if !s.Enabled {
		return nil
	}
	if s.Path == "" {
		return errors.New("snapshot.path must not be empty when snapshots are enabled")
	}
	if s.KeepLast < 0 {
		return errors.New("snapshot.keep_last must be >= 0")
	}
	return nil
}
