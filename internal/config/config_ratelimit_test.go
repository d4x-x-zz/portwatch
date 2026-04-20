package config

import (
	"testing"
	"time"
)

func TestDefaultRateLimit_Values(t *testing.T) {
	r := DefaultRateLimit()
	if !r.Enabled {
		t.Error("expected Enabled to be true")
	}
	if r.Cooldown != 5*time.Minute {
		t.Errorf("expected 5m cooldown, got %s", r.Cooldown)
	}
}

func TestValidateRateLimit_Valid(t *testing.T) {
	r := DefaultRateLimit()
	if err := validateRateLimit(r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateRateLimit_Disabled(t *testing.T) {
	r := RateLimitConfig{Enabled: false, Cooldown: 0}
	if err := validateRateLimit(r); err != nil {
		t.Fatalf("disabled ratelimit should not error, got: %v", err)
	}
}

func TestValidateRateLimit_ZeroCooldown(t *testing.T) {
	r := RateLimitConfig{Enabled: true, Cooldown: 0}
	if err := validateRateLimit(r); err == nil {
		t.Fatal("expected error for zero cooldown")
	}
}

func TestValidateRateLimit_NegativeCooldown(t *testing.T) {
	r := RateLimitConfig{Enabled: true, Cooldown: -time.Second}
	if err := validateRateLimit(r); err == nil {
		t.Fatal("expected error for negative cooldown")
	}
}
