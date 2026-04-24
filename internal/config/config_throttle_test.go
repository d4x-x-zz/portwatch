package config

import (
	"testing"
	"time"
)

func TestDefaultThrottle_Values(t *testing.T) {
	cfg := DefaultThrottle()
	if cfg.Window != 60*time.Second {
		t.Errorf("expected 60s window, got %v", cfg.Window)
	}
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
}

func TestValidateThrottle_Disabled(t *testing.T) {
	cfg := DefaultThrottle()
	cfg.Enabled = false
	if err := validateThrottle(cfg); err != nil {
		t.Errorf("unexpected error for disabled throttle: %v", err)
	}
}

func TestValidateThrottle_Valid(t *testing.T) {
	cfg := DefaultThrottle()
	cfg.Enabled = true
	if err := validateThrottle(cfg); err != nil {
		t.Errorf("unexpected error for valid throttle: %v", err)
	}
}

func TestValidateThrottle_ZeroWindow(t *testing.T) {
	cfg := DefaultThrottle()
	cfg.Enabled = true
	cfg.Window = 0
	if err := validateThrottle(cfg); err == nil {
		t.Error("expected error for zero window")
	}
}

func TestValidateThrottle_NegativeWindow(t *testing.T) {
	cfg := DefaultThrottle()
	cfg.Enabled = true
	cfg.Window = -1 * time.Second
	if err := validateThrottle(cfg); err == nil {
		t.Error("expected error for negative window")
	}
}
