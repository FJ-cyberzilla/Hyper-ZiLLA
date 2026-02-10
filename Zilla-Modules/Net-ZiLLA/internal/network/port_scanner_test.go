package network

import (
	"context"
	"fmt"
	"net"
	"testing"

	"net-zilla/pkg/logger"
)

func TestPortScanner_Scan(t *testing.T) {
	l := logger.NewLogger()
	ps := NewPortScanner(l)

	// Start a local listener to simulate an open port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	defer ln.Close()

	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	var port int
	_, err = fmt.Sscanf(portStr, "%d", &port)
	if err != nil {
		t.Fatalf("failed to parse port: %v", err)
	}

	results := ps.Scan(context.Background(), "127.0.0.1", []int{port, 12345})

	foundOpen := false
	for _, res := range results {
		if res.Port == port {
			if !res.Open {
				t.Errorf("expected port %d to be open", port)
			}
			foundOpen = true
		}
		if res.Port == 12345 && res.Open {
			t.Errorf("expected port 12345 to be closed")
		}
	}

	if !foundOpen {
		t.Errorf("did not find result for port %d", port)
	}
}
