package notify

import (
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

func TestWindowLimiter_AllowsUpToMax(t *testing.T) {
	called := 0
	base := alerterFunc(func(_ monitor.Diff) error { called++; return nil })
	lim := NewWindowLimiter(base, WindowOptions{Window: time.Minute, MaxCalls: 3})

	d := monitor.Diff{Opened: []int{80}}
	for i := 0; i < 3; i++ {
		if err := lim.Notify(d); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if called != 3 {
		t.Fatalf("expected 3 calls, got %d", called)
	}
}

func TestWindowLimiter_DropsOverMax(t *testing.T) {
	called := 0
	base := alerterFunc(func(_ monitor.Diff) error { called++; return nil })
	lim := NewWindowLimiter(base, WindowOptions{Window: time.Minute, MaxCalls: 2})

	d := monitor.Diff{Opened: []int{443}}
	for i := 0; i < 5; i++ {
		_ = lim.Notify(d)
	}
	if called != 2 {
		t.Fatalf("expected 2 calls, got %d", called)
	}
}

func TestWindowLimiter_ResetsAfterWindow(t *testing.T) {
	called := 0
	base := alerterFunc(func(_ monitor.Diff) error { called++; return nil })

	now := time.Now()
	wl := &windowLimiter{
		next: base,
		opts: WindowOptions{Window: time.Second, MaxCalls: 1},
		now:  func() time.Time { return now },
	}

	d := monitor.Diff{Opened: []int{22}}
	_ = wl.Notify(d) // allowed
	_ = wl.Notify(d) // dropped

	now = now.Add(2 * time.Second) // advance past window
	_ = wl.Notify(d)              // allowed again

	if called != 2 {
		t.Fatalf("expected 2 calls after reset, got %d", called)
	}
}

func TestWindowLimiter_PropagatesError(t *testing.T) {
	sentinel := errors.New("boom")
	base := alerterFunc(func(_ monitor.Diff) error { return sentinel })
	lim := NewWindowLimiter(base, WindowOptions{Window: time.Minute, MaxCalls: 5})

	err := lim.Notify(monitor.Diff{Opened: []int{8080}})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}
