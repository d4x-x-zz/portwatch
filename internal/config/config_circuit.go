package config

import (
	"errors"
	"time"
)

// CircuitBreaker holds configuration for the circuit-breaker alerting middleware.
type CircuitBreaker struct {
	// Enabled controls whether the circuit breaker is active.
	Enabled bool `toml:"enabled"`

	// Threshold is the number of consecutive failures before the circuit opens.
	Threshold int `toml:"threshold"`

	// Cooldown is the duration the circuit stays open before moving to half-open.
	Cooldown time.Duration `toml:"cooldown"`
}

// DefaultCircuitBreaker returns sensible production defaults.
func DefaultCircuitBreaker() CircuitBreaker {
	return CircuitBreaker{
		Enabled:   true,
		Threshold: 5,
		Cooldown:  30 * time.Second,
	}
}

func validateCircuitBreaker(c CircuitBreaker) error {
	if !c.Enabled {
		return nil
	}
	if c.Threshold <= 0 {
		return errors.New("circuit_breaker.threshold must be greater than zero")
	}
	if c.Cooldown <= 0 {
		return errors.New("circuit_breaker.cooldown must be greater than zero")
	}
	return nil
}
