package notify

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// State represents the circuit breaker state.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker wraps an Alerter and stops calling it after repeated failures.
type CircuitBreaker struct {
	mu          sync.Mutex
	inner       Alerter
	failures    int
	threshold   int
	cooldown    time.Duration
	openedAt    time.Time
	state       State
}

type Alerter interface {
	Notify(diff monitor.Diff) error
}

// NewCircuitBreaker returns a CircuitBreaker that opens after threshold failures
// and retries after cooldown.
func NewCircuitBreaker(inner Alerter, threshold int, cooldown time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		inner:     inner,
		threshold: threshold,
		cooldown:  cooldown,
		state:     StateClosed,
	}
}

func (cb *CircuitBreaker) Notify(diff monitor.Diff) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateOpen:
		if time.Since(cb.openedAt) >= cb.cooldown {
			cb.state = StateHalfOpen
		} else {
			return fmt.Errorf("circuit open: alerter unavailable")
		}
	case StateClosed, StateHalfOpen:
		// proceed
	}

	err := cb.inner.Notify(diff)
	if err != nil {
		cb.failures++
		if cb.failures >= cb.threshold {
			cb.state = StateOpen
			cb.openedAt = time.Now()
		}
		return err
	}

	cb.failures = 0
	cb.state = StateClosed
	return nil
}

func (cb *CircuitBreaker) State() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}
