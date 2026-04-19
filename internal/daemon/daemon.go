// Package daemon wires together scanning, monitoring, alerting, and history
// into a single run loop.
package daemon

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

// Daemon runs the port-watch loop.
type Daemon struct {
	cfg     *config.Config
	store   *history.Store
	alerter alert.Alerter
}

// New creates a Daemon from the given config, history store, and alerter.
func New(cfg *config.Config, store *history.Store, alerter alert.Alerter) *Daemon {
	return &Daemon{cfg: cfg, store: store, alerter: alerter}
}

// Run starts the monitoring loop and blocks until ctx is cancelled.
func (d *Daemon) Run(ctx context.Context) error {
	s := scanner.New(d.cfg.Timeout)
	ticker := time.NewTicker(d.cfg.Interval)
	defer ticker.Stop()

	log.Printf("portwatch starting — range %d-%d, interval %s",
		d.cfg.StartPort, d.cfg.EndPort, d.cfg.Interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("portwatch shutting down")
			return ctx.Err()
		case <-ticker.C:
			d.tick(s)
		}
	}
}

func (d *Daemon) tick(s *scanner.Scanner) {
	ports, err := s.Scan(d.cfg.StartPort, d.cfg.EndPort)
	if err != nil {
		log.Printf("scan error: %v", err)
		return
	}

	current := monitor.NewPortSet(ports)

	previous, err := d.store.Load()
	if err != nil {
		log.Printf("history load error (treating as empty): %v", err)
		previous = monitor.NewPortSet(nil)
	}

	diff := monitor.Compare(previous, current)

	if err := d.alerter.Notify(diff); err != nil {
		log.Printf("alert error: %v", err)
	}

	if err := d.store.Save(current); err != nil {
		log.Printf("history save error: %v", err)
	}
}
