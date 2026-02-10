package models

import "time"

// AdvancedReport is the comprehensive production-ready report structure
type AdvancedReport struct {
	ReportID  string    `json:"report_id"`
	Timestamp time.Time `json:"timestamp"`
	Target    string    `json:"target"`

	// Comprehensive Analysis Sections
	ThreatIntelligence *IOCRegistry       `json:"threat_intelligence"`
	Reputation         *ReputationSummary `json:"reputation"`
	BehavioralAnalysis *BehaviorAnalysis  `json:"behavioral_analysis"`
	RiskAssessment     *RiskAssessment    `json:"risk_assessment"`

	// Legacy Support Integration
	BasicAnalysis *ThreatAnalysis `json:"basic_analysis"`

	Findings   []string               `json:"findings"`
	InstanceID string                 `json:"instance_id"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// AnalysisReport is a simplified report used for correlation and quick results
type AnalysisReport struct {
	RiskScore   int      `json:"risk_score"`
	Findings    []string `json:"findings"`
	SandboxUsed bool     `json:"sandbox_used"`
}
