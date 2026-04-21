package config

import (
	"testing"
	"time"
)

func TestDefaultCircuitBreaker_Values(t *testing.T) {
	c := DefaultCircuitBreaker()
	if !c.Enabled {
		t.Error("expected Enabled to be true")
	}
	if c.Threshold != 5 {
		t.Errorf("expected Threshold 5, got %d", c.Threshold)
	}
	if c.Cooldown != 30*time.Second {
		t.Errorf("expected Cooldown 30s, got %v", c.Cooldown)
	}
}

func TestValidateCircuitBreaker_Disabled(t *testing.T) {
	c := CircuitBreaker{Enabled: false}
	if err := validateCircuitBreaker(c); err != nil {
		t.Errorf("expected no error for disabled circuit breaker, got %v", err)
	}
}

func TestValidateCircuitBreaker_Valid(t *testing.T) {
	c := DefaultCircuitBreaker()
	if err := validateCircuitBreaker(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateCircuitBreaker_ZeroThreshold(t *testing.T) {
	c := DefaultCircuitBreaker()
	c.Threshold = 0
	if err := validateCircuitBreaker(c); err == nil {
		t.Error("expected error for zero threshold")
	}
}

func TestValidateCircuitBreaker_ZeroCooldown(t *testing.T) {
	c := DefaultCircuitBreaker()
	c.Cooldown = 0
	if err := validateCircuitBreaker(c); err == nil {
		t.Error("expected error for zero cooldown")
	}
}

func TestValidateCircuitBreaker_NegativeCooldown(t *testing.T) {
	c := DefaultCircuitBreaker()
	c.Cooldown = -1 * time.Second
	if err := validateCircuitBreaker(c); err == nil {
		t.Error("expected error for negative cooldown")
	}
}
