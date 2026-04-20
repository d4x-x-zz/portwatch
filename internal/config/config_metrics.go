package config

import "fmt"

// MetricsConfig controls the built-in HTTP metrics endpoint.
type MetricsConfig struct {
	Enabled bool   `toml:"enabled"`
	Addr    string `toml:"addr"`
	Path    string `toml:"path"`
	MaxHistory int `toml:"max_history"`
}

// DefaultMetrics returns sensible defaults for the metrics endpoint.
func DefaultMetrics() MetricsConfig {
	return MetricsConfig{
		Enabled:    false,
		Addr:       ":9090",
		Path:       "/metrics",
		MaxHistory: 100,
	}
}

func validateMetrics(m MetricsConfig) error {
	if !m.Enabled {
		return nil
	}
	if m.Addr == "" {
		return fmt.Errorf("metrics.addr must not be empty")
	}
	if m.Path == "" {
		return fmt.Errorf("metrics.path must not be empty")
	}
	if m.MaxHistory < 1 {
		return fmt.Errorf("metrics.max_history must be at least 1")
	}
	return nil
}
