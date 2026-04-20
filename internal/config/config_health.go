package config

import "fmt"

// HealthConfig controls the built-in health-check HTTP server.
type HealthConfig struct {
	Enabled bool   `toml:"enabled"`
	Addr    string `toml:"addr"`
	Path    string `toml:"path"`
}

// DefaultHealth returns a HealthConfig with sensible defaults.
func DefaultHealth() HealthConfig {
	return HealthConfig{
		Enabled: true,
		Addr:    ":9110",
		Path:    "/health",
	}
}

func validateHealth(h HealthConfig) error {
	if !h.Enabled {
		return nil
	}
	if h.Addr == "" {
		return fmt.Errorf("health.addr must not be empty")
	}
	if len(h.Path) == 0 || h.Path[0] != '/' {
		return fmt.Errorf("health.path must start with '/'")
	}
	return nil
}
