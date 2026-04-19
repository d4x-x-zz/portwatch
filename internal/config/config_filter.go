package config

import (
	"fmt"
	"strings"
)

// FilterConfig holds port filter rules loaded from configuration.
type FilterConfig struct {
	Rules []string `toml:"rules"`
}

// validateFilter checks that filter rules are syntactically valid.
func validateFilter(fc FilterConfig) error {
	for _, rule := range fc.Rules {
		if err := validateFilterRule(rule); err != nil {
			return err
		}
	}
	return nil
}

func validateFilterRule(s string) error {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("config: filter rule %q missing colon separator", s)
	}
	kind := parts[0]
	if kind != "allow" && kind != "deny" {
		return fmt.Errorf("config: filter rule %q has unknown type %q", s, kind)
	}
	rangeStr := parts[1]
	var start, end int
	if _, err := fmt.Sscanf(rangeStr, "%d-%d", &start, &end); err != nil {
		if _, err2 := fmt.Sscanf(rangeStr, "%d", &start); err2 != nil {
			return fmt.Errorf("config: filter rule %q has invalid port %q", s, rangeStr)
		}
		end = start
	}
	if start < 1 || end > 65535 || start > end {
		return fmt.Errorf("config: filter rule %q port range %d-%d out of bounds", s, start, end)
	}
	return nil
}
