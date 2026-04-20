package config

import (
	"errors"
	"time"
)

// BackoffConfig controls the back-off strategy used by the retry alerter.
type BackoffConfig struct {
	// Strategy is either "exponential" or "constant".
	Strategy string        `toml:"strategy"`
	Base     time.Duration `toml:"base"`
	Max      time.Duration `toml:"max"`
	Factor   float64       `toml:"factor"`
}

// DefaultBackoff returns conservative defaults for exponential back-off.
func DefaultBackoff() BackoffConfig {
	return BackoffConfig{
		Strategy: "exponential",
		Base:     500 * time.Millisecond,
		Max:      30 * time.Second,
		Factor:   2.0,
	}
}

func validateBackoff(b BackoffConfig) error {
	if b.Strategy != "exponential" && b.Strategy != "constant" {
		return errors.New("backoff.strategy must be \"exponential\" or \"constant\"")
	}
	if b.Base <= 0 {
		return errors.New("backoff.base must be positive")
	}
	if b.Strategy == "exponential" {
		if b.Max < b.Base {
			return errors.New("backoff.max must be >= backoff.base")
		}
		if b.Factor <= 1.0 {
			return errors.New("backoff.factor must be greater than 1")
		}
	}
	return nil
}
