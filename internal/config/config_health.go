package config

import "fmt"

// Health holds configuration for the health-check HTTP server.
type Health struct {
	Enabled bool   `toml:"enabled"`
	Addr    string `toml:"addr"`
	Path    string `toml:"path"`
}

// DefaultHealth returns a Health config with sensible defaults.
func DefaultHealth() Health {
	return Health{
		Enabled: true,
		Addr:    ":9091",
		Path:    "/healthz",
	}
}

func validateHealth(h Health) error {
	if !h.Enabled {
		return nil
	}
	if h.Addr == "" {
		return fmt.Errorf("health.addr must not be empty when health is enabled")
	}
	if h.Path == "" {
		return fmt.Errorf("health.path must not be empty when health is enabled")
	}
	return nil
}
