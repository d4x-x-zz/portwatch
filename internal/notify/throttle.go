// Package notify provides alerting middleware such as rate limiting,
// debouncing, and throttling.
package notify

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// ThrottledAlerter wraps an Alerter and ensures Notify is called at most once
// per interval, regardless of how many diffs arrive.
type ThrottledAlerter struct {
	mu       sync.Mutex
	inner    Alerter
	interval time.Duration
	lastFire time.Time
	pending  *monitor.Diff
	stop     chan struct{}
}

// Alerter is the interface satisfied by all alert backends.
type Alerter interface {
	Notify(d monitor.Diff) error
}

// NewThrottledAlerter returns a ThrottledAlerter that forwards to inner at
// most once per interval. Pending diffs are merged; the latest wins.
func NewThrottledAlerter(inner Alerter, interval time.Duration) *ThrottledAlerter {
	return &ThrottledAlerter{
		inner:    inner,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Notify queues or immediately forwards d depending on the throttle window.
func (t *ThrottledAlerter) Notify(d monitor.Diff) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	if t.lastFire.IsZero() || now.Sub(t.lastFire) >= t.interval {
		t.lastFire = now
		t.pending = nil
		return t.inner.Notify(d)
	}

	// Store latest pending diff to send after the window expires.
	copy := d
	t.pending = &copy
	return nil
}

// Flush sends any pending diff immediately, resetting the throttle window.
func (t *ThrottledAlerter) Flush() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.pending == nil {
		return nil
	}
	d := *t.pending
	t.pending = nil
	t.lastFire = time.Now()
	return t.inner.Notify(d)
}
