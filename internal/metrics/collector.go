package metrics

import (
	"sync"
	"time"
)

// ScanRecord holds the result of a single scan cycle.
type ScanRecord struct {
	At        time.Time
	PortsOpen int
	Duration  time.Duration
}

// Collector keeps a rolling window of recent scan records.
type Collector struct {
	mu      sync.Mutex
	records []ScanRecord
	maxSize int
}

// NewCollector returns a Collector that retains up to maxSize records.
func NewCollector(maxSize int) *Collector {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &Collector{maxSize: maxSize}
}

// Record appends a new scan result, evicting the oldest if at capacity.
func (c *Collector) Record(r ScanRecord) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.records) >= c.maxSize {
		c.records = c.records[1:]
	}
	c.records = append(c.records, r)
}

// Latest returns the most recent ScanRecord and true, or false if empty.
func (c *Collector) Latest() (ScanRecord, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.records) == 0 {
		return ScanRecord{}, false
	}
	return c.records[len(c.records)-1], true
}

// All returns a copy of all stored records.
func (c *Collector) All() []ScanRecord {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]ScanRecord, len(c.records))
	copy(out, c.records)
	return out
}

// Len returns the current number of stored records.
func (c *Collector) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.records)
}
