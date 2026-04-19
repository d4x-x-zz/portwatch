package notify_test

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/notify"
)

type countAlerter struct {
	calls atomic.Int32
	err   error
}

func (c *countAlerter) Notify(_ monitor.Diff) error {
	c.calls.Add(1)
	return c.err
}

func TestRateLimiter_AllowsFirstCall(t *testing.T) {
	inner := &countAlerter{}
	rl := notify.NewRateLimiter(inner, 5*time.Second)
	d := monitor.Diff{Opened: []int{8080}}

	if err := rl.Notify(d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inner.calls.Load() != 1 {
		t.Fatalf("expected 1 call, got %d", inner.calls.Load())
	}
}

func TestRateLimiter_SuppressesDuplicate(t *testing.T) {
	inner := &countAlerter{}
	rl := notify.NewRateLimiter(inner, 5*time.Second)
	d := monitor.Diff{Opened: []int{8080}}

	_ = rl.Notify(d)
	_ = rl.Notify(d)

	if inner.calls.Load() != 1 {
		t.Fatalf("expected 1 call, got %d", inner.calls.Load())
	}
}

func TestRateLimiter_AllowsDifferentDiff(t *testing.T) {
	inner := &countAlerter{}
	rl := notify.NewRateLimiter(inner, 5*time.Second)

	_ = rl.Notify(monitor.Diff{Opened: []int{8080}})
	_ = rl.Notify(monitor.Diff{Opened: []int{9090}})

	if inner.calls.Load() != 2 {
		t.Fatalf("expected 2 calls, got %d", inner.calls.Load())
	}
}

func TestRateLimiter_AllowsAfterCooldown(t *testing.T) {
	inner := &countAlerter{}
	rl := notify.NewRateLimiter(inner, 10*time.Millisecond)
	d := monitor.Diff{Opened: []int{8080}}

	_ = rl.Notify(d)
	time.Sleep(20 * time.Millisecond)
	_ = rl.Notify(d)

	if inner.calls.Load() != 2 {
		t.Fatalf("expected 2 calls, got %d", inner.calls.Load())
	}
}

func TestRateLimiter_PropagatesError(t *testing.T) {
	expected := errors.New("boom")
	inner := &countAlerter{err: expected}
	rl := notify.NewRateLimiter(inner, 5*time.Second)

	err := rl.Notify(monitor.Diff{Opened: []int{1234}})
	if !errors.Is(err, expected) {
		t.Fatalf("expected error %v, got %v", expected, err)
	}
}
