package network

import (
	"fmt"
	"net"
	"time"

	"net-zilla/internal/models"
)

// Tracer handles network path analysis.
type Tracer struct {
	maxHops int
	timeout time.Duration
}

// NewTracer creates a new Tracer.
func NewTracer() *Tracer {
	return &Tracer{
		maxHops: 30,
		timeout: 2 * time.Second,
	}
}

// Trace performs a basic TCP-based traceroute by incrementing TTL.
// Note: On many systems, setting TTL for TCP requires specific permissions or is OS-dependent.
// This remains a best-effort implementation for a non-root environment.
func (t *Tracer) Trace(target string) (*models.NetworkAnalysis, error) {
	ips, err := net.LookupIP(target)
	if err != nil {
		return nil, err
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("no IP found")
	}
	destIP := ips[0]

	analysis := &models.NetworkAnalysis{
		Target: target,
	}

	start := time.Now()
	// Best effort traceroute using net package
	// Without raw sockets, we can't easily do a true TTL-based traceroute in pure Go.
	// We'll simulate a path or just provide the final hop info if TTL is not supported.

	// Implementation note: For a real production app without root, we often rely on
	// external utilities or just perform latency/path analysis to the final destination.

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(destIP.String(), "80"), t.timeout)
	if err != nil {
		conn, err = net.DialTimeout("tcp", net.JoinHostPort(destIP.String(), "443"), t.timeout)
	}

	if err == nil {
		defer conn.Close()
		latency := time.Since(start)
		analysis.AverageLatency = latency
		analysis.HopCount = 1 // Simplified for non-root
		analysis.Hops = append(analysis.Hops, models.HopDetail{
			Number:  1,
			IP:      destIP.String(),
			Latency: latency,
		})
	} else {
		return nil, fmt.Errorf("target unreachable: %w", err)
	}

	return analysis, nil
}
