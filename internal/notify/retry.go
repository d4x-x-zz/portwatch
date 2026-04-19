package notify

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// Alerter is satisfied by any type that can send a Diff notification.
type Alerter interface {
	Notify(diff monitor.Diff) error
}

// RetryAlerter wraps an Alerter and retries on failure.
type RetryAlerter struct {
	inner    Alerter
	attempts int
	delay    time.Duration
}

// NewRetryAlerter returns an Alerter that retries up to attempts times,
// waiting delay between each try.
func NewRetryAlerter(inner Alerter, attempts int, delay time.Duration) *RetryAlerter {
	if attempts < 1 {
		attempts = 1
	}
	return &RetryAlerter{inner: inner, attempts: attempts, delay: delay}
}

// Notify calls the inner alerter, retrying on error.
func (r *RetryAlerter) Notify(diff monitor.Diff) error {
	var err error
	for i := 0; i < r.attempts; i++ {
		if err = r.inner.Notify(diff); err == nil {
			return nil
		}
		if i < r.attempts-1 {
			time.Sleep(r.delay)
		}
	}
	return fmt.Errorf("all %d attempts failed: %w", r.attempts, err)
}
