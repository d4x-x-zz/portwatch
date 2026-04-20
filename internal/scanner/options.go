package scanner

import "time"

// Options configures a Scanner.
type Options struct {
	// Host is the target host.
	Host string
	// Start is the first port in the scan range (inclusive).
	Start int
	// End is the last port in the scan range (inclusive).
	End int
	// Concurrency controls how many ports are dialled in parallel.
	Concurrency int
	// Timeout is the per-port dial timeout.
	Timeout time.Duration
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		Host:        "127.0.0.1",
		Start:       1,
		End:         1024,
		Concurrency: 100,
		Timeout:     500 * time.Millisecond,
	}
}

// Validate returns an error if the options are not usable.
func (o Options) Validate() error {
	if o.Host == "" {
		return errorf("host must not be empty")
	}
	if o.Start < 1 || o.Start > 65535 {
		return errorf("start port %d out of range", o.Start)
	}
	if o.End < o.Start || o.End > 65535 {
		return errorf("end port %d out of range or less than start", o.End)
	}
	if o.Concurrency <= 0 {
		return errorf("concurrency must be greater than zero")
	}
	if o.Timeout <= 0 {
		return errorf("timeout must be greater than zero")
	}
	return nil
}

func errorf(format string, args ...interface{}) error {
	return fmt.Errorf("scanner: "+format, args...)
}
