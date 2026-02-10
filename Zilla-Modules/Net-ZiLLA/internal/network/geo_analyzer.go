package network

import (
	"context"
	"net-zilla/internal/models"
	"net-zilla/pkg/logger"
)

type GeoAnalyzer struct {
	logger *logger.Logger
}

func NewGeoAnalyzer(logger *logger.Logger) *GeoAnalyzer {
	return &GeoAnalyzer{logger: logger}
}

// AnalyzeRegion assesses the risk of a specific geographic region
func (ga *GeoAnalyzer) AnalyzeRegion(ctx context.Context, geo *models.GeoAnalysis) int {
	riskScore := 0

	// High risk jurisdictions (Example list)
	highRiskCountries := map[string]bool{
		"Unknown": true, // Often used by TOR/Proxies
	}

	if highRiskCountries[geo.Country] {
		riskScore += 40
	}

	// Proxy/VPN check
	if geo.IsProxy {
		riskScore += 30
	}

	// Hosting center check (Cloud providers are often used for short-lived phishing)
	if geo.HostingType == "Hosting Provider" {
		riskScore += 15
	}

	return riskScore
}
