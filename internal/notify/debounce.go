package notify

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// Debouncer suppresses rapid successive alerts by waiting for a quiet period
// before forwarding a diff to the wrapped notify function.
type Debouncer struct {
	mu      sync.Mutex
	wait    time.Duration
	timer   *time.Timer
	pending *monitor.Diff
	notify  func(monitor.Diff) error
}

// NewDebouncer returns a Debouncer that waits for the given quiet period before
// calling notify with the most recent diff.
func NewDebouncer(wait time.Duration, notify func(monitor.Diff) error) *Debouncer {
	return &Debouncer{
		wait:   wait,
		notify: notify,
	}
}

// Submit schedules a notification. If another diff arrives before the quiet
// period elapses the timer is reset and only the latest diff is forwarded.
func (d *Debouncer) Submit(diff monitor.Diff) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.pending = &diff

	if d.timer != nil {
		d.timer.Stop()
	}

	d.timer = time.AfterFunc(d.wait, func() {
		d.mu.Lock()
		pending := d.pending
		d.pending = nil
		d.mu.Unlock()

		if pending != nil {
			_ = d.notify(*pending)
		}
	})
}

// Flush immediately fires any pending notification, cancelling the timer.
func (d *Debouncer) Flush() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}

	if d.pending == nil {
		return nil
	}

	pending := d.pending
	d.pending = nil
	return d.notify(*pending)
}
