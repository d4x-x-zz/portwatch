package config

import "fmt"

// ServerConfig holds settings for the optional HTTP admin server.
type ServerConfig struct {
	Enabled bool   `toml:"enabled"`
	Addr    string `toml:"addr"`
	Path    string `toml:"path"`
}

// DefaultServer returns a ServerConfig with sensible defaults.
func DefaultServer() ServerConfig {
	return ServerConfig{
		Enabled: false,
		Addr:    ":9091",
		Path:    "/status",
	}
}

func validateServer(s ServerConfig) error {
	if !s.Enabled {
		return nil
	}
	if s.Addr == "" {
		return fmt.Errorf("server.addr must not be empty")
	}
	if s.Path == "" {
		return fmt.Errorf("server.path must not be empty")
	}
	if s.Path[0] != '/' {
		return fmt.Errorf("server.path must start with /")
	}
	return nil
}
