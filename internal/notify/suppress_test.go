package notify

import (
	"errors"
	"testing"

	"github.com/user/portwatch/internal/monitor"
)

func TestSuppressAlerter_ForwardsUnderLimit(t *testing.T) {
	calls := 0
	stub := alerterFunc(func(_ monitor.Diff) error { calls++; return nil })
	s := NewSuppressAlerter(stub, 3, false)
	d := makeDiff([]uint16{80}, nil)

	for i := 0; i < 3; i++ {
		if err := s.Notify(d); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestSuppressAlerter_SuppressesAfterLimit(t *testing.T) {
	calls := 0
	stub := alerterFunc(func(_ monitor.Diff) error { calls++; return nil })
	s := NewSuppressAlerter(stub, 2, false)
	d := makeDiff([]uint16{443}, nil)

	for i := 0; i < 5; i++ {
		_ = s.Notify(d)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestSuppressAlerter_ResetOnChange(t *testing.T) {
	calls := 0
	stub := alerterFunc(func(_ monitor.Diff) error { calls++; return nil })
	s := NewSuppressAlerter(stub, 1, true)

	d1 := makeDiff([]uint16{80}, nil)
	d2 := makeDiff([]uint16{443}, nil)

	_ = s.Notify(d1) // count=1, allowed
	_ = s.Notify(d1) // count=2, suppressed
	_ = s.Notify(d2) // different → reset, count=1, allowed
	_ = s.Notify(d2) // count=2, suppressed

	if calls != 2 {
		t.Fatalf("expected 2 calls after resets, got %d", calls)
	}
}

func TestSuppressAlerter_Reset(t *testing.T) {
	calls := 0
	stub := alerterFunc(func(_ monitor.Diff) error { calls++; return nil })
	s := NewSuppressAlerter(stub, 1, false)
	d := makeDiff([]uint16{22}, nil)

	_ = s.Notify(d) // allowed
	_ = s.Notify(d) // suppressed
	s.Reset()
	_ = s.Notify(d) // allowed again after reset

	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestSuppressAlerter_PropagatesError(t *testing.T) {
	want := errors.New("boom")
	stub := alerterFunc(func(_ monitor.Diff) error { return want })
	s := NewSuppressAlerter(stub, 5, false)
	d := makeDiff([]uint16{8080}, nil)

	if err := s.Notify(d); !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}
