package models

// RiskMetric defines a single vector of risk
type RiskMetric struct {
	Vector string `json:"vector"`
	Value  int    `json:"value"` // 0-10
	Impact string `json:"impact"`
}

// RiskAssessment provides the final security posture evaluation
type RiskAssessment struct {
	OverallRiskLevel string       `json:"overall_risk_level"`
	RiskScore        float64      `json:"risk_score"` // 0.0 - 1.0
	Metrics          []RiskMetric `json:"metrics"`
	Summary          string       `json:"summary"`
	CriticalFindings int          `json:"critical_findings"`
}
