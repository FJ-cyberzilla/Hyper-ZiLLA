package shared_models

import (
	"testing"
)

func TestAIAnalysisResult_GetRiskScore(t *testing.T) {
	res := AIAnalysisResult{
		RiskLevel:  "HIGH",
		Confidence: 0.9,
	}
	score := res.GetRiskScore()
	if score < 70 || score > 100 {
		t.Errorf("unexpected risk score: %d", score)
	}
}

func TestOrchestrationResult_GetHealthScore(t *testing.T) {
	res := OrchestrationResult{
		Success: true,
		PerformanceMetrics: map[string]float64{
			"success_rate": 0.95,
		},
	}
	score := res.GetHealthScore()
	if score != 0.95 {
		t.Errorf("expected 0.95, got %f", score)
	}
}

func TestAIAnalysisResult_Summary(t *testing.T) {
	res := AIAnalysisResult{
		IsSafe:     false,
		Confidence: 0.8,
		RiskLevel:  "HIGH",
		Threats:    []string{"Phishing"},
	}
	summary := res.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestOrchestrationResult_Helpers(t *testing.T) {
	res := NewOrchestrationResult()
	res.AddError("test error")
	res.AddTask("test task")
	res.AddFinding("test finding")
	res.AddMetric("test metric", 1.0)
	res.AddRecommendation("test rec")
	res.AddNextAction("test action")
	res.AddAIEnhancement(AIAnalysisResult{RiskLevel: "LOW"})

	if res.Success {
		t.Error("expected failure after AddError")
	}
	if len(res.Errors) != 1 {
		t.Error("expected 1 error")
	}
}

func TestAIAnalysisResult_RiskChecks(t *testing.T) {
	res := AIAnalysisResult{RiskLevel: "CRITICAL"}
	if !res.IsHighRisk() {
		t.Error("expected high risk")
	}
	
	res.RiskLevel = "MEDIUM"
	if !res.IsMediumRisk() {
		t.Error("expected medium risk")
	}
	
	res.RiskLevel = "SAFE"
	if !res.IsLowRisk() {
		t.Error("expected low risk")
	}
}

func TestAIAnalysisResult_Metadata(t *testing.T) {
	res := NewAIAnalysisResult()
	res.AddMetadata("key", "value")
	val, ok := res.GetMetadata("key")
	if !ok || val != "value" {
		t.Error("metadata retrieval failed")
	}
}

func TestOrchestrationResult_GetDescription(t *testing.T) {
	res := OrchestrationResult{
		Success: true,
		Target:  "example.com",
		TasksExecuted: []string{"task1"},
	}
	desc := res.GetDescription()
	if desc == "" {
		t.Error("expected non-empty description")
	}
}

func TestAIAnalysisResult_ThreatsAndRecs(t *testing.T) {
	res := NewAIAnalysisResult()
	res.AddThreat("Malware")
	res.AddRecommendation("Block")
	
	if len(res.Threats) != 1 || res.Threats[0] != "Malware" {
		t.Error("threat addition failed")
	}
	if len(res.Recommendations) != 1 || res.Recommendations[0] != "Block" {
		t.Error("recommendation addition failed")
	}
}

func TestOrchestrationResult_AddMethods(t *testing.T) {
	res := NewOrchestrationResult()
	res.AddFinding("finding")
	res.AddNextAction("action")
	
	if len(res.Findings) != 1 || res.Findings[0] != "finding" {
		t.Error("finding addition failed")
	}
	if len(res.NextActions) != 1 || res.NextActions[0] != "action" {
		t.Error("next action addition failed")
	}
}

func TestAIAnalysisResult_Validate(t *testing.T) {
	res := NewAIAnalysisResult()
	res.RiskLevel = "LOW"
	if err := res.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}
	
	res.Confidence = 1.5
	if err := res.Validate(); err == nil {
		t.Error("expected validation error for confidence > 1")
	}
}
