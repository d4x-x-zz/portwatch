package config

import (
	"fmt"
	"time"
)

// AuditRotateConfig controls automatic rotation of the audit log file.
type AuditRotateConfig struct {
	// Enabled turns rotation on or off.
	Enabled bool `toml:"enabled"`

	// MaxSizeBytes is the maximum file size before rotation is triggered.
	MaxSizeBytes int64 `toml:"max_size_bytes"`

	// MaxAge is the maximum age of the log file before rotation.
	MaxAge time.Duration `toml:"max_age"`

	// KeepLast is the number of rotated files to retain.
	KeepLast int `toml:"keep_last"`
}

// DefaultAuditRotate returns sensible rotation defaults.
func DefaultAuditRotate() AuditRotateConfig {
	return AuditRotateConfig{
		Enabled:      true,
		MaxSizeBytes: 10 * 1024 * 1024, // 10 MiB
		MaxAge:       7 * 24 * time.Hour,
		KeepLast:     5,
	}
}

// validateAuditRotate checks that rotation settings are coherent.
func validateAuditRotate(r AuditRotateConfig) error {
	if !r.Enabled {
		return nil
	}
	if r.MaxSizeBytes <= 0 {
		return fmt.Errorf("audit.rotate.max_size_bytes must be positive")
	}
	if r.MaxAge <= 0 {
		return fmt.Errorf("audit.rotate.max_age must be positive")
	}
	if r.KeepLast <= 0 {
		return fmt.Errorf("audit.rotate.keep_last must be positive")
	}
	return nil
}
