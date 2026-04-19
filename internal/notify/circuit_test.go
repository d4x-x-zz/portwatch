package notify

import (
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

type stubAlerter struct {
	calls int
	err   error
}

func (s *stubAlerter) Notify(_ monitor.Diff) error {
	s.calls++
	return s.err
}

func TestCircuitBreaker_ClosedOnSuccess(t *testing.T) {
	stub := &stubAlerter{}
	cb := NewCircuitBreaker(stub, 3, time.Second)

	if err := cb.Notify(monitor.Diff{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cb.State() != StateClosed {
		t.Errorf("expected closed state")
	}
}

func TestCircuitBreaker_OpensAfterThreshold(t *testing.T) {
	stub := &stubAlerter{err: errors.New("fail")}
	cb := NewCircuitBreaker(stub, 3, time.Minute)

	for i := 0; i < 3; i++ {
		_ = cb.Notify(monitor.Diff{})
	}

	if cb.State() != StateOpen {
		t.Errorf("expected open state after threshold failures")
	}
}

func TestCircuitBreaker_BlocksWhenOpen(t *testing.T) {
	stub := &stubAlerter{err: errors.New("fail")}
	cb := NewCircuitBreaker(stub, 2, time.Minute)

	for i := 0; i < 2; i++ {
		_ = cb.Notify(monitor.Diff{})
	}

	callsBefore := stub.calls
	err := cb.Notify(monitor.Diff{})
	if err == nil {
		t.Error("expected error when circuit open")
	}
	if stub.calls != callsBefore {
		t.Error("inner alerter should not be called when circuit is open")
	}
}

func TestCircuitBreaker_HalfOpenAfterCooldown(t *testing.T) {
	stub := &stubAlerter{err: errors.New("fail")}
	cb := NewCircuitBreaker(stub, 2, 10*time.Millisecond)

	for i := 0; i < 2; i++ {
		_ = cb.Notify(monitor.Diff{})
	}

	time.Sleep(20 * time.Millisecond)
	stub.err = nil

	if err := cb.Notify(monitor.Diff{}); err != nil {
		t.Fatalf("unexpected error after cooldown: %v", err)
	}
	if cb.State() != StateClosed {
		t.Errorf("expected closed after successful half-open attempt")
	}
}
