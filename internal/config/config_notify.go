package config

import "errors"

// Notify groups all notification pipeline configuration.
type Notify struct {
	Retry          Retry          `toml:"retry"`
	Backoff        Backoff        `toml:"backoff"`
	RateLimit      RateLimit      `toml:"rate_limit"`
	Dedupe         Dedupe         `toml:"dedupe"`
	Suppress       Suppress       `toml:"suppress"`
	Debounce       Debounce       `toml:"debounce"`
	CircuitBreaker CircuitBreaker `toml:"circuit_breaker"`
}

// DefaultNotify returns a Notify config populated with all sub-defaults.
func DefaultNotify() Notify {
	return Notify{
		Retry:          DefaultRetry(),
		Backoff:        DefaultBackoff(),
		RateLimit:      DefaultRateLimit(),
		Dedupe:         DefaultDedupe(),
		Suppress:       DefaultSuppress(),
		Debounce:       DefaultDebounce(),
		CircuitBreaker: DefaultCircuitBreaker(),
	}
}

func validateNotify(n Notify) error {
	if err := validateRetry(n.Retry); err != nil {
		return err
	}
	if err := validateBackoff(n.Backoff); err != nil {
		return err
	}
	if err := validateRateLimit(n.RateLimit); err != nil {
		return err
	}
	if err := validateDedupe(n.Dedupe); err != nil {
		return err
	}
	if err := validateSuppress(n.Suppress); err != nil {
		return err
	}
	if err := validateDebounce(n.Debounce); err != nil {
		return err
	}
	if err := validateCircuitBreaker(n.CircuitBreaker); err != nil {
		return errors.New("notify.circuit_breaker: " + err.Error())
	}
	return nil
}
