package snapshot_test

import (
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func tempDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp("", "snapshot-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	return d
}

var epoch = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func TestSave_CreatesFile(t *testing.T) {
	dir := tempDir(t)
	path, err := snapshot.Save(dir, []int{80, 443}, epoch)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file at %s: %v", path, err)
	}
}

func TestList_ReturnsSortedEntries(t *testing.T) {
	dir := tempDir(t)
	t1 := epoch
	t2 := epoch.Add(time.Hour)

	if _, err := snapshot.Save(dir, []int{22}, t1); err != nil {
		t.Fatalf("Save t1: %v", err)
	}
	if _, err := snapshot.Save(dir, []int{80, 443}, t2); err != nil {
		t.Fatalf("Save t2: %v", err)
	}

	entries, err := snapshot.List(dir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Ports[0] != 22 {
		t.Errorf("first entry should have port 22, got %v", entries[0].Ports)
	}
	if len(entries[1].Ports) != 2 {
		t.Errorf("second entry should have 2 ports, got %v", entries[1].Ports)
	}
}

func TestSave_SortsPorts(t *testing.T) {
	dir := tempDir(t)
	_, err := snapshot.Save(dir, []int{443, 22, 80}, epoch)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	entries, _ := snapshot.List(dir)
	if entries[0].Ports[0] != 22 {
		t.Errorf("expected ports sorted, got %v", entries[0].Ports)
	}
}

func TestPrune_KeepsLatest(t *testing.T) {
	dir := tempDir(t)
	for i := 0; i < 5; i++ {
		ts := epoch.Add(time.Duration(i) * time.Minute)
		if _, err := snapshot.Save(dir, []int{i + 1}, ts); err != nil {
			t.Fatalf("Save %d: %v", i, err)
		}
	}
	if err := snapshot.Prune(dir, 3); err != nil {
		t.Fatalf("Prune: %v", err)
	}
	entries, _ := snapshot.List(dir)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries after prune, got %d", len(entries))
	}
	// Should keep the newest three (ports 3,4,5)
	if entries[0].Ports[0] != 3 {
		t.Errorf("expected oldest remaining to have port 3, got %v", entries[0].Ports)
	}
}

func TestList_EmptyDir(t *testing.T) {
	dir := tempDir(t)
	entries, err := snapshot.List(dir)
	if err != nil {
		t.Fatalf("List on empty dir: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}
