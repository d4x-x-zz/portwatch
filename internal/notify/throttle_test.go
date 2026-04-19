package notify

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

type captureAlerter struct {
	mu    sync.Mutex
	calls []monitor.Diff
	err   error
}

func (c *captureAlerter) Notify(d monitor.Diff) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.calls = append(c.calls, d)
	return c.err
}

func (c *captureAlerter) count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.calls)
}

func TestThrottledAlerter_AllowsFirst(t *testing.T) {
	cap := &captureAlerter{}
	th := NewThrottledAlerter(cap, 1*time.Second)
	d := monitor.Diff{Opened: monitor.NewPortSet([]int{80})}

	if err := th.Notify(d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cap.count() != 1 {
		t.Fatalf("expected 1 call, got %d", cap.count())
	}
}

func TestThrottledAlerter_SuppressesWithinWindow(t *testing.T) {
	cap := &captureAlerter{}
	th := NewThrottledAlerter(cap, 5*time.Second)
	d := monitor.Diff{Opened: monitor.NewPortSet([]int{443})}

	_ = th.Notify(d)
	_ = th.Notify(d)
	_ = th.Notify(d)

	if cap.count() != 1 {
		t.Fatalf("expected 1 call, got %d", cap.count())
	}
}

func TestThrottledAlerter_FlushSendsPending(t *testing.T) {
	cap := &captureAlerter{}
	th := NewThrottledAlerter(cap, 5*time.Second)
	d := monitor.Diff{Opened: monitor.NewPortSet([]int{8080})}

	_ = th.Notify(d) // fires
	_ = th.Notify(d) // queued

	if err := th.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}
	if cap.count() != 2 {
		t.Fatalf("expected 2 calls after flush, got %d", cap.count())
	}
}

func TestThrottledAlerter_FlushEmptyIsNoop(t *testing.T) {
	cap := &captureAlerter{}
	th := NewThrottledAlerter(cap, 1*time.Second)

	if err := th.Flush(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cap.count() != 0 {
		t.Fatalf("expected 0 calls, got %d", cap.count())
	}
}

func TestThrottledAlerter_PropagatesError(t *testing.T) {
	want := errors.New("backend down")
	cap := &captureAlerter{err: want}
	th := NewThrottledAlerter(cap, 1*time.Second)
	d := monitor.Diff{Opened: monitor.NewPortSet([]int{22})}

	if err := th.Notify(d); !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}
