package notify

import (
	"context"

	"github.com/user/portwatch/internal/monitor"
)

// Alerter is the interface for sending diff notifications.
type Alerter interface {
	Notify(ctx context.Context, diff monitor.Diff) error
}

// Pipeline chains multiple alerter middlewares around a base alerter.
// Middlewares are applied in order, outermost first.
type Pipeline struct {
	alerter Alerter
}

// NewPipeline constructs a Pipeline wrapping base with the given middlewares.
// Each middleware wraps the previous result, so middlewares[0] is outermost.
func NewPipeline(base Alerter, middlewares ...func(Alerter) Alerter) *Pipeline {
	a := base
	for i := len(middlewares) - 1; i >= 0; i-- {
		a = middlewares[i](a)
	}
	return &Pipeline{alerter: a}
}

// Notify dispatches the diff through the pipeline.
func (p *Pipeline) Notify(ctx context.Context, diff monitor.Diff) error {
	return p.alerter.Notify(ctx, diff)
}
