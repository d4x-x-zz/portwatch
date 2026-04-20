package config

import "testing"

func TestDefaultServer(t *testing.T) {
	s := DefaultServer()
	if s.Enabled {
		t.Error("expected disabled by default")
	}
	if s.Addr != ":9091" {
		t.Errorf("unexpected addr: %s", s.Addr)
	}
	if s.Path != "/status" {
		t.Errorf("unexpected path: %s", s.Path)
	}
}

func TestValidateServer_Disabled(t *testing.T) {
	if err := validateServer(ServerConfig{Enabled: false}); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestValidateServer_EmptyAddr(t *testing.T) {
	err := validateServer(ServerConfig{Enabled: true, Addr: "", Path: "/status"})
	if err == nil {
		t.Error("expected error for empty addr")
	}
}

func TestValidateServer_EmptyPath(t *testing.T) {
	err := validateServer(ServerConfig{Enabled: true, Addr: ":9091", Path: ""})
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestValidateServer_BadPath(t *testing.T) {
	err := validateServer(ServerConfig{Enabled: true, Addr: ":9091", Path: "status"})
	if err == nil {
		t.Error("expected error for path not starting with /")
	}
}

func TestValidateServer_Valid(t *testing.T) {
	err := validateServer(ServerConfig{Enabled: true, Addr: ":9091", Path: "/status"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
