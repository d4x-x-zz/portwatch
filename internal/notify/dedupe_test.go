package notify

import (
	"errors"
	"testing"

	"github.com/patrickward/portwatch/internal/monitor"
)

type recordingAlerter struct {
	calls []monitor.Diff
	err   error
}

func (r *recordingAlerter) Notify(d monitor.Diff) error {
	r.calls = append(r.calls, d)
	return r.err
}

func TestDedupeAlerter_ForwardsFirst(t *testing.T) {
	rec := &recordingAlerter{}
	da := NewDedupeAlerter(rec)

	d := makeDiff([]int{8080}, nil)
	if err := da.Notify(d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(rec.calls))
	}
}

func TestDedupeAlerter_SuppressesDuplicate(t *testing.T) {
	rec := &recordingAlerter{}
	da := NewDedupeAlerter(rec)

	d := makeDiff([]int{9090}, nil)
	_ = da.Notify(d)
	_ = da.Notify(d) // same diff — should be suppressed

	if len(rec.calls) != 1 {
		t.Fatalf("expected 1 call after duplicate, got %d", len(rec.calls))
	}
}

func TestDedupeAlerter_ForwardsDifferentDiff(t *testing.T) {
	rec := &recordingAlerter{}
	da := NewDedupeAlerter(rec)

	_ = da.Notify(makeDiff([]int{8080}, nil))
	_ = da.Notify(makeDiff([]int{9090}, nil))

	if len(rec.calls) != 2 {
		t.Fatalf("expected 2 calls for different diffs, got %d", len(rec.calls))
	}
}

func TestDedupeAlerter_Reset(t *testing.T) {
	rec := &recordingAlerter{}
	da := NewDedupeAlerter(rec)

	d := makeDiff([]int{7070}, nil)
	_ = da.Notify(d)
	da.Reset()
	_ = da.Notify(d) // should be forwarded again after reset

	if len(rec.calls) != 2 {
		t.Fatalf("expected 2 calls after reset, got %d", len(rec.calls))
	}
}

func TestDedupeAlerter_PropagatesError(t *testing.T) {
	want := errors.New("downstream failure")
	rec := &recordingAlerter{err: want}
	da := NewDedupeAlerter(rec)

	err := da.Notify(makeDiff([]int{8080}, nil))
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestDedupeAlerter_EmptyDiffAlwaysForwarded(t *testing.T) {
	rec := &recordingAlerter{}
	da := NewDedupeAlerter(rec)

	empty := monitor.Diff{}
	_ = da.Notify(empty)
	_ = da.Notify(empty)

	if len(rec.calls) != 2 {
		t.Fatalf("expected empty diffs to always be forwarded, got %d calls", len(rec.calls))
	}
}
