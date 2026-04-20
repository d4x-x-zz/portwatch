package daemon

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// WaitForSignal blocks until SIGINT or SIGTERM is received, then cancels the
// provided context. It returns the signal that triggered the shutdown.
func WaitForSignal(cancel context.CancelFunc) os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)

	sig := <-ch
	cancel()
	return sig
}

// WaitForSignalCtx blocks until a signal is received or the context is done.
// Returns the signal if one was received, or nil if the context expired first.
func WaitForSignalCtx(ctx context.Context, cancel context.CancelFunc) os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)

	select {
	case sig := <-ch:
		cancel()
		return sig
	case <-ctx.Done():
		return nil
	}
}
