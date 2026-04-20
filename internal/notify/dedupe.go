package notify

import (
	"sync"

	"github.com/patrickward/portwatch/internal/monitor"
)

// DedupeAlerter wraps an Alerter and suppresses consecutive identical diffs.
// Unlike RateLimiter (which uses a cooldown window), DedupeAlerter simply
// skips a notification if the diff is byte-for-byte identical to the last
// one that was successfully forwarded.
type DedupeAlerter struct {
	mu   sync.Mutex
	next monitor.Alerter
	last string
}

// NewDedupeAlerter returns a DedupeAlerter wrapping next.
func NewDedupeAlerter(next monitor.Alerter) *DedupeAlerter {
	return &DedupeAlerter{next: next}
}

// Notify forwards diff to the wrapped Alerter only when it differs from the
// most recently forwarded diff. If the diff is empty (no changes) it is
// always forwarded so that downstream alerters can decide what to do.
func (d *DedupeAlerter) Notify(diff monitor.Diff) error {
	d.mu.Lock()
	key := diff.String()
	if key != "" && key == d.last {
		d.mu.Unlock()
		return nil
	}
	d.mu.Unlock()

	if err := d.next.Notify(diff); err != nil {
		return err
	}

	d.mu.Lock()
	d.last = key
	d.mu.Unlock()
	return nil
}

// Reset clears the remembered last diff, so the next notification will
// always be forwarded regardless of content.
func (d *DedupeAlerter) Reset() {
	d.mu.Lock()
	d.last = ""
	d.mu.Unlock()
}
