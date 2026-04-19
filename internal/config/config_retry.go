package config

import (
	"fmt"
	"time"
)

// RetryConfig controls retry behaviour for alert delivery.
type RetryConfig struct {
	Attempts int           `toml:"attempts"`
	Delay    time.Duration `toml:"delay"`
}

// DefaultRetry returns sensible retry defaults.
func DefaultRetry() RetryConfig {
	return RetryConfig{
		Attempts: 3,
		Delay:    2 * time.Second,
	}
}

func validateRetry(r RetryConfig) error {
	if r.Attempts < 1 {
		return fmt.Errorf("retry.attempts must be at least 1")
	}
	if r.Delay < 0 {
		return fmt.Errorf("retry.delay must not be negative")
	}
	return nil
}
