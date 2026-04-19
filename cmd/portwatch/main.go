package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
)

func main() {
	cfgPath := flag.String("config", "portwatch.toml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Printf("could not load config, using defaults: %v", err)
		cfg = config.Default()
	}

	sc := scanner.New(cfg.Timeout)
	alerter := alert.NewLogAlerter(os.Stdout)

	log.Printf("starting portwatch: range=%d-%d interval=%s", cfg.StartPort, cfg.EndPort, cfg.Interval)

	initial, err := sc.Scan(cfg.StartPort, cfg.EndPort)
	if err != nil {
		log.Fatalf("initial scan failed: %v", err)
	}
	prev := monitor.NewPortSet(initial)
	log.Printf("initial scan complete: %d open ports", len(initial)icker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			current, err := sc.Scan(cfg.StartPort, cfg.EndPort)
			if err != nil {
				log.Printf("scan error: %v", err)
				continue
			}
			curSet := monitor.NewPortSet(current)
			diff := monitor.Compare(prev, curSet)
			if diff.HasChanges() {
				if err := alerter.Notify(diff); err != nil {
					log.Printf("alert error: %v", err)
				}
			}
			prev = curSet
		 case <-sigs:
			log.Println("shutting down")
			return
		}
	}
}
