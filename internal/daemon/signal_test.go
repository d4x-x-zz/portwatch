package daemon

import (
	"context"
	"syscall"
	"testing"
	"time"
)

func TestWaitForSignalCtx_CancelFromContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// wrap cancel so WaitForSignalCtx gets its own cancel
	inner, innerCancel := context.WithCancel(ctx)
	defer innerCancel()

	done := make(chan struct{})
	var got interface{}
	go func() {
		got = WaitForSignalCtx(inner, innerCancel)
		close(done)
	}()

	// let the context time out naturally
	select {
	case <-done:
		if got != nil {
			t.Errorf("expected nil signal, got %v", got)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for WaitForSignalCtx to return")
	}
}

func TestWaitForSignalCtx_SignalCancels(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	var got interface{}
	go func() {
		got = WaitForSignalCtx(ctx, cancel)
		close(done)
	}()

	// give the goroutine time to register
	time.Sleep(20 * time.Millisecond)
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM) //nolint:errcheck

	select {
	case <-done:
		if got == nil {
			t.Error("expected a signal, got nil")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for signal delivery")
	}
}
