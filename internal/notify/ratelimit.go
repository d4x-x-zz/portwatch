// Package notify provides rate-limiting and deduplication for alerters.
package notify

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
)

// RateLimiter wraps an Alerter and suppresses duplicate notifications
// within a configurable cooldown window.
type RateLimiter struct {
	mu       sync.Mutex
	inner    alert.Alerter
	cooldown time.Duration
	lastSent time.Time
	lastDiff monitor.Diff
}

// NewRateLimiter returns a RateLimiter that forwards to inner at most once
// per cooldown duration for identical diffs.
func NewRateLimiter(inner alert.Alerter, cooldown time.Duration) *RateLimiter {
	return &RateLimiter{
		inner:    inner,
		cooldown: cooldown,
	}
}

// Notify forwards the diff to the inner alerter unless the same diff was
// already sent within the cooldown window.
func (r *RateLimiter) Notify(d monitor.Diff) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	sameContent := diffEqual(r.lastDiff, d)
	withinWindow := now.Sub(r.lastSent) < r.cooldown

	if sameContent && withinWindow {
		return nil
	}

	if err := r.inner.Notify(d); err != nil {
		return err
	}

	r.lastSent = now
	r.lastDiff = d
	return nil
}

func diffEqual(a, b monitor.Diff) bool {
	if len(a.Opened) != len(b.Opened) || len(a.Closed) != len(b.Closed) {
		return false
	}
	for i := range a.Opened {
		if a.Opened[i] != b.Opened[i] {
			return false
		}
	}
	for i := range a.Closed {
		if a.Closed[i] != b.Closed[i] {
			return false
		}
	}
	return true
}
