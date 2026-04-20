package config

import "testing"

func TestDefaultDedupe_Values(t *testing.T) {
	d := DefaultDedupe()
	if !d.Enabled {
		t.Error("expected Enabled to be true by default")
	}
	if !d.ResetOnDiff {
		t.Error("expected ResetOnDiff to be true by default")
	}
}

func TestValidateDedupe_Enabled(t *testing.T) {
	d := DefaultDedupe()
	if err := validateDedupe(d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateDedupe_Disabled(t *testing.T) {
	d := DedupeConfig{Enabled: false}
	if err := validateDedupe(d); err != nil {
		t.Fatalf("unexpected error when disabled: %v", err)
	}
}

func TestValidateDedupe_ResetOnDiff_False(t *testing.T) {
	d := DedupeConfig{Enabled: true, ResetOnDiff: false}
	if err := validateDedupe(d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
