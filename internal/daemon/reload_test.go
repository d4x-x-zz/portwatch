package daemon

import (
	"context"
	"errors"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestWatchReload_CancelFromContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	called := 0
	done := make(chan error, 1)
	go func() {
		done <- WatchReload(ctx, func() error {
			called++
			return nil
		})
	}()

	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("WatchReload did not return after context cancel")
	}

	if called != 0 {
		t.Errorf("expected fn not to be called, got %d", called)
	}
}

func TestWatchReload_CallsFnOnSIGHUP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	called := make(chan struct{}, 1)
	done := make(chan error, 1)
	go func() {
		done <- WatchReload(ctx, func() error {
			called <- struct{}{}
			return nil
		})
	}()

	time.Sleep(20 * time.Millisecond)
	p, _ := os.FindProcess(os.Getpid())
	p.Signal(syscall.SIGHUP)

	select {
	case <-called:
		// success
	case <-time.After(time.Second):
		t.Fatal("fn was not called after SIGHUP")
	}

	cancel()
	<-done
}

func TestWatchReload_StopsOnFnError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	expected := errors.New("reload failed")
	done := make(chan error, 1)
	go func() {
		done <- WatchReload(ctx, func() error {
			return expected
		})
	}()

	time.Sleep(20 * time.Millisecond)
	p, _ := os.FindProcess(os.Getpid())
	p.Signal(syscall.SIGHUP)

	select {
	case err := <-done:
		if !errors.Is(err, expected) {
			t.Fatalf("expected %v, got %v", expected, err)
		}
	case <-time.After(time.Second):
		t.Fatal("WatchReload did not return after fn error")
	}
}
