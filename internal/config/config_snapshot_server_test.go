package config

import (
	"testing"
)

func TestDefaultSnapshotServer_Values(t *testing.T) {
	s := DefaultSnapshotServer()
	if s.Addr == "" {
		t.Error("expected non-empty default addr")
	}
	if s.Path == "" {
		t.Error("expected non-empty default path")
	}
}

func TestValidateSnapshotServer_Disabled(t *testing.T) {
	s := DefaultSnapshotServer()
	s.Enabled = false
	if err := validateSnapshotServer(s); err != nil {
		t.Fatalf("unexpected error for disabled server: %v", err)
	}
}

func TestValidateSnapshotServer_Valid(t *testing.T) {
	s := DefaultSnapshotServer()
	s.Enabled = true
	if err := validateSnapshotServer(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateSnapshotServer_EmptyAddr(t *testing.T) {
	s := DefaultSnapshotServer()
	s.Enabled = true
	s.Addr = ""
	if err := validateSnapshotServer(s); err == nil {
		t.Fatal("expected error for empty addr")
	}
}

func TestValidateSnapshotServer_EmptyPath(t *testing.T) {
	s := DefaultSnapshotServer()
	s.Enabled = true
	s.Path = ""
	if err := validateSnapshotServer(s); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestValidateSnapshotServer_BadPath(t *testing.T) {
	s := DefaultSnapshotServer()
	s.Enabled = true
	s.Path = "no-leading-slash"
	if err := validateSnapshotServer(s); err == nil {
		t.Fatal("expected error for path missing leading slash")
	}
}
