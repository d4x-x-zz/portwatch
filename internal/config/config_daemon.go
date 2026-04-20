package config

import (
	"fmt"
	"time"
)

// Daemon holds configuration for the main daemon loop.
type Daemon struct {
	Interval    time.Duration `toml:"interval"`
	StateFile   string        `toml:"state_file"`
	PIDFile     string        `toml:"pid_file"`
}

// DefaultDaemon returns a Daemon config with sensible defaults.
func DefaultDaemon() Daemon {
	return Daemon{
		Interval:  30 * time.Second,
		StateFile: "/var/lib/portwatch/state.json",
		PIDFile:   "/var/run/portwatch.pid",
	}
}

func validateDaemon(d Daemon) error {
	if d.Interval <= 0 {
		return fmt.Errorf("daemon.interval must be positive")
	}
	if d.StateFile == "" {
		return fmt.Errorf("daemon.state_file must not be empty")
	}
	if d.PIDFile == "" {
		return fmt.Errorf("daemon.pid_file must not be empty")
	}
	return nil
}
