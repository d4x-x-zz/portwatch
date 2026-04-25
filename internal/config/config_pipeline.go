package config

import "fmt"

// PipelineConfig controls the order in which notify middleware is applied.
// Stages are applied innermost-first, so the first entry in Stages wraps
// closest to the base alerter.
type PipelineConfig struct {
	Stages []string `toml:"stages"`
}

// knownStages is the set of valid stage names.
var knownStages = map[string]bool{
	"dedupe":   true,
	"suppress": true,
	"debounce": true,
	"throttle": true,
	"ratelimit": true,
	"circuit":  true,
	"retry":    true,
}

// DefaultPipeline returns a sensible default ordering of notify stages.
func DefaultPipeline() PipelineConfig {
	return PipelineConfig{
		Stages: []string{
			"dedupe",
			"suppress",
			"throttle",
			"circuit",
			"retry",
		},
	}
}

// validatePipeline checks that every declared stage name is recognised and
// that no stage appears more than once.
func validatePipeline(p PipelineConfig) error {
	seen := make(map[string]bool, len(p.Stages))
	for i, s := range p.Stages {
		if !knownStages[s] {
			return fmt.Errorf("pipeline: unknown stage %q at index %d", s, i)
		}
		if seen[s] {
			return fmt.Errorf("pipeline: duplicate stage %q at index %d", s, i)
		}
		seen[s] = true
	}
	return nil
}
