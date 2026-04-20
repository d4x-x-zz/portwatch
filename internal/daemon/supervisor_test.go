package daemon

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

var errBoom = errors.New("boom")

func TestSupervise_FnSucceeds(t *testing.T) {
	ctx := context.Background()
	opts := DefaultSupervisorOptions()
	opts.RestartDelay = time.Millisecond

	calls := 0
	err := Supervise(ctx, opts, func(_ context.Context) error {
		calls++
		return nil
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestSupervise_RetriesOnFailure(t *testing.T) {
	ctx := context.Background()
	opts := DefaultSupervisorOptions()
	opts.MaxRestarts = 3
	opts.RestartDelay = time.Millisecond

	var calls atomic.Int32
	err := Supervise(ctx, opts, func(_ context.Context) error {
		calls.Add(1)
		return errBoom
	})

	if err == nil {
		t.Fatal("expected error after max restarts")
	}
	if calls.Load() != 3 {
		t.Fatalf("expected 3 calls, got %d", calls.Load())
	}
}

func TestSupervise_CancelStops(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	opts := DefaultSupervisorOptions()
	opts.MaxRestarts = 0 // unlimited
	opts.RestartDelay = 5 * time.Millisecond

	var calls atomic.Int32
	done := make(chan error, 1)
	go func() {
		done <- Supervise(ctx, opts, func(_ context.Context) error {
			if calls.Add(1) >= 2 {
				cancel()
			}
			return errBoom
		})
	}()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("supervisor did not stop after context cancel")
	}
}

func TestDefaultSupervisorOptions_Sensible(t *testing.T) {
	opts := DefaultSupervisorOptions()
	if opts.MaxRestarts <= 0 {
		t.Error("expected positive MaxRestarts")
	}
	if opts.RestartDelay <= 0 {
		t.Error("expected positive RestartDelay")
	}
}
