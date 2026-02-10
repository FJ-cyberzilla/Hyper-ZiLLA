package network

import (
	"context"
	"testing"

	"net-zilla/pkg/logger"
)

func TestNewDNSClient(t *testing.T) {
	l := logger.NewLogger()
	dc := NewDNSClient(l)
	if dc == nil {
		t.Fatal("NewDNSClient returned nil")
	}
	if dc.logger != l {
		t.Error("DNSClient logger mismatch")
	}
}

func TestDNSClient_ReverseDNSLookup_Empty(t *testing.T) {
	l := logger.NewLogger()
	dc := NewDNSClient(l)
	_, err := dc.ReverseDNSLookup(context.Background(), "")
	if err == nil {
		t.Error("expected error for empty IP in ReverseDNSLookup")
	}
}
