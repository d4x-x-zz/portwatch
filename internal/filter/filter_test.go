package filter

import (
	"testing"
)

func TestFilter_NoRules_AllowsAll(t *testing.T) {
	f, err := New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, p := range []int{1, 80, 443, 65535} {
		if !f.Allow(p) {
			t.Errorf("expected port %d to be allowed", p)
		}
	}
}

func TestFilter_AllowSingle(t *testing.T) {
	f, err := New([]string{"allow:80", "deny:1-65535"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Allow(80) {
		t.Error("expected port 80 to be allowed")
	}
	if f.Allow(443) {
		t.Error("expected port 443 to be denied")
	}
}

func TestFilter_DenyRange(t *testing.T) {
	f, err := New([]string{"deny:8000-9000"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Allow(8080) {
		t.Error("expected 8080 to be denied")
	}
	if !f.Allow(80) {
		t.Error("expected 80 to be allowed (no matching rule)")
	}
}

func TestFilter_InvalidRule(t *testing.T) {
	cases := []string{
		"badformat",
		"unknown:80",
		"allow:0",
		"deny:70000",
		"allow:900-80",
	}
	for _, c := range cases {
		_, err := New([]string{c})
		if err == nil {
			t.Errorf("expected error for rule %q", c)
		}
	}
}

func TestParseRule_Range(t *testing.T) {
	r, err := parseRule("allow:1024-2048")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Start != 1024 || r.End != 2048 || !r.Allow {
		t.Errorf("unexpected rule: %+v", r)
	}
}
