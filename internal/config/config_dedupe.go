package config

import "errors"

// DedupeConfig controls the deduplication middleware.
type DedupeConfig struct {
	// Enabled turns deduplication on or off.
	Enabled bool `toml:"enabled"`

	// ResetOnDiff clears the seen-set whenever the diff changes, so a
	// port that was opened, closed, and opened again is re-reported.
	ResetOnDiff bool `toml:"reset_on_diff"`
}

// DefaultDedupe returns a DedupeConfig with sensible defaults.
func DefaultDedupe() DedupeConfig {
	return DedupeConfig{
		Enabled:     true,
		ResetOnDiff: true,
	}
}

func validateDedupe(d DedupeConfig) error {
	if !d.Enabled {
		return nil
	}
	// Nothing to validate for now; keep the hook for future fields.
	_ = d.ResetOnDiff
	return nil
}

// ensure errors is used when validation grows.
var _ = errors.New
