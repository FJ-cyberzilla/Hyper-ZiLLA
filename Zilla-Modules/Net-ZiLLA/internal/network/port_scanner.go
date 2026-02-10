package network

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"net-zilla/pkg/logger"
)

type PortScanner struct {
	logger  *logger.Logger
	timeout time.Duration
	workers int
}

func NewPortScanner(logger *logger.Logger) *PortScanner {
	return &PortScanner{
		logger:  logger,
		timeout: 1 * time.Second,
		workers: 100,
	}
}

type ScanResult struct {
	Port    int
	Open    bool
	Service string
}

func (ps *PortScanner) Scan(ctx context.Context, target string, ports []int) []ScanResult {
	var results []ScanResult
	var wg sync.WaitGroup

	portsChan := make(chan int, ps.workers)
	resultsChan := make(chan ScanResult, len(ports))

	// Start workers
	for i := 0; i < ps.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range portsChan {
				address := net.JoinHostPort(target, fmt.Sprintf("%d", port))
				conn, err := net.DialTimeout("tcp", address, ps.timeout)
				if err == nil {
					conn.Close()
					resultsChan <- ScanResult{Port: port, Open: true}
				}
			}
		}()
	}

	// Send ports to scan
	go func() {
		for _, port := range ports {
			portsChan <- port
		}
		close(portsChan)
	}()

	// Close results channel when done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for res := range resultsChan {
		results = append(results, res)
	}

	return results
}
