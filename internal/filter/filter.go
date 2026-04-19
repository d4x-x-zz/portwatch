// Package filter provides port filtering based on allow/deny rules.
package filter

import "fmt"

// Rule represents a single port filter rule.
type Rule struct {
	Allow bool
	Start int
	End   int
}

// Filter decides which ports should trigger alerts.
type Filter struct {
	rules []Rule
}

// New creates a Filter from a slice of rule strings.
// Format: "allow:80", "deny:8000-9000"
func New(specs []string) (*Filter, error) {
	f := &Filter{}
	for _, s := range specs {
		r, err := parseRule(s)
		if err != nil {
			return nil, err
		}
		f.rules = append(f.rules, r)
	}
	return f, nil
}

// Allow returns true if the port should be included in alerting.
// If no rules are defined, all ports are allowed.
func (f *Filter) Allow(port int) bool {
	if len(f.rules) == 0 {
		return true
	}
	for _, r := range f.rules {
		if port >= r.Start && port <= r.End {
			return r.Allow
		}
	}
	return true
}

func parseRule(s string) (Rule, error) {
	var kind string
	var rangeStr string
	if _, err := fmt.Sscanf(s, "%3s:%s", &kind, &rangeStr); err != nil {
		// try longer keyword
		if _, err2 := fmt.Sscanf(s, "%4s:%s", &kind, &rangeStr); err2 != nil {
			return Rule{}, fmt.Errorf("filter: invalid rule %q", s)
		}
	}
	allow := kind == "allow"
	if kind != "allow" && kind != "deny" {
		return Rule{}, fmt.Errorf("filter: unknown rule type %q", kind)
	}
	var start, end int
	if _, err := fmt.Sscanf(rangeStr, "%d-%d", &start, &end); err != nil {
		if _, err2 := fmt.Sscanf(rangeStr, "%d", &start); err2 != nil {
			return Rule{}, fmt.Errorf("filter: invalid port range %q", rangeStr)
		}
		end = start
	}
	if start < 1 || end > 65535 || start > end {
		return Rule{}, fmt.Errorf("filter: port range %d-%d out of bounds", start, end)
	}
	return Rule{Allow: allow, Start: start, End: end}, nil
}
