package config

import "testing"

func TestDefaultMetrics(t *testing.T) {
	m := DefaultMetrics()
	if m.Enabled {
		t.Error("expected disabled by default")
	}
	if m.Addr != ":9090" {
		t.Errorf("unexpected addr: %s", m.Addr)
	}
	if m.Path != "/metrics" {
		t.Errorf("unexpected path: %s", m.Path)
	}
	if m.MaxHistory != 100 {
		t.Errorf("unexpected max_history: %d", m.MaxHistory)
	}
}

func TestValidateMetrics_Disabled(t *testing.T) {
	m := MetricsConfig{Enabled: false}
	if err := validateMetrics(m); err != nil {
		t.Errorf("expected no error for disabled metrics, got %v", err)
	}
}

func TestValidateMetrics_EmptyAddr(t *testing.T) {
	m := MetricsConfig{Enabled: true, Path: "/metrics", MaxHistory: 10}
	if err := validateMetrics(m); err == nil {
		t.Error("expected error for empty addr")
	}
}

func TestValidateMetrics_EmptyPath(t *testing.T) {
	m := MetricsConfig{Enabled: true, Addr: ":9090", MaxHistory: 10}
	if err := validateMetrics(m); err == nil {
		t.Error("expected error for empty path")
	}
}

func TestValidateMetrics_ZeroMaxHistory(t *testing.T) {
	m := MetricsConfig{Enabled: true, Addr: ":9090", Path: "/metrics", MaxHistory: 0}
	if err := validateMetrics(m); err == nil {
		t.Error("expected error for zero max_history")
	}
}

func TestValidateMetrics_Valid(t *testing.T) {
	m := DefaultMetrics()
	m.Enabled = true
	if err := validateMetrics(m); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
