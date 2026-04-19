package metrics

import (
	"testing"
	"time"
)

func TestNew_InitialisesUpSince(t *testing.T) {
	before := time.Now()
	c := New()
	after := time.Now()
	s := c.Snapshot()
	if s.UpSince.Before(before) || s.UpSince.After(after) {
		t.Errorf("unexpected UpSince: %v", s.UpSince)
	}
}

func TestRecordScan(t *testing.T) {
	c := New()
	c.RecordScan(42)
	c.RecordScan(7)
	s := c.Snapshot()
	if s.ScansTotal != 2 {
		t.Errorf("expected 2 scans, got %d", s.ScansTotal)
	}
	if s.LastScanPorts != 7 {
		t.Errorf("expected last scan ports 7, got %d", s.LastScanPorts)
	}
	if s.LastScanAt.IsZero() {
		t.Error("LastScanAt should not be zero")
	}
}

func TestRecordAlert(t *testing.T) {
	c := New()
	c.RecordAlert()
	c.RecordAlert()
	c.RecordAlert()
	s := c.Snapshot()
	if s.AlertsTotal != 3 {
		t.Errorf("expected 3 alerts, got %d", s.AlertsTotal)
	}
}

func TestSnapshot_IsCopy(t *testing.T) {
	c := New()
	c.RecordScan(10)
	s1 := c.Snapshot()
	c.RecordScan(20)
	s2 := c.Snapshot()
	if s1.ScansTotal == s2.ScansTotal {
		t.Error("snapshots should be independent")
	}
}
