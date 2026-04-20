package daemon

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// ReloadFunc is called when a reload signal (SIGHUP) is received.
type ReloadFunc func() error

// WatchReload blocks until SIGHUP is received or ctx is cancelled.
// On SIGHUP it calls fn and continues watching. On ctx cancellation it returns.
func WatchReload(ctx context.Context, fn ReloadFunc) error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP)
	defer signal.Stop(ch)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ch:
			if err := fn(); err != nil {
				return err
			}
		}
	}
}

// ReloadChannel returns a channel that emits on SIGHUP.
// The caller is responsible for calling the returned stop function.
func ReloadChannel() (ch <-chan os.Signal, stop func()) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP)
	return sig, func() { signal.Stop(sig) }
}
