package config

import (
	"testing"
	"time"
)

func TestDefaultNotify_Values(t *testing.T) {
	n := DefaultNotify()
	if n.Retry.Attempts == 0 {
		t.Error("expected non-zero retry attempts")
	}
	if n.Throttle.Window == 0 {
		t.Error("expected non-zero throttle window")
	}
	if n.RateLimit.Cooldown == 0 {
		t.Error("expected non-zero rate limit cooldown")
	}
}

func TestValidateNotify_Valid(t *testing.T) {
	n := DefaultNotify()
	if err := validateNotify(n); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateNotify_InvalidRetry(t *testing.T) {
	n := DefaultNotify()
	n.Retry.Attempts = 0
	if err := validateNotify(n); err == nil {
		t.Error("expected error for zero retry attempts")
	}
}

func TestValidateNotify_InvalidThrottle(t *testing.T) {
	n := DefaultNotify()
	n.Throttle.Enabled = true
	n.Throttle.Window = -1 * time.Second
	if err := validateNotify(n); err == nil {
		t.Error("expected error for negative throttle window")
	}
}

func TestValidateNotify_InvalidCircuitBreaker(t *testing.T) {
	n := DefaultNotify()
	n.CircuitBreaker.Enabled = true
	n.CircuitBreaker.Threshold = 0
	if err := validateNotify(n); err == nil {
		t.Error("expected error for zero circuit breaker threshold")
	}
}
