package notify

import (
	"math"
	"time"
)

// BackoffStrategy computes the delay before the nth retry attempt.
type BackoffStrategy interface {
	Delay(attempt int) time.Duration
}

// ExponentialBackoff implements exponential back-off with optional jitter cap.
type ExponentialBackoff struct {
	Base    time.Duration
	Max     time.Duration
	Factor  float64
}

// DefaultExponentialBackoff returns a sensible default strategy.
func DefaultExponentialBackoff() ExponentialBackoff {
	return ExponentialBackoff{
		Base:   500 * time.Millisecond,
		Max:    30 * time.Second,
		Factor: 2.0,
	}
}

// Delay returns the capped exponential delay for the given attempt (0-indexed).
func (e ExponentialBackoff) Delay(attempt int) time.Duration {
	if attempt <= 0 {
		return e.Base
	}
	mult := math.Pow(e.Factor, float64(attempt))
	d := time.Duration(float64(e.Base) * mult)
	if d > e.Max {
		d = e.Max
	}
	return d
}

// ConstantBackoff always returns the same delay regardless of attempt.
type ConstantBackoff struct {
	Interval time.Duration
}

// Delay returns the fixed interval.
func (c ConstantBackoff) Delay(_ int) time.Duration {
	return c.Interval
}
