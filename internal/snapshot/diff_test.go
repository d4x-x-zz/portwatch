package snapshot_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func makeEntry(ports []int, t time.Time) snapshot.Entry {
	return snapshot.Entry{Ports: ports, Time: t}
}

func TestDiff_NoChanges(t *testing.T) {
	old := makeEntry([]int{80, 443}, time.Now().Add(-time.Hour))
	new := makeEntry([]int{80, 443}, time.Now())

	d := snapshot.Diff(old, new)
	if !d.Empty() {
		t.Errorf("expected empty diff, got opened=%v closed=%v", d.Opened, d.Closed)
	}
}

func TestDiff_OpenedPorts(t *testing.T) {
	old := makeEntry([]int{80}, time.Now().Add(-time.Hour))
	new := makeEntry([]int{80, 443, 8080}, time.Now())

	d := snapshot.Diff(old, new)
	if len(d.Opened) != 2 {
		t.Fatalf("expected 2 opened ports, got %d", len(d.Opened))
	}
	if d.Opened[0] != 443 || d.Opened[1] != 8080 {
		t.Errorf("unexpected opened ports: %v", d.Opened)
	}
	if len(d.Closed) != 0 {
		t.Errorf("expected no closed ports, got %v", d.Closed)
	}
}

func TestDiff_ClosedPorts(t *testing.T) {
	old := makeEntry([]int{80, 443, 8080}, time.Now().Add(-time.Hour))
	new := makeEntry([]int{80}, time.Now())

	d := snapshot.Diff(old, new)
	if len(d.Closed) != 2 {
		t.Fatalf("expected 2 closed ports, got %d", len(d.Closed))
	}
	if d.Closed[0] != 443 || d.Closed[1] != 8080 {
		t.Errorf("unexpected closed ports: %v", d.Closed)
	}
}

func TestDiff_String_WithChanges(t *testing.T) {
	old := makeEntry([]int{80}, time.Now().Add(-time.Hour))
	new := makeEntry([]int{443}, time.Now())

	d := snapshot.Diff(old, new)
	s := d.String()
	if s == "no changes between snapshots" {
		t.Error("expected non-empty diff string")
	}
}

func TestDiff_String_NoChanges(t *testing.T) {
	old := makeEntry([]int{80}, time.Now().Add(-time.Hour))
	new := makeEntry([]int{80}, time.Now())

	d := snapshot.Diff(old, new)
	if d.String() != "no changes between snapshots" {
		t.Errorf("unexpected string: %s", d.String())
	}
}

func TestDiff_Empty_True(t *testing.T) {
	d := snapshot.SnapshotDiff{}
	if !d.Empty() {
		t.Error("zero-value SnapshotDiff should be empty")
	}
}
