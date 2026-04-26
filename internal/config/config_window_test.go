package config

import (
	"testing"
	"time"
)

func TestDefaultWindowLimit_Values(t *testing.T) {
	cfg := DefaultWindowLimit()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.Max <= 0 {
		t.Errorf("expected positive Max, got %d", cfg.Max)
	}
	if cfg.Window <= 0 {
		t.Errorf("expected positive Window, got %v", cfg.Window)
	}
}

func TestValidateWindowLimit_Disabled(t *testing.T) {
	cfg := WindowLimit{Enabled: false}
	if err := validateWindowLimit(cfg); err != nil {
		t.Errorf("unexpected error for disabled config: %v", err)
	}
}

func TestValidateWindowLimit_Valid(t *testing.T) {
	cfg := WindowLimit{
		Enabled: true,
		Max:     10,
		Window:  time.Minute,
	}
	if err := validateWindowLimit(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateWindowLimit_ZeroMax(t *testing.T) {
	cfg := WindowLimit{
		Enabled: true,
		Max:     0,
		Window:  time.Minute,
	}
	if err := validateWindowLimit(cfg); err == nil {
		t.Error("expected error for zero Max")
	}
}

func TestValidateWindowLimit_NegativeMax(t *testing.T) {
	cfg := WindowLimit{
		Enabled: true,
		Max:     -5,
		Window:  time.Minute,
	}
	if err := validateWindowLimit(cfg); err == nil {
		t.Error("expected error for negative Max")
	}
}

func TestValidateWindowLimit_ZeroWindow(t *testing.T) {
	cfg := WindowLimit{
		Enabled: true,
		Max:     10,
		Window:  0,
	}
	if err := validateWindowLimit(cfg); err == nil {
		t.Error("expected error for zero Window")
	}
}
