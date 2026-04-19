package history_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/portwatch/internal/history"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snapshot.json")
}

func TestStore_SaveAndLoad(t *testing.T) {
	store := history.NewStore(tempPath(t))

	ports := []uint16{22, 80, 443}
	if err := store.Save(ports); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if snap == nil {
		t.Fatal("expected snapshot, got nil")
	}
	if len(snap.Ports) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(snap.Ports))
	}
	for i, p := range ports {
		if snap.Ports[i] != p {
			t.Errorf("port[%d]: want %d got %d", i, p, snap.Ports[i])
		}
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestStore_Load_Missing(t *testing.T) {
	store := history.NewStore(filepath.Join(t.TempDir(), "nope.json"))
	snap, err := store.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap != nil {
		t.Fatal("expected nil snapshot for missing file")
	}
}

func TestStore_Load_Corrupt(t *testing.T) {
	p := tempPath(t)
	_ = os.WriteFile(p, []byte("not json{"), 0o644)
	store := history.NewStore(p)
	_, err := store.Load()
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}

func TestStore_Save_Overwrites(t *testing.T) {
	store := history.NewStore(tempPath(t))
	_ = store.Save([]uint16{22})
	_ = store.Save([]uint16{80, 443})
	snap, _ := store.Load()
	if len(snap.Ports) != 2 {
		t.Fatalf("expected 2 ports after overwrite, got %d", len(snap.Ports))
	}
}
