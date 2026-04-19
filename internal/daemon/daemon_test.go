package daemon_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/daemon"
	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/monitor"
)

type captureAlerter struct {
	diffs []monitor.Diff
}

func (c *captureAlerter) Notify(d monitor.Diff) error {
	c.diffs = append(c.diffs, d)
	return nil
}

func TestDaemon_RunTicksAndCancels(t *testing.T) {
	cfg := config.Default()
	cfg.Interval = 50 * time.Millisecond
	cfg.StartPort = 1
	cfg.EndPort = 1 // tiny range, likely nothing open

	tmpDir := t.TempDir()
	store, err := history.NewStore(filepath.Join(tmpDir, "state.json"))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	cap := &captureAlerter{}
	d := daemon.New(cfg, store, cap)

	ctx, cancel := context.WithTimeout(context.Background(), 160*time.Millisecond)
	defer cancel()

	err = d.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}

	// At least 2 ticks should have fired in 160 ms with a 50 ms interval.
	if len(cap.diffs) < 2 {
		t.Errorf("expected >=2 ticks, got %d", len(cap.diffs))
	}
}

func TestDaemon_PersistsState(t *testing.T) {
	cfg := config.Default()
	cfg.Interval = 40 * time.Millisecond
	cfg.StartPort = 1
	cfg.EndPort = 1

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "state.json")
	store, _ := history.NewStore(path)

	d := daemon.New(cfg, store, &captureAlerter{})
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Millisecond)
	defer cancel()
	d.Run(ctx) //nolint:errcheck

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected state file to exist after run")
	}
}

// Ensure alert.NewLogAlerter satisfies the Alerter interface used by daemon.
func TestDaemon_AcceptsLogAlerter(t *testing.T) {
	cfg := config.Default()
	cfg.Interval = 30 * time.Millisecond
	cfg.StartPort = 1
	cfg.EndPort = 1

	store, _ := history.NewStore(filepath.Join(t.TempDir(), "s.json"))
	a := alert.NewLogAlerter(nil)
	d := daemon.New(cfg, store, a)

	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()
	d.Run(ctx) //nolint:errcheck
}
