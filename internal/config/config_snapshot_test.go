package config

import "testing"

func TestDefaultSnapshot_Values(t *testing.T) {
	s := DefaultSnapshot()
	if !s.Enabled {
		t.Error("expected Enabled to be true")
	}
	if s.Path == "" {
		t.Error("expected a non-empty default Path")
	}
	if s.KeepLast <= 0 {
		t.Errorf("expected KeepLast > 0, got %d", s.KeepLast)
	}
}

func TestValidateSnapshot_Disabled(t *testing.T) {
	s := Snapshot{Enabled: false}
	if err := validateSnapshot(s); err != nil {
		t.Errorf("unexpected error for disabled snapshot: %v", err)
	}
}

func TestValidateSnapshot_Valid(t *testing.T) {
	s := DefaultSnapshot()
	if err := validateSnapshot(s); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateSnapshot_EmptyPath(t *testing.T) {
	s := DefaultSnapshot()
	s.Path = ""
	if err := validateSnapshot(s); err == nil {
		t.Error("expected error for empty path")
	}
}

func TestValidateSnapshot_NegativeKeepLast(t *testing.T) {
	s := DefaultSnapshot()
	s.KeepLast = -1
	if err := validateSnapshot(s); err == nil {
		t.Error("expected error for negative keep_last")
	}
}

func TestValidateSnapshot_ZeroKeepLast(t *testing.T) {
	s := DefaultSnapshot()
	s.KeepLast = 0
	if err := validateSnapshot(s); err != nil {
		t.Errorf("keep_last=0 should be valid, got: %v", err)
	}
}
