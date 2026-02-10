package analyzer_interface

import (
	"context"
	"net-zilla/internal/models"
)

// ThreatAnalyzerInterface defines the methods that the Menu needs from the ThreatAnalyzer.
type ThreatAnalyzerInterface interface {
	ComprehensiveAnalysis(ctx context.Context, targetURL string) (*models.ThreatAnalysis, error)
	PerformDNSLookup(ctx context.Context, target string) (*models.DNSAnalysis, error)
	PerformWhoisLookup(ctx context.Context, target string) (*models.WhoisAnalysis, error)
	PerformIPGeolocation(ctx context.Context, target string) (*models.GeoAnalysis, error)
	PerformTLSAnalysis(ctx context.Context, target string) (*models.TLSAnalysis, error)
	PerformTraceroute(ctx context.Context, target string) (*models.NetworkAnalysis, error)
}
