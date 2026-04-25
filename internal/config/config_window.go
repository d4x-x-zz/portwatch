package config

import (
	"fmt"
	"time"
)

// WindowLimit configures the sliding-window notification limiter.
type WindowLimit struct {
	Enabled  bool          `toml:"enabled"`
	Window   time.Duration `toml:"window"`
	MaxCalls int           `toml:"max_calls"`
}

// DefaultWindowLimit returns sensible defaults: 10 calls per minute.
func DefaultWindowLimit() WindowLimit {
	return WindowLimit{
		Enabled:  false,
		Window:   time.Minute,
		MaxCalls: 10,
	}
}

func validateWindowLimit(w WindowLimit) error {
	if !w.Enabled {
		return nil
	}
	if w.Window <= 0 {
		return fmt.Errorf("window_limit.window must be positive, got %s", w.Window)
	}
	if w.MaxCalls <= 0 {
		return fmt.Errorf("window_limit.max_calls must be positive, got %d", w.MaxCalls)
	}
	return nil
}
