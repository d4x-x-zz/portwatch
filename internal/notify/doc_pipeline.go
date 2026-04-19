// Package notify provides alerter middleware and pipeline utilities.
//
// A Pipeline chains middlewares (rate limiting, throttling, retry, circuit
// breaking, debouncing) around a base Alerter. Build one with NewPipeline:
//
//	pipeline := notify.NewPipeline(
//		baseAlerter,
//		func(a notify.Alerter) notify.Alerter { return notify.NewRateLimiter(a, cooldown) },
//		func(a notify.Alerter) notify.Alerter { return notify.NewRetryAlerter(a, cfg) },
//	)
//
// Middlewares are applied outermost-first, so the first middleware in the
// slice is the first to handle each Notify call.
package notify
