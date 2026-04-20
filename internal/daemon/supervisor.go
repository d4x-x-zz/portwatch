package daemon

import (
	"context"
	"log"
	"time"
)

// SupervisorOptions controls restart behaviour for the supervised function.
type SupervisorOptions struct {
	// MaxRestarts is the maximum number of consecutive restarts before giving up.
	// Zero means unlimited.
	MaxRestarts int
	// RestartDelay is the wait between restart attempts.
	RestartDelay time.Duration
	// Logger receives restart notifications. Defaults to the standard logger.
	Logger *log.Logger
}

// DefaultSupervisorOptions returns sensible defaults.
func DefaultSupervisorOptions() SupervisorOptions {
	return SupervisorOptions{
		MaxRestarts:  5,
		RestartDelay: 2 * time.Second,
	}
}

// Supervise runs fn repeatedly until ctx is cancelled, fn returns nil, or
// MaxRestarts consecutive failures occur. Each failure is logged and the
// supervisor waits RestartDelay before the next attempt.
func Supervise(ctx context.Context, opts SupervisorOptions, fn func(ctx context.Context) error) error {
	logger := opts.Logger
	if logger == nil {
		logger = log.Default()
	}

	attempts := 0
	for {
		err := fn(ctx)
		if err == nil {
			return nil
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		attempts++
		logger.Printf("supervisor: fn exited with error (attempt %d): %v", attempts, err)

		if opts.MaxRestarts > 0 && attempts >= opts.MaxRestarts {
			logger.Printf("supervisor: max restarts (%d) reached, giving up", opts.MaxRestarts)
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(opts.RestartDelay):
		}
	}
}
