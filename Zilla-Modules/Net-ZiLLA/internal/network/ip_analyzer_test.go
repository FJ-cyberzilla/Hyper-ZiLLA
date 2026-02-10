package network

import (
	"net"
	"testing"

	"net-zilla/internal/models"
)

func TestIsPublicIP(t *testing.T) {
	tests := []struct {
		ip   string
		want bool
	}{
		{"8.8.8.8", true},
		{"1.1.1.1", true},
		{"127.0.0.1", false},
		{"10.0.0.1", false},
		{"192.168.1.1", false},
		{"172.16.0.1", false},
		{"172.31.255.255", false},
		{"172.32.0.1", true},
		{"100.64.0.1", false},
	}

	for _, tt := range tests {
		if got := isPublicIP(net.ParseIP(tt.ip)); got != tt.want {
			t.Errorf("isPublicIP(%s) = %v, want %v", tt.ip, got, tt.want)
		}
	}
}

func TestIsReservedIP(t *testing.T) {
	tests := []struct {
		ip   string
		want bool
	}{
		{"8.8.8.8", false},
		{"0.0.0.0", true},
		{"240.0.0.1", true},
		{"127.0.0.1", true},
	}

	for _, tt := range tests {
		if got := isReservedIP(net.ParseIP(tt.ip)); got != tt.want {
			t.Errorf("isReservedIP(%s) = %v, want %v", tt.ip, got, tt.want)
		}
	}
}

func TestDetectHostingType(t *testing.T) {
	ipa := &IPAnalyzer{}
	tests := []struct {
		isp  string
		want string
	}{
		{"Amazon.com", "Hosting Provider"},
		{"Google LLC", "Hosting Provider"},
		{"Comcast Cable", "ISP/Residential"},
		{"DigitalOcean", "Hosting Provider"},
		{"Some Random ISP", "ISP/Residential"},
	}

	for _, tt := range tests {
		analysis := &models.GeoAnalysis{ISP: tt.isp}
		if got := ipa.detectHostingType(analysis); got != tt.want {
			t.Errorf("detectHostingType(%s) = %v, want %v", tt.isp, got, tt.want)
		}
	}
}
