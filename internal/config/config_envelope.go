package config

import "fmt"

// Envelope holds optional metadata attached to every alert payload.
type Envelope struct {
	// Host is the hostname included in alert messages.
	// Defaults to the system hostname when empty.
	Host string `toml:"host"`

	// Tags are arbitrary key=value pairs forwarded to alerters.
	Tags map[string]string `toml:"tags"`

	// Environment labels the deployment context (e.g. "prod", "staging").
	Environment string `toml:"environment"`
}

// DefaultEnvelope returns a sensible zero-value Envelope.
func DefaultEnvelope() Envelope {
	return Envelope{
		Tags:        map[string]string{},
		Environment: "production",
	}
}

func validateEnvelope(e Envelope) error {
	for k, v := range e.Tags {
		if k == "" {
			return fmt.Errorf("envelope: tag key must not be empty (value %q)", v)
		}
	}
	return nil
}
