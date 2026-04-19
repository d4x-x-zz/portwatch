package config

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, cfg *Config) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-config-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if err := json.NewEncoder(f).Encode(cfg); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	cfg := &Config{
		PortRange: PortRange{From: 80, To: 443},
		Interval:  10 * time.Second,
	}
	path := writeTempConfig(t, cfg)
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.PortRange.From != 80 || loaded.PortRange.To != 443 {
		t.Errorf("unexpected port range: %+v", loaded.PortRange)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestValidate_InvalidRange(t *testing.T) {
	cfg := &Config{
		PortRange: PortRange{From: 500, To: 100},
		Interval:  5 * time.Second,
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for inverted range")
	}
}

func TestValidate_ZeroInterval(t *testing.T) {
	cfg := &Config{
		PortRange: PortRange{From: 1, To: 1024},
		Interval:  0,
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for zero interval")
	}
}

func TestDefault(t *testing.T) {
	cfg := Default()
	if err := cfg.Validate(); err != nil {
		t.Errorf("default config should be valid, got: %v", err)
	}
}
