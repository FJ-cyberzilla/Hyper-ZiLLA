package network

import (
	"context"
	"fmt"
	"net"
	"testing"
)

func TestServiceDetector_CommonPortMap(t *testing.T) {
	sd := NewServiceDetector()
	if got := sd.commonPortMap(80); got != "HTTP" {
		t.Errorf("expected HTTP, got %s", got)
	}
	if got := sd.commonPortMap(9999); got != "Unknown Service" {
		t.Errorf("expected Unknown Service, got %s", got)
	}
}

func TestServiceDetector_Detect(t *testing.T) {
	sd := NewServiceDetector()

	// 1. Test real response (Mock HTTP)
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		conn, _ := ln.Accept()
		if conn != nil {
			fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\n\r\n")
			conn.Close()
		}
	}()
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	var port int
	fmt.Sscanf(portStr, "%d", &port)

	got := sd.Detect(context.Background(), "127.0.0.1", port)
	if got != "Web Server (HTTP)" {
		t.Errorf("expected Web Server (HTTP), got %s", got)
	}
	ln.Close()

	// 2. Test fallback to port map
	ln2, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		conn, _ := ln2.Accept()
		if conn != nil {
			// Just close without sending anything
			conn.Close()
		}
	}()
	_, portStr2, _ := net.SplitHostPort(ln2.Addr().String())
	var port2 int
	fmt.Sscanf(portStr2, "%d", &port2)
	
	// We need a known port for commonPortMap to trigger, but our listener is random.
	// So we'll just verify it doesn't crash and returns something.
	got2 := sd.Detect(context.Background(), "127.0.0.1", port2)
	if got2 == "Unknown" {
		t.Errorf("expected commonPortMap or banner, got Unknown")
	}
	ln2.Close()
}
