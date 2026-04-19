package alert

// MultiAlerter fans out notifications to multiple Alerter implementations.
type MultiAlerter struct {
	alerters []Alerter
}

// NewMultiAlerter creates a MultiAlerter that notifies all provided alerters.
func NewMultiAlerter(alerters ...Alerter) *MultiAlerter {
	return &MultiAlerter{alerters: alerters}
}

// Notify calls Notify on every registered alerter, collecting any errors.
// It continues even if one alerter fails and returns the last non-nil error.
func (m *MultiAlerter) Notify(diff interface{ String() string }) error {
	var lastErr error
	for _, a := range m.alerters {
		if err := a.Notify(diff); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// Add appends an alerter to the list at runtime.
func (m *MultiAlerter) Add(a Alerter) {
	m.alerters = append(m.alerters, a)
}

// Len returns the number of registered alerters.
func (m *MultiAlerter) Len() int {
	return len(m.alerters)
}
