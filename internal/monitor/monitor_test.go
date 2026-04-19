package monitor_test

import (
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func TestCompare_NoChanges(t *testing.T) {
	prev := monitor.NewPortSet([]int{80, 443})
	curr := monitor.NewPortSet([]int{80, 443})
	d := monitor.Compare(prev, curr)
	if d.HasChanges() {
		t.Errorf("expected no changes, got %s", d)
	}
}

func TestCompare_OpenedPort(t *testing.T) {
	prev := monitor.NewPortSet([]int{80})
	curr := monitor.NewPortSet([]int{80, 8080})
	d := monitor.Compare(prev, curr)
	if !d.HasChanges() {
		t.Fatal("expected changes")
	}
	if len(d.Opened) != 1 || d.Opened[0] != 8080 {
		t.Errorf("expected opened [8080], got %v", d.Opened)
	}
	if len(d.Closed) != 0 {
		t.Errorf("expected no closed ports, got %v", d.Closed)
	}
}

func TestCompare_ClosedPort(t *testing.T) {
	prev := monitor.NewPortSet([]int{80, 443})
	curr := monitor.NewPortSet([]int{80})
	d := monitor.Compare(prev, curr)
	if !d.HasChanges() {
		t.Fatal("expected changes")
	}
	if len(d.Closed) != 1 || d.Closed[0] != 443 {
		t.Errorf("expected closed [443], got %v", d.Closed)
	}
}

func TestDiff_String(t *testing.T) {
	prev := monitor.NewPortSet([]int{80})
	curr := monitor.NewPortSet([]int{443})
	d := monitor.Compare(prev, curr)
	s := d.String()
	if s == "no changes" {
		t.Error("expected non-empty diff string")
	}
}

func TestDiff_StringNoChanges(t *testing.T) {
	d := monitor.Diff{}
	if d.String() != "no changes" {
		t.Errorf("expected 'no changes', got %q", d.String())
	}
}
