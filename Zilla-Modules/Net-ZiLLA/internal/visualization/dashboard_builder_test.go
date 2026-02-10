package visualization

import (
	"testing"

	"net-zilla/internal/models"
)

func TestDashboardBuilder_BuildSummary(t *testing.T) {
	db := NewDashboardBuilder()
	report := &models.AdvancedReport{
		Target: "http://test.com",
		RiskAssessment: &models.RiskAssessment{
			OverallRiskLevel: "MEDIUM",
			RiskScore:        0.5,
		},
		BasicAnalysis: &models.ThreatAnalysis{
			RedirectCount: 2,
			GeoAnalysis:   &models.GeoAnalysis{Country: "USA"},
		},
	}

	summary := db.BuildSummary(report)

	if summary["target"] != "http://test.com" {
		t.Error("target mismatch")
	}
	if summary["risk_level"] != "MEDIUM" {
		t.Error("risk level mismatch")
	}
	if summary["origin_country"] != "USA" {
		t.Error("country mismatch")
	}
}
