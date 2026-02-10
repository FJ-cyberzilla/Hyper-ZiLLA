package network

import (
	"context"
	"strings"
	"testing"
	"time"

	"net-zilla/pkg/logger"
)

func TestNewWhoisClient(t *testing.T) {
	l := logger.NewLogger()
	wc := NewWhoisClient(l)
	if wc == nil {
		t.Fatal("NewWhoisClient returned nil")
	}
	if wc.logger != l {
		t.Error("WhoisClient logger mismatch")
	}
}

func TestWhoisClient_GetWhoisServer(t *testing.T) {
	wc := &WhoisClient{
		servers: map[string]string{"com": "whois.verisign-grs.com"},
		logger:  logger.NewLogger(),
	}
	
	tests := []struct {
		domain string
		want   string
	}{
		{"example.com", "whois.verisign-grs.com"},
		{"example.unknown", "whois.iana.org"},
	}

	for _, tt := range tests {
		got, err := wc.getWhoisServer(context.Background(), tt.domain)
		if err != nil {
			t.Errorf("getWhoisServer(%s) error: %v", tt.domain, err)
		}
		if got != tt.want {
			t.Errorf("getWhoisServer(%s) = %v, want %v", tt.domain, got, tt.want)
		}
	}
}

func TestCalculateDomainAge(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		created time.Time
		want    string
	}{
		{time.Time{}, "Unknown"},
		{now.AddDate(-5, 0, 0), "5 years"},
		{now.AddDate(0, -3, 0), "3 months"},
		{now.AddDate(0, 0, -10), "10 days"},
	}

	for _, tt := range tests {
		if got := calculateDomainAge(tt.created); !strings.Contains(got, tt.want) {
			t.Errorf("calculateDomainAge(%v) = %v, want %v", tt.created, got, tt.want)
		}
	}
}

func TestWhoisClient_ParseWhoisResponse(t *testing.T) {
	wc := &WhoisClient{}
	info := &WhoisInfo{}
	raw := `
Domain Name: EXAMPLE.COM
Registrar Name: Safe Registrar LLC
Creation Date: 2020-01-01T00:00:00Z
Registry Expiry Date: 2025-01-01T00:00:00Z
Name Server: NS1.EXAMPLE.COM
Name Server: NS2.EXAMPLE.COM
Domain Status: clientTransferProhibited
`
	wc.parseWhoisResponse(raw, info)

	if info.Registrar != "Safe Registrar LLC" {
		t.Errorf("expected Registrar Safe Registrar LLC, got %s", info.Registrar)
	}
	if info.CreatedDate.Year() != 2020 {
		t.Errorf("expected Year 2020, got %d", info.CreatedDate.Year())
	}
	if len(info.NameServers) != 2 {
		t.Errorf("expected 2 NameServers, got %d", len(info.NameServers))
	}
	if len(info.Status) != 1 {
		t.Errorf("expected 1 Status, got %d", len(info.Status))
	}
}
