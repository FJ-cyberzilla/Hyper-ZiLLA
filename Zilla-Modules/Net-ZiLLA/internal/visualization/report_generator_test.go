package visualization

import (
	"strings"
	"testing"
	"time"

	"net-zilla/internal/models"
)

func TestReportGenerator_GenerateTextReport(t *testing.T) {
	rg := NewReportGenerator()

	report := &models.AdvancedReport{
		ReportID:  "TEST-123",
		Target:    "http://example.com",
		Timestamp: time.Now(),
		RiskAssessment: &models.RiskAssessment{
			OverallRiskLevel: "HIGH",
			RiskScore:        0.75,
			Summary:          "Test summary",
		},
		BasicAnalysis: &models.ThreatAnalysis{
			GeoAnalysis: &models.GeoAnalysis{
				IP:      "1.2.3.4",
				Country: "USA",
				ISP:     "Test ISP",
			},
		},
		BehavioralAnalysis: &models.BehaviorAnalysis{
			Patterns: []models.BehavioralPattern{
				{Name: "Test Pattern", Type: "Malware", Weight: 50, Description: "Test description"},
			},
		},
	}

	got := rg.GenerateTextReport(report)

	if !strings.Contains(got, "TEST-123") {
		t.Error("report ID missing")
	}
	if !strings.Contains(got, "HIGH") {
		t.Error("risk level missing")
	}
	if !strings.Contains(got, "Test ISP") {
		t.Error("ISP missing")
	}
	if !strings.Contains(got, "Test Pattern") {
		t.Error("behavioral pattern missing")
	}
}

func TestReportGenerator_GenerateTextReport_Minimal(t *testing.T) {
	rg := NewReportGenerator()
	report := &models.AdvancedReport{
		ReportID: "MIN-1",
		Target:   "http://example.com",
	}

	got := rg.GenerateTextReport(report)
	if !strings.Contains(got, "MIN-1") {
		t.Error("report ID missing")
	}
}