package config

import (
	"errors"
	"time"
)

// RateLimitConfig controls the rate-limit middleware in the notify pipeline.
type RateLimitConfig struct {
	// Enabled toggles the rate-limiter. When false the middleware is a no-op.
	Enabled bool `toml:"enabled"`

	// Cooldown is the minimum duration between identical notifications.
	Cooldown time.Duration `toml:"cooldown"`
}

// DefaultRateLimit returns a RateLimitConfig with sensible defaults.
func DefaultRateLimit() RateLimitConfig {
	return RateLimitConfig{
		Enabled:  true,
		Cooldown: 5 * time.Minute,
	}
}

func validateRateLimit(r RateLimitConfig) error {
	if !r.Enabled {
		return nil
	}
	if r.Cooldown <= 0 {
		return errors.New("ratelimit: cooldown must be positive")
	}
	return nil
}
