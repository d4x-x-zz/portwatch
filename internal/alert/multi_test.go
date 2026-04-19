package alert

import (
	"errors"
	"testing"
)

type stubAlerter struct {
	called bool
	err    error
}

func (s *stubAlerter) Notify(diff interface{ String() string }) error {
	s.called = true
	return s.err
}

type stubDiff struct{ msg string }

func (d stubDiff) String() string { return d.msg }

func TestMultiAlerter_NotifiesAll(t *testing.T) {
	a1 := &stubAlerter{}
	a2 := &stubAlerter{}
	m := NewMultiAlerter(a1, a2)

	if err := m.Notify(stubDiff{"changes"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !a1.called || !a2.called {
		t.Error("expected both alerters to be called")
	}
}

func TestMultiAlerter_ReturnsLastError(t *testing.T) {
	errBoom := errors.New("boom")
	a1 := &stubAlerter{err: errors.New("first")}
	a2 := &stubAlerter{err: errBoom}
	m := NewMultiAlerter(a1, a2)

	err := m.Notify(stubDiff{})
	if err != errBoom {
		t.Fatalf("expected last error, got %v", err)
	}
	if !a1.called {
		t.Error("first alerter should still be called despite error")
	}
}

func TestMultiAlerter_Add(t *testing.T) {
	m := NewMultiAlerter()
	if m.Len() != 0 {
		t.Fatal("expected empty")
	}
	m.Add(&stubAlerter{})
	if m.Len() != 1 {
		t.Fatal("expected one alerter")
	}
}
