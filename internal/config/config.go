// Package config handles loading and validation of portwatch configuration.
package config

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// Config holds the runtime configuration for portwatch.
type Config struct {
	// PortRange defines the inclusive range of ports to scan.
	PortRange PortRange `json:"port_range"`
	// Interval is how often to run a scan.
	Interval time.Duration `json:"interval"`
	// LogFile is the path to write alerts. Defaults to stdout if empty.
	LogFile string `json:"log_file"`
}

// PortRange defines a start and end port (inclusive).
type PortRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// Validate returns an error if the config is invalid.
func (c *Config) Validate() error {
	if c.PortRange.From < 1 || c.PortRange.From > 65535 {
		return errors.New("port_range.from must be between 1 and 65535")
	}
	if c.PortRange.To < 1 || c.PortRange.To > 65535 {
		return errors.New("port_range.to must be between 1 and 65535")
	}
	if c.PortRange.From > c.PortRange.To {
		return errors.New("port_range.from must be <= port_range.to")
	}
	if c.Interval <= 0 {
		return errors.New("interval must be positive")
	}
	return nil
}

// Load reads and parses a JSON config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Default returns a Config with sensible defaults.
func Default() *Config {
	return &Config{
		PortRange: PortRange{From: 1, To: 1024},
		Interval:  30 * time.Second,
	}
}
