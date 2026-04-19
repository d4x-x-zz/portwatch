// Package scanner provides TCP port scanning functionality for portwatch.
//
// Usage:
//
//	s := scanner.New("localhost", time.Second)
//	ports, err := s.Scan(1, 1024)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, p := range ports {
//		fmt.Printf("open: %d/%s\n", p.Number, p.Protocol)
//	}
//
// Scan performs a synchronous, sequential TCP dial for each port in the
// specified range. The Timeout field controls how long each dial attempt
// waits before moving on. For large ranges consider splitting the range
// across multiple goroutines in the calling code.
package scanner
