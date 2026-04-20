package config

import (
	"testing"
	"time"
)

func TestDefaultBackoff_Values(t *testing.T) {
	b := DefaultBackoff()
	if b.Strategy != "exponential" {
		t.Fatalf("expected exponential, got %q", b.Strategy)
	}
	if b.Base <= 0 {
		t.Fatal("base must be positive")
	}
	if b.Max < b.Base {
		t.Fatal("max must be >= base")
	}
	if b.Factor <= 1.0 {
		t.Fatal("factor must be > 1")
	}
}

func TestValidateBackoff_Valid(t *testing.T) {
	if err := validateBackoff(DefaultBackoff()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateBackoff_Constant(t *testing.T) {
	b := BackoffConfig{Strategy: "constant", Base: time.Second}
	if err := validateBackoff(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateBackoff_BadStrategy(t *testing.T) {
	b := DefaultBackoff()
	b.Strategy = "linear"
	if err := validateBackoff(b); err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestValidateBackoff_ZeroBase(t *testing.T) {
	b := DefaultBackoff()
	b.Base = 0
	if err := validateBackoff(b); err == nil {
		t.Fatal("expected error for zero base")
	}
}

func TestValidateBackoff_MaxLessThanBase(t *testing.T) {
	b := DefaultBackoff()
	b.Max = b.Base / 2
	if err := validateBackoff(b); err == nil {
		t.Fatal("expected error when max < base")
	}
}

func TestValidateBackoff_FactorTooLow(t *testing.T) {
	b := DefaultBackoff()
	b.Factor = 1.0
	if err := validateBackoff(b); err == nil {
		t.Fatal("expected error when factor <= 1")
	}
}
