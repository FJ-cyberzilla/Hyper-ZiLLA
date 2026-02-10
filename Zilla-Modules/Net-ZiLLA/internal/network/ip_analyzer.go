package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"net-zilla/internal/models"
	"net-zilla/pkg/logger"
)

// IPAnalyzer performs IP address analysis, including geolocation.
type IPAnalyzer struct {
	logger    *logger.Logger
	dnsClient *DNSClient
	client    *http.Client
}

// NewIPAnalyzer creates and initializes a new IPAnalyzer.
func NewIPAnalyzer(logger *logger.Logger) *IPAnalyzer {
	return &IPAnalyzer{
		logger:    logger,
		dnsClient: NewDNSClient(logger),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetGeolocation performs IP geolocation using ip-api.com.
func (ipa *IPAnalyzer) GetGeolocation(ctx context.Context, ip string) (*models.GeoAnalysis, error) {
	analysis := &models.GeoAnalysis{
		IP: ip,
	}

	// Basic IP validation
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		// If not an IP, try to resolve it
		ips, err := net.LookupIP(ip)
		if err != nil || len(ips) == 0 {
			return nil, fmt.Errorf("invalid IP address or unresolvable host: %s", ip)
		}
		ip = ips[0].String()
		parsedIP = ips[0]
		analysis.IP = ip
	}

	// Geolocation using ip-api.com (Free for non-commercial)
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,message,country,city,isp,as,lat,lon,proxy", ip)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := ipa.client.Do(req)
	if err != nil {
		ipa.logger.Warn("Geolocation lookup failed for %s: %v", ip, err)
	} else {
		defer resp.Body.Close()
		var result struct {
			Status  string  `json:"status"`
			Country string  `json:"country"`
			City    string  `json:"city"`
			ISP     string  `json:"isp"`
			AS      string  `json:"as"`
			Lat     float64 `json:"lat"`
			Lon     float64 `json:"lon"`
			Proxy   bool    `json:"proxy"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err == nil && result.Status == "success" {
			analysis.Country = result.Country
			analysis.City = result.City
			analysis.ISP = result.ISP
			analysis.ASN = result.AS
			analysis.Latitude = result.Lat
			analysis.Longitude = result.Lon
			analysis.IsProxy = result.Proxy
		}
	}

	analysis.IsPublic = isPublicIP(parsedIP)
	analysis.IsReserved = isReservedIP(parsedIP)
	analysis.HostingType = ipa.detectHostingType(analysis)

	return analysis, nil
}

func isPublicIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := ip.To4(); ip4 != nil {
		if ip4[0] == 10 ||
			(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) ||
			(ip4[0] == 192 && ip4[1] == 168) ||
			(ip4[0] == 100 && ip4[1] >= 64 && ip4[1] <= 127) {
			return false
		}
	}
	return true
}

func isReservedIP(ip net.IP) bool {
	if ip.IsUnspecified() || ip.IsMulticast() || ip.IsInterfaceLocalMulticast() || ip.IsLoopback() {
		return true
	}
	if ip4 := ip.To4(); ip4 != nil {
		if ip4[0] == 0 ||
			(ip4[0] == 192 && ip4[1] == 0 && ip4[2] == 0) ||
			(ip4[0] == 192 && ip4[1] == 0 && ip4[2] == 2) ||
			(ip4[0] == 203 && ip4[1] == 0 && ip4[2] == 113) ||
			(ip4[0] >= 240 && ip4[0] <= 255) {
			return true
		}
	}
	return false
}

func (ipa *IPAnalyzer) detectHostingType(analysis *models.GeoAnalysis) string {
	ispLower := strings.ToLower(analysis.ISP)
	hostingKeywords := []string{"cloud", "host", "server", "data center", "amazon", "google", "microsoft", "digitalocean", "linode", "akamai"}
	for _, keyword := range hostingKeywords {
		if strings.Contains(ispLower, keyword) {
			return "Hosting Provider"
		}
	}
	return "ISP/Residential"
}
