package daemon

import (
	"os"
	"path/filepath"
	"testing"
)

func tempPIDPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "portwatch.pid")
}

func TestWritePID_CreatesFile(t *testing.T) {
	path := tempPIDPath(t)
	if err := WritePID(path); err != nil {
		t.Fatalf("WritePID: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected pid file to exist: %v", err)
	}
}

func TestReadPID_MatchesCurrentProcess(t *testing.T) {
	path := tempPIDPath(t)
	if err := WritePID(path); err != nil {
		t.Fatalf("WritePID: %v", err)
	}
	pid, err := ReadPID(path)
	if err != nil {
		t.Fatalf("ReadPID: %v", err)
	}
	if pid != os.Getpid() {
		t.Errorf("got pid %d, want %d", pid, os.Getpid())
	}
}

func TestRemovePID_DeletesFile(t *testing.T) {
	path := tempPIDPath(t)
	if err := WritePID(path); err != nil {
		t.Fatalf("WritePID: %v", err)
	}
	if err := RemovePID(path); err != nil {
		t.Fatalf("RemovePID: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected pid file to be removed")
	}
}

func TestRemovePID_MissingFileIsOK(t *testing.T) {
	path := tempPIDPath(t)
	if err := RemovePID(path); err != nil {
		t.Errorf("RemovePID on missing file: %v", err)
	}
}

func TestReadPID_CorruptFile(t *testing.T) {
	path := tempPIDPath(t)
	if err := os.WriteFile(path, []byte("notanumber\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadPID(path); err == nil {
		t.Error("expected error for corrupt pid file")
	}
}
