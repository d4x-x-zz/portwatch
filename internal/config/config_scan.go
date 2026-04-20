package config

import "fmt"

// ScanConfig holds port scanning configuration.
type ScanConfig struct {
	// Host is the target host to scan.
	Host string `toml:"host"`
	// PortRange is the range of ports to scan, e.g. "1-1024".
	PortRange string `toml:"port_range"`
	// Concurrency is the number of parallel scan workers.
	Concurrency int `toml:"concurrency"`
	// TimeoutMs is the per-port dial timeout in milliseconds.
	TimeoutMs int `toml:"timeout_ms"`
}

// DefaultScan returns a ScanConfig with sensible defaults.
func DefaultScan() ScanConfig {
	return ScanConfig{
		Host:        "127.0.0.1",
		PortRange:   "1-1024",
		Concurrency: 100,
		TimeoutMs:   500,
	}
}

func validateScan(s ScanConfig) error {
	if s.Host == "" {
		return fmt.Errorf("scan.host must not be empty")
	}
	if s.PortRange == "" {
		return fmt.Errorf("scan.port_range must not be empty")
	}
	if s.Concurrency <= 0 {
		return fmt.Errorf("scan.concurrency must be greater than zero")
	}
	if s.TimeoutMs <= 0 {
		return fmt.Errorf("scan.timeout_ms must be greater than zero")
	}
	return nil
}
