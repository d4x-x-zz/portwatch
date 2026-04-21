package notify

import (
	"sync"

	"github.com/user/portwatch/internal/monitor"
)

// SuppressAlerter wraps an Alerter and suppresses repeated identical diffs
// after a configurable number of consecutive occurrences.
type SuppressAlerter struct {
	mu           sync.Mutex
	next         Alerter
	maxRepeats   int
	resetOnChange bool
	last         *monitor.Diff
	count        int
}

// NewSuppressAlerter returns a SuppressAlerter that forwards to next.
// After maxRepeats consecutive identical diffs the alert is suppressed.
// If resetOnChange is true the repeat counter resets when a different diff
// is seen.
func NewSuppressAlerter(next Alerter, maxRepeats int, resetOnChange bool) *SuppressAlerter {
	return &SuppressAlerter{
		next:          next,
		maxRepeats:    maxRepeats,
		resetOnChange: resetOnChange,
	}
}

// Notify forwards the diff unless it has been seen maxRepeats times in a row.
func (s *SuppressAlerter) Notify(d monitor.Diff) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.last != nil && diffEqual(d, *s.last) {
		s.count++
	} else {
		if s.resetOnChange {
			s.count = 1
		} else {
			s.count++
		}
		copy := d
		s.last = &copy
	}

	if s.last == nil {
		copy := d
		s.last = &copy
		s.count = 1
	}

	if s.count > s.maxRepeats {
		return nil
	}

	return s.next.Notify(d)
}

// Reset clears the internal state, allowing the next diff to pass through.
func (s *SuppressAlerter) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.last = nil
	s.count = 0
}
