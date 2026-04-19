package notify_test

import (
	"context"
	"errors"
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/notify"
)

type recordingAlerter struct {
	calls []monitor.Diff
	err   error
}

func (r *recordingAlerter) Notify(_ context.Context, d monitor.Diff) error {
	r.calls = append(r.calls, d)
	return r.err
}

func wrapTag(tag string, next notify.Alerter) notify.Alerter {
	return &taggingAlerter{tag: tag, next: next}
}

type taggingAlerter struct {
	tag  string
	next notify.Alerter
	order *[]string
}

func (t *taggingAlerter) Notify(ctx context.Context, d monitor.Diff) error {
	return t.next.Notify(ctx, d)
}

func TestPipeline_CallsBase(t *testing.T) {
	base := &recordingAlerter{}
	p := notify.NewPipeline(base)
	d := monitor.Diff{Opened: []int{8080}}
	if err := p.Notify(context.Background(), d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(base.calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(base.calls))
	}
}

func TestPipeline_PropagatesError(t *testing.T) {
	expected := errors.New("boom")
	base := &recordingAlerter{err: expected}
	p := notify.NewPipeline(base)
	err := p.Notify(context.Background(), monitor.Diff{})
	if !errors.Is(err, expected) {
		t.Fatalf("expected %v, got %v", expected, err)
	}
}

func TestPipeline_WrapsMiddlewares(t *testing.T) {
	base := &recordingAlerter{}
	// middleware that just passes through
	pass := func(next notify.Alerter) notify.Alerter { return next }
	p := notify.NewPipeline(base, pass, pass)
	if err := p.Notify(context.Background(), monitor.Diff{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(base.calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(base.calls))
	}
}
