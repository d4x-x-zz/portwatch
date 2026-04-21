package config

import (
	"testing"
	"time"
)

func TestDefaultDebounce_Values(t *testing.T) {
	d := DefaultDebounce()
	if d.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if d.QuietPeriod != 2*time.Second {
		t.Errorf("unexpected QuietPeriod: %s", d.QuietPeriod)
	}
	if d.MaxDelay != 30*time.Second {
		t.Errorf("unexpected MaxDelay: %s", d.MaxDelay)
	}
}

func TestValidateDebounce_Disabled(t *testing.T) {
	// disabled with bad values should still pass
	d := Debounce{Enabled: false, QuietPeriod: -1, MaxDelay: -1}
	if err := validateDebounce(d); err != nil {
		t.Fatalf("expected no error when disabled, got: %v", err)
	}
}

func TestValidateDebounce_Valid(t *testing.T) {
	d := Debounce{Enabled: true, QuietPeriod: 1 * time.Second, MaxDelay: 10 * time.Second}
	if err := validateDebounce(d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateDebounce_ZeroQuietPeriod(t *testing.T) {
	d := Debounce{Enabled: true, QuietPeriod: 0, MaxDelay: 10 * time.Second}
	if err := validateDebounce(d); err == nil {
		t.Fatal("expected error for zero quiet_period")
	}
}

func TestValidateDebounce_ZeroMaxDelay(t *testing.T) {
	d := Debounce{Enabled: true, QuietPeriod: 1 * time.Second, MaxDelay: 0}
	if err := validateDebounce(d); err == nil {
		t.Fatal("expected error for zero max_delay")
	}
}

func TestValidateDebounce_QuietPeriodExceedsMaxDelay(t *testing.T) {
	d := Debounce{Enabled: true, QuietPeriod: 10 * time.Second, MaxDelay: 5 * time.Second}
	if err := validateDebounce(d); err == nil {
		t.Fatal("expected error when quiet_period >= max_delay")
	}
}

func TestValidateDebounce_QuietPeriodEqualsMaxDelay(t *testing.T) {
	d := Debounce{Enabled: true, QuietPeriod: 5 * time.Second, MaxDelay: 5 * time.Second}
	if err := validateDebounce(d); err == nil {
		t.Fatal("expected error when quiet_period == max_delay")
	}
}
