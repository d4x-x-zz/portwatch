package config

import (
	"fmt"
	"time"
)

// Debounce holds configuration for the debounce alerting middleware.
type Debounce struct {
	Enabled     bool          `toml:"enabled"`
	QuietPeriod time.Duration `toml:"quiet_period"`
	MaxDelay    time.Duration `toml:"max_delay"`
}

// DefaultDebounce returns sensible debounce defaults.
func DefaultDebounce() Debounce {
	return Debounce{
		Enabled:     false,
		QuietPeriod: 2 * time.Second,
		MaxDelay:    30 * time.Second,
	}
}

func validateDebounce(d Debounce) error {
	if !d.Enabled {
		return nil
	}
	if d.QuietPeriod <= 0 {
		return fmt.Errorf("debounce: quiet_period must be positive, got %s", d.QuietPeriod)
	}
	if d.MaxDelay <= 0 {
		return fmt.Errorf("debounce: max_delay must be positive, got %s", d.MaxDelay)
	}
	if d.QuietPeriod >= d.MaxDelay {
		return fmt.Errorf("debounce: quiet_period (%s) must be less than max_delay (%s)", d.QuietPeriod, d.MaxDelay)
	}
	return nil
}
