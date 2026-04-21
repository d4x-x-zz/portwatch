package config

import "testing"

func TestDefaultSuppress_Values(t *testing.T) {
	s := DefaultSuppress()
	if s.MaxRepeats != 1 {
		t.Errorf("expected MaxRepeats=1, got %d", s.MaxRepeats)
	}
	if !s.ResetOnChange {
		t.Error("expected ResetOnChange=true")
	}
}

func TestValidateSuppress_Valid(t *testing.T) {
	s := DefaultSuppress()
	if err := validateSuppress(s); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateSuppress_ZeroMaxRepeats(t *testing.T) {
	s := Suppress{MaxRepeats: 0, ResetOnChange: true}
	if err := validateSuppress(s); err != nil {
		t.Errorf("zero max_repeats should be valid, got: %v", err)
	}
}

func TestValidateSuppress_NegativeMaxRepeats(t *testing.T) {
	s := Suppress{MaxRepeats: -1}
	if err := validateSuppress(s); err == nil {
		t.Error("expected error for negative max_repeats")
	}
}

func TestValidateSuppress_ResetOnChangeFalse(t *testing.T) {
	s := Suppress{MaxRepeats: 3, ResetOnChange: false}
	if err := validateSuppress(s); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
