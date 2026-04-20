package config

import (
	"errors"
	"time"
)

// Audit holds configuration for the audit log.
type Audit struct {
	// Enabled controls whether audit logging is active.
	Enabled bool `toml:"enabled"`
	// Path is the file path for the audit log.
	Path string `toml:"path"`
	// RotateMaxBytes triggers rotation when the file exceeds this size (bytes). 0 = disabled.
	RotateMaxBytes int64 `toml:"rotate_max_bytes"`
	// RotateMaxAge triggers rotation when the file is older than this duration. 0 = disabled.
	RotateMaxAge duration `toml:"rotate_max_age"`
}

type duration struct{ time.Duration }

func (d *duration) UnmarshalText(b []byte) error {
	v, err := time.ParseDuration(string(b))
	if err != nil {
		return err
	}
	d.Duration = v
	return nil
}

// DefaultAudit returns a sensible audit configuration.
func DefaultAudit() Audit {
	return Audit{
		Enabled:        true,
		Path:           "/var/lib/portwatch/audit.log",
		RotateMaxBytes: 10 * 1024 * 1024,
		RotateMaxAge:   duration{24 * time.Hour},
	}
}

func validateAudit(a Audit) error {
	if !a.Enabled {
		return nil
	}
	if a.Path == "" {
		return errors.New("audit.path must not be empty when audit is enabled")
	}
	if a.RotateMaxBytes < 0 {
		return errors.New("audit.rotate_max_bytes must be >= 0")
	}
	if a.RotateMaxAge.Duration < 0 {
		return errors.New("audit.rotate_max_age must be >= 0")
	}
	return nil
}
