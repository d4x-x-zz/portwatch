package config

import (
	"fmt"
	"time"
)

// Throttle holds configuration for the ThrottledAlerter middleware.
type Throttle struct {
	// Enabled controls whether throttling is applied.
	Enabled bool `toml:"enabled"`

	// Window is the minimum duration between successive alerts.
	Window time.Duration `toml:"window"`

	// FlushOnExit sends any pending alert when the daemon shuts down.
	FlushOnExit bool `toml:"flush_on_exit"`
}

// DefaultThrottle returns a Throttle with sensible defaults.
func DefaultThrottle() Throttle {
	return Throttle{
		Enabled:     false,
		Window:      30 * time.Second,
		FlushOnExit: true,
	}
}

func validateThrottle(t Throttle) error {
	if !t.Enabled {
		return nil
	}
	if t.Window <= 0 {
		return fmt.Errorf("throttle: window must be positive, got %s", t.Window)
	}
	return nil
}
