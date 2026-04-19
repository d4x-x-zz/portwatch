package notify

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

func makeDiff(opened, closed []int) monitor.Diff {
	ops := monitor.NewPortSet(opened...)
	cls := monitor.NewPortSet(closed...)
	return monitor.Diff{Opened: ops, Closed: cls}
}

func TestDebouncer_FlushSendsLatest(t *testing.T) {
	var calls int32
	var last monitor.Diff

	d := NewDebouncer(200*time.Millisecond, func(diff monitor.Diff) error {
		atomic.AddInt32(&calls, 1)
		last = diff
		return nil
	})

	d.Submit(makeDiff([]int{80}, nil))
	d.Submit(makeDiff([]int{443}, nil))

	if err := d.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}

	if atomic.LoadInt32(&calls) != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}

	if !last.Opened.Contains(443) {
		t.Error("expected latest diff (port 443) to be forwarded")
	}
}

func TestDebouncer_FiresAfterQuietPeriod(t *testing.T) {
	var calls int32

	d := NewDebouncer(50*time.Millisecond, func(diff monitor.Diff) error {
		atomic.AddInt32(&calls, 1)
		return nil
	})

	d.Submit(makeDiff([]int{8080}, nil))
	time.Sleep(120 * time.Millisecond)

	if atomic.LoadInt32(&calls) != 1 {
		t.Fatalf("expected 1 call after quiet period, got %d", calls)
	}
}

func TestDebouncer_FlushOnEmpty(t *testing.T) {
	d := NewDebouncer(100*time.Millisecond, func(diff monitor.Diff) error {
		t.Error("notify should not be called when nothing pending")
		return nil
	})

	if err := d.Flush(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
