package config

import "fmt"

// Suppress controls alert suppression behaviour.
type Suppress struct {
	// MaxRepeats is the number of times the same diff may be alerted
	// before it is suppressed. 0 means suppress after the first alert.
	MaxRepeats int `toml:"max_repeats"`

	// ResetOnChange re-enables alerting when the diff changes.
	ResetOnChange bool `toml:"reset_on_change"`
}

// DefaultSuppress returns sensible defaults for suppression.
func DefaultSuppress() Suppress {
	return Suppress{
		MaxRepeats:    1,
		ResetOnChange: true,
	}
}

func validateSuppress(s Suppress) error {
	if s.MaxRepeats < 0 {
		return fmt.Errorf("suppress.max_repeats must be >= 0, got %d", s.MaxRepeats)
	}
	return nil
}
