package correlation

import (
	"testing"

	"net-zilla/internal/models"
)

func TestEventCorrelator_Correlate(t *testing.T) {
	ec := NewEventCorrelator()

	tests := []struct {
		name    string
		report  *models.AdvancedReport
		wantKey string
	}{
		{
			name: "New Domain Invalid SSL",
			report: &models.AdvancedReport{
				BasicAnalysis: &models.ThreatAnalysis{
					WhoisInfo: &models.WhoisAnalysis{DomainAge: "New"},
					TLSInfo:   &models.TLSAnalysis{CertificateValid: false},
				},
			},
			wantKey: "correlation_insight_0",
		},
		{
			name: "Evasion Pattern",
			report: &models.AdvancedReport{
				BasicAnalysis: &models.ThreatAnalysis{
					RedirectChain: make([]models.RedirectDetail, 4),
					RedirectCount: 4,
					GeoAnalysis:   &models.GeoAnalysis{IsProxy: true},
				},
			},
			wantKey: "correlation_insight_0",
		},
		{
			name: "Infrastructure Anomaly",
			report: &models.AdvancedReport{
				BasicAnalysis: &models.ThreatAnalysis{
					DNSInfo:     &models.DNSAnalysis{ARecords: []string{"1.1.1.1", "2.2.2.2"}},
					GeoAnalysis: &models.GeoAnalysis{Country: "Unknown"},
				},
			},
			wantKey: "correlation_insight_0",
		},
		{
			name:   "Nil Analysis",
			report: &models.AdvancedReport{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ec.Correlate(tt.report)
			if tt.wantKey != "" {
				if tt.report.Metadata[tt.wantKey] == "" {
					t.Errorf("expected insight at %s", tt.wantKey)
				}
			}
		})
	}
}