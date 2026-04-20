package config

import "fmt"

// NotifyConfig holds configuration for the notification pipeline.
type NotifyConfig struct {
	// Throttle suppresses repeated alerts within a time window.
	Throttle ThrottleConfig `toml:"throttle"`

	// RateLimit prevents identical diffs from firing repeatedly.
	RateLimit RateLimitConfig `toml:"rate_limit"`

	// Circuit breaker stops alerting after repeated failures.
	Circuit CircuitConfig `toml:"circuit"`
}

// ThrottleConfig controls how often alerts may fire.
type ThrottleConfig struct {
	Enabled  bool   `toml:"enabled"`
	WindowSec int   `toml:"window_sec"`
}

// RateLimitConfig suppresses duplicate diffs within a cooldown period.
type RateLimitConfig struct {
	Enabled     bool  `toml:"enabled"`
	CooldownSec int   `toml:"cooldown_sec"`
}

// CircuitConfig configures the circuit breaker for alerters.
type CircuitConfig struct {
	Enabled        bool  `toml:"enabled"`
	Threshold      int   `toml:"threshold"`
	CooldownSec    int   `toml:"cooldown_sec"`
}

// DefaultNotify returns a NotifyConfig with sensible defaults.
func DefaultNotify() NotifyConfig {
	return NotifyConfig{
		Throttle: ThrottleConfig{
			Enabled:   false,
			WindowSec: 60,
		},
		RateLimit: RateLimitConfig{
			Enabled:     true,
			CooldownSec: 300,
		},
		Circuit: CircuitConfig{
			Enabled:     true,
			Threshold:   5,
			CooldownSec: 120,
		},
	}
}

// validateNotify checks that notify config values are in range.
func validateNotify(n NotifyConfig) error {
	if n.Throttle.Enabled && n.Throttle.WindowSec <= 0 {
		return fmt.Errorf("notify.throttle.window_sec must be positive")
	}
	if n.RateLimit.Enabled && n.RateLimit.CooldownSec <= 0 {
		return fmt.Errorf("notify.rate_limit.cooldown_sec must be positive")
	}
	if n.Circuit.Enabled {
		if n.Circuit.Threshold <= 0 {
			return fmt.Errorf("notify.circuit.threshold must be positive")
		}
		if n.Circuit.CooldownSec <= 0 {
			return fmt.Errorf("notify.circuit.cooldown_sec must be positive")
		}
	}
	return nil
}
