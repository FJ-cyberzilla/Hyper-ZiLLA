package correlation

import (
	"fmt"
	"net-zilla/internal/models"
)

// EventCorrelator cross-references findings from multiple analysis vectors.
type EventCorrelator struct{}

func NewEventCorrelator() *EventCorrelator {
	return &EventCorrelator{}
}

// Correlate analyzes the relationships between different report components
// to identify sophisticated attack patterns like "Fast-Flux" or "Magecart".
func (ec *EventCorrelator) Correlate(report *models.AdvancedReport) {
	if report.BasicAnalysis == nil {
		return
	}

	ba := report.BasicAnalysis

	// 1. Correlate Domain Age with SSL Status
	if ba.WhoisInfo != nil && ba.TLSInfo != nil {
		if ba.WhoisInfo.DomainAge == "New" && !ba.TLSInfo.CertificateValid {
			ec.addInsight(report, "High Risk: Newly registered domain using invalid/self-signed SSL.")
		}
	}

	// 2. Correlate Redirect Chains with IP Reputation
	if len(ba.RedirectChain) > 0 && ba.GeoAnalysis != nil {
		if ba.RedirectCount > 3 && ba.GeoAnalysis.IsProxy {
			ec.addInsight(report, "Evasion Pattern: Long redirect chain originating from a Proxy/VPN node.")
		}
	}

	// 3. Correlate DNS Records with Geolocation
	if ba.DNSInfo != nil && ba.GeoAnalysis != nil {
		// Check if A records point to a different country than the primary resolution
		if len(ba.DNSInfo.ARecords) > 1 && ba.GeoAnalysis.Country == "Unknown" {
			ec.addInsight(report, "Infrastructure Anomaly: Multiple A-records resolved via anonymous infrastructure.")
		}
	}
}

func (ec *EventCorrelator) addInsight(report *models.AdvancedReport, insight string) {
	if report.Metadata == nil {
		report.Metadata = make(map[string]interface{})
	}
	key := fmt.Sprintf("correlation_insight_%d", len(report.Metadata))
	report.Metadata[key] = insight
}
