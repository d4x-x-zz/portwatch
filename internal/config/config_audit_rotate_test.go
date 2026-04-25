package config

import (
	"testing"
	"time"
)

func TestDefaultAuditRotate_Values(t *testing.T) {
	r := DefaultAuditRotate()
	if r.MaxSizeBytes <= 0 {
		t.Errorf("expected positive MaxSizeBytes, got %d", r.MaxSizeBytes)
	}
	if r.MaxAge <= 0 {
		t.Errorf("expected positive MaxAge, got %v", r.MaxAge)
	}
	if r.KeepLast <= 0 {
		t.Errorf("expected positive KeepLast, got %d", r.KeepLast)
	}
}

func TestValidateAuditRotate_Disabled(t *testing.T) {
	r := AuditRotateConfig{Enabled: false}
	if err := validateAuditRotate(r); err != nil {
		t.Errorf("unexpected error for disabled rotate: %v", err)
	}
}

func TestValidateAuditRotate_Valid(t *testing.T) {
	r := AuditRotateConfig{
		Enabled:      true,
		MaxSizeBytes: 1024 * 1024,
		MaxAge:       24 * time.Hour,
		KeepLast:     5,
	}
	if err := validateAuditRotate(r); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateAuditRotate_ZeroMaxSize(t *testing.T) {
	r := AuditRotateConfig{
		Enabled:      true,
		MaxSizeBytes: 0,
		MaxAge:       time.Hour,
		KeepLast:     3,
	}
	if err := validateAuditRotate(r); err == nil {
		t.Error("expected error for zero MaxSizeBytes")
	}
}

func TestValidateAuditRotate_ZeroMaxAge(t *testing.T) {
	r := AuditRotateConfig{
		Enabled:      true,
		MaxSizeBytes: 512,
		MaxAge:       0,
		KeepLast:     3,
	}
	if err := validateAuditRotate(r); err == nil {
		t.Error("expected error for zero MaxAge")
	}
}

func TestValidateAuditRotate_ZeroKeepLast(t *testing.T) {
	r := AuditRotateConfig{
		Enabled:      true,
		MaxSizeBytes: 512,
		MaxAge:       time.Hour,
		KeepLast:     0,
	}
	if err := validateAuditRotate(r); err == nil {
		t.Error("expected error for zero KeepLast")
	}
}
