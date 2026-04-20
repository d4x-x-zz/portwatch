package notify

import (
	"testing"
	"time"
)

func TestExponentialBackoff_FirstAttempt(t *testing.T) {
	b := DefaultExponentialBackoff()
	got := b.Delay(0)
	if got != b.Base {
		t.Fatalf("expected %v, got %v", b.Base, got)
	}
}

func TestExponentialBackoff_Grows(t *testing.T) {
	b := ExponentialBackoff{
		Base:   100 * time.Millisecond,
		Max:    10 * time.Second,
		Factor: 2.0,
	}
	d0 := b.Delay(0)
	d1 := b.Delay(1)
	d2 := b.Delay(2)
	if d1 <= d0 {
		t.Fatalf("expected delay to grow: d0=%v d1=%v", d0, d1)
	}
	if d2 <= d1 {
		t.Fatalf("expected delay to grow: d1=%v d2=%v", d1, d2)
	}
}

func TestExponentialBackoff_CapsAtMax(t *testing.T) {
	b := ExponentialBackoff{
		Base:   1 * time.Second,
		Max:    2 * time.Second,
		Factor: 10.0,
	}
	for attempt := 0; attempt < 10; attempt++ {
		d := b.Delay(attempt)
		if d > b.Max {
			t.Fatalf("attempt %d: delay %v exceeds max %v", attempt, d, b.Max)
		}
	}
}

func TestConstantBackoff_AlwaysSame(t *testing.T) {
	c := ConstantBackoff{Interval: 5 * time.Second}
	for attempt := 0; attempt < 5; attempt++ {
		d := c.Delay(attempt)
		if d != 5*time.Second {
			t.Fatalf("attempt %d: expected 5s, got %v", attempt, d)
		}
	}
}

func TestDefaultExponentialBackoff_Sensible(t *testing.T) {
	b := DefaultExponentialBackoff()
	if b.Base <= 0 {
		t.Fatal("base must be positive")
	}
	if b.Max < b.Base {
		t.Fatal("max must be >= base")
	}
	if b.Factor <= 1.0 {
		t.Fatal("factor must be > 1")
	}
}
