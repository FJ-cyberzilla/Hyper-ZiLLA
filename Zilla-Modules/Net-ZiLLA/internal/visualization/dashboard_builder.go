package visualization

import (
	"net-zilla/internal/models"
)

// DashboardBuilder constructs high-level data summaries.
type DashboardBuilder struct{}

func NewDashboardBuilder() *DashboardBuilder {
	return &DashboardBuilder{}
}

// BuildSummary provides a "one-glance" overview of the analysis.
func (db *DashboardBuilder) BuildSummary(report *models.AdvancedReport) map[string]interface{} {
	summary := make(map[string]interface{})

	summary["target"] = report.Target
	summary["risk_level"] = "UNKNOWN"
	if report.RiskAssessment != nil {
		summary["risk_level"] = report.RiskAssessment.OverallRiskLevel
		summary["score"] = report.RiskAssessment.RiskScore
	}

	if report.BasicAnalysis != nil {
		summary["hops"] = report.BasicAnalysis.RedirectCount
		if report.BasicAnalysis.GeoAnalysis != nil {
			summary["origin_country"] = report.BasicAnalysis.GeoAnalysis.Country
		}
	}

	return summary
}
