package config

import "errors"

// SnapshotServer holds HTTP server settings for snapshot endpoints.
type SnapshotServer struct {
	Enabled bool   `toml:"enabled"`
	Addr    string `toml:"addr"`
}

// DefaultSnapshotServer returns sensible defaults for the snapshot server.
func DefaultSnapshotServer() SnapshotServer {
	return SnapshotServer{
		Enabled: false,
		Addr:    "127.0.0.1:9092",
	}
}

func validateSnapshotServer(s SnapshotServer) error {
	if !s.Enabled {
		return nil
	}
	if s.Addr == "" {
		return errors.New("snapshot_server.addr must not be empty when enabled")
	}
	return nil
}
