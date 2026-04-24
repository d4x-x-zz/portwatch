package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/wjam/portwatch/internal/snapshot"
)

func TestLoad_ValidSnapshot(t *testing.T) {
	dir := t.TempDir()
	ports := []int{22, 80, 443}
	if err := snapshot.Save(dir, ports, 10); err != nil {
		t.Fatalf("Save: %v", err)
	}

	entries, err := snapshot.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one snapshot entry")
	}

	path := filepath.Join(dir, entries[0].Name())
	e, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(e.Ports) != len(ports) {
		t.Errorf("ports len: got %d, want %d", len(e.Ports), len(ports))
	}
	if e.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/snapshot.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_CorruptFile(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(bad, []byte("not-json{"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	_, err := snapshot.Load(bad)
	if err == nil {
		t.Fatal("expected decode error")
	}
}

func TestLatest_ReturnsNewest(t *testing.T) {
	dir := t.TempDir()

	if err := snapshot.Save(dir, []int{22}, 10); err != nil {
		t.Fatalf("Save first: %v", err)
	}
	time.Sleep(10 * time.Millisecond)
	if err := snapshot.Save(dir, []int{22, 80}, 10); err != nil {
		t.Fatalf("Save second: %v", err)
	}

	e, err := snapshot.Latest(dir)
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil entry")
	}
	if len(e.Ports) != 2 {
		t.Errorf("expected 2 ports in latest snapshot, got %d", len(e.Ports))
	}
}

func TestLatest_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	e, err := snapshot.Latest(dir)
	if err != nil {
		t.Fatalf("Latest on empty dir: %v", err)
	}
	if e != nil {
		t.Errorf("expected nil entry for empty dir, got %+v", e)
	}
}
