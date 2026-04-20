package config

import "testing"

func TestDefaultScan(t *testing.T) {
	s := DefaultScan()
	if s.Host == "" {
		t.Error("expected non-empty host")
	}
	if s.PortRange == "" {
		t.Error("expected non-empty port_range")
	}
	if s.Concurrency <= 0 {
		t.Errorf("expected positive concurrency, got %d", s.Concurrency)
	}
	if s.TimeoutMs <= 0 {
		t.Errorf("expected positive timeout_ms, got %d", s.TimeoutMs)
	}
}

func TestValidateScan_Valid(t *testing.T) {
	if err := validateScan(DefaultScan()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateScan_EmptyHost(t *testing.T) {
	s := DefaultScan()
	s.Host = ""
	if err := validateScan(s); err == nil {
		t.Error("expected error for empty host")
	}
}

func TestValidateScan_EmptyRange(t *testing.T) {
	s := DefaultScan()
	s.PortRange = ""
	if err := validateScan(s); err == nil {
		t.Error("expected error for empty port_range")
	}
}

func TestValidateScan_ZeroConcurrency(t *testing.T) {
	s := DefaultScan()
	s.Concurrency = 0
	if err := validateScan(s); err == nil {
		t.Error("expected error for zero concurrency")
	}
}

func TestValidateScan_ZeroTimeout(t *testing.T) {
	s := DefaultScan()
	s.TimeoutMs = 0
	if err := validateScan(s); err == nil {
		t.Error("expected error for zero timeout_ms")
	}
}
