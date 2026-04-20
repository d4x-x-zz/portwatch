package metrics

import (
	"testing"
	"time"
)

func makeRecord(ports int) ScanRecord {
	return ScanRecord{
		At:        time.Now(),
		PortsOpen: ports,
		Duration:  10 * time.Millisecond,
	}
}

func TestCollector_RecordAndLatest(t *testing.T) {
	c := NewCollector(10)
	if _, ok := c.Latest(); ok {
		t.Fatal("expected empty collector")
	}
	c.Record(makeRecord(5))
	c.Record(makeRecord(7))
	r, ok := c.Latest()
	if !ok {
		t.Fatal("expected a record")
	}
	if r.PortsOpen != 7 {
		t.Fatalf("want 7 open ports, got %d", r.PortsOpen)
	}
}

func TestCollector_EvictsOldest(t *testing.T) {
	c := NewCollector(3)
	for i := 1; i <= 4; i++ {
		c.Record(makeRecord(i))
	}
	if c.Len() != 3 {
		t.Fatalf("want 3 records, got %d", c.Len())
	}
	all := c.All()
	if all[0].PortsOpen != 2 {
		t.Fatalf("want oldest=2, got %d", all[0].PortsOpen)
	}
}

func TestCollector_AllReturnsCopy(t *testing.T) {
	c := NewCollector(10)
	c.Record(makeRecord(3))
	all := c.All()
	all[0].PortsOpen = 999
	r, _ := c.Latest()
	if r.PortsOpen == 999 {
		t.Fatal("All() should return a copy, not a reference")
	}
}

func TestCollector_DefaultMaxSize(t *testing.T) {
	c := NewCollector(0)
	if c.maxSize != 100 {
		t.Fatalf("expected default maxSize 100, got %d", c.maxSize)
	}
}
