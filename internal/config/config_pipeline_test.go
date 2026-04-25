package config

import "testing"

func TestDefaultPipeline_Values(t *testing.T) {
	p := DefaultPipeline()
	if len(p.Stages) == 0 {
		t.Fatal("expected non-empty default stages")
	}
	for _, s := range p.Stages {
		if !knownStages[s] {
			t.Errorf("default stage %q is not in knownStages", s)
		}
	}
}

func TestValidatePipeline_Valid(t *testing.T) {
	p := DefaultPipeline()
	if err := validatePipeline(p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidatePipeline_Empty(t *testing.T) {
	p := PipelineConfig{Stages: []string{}}
	if err := validatePipeline(p); err != nil {
		t.Fatalf("empty stages should be valid, got: %v", err)
	}
}

func TestValidatePipeline_UnknownStage(t *testing.T) {
	p := PipelineConfig{Stages: []string{"dedupe", "magic"}}
	if err := validatePipeline(p); err == nil {
		t.Fatal("expected error for unknown stage")
	}
}

func TestValidatePipeline_DuplicateStage(t *testing.T) {
	p := PipelineConfig{Stages: []string{"retry", "dedupe", "retry"}}
	if err := validatePipeline(p); err == nil {
		t.Fatal("expected error for duplicate stage")
	}
}

func TestValidatePipeline_SingleStage(t *testing.T) {
	p := PipelineConfig{Stages: []string{"circuit"}}
	if err := validatePipeline(p); err != nil {
		t.Fatalf("single valid stage should pass, got: %v", err)
	}
}
