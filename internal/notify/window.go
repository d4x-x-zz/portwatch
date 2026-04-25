package notify

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// WindowOptions configures the sliding-window rate limiter.
type WindowOptions struct {
	// Window is the duration of the sliding window.
	Window time.Duration
	// MaxCalls is the maximum number of Notify calls allowed within Window.
	MaxCalls int
}

type windowLimiter struct {
	mu       sync.Mutex
	next     Alerter
	opts     WindowOptions
	timestamps []time.Time
	now      func() time.Time
}

// NewWindowLimiter returns an Alerter that forwards at most opts.MaxCalls
// notifications per opts.Window duration, dropping excess calls silently.
func NewWindowLimiter(next Alerter, opts WindowOptions) Alerter {
	return &windowLimiter{
		next: next,
		opts: opts,
		now:  time.Now,
	}
}

func (w *windowLimiter) Notify(d monitor.Diff) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := w.now()
	cutoff := now.Add(-w.opts.Window)

	// evict timestamps outside the window
	valid := w.timestamps[:0]
	for _, t := range w.timestamps {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	w.timestamps = valid

	if len(w.timestamps) >= w.opts.MaxCalls {
		return nil // silently drop
	}

	w.timestamps = append(w.timestamps, now)
	return w.next.Notify(d)
}
