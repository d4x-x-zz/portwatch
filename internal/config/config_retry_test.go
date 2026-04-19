package config

import (
	"testing"
	"time"
)

func TestDefaultRetry(t *testing.T) {
	r := DefaultRetry()
	if r.Attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", r.Attempts)
	}
	if r.Delay != 2*time.Second {
		t.Errorf("expected 2s delay, got %v", r.Delay)
	}
}

func TestValidateRetry_Valid(t *testing.T) {
	if err := validateRetry(DefaultRetry()); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateRetry_ZeroAttempts(t *testing.T) {
	err := validateRetry(RetryConfig{Attempts: 0, Delay: time.Second})
	if err == nil {
		t.Fatal("expected error for zero attempts")
	}
}

func TestValidateRetry_NegativeDelay(t *testing.T) {
	err := validateRetry(RetryConfig{Attempts: 1, Delay: -time.Second})
	if err == nil {
		t.Fatal("expected error for negative delay")
	}
}
