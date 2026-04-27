package config

import "testing"

func TestDefaultEnvelope_Values(t *testing.T) {
	e := DefaultEnvelope()
	if e.Environment != "production" {
		t.Errorf("expected environment=production, got %q", e.Environment)
	}
	if e.Tags == nil {
		t.Error("expected non-nil Tags map")
	}
	if e.Host != "" {
		t.Errorf("expected empty host, got %q", e.Host)
	}
}

func TestValidateEnvelope_Valid(t *testing.T) {
	e := DefaultEnvelope()
	e.Host = "myhost"
	e.Tags["team"] = "infra"
	if err := validateEnvelope(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateEnvelope_EmptyTagKey(t *testing.T) {
	e := DefaultEnvelope()
	e.Tags[""] = "oops"
	if err := validateEnvelope(e); err == nil {
		t.Fatal("expected error for empty tag key, got nil")
	}
}

func TestValidateEnvelope_NoTags(t *testing.T) {
	e := Envelope{Environment: "staging", Tags: map[string]string{}}
	if err := validateEnvelope(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateEnvelope_MultipleTagsAllValid(t *testing.T) {
	e := DefaultEnvelope()
	e.Tags["region"] = "us-east-1"
	e.Tags["env"] = "prod"
	if err := validateEnvelope(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
