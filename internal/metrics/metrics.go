// Package metrics tracks runtime statistics for portwatch.
package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time view of daemon metrics.
type Snapshot struct {
	ScansTotal    int64
	AlertsTotal   int64
	LastScanAt    time.Time
	LastScanPorts int
	UpSince       time.Time
}

// Collector accumulates runtime metrics.
type Collector struct {
	mu            sync.RWMutex
	scansTotal    int64
	alertsTotal   int64
	lastScanAt    time.Time
	lastScanPorts int
	upSince       time.Time
}

// New returns a new Collector initialised with the current time.
func New() *Collector {
	return &Collector{upSince: time.Now()}
}

// RecordScan records a completed scan result.
func (c *Collector) RecordScan(portCount int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.scansTotal++
	c.lastScanAt = time.Now()
	c.lastScanPorts = portCount
}

// RecordAlert increments the alert counter.
func (c *Collector) RecordAlert() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.alertsTotal++
}

// Snapshot returns a copy of the current metrics.
func (c *Collector) Snapshot() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Snapshot{
		ScansTotal:    c.scansTotal,
		AlertsTotal:   c.alertsTotal,
		LastScanAt:    c.lastScanAt,
		LastScanPorts: c.lastScanPorts,
		UpSince:       c.upSince,
	}
}
