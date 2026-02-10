package visualization

import (
	"fmt"
	"net-zilla/internal/models"
	"strings"
)

// ReportGenerator creates human-readable security reports.
type ReportGenerator struct{}

func NewReportGenerator() *ReportGenerator {
	return &ReportGenerator{}
}

// GenerateTextReport produces a professional ASCII-formatted report.
func (rg *ReportGenerator) GenerateTextReport(report *models.AdvancedReport) string {
	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString("================================================================================\n")
	sb.WriteString("                       NET-ZILLA ADVANCED SECURITY REPORT                       \n")
	sb.WriteString("================================================================================\n")
	sb.WriteString(fmt.Sprintf("Report ID: %s\n", report.ReportID))
	sb.WriteString(fmt.Sprintf("Target:    %s\n", report.Target))
	sb.WriteString(fmt.Sprintf("Timestamp: %s\n", report.Timestamp.Format("2006-01-02 15:04:05")))
	sb.WriteString("================================================================================\n\n")

	if report.RiskAssessment != nil {
		sb.WriteString("1. RISK ASSESSMENT SUMMARY\n")
		sb.WriteString("--------------------------\n")
		sb.WriteString(fmt.Sprintf("Overall Risk: %s\n", report.RiskAssessment.OverallRiskLevel))
		sb.WriteString(fmt.Sprintf("Risk Score:   %.2f/1.0\n", report.RiskAssessment.RiskScore))
		sb.WriteString(fmt.Sprintf("Findings:     %d critical issues detected\n", report.RiskAssessment.CriticalFindings))
		sb.WriteString(fmt.Sprintf("Summary:      %s\n\n", report.RiskAssessment.Summary))
	}

	if report.BasicAnalysis != nil {
		sb.WriteString("2. INFRASTRUCTURE INTELLIGENCE\n")
		sb.WriteString("------------------------------\n")
		if report.BasicAnalysis.GeoAnalysis != nil {
			geo := report.BasicAnalysis.GeoAnalysis
			sb.WriteString(fmt.Sprintf("IP Address:   %s\n", geo.IP))
			sb.WriteString(fmt.Sprintf("Location:     %s, %s\n", geo.City, geo.Country))
			sb.WriteString(fmt.Sprintf("ASN/ISP:      %s (%s)\n", geo.ASN, geo.ISP))
			sb.WriteString(fmt.Sprintf("Privacy:      Proxy/VPN: %v\n", geo.IsProxy))
		}
		sb.WriteString("\n")
	}

	if report.BehavioralAnalysis != nil && len(report.BehavioralAnalysis.Patterns) > 0 {
		sb.WriteString("3. BEHAVIORAL PATTERNS\n")
		sb.WriteString("----------------------\n")
		for _, p := range report.BehavioralAnalysis.Patterns {
			sb.WriteString(fmt.Sprintf("[!] %-20s | Type: %-10s | Weight: %d\n", p.Name, p.Type, p.Weight))
			sb.WriteString(fmt.Sprintf("    Description: %s\n", p.Description))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("================================================================================\n")
	sb.WriteString("                         END OF SECURITY ASSESSMENT                             \n")
	sb.WriteString("================================================================================\n")

	return sb.String()
}
