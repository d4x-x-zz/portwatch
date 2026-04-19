package notify

import (
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

type countAlerter struct {
	calls  int
	failOn int // fail this many times before succeeding
}

func (c *countAlerter) Notify(_ monitor.Diff) error {
	c.calls++
	if c.calls <= c.failOn {
		return errors.New("transient error")
	}
	return nil
}

func TestRetryAlerter_SucceedsFirstTry(t *testing.T) {
	inner := &countAlerter{}
	r := NewRetryAlerter(inner, 3, 0)
	if err := r.Notify(monitor.Diff{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inner.calls != 1 {
		t.Fatalf("expected 1 call, got %d", inner.calls)
	}
}

func TestRetryAlerter_RetriesOnFailure(t *testing.T) {
	inner := &countAlerter{failOn: 2}
	r := NewRetryAlerter(inner, 3, 0)
	if err := r.Notify(monitor.Diff{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inner.calls != 3 {
		t.Fatalf("expected 3 calls, got %d", inner.calls)
	}
}

func TestRetryAlerter_ReturnsErrorAfterExhaustion(t *testing.T) {
	inner := &countAlerter{failOn: 10}
	r := NewRetryAlerter(inner, 3, 0)
	if err := r.Notify(monitor.Diff{}); err == nil {
		t.Fatal("expected error, got nil")
	}
	if inner.calls != 3 {
		t.Fatalf("expected 3 calls, got %d", inner.calls)
	}
}

func TestRetryAlerter_RespectsDelay(t *testing.T) {
	inner := &countAlerter{failOn: 1}
	r := NewRetryAlerter(inner, 2, 20*time.Millisecond)
	start := time.Now()
	_ = r.Notify(monitor.Diff{})
	if time.Since(start) < 20*time.Millisecond {
		t.Fatal("expected delay between retries")
	}
}

func TestNewRetryAlerter_MinAttemptsOne(t *testing.T) {
	inner := &countAlerter{failOn: 10}
	r := NewRetryAlerter(inner, 0, 0)
	_ = r.Notify(monitor.Diff{})
	if inner.calls != 1 {
		t.Fatalf("expected 1 call, got %d", inner.calls)
	}
}
