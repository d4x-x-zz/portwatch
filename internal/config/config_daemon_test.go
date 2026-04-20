package config

import (
	"testing"
	"time"
)

func TestDefaultDaemon(t *testing.T) {
	d := DefaultDaemon()
	if d.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", d.Interval)
	}
	if d.StateFile == "" {
		t.Error("expected non-empty state_file")
	}
	if d.PIDFile == "" {
		t.Error("expected non-empty pid_file")
	}
}

func TestValidateDaemon_Valid(t *testing.T) {
	if err := validateDaemon(DefaultDaemon()); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateDaemon_ZeroInterval(t *testing.T) {
	d := DefaultDaemon()
	d.Interval = 0
	if err := validateDaemon(d); err == nil {
		t.Error("expected error for zero interval")
	}
}

func TestValidateDaemon_EmptyStateFile(t *testing.T) {
	d := DefaultDaemon()
	d.StateFile = ""
	if err := validateDaemon(d); err == nil {
		t.Error("expected error for empty state_file")
	}
}

func TestValidateDaemon_EmptyPIDFile(t *testing.T) {
	d := DefaultDaemon()
	d.PIDFile = ""
	if err := validateDaemon(d); err == nil {
		t.Error("expected error for empty pid_file")
	}
}
