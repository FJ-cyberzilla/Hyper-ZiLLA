package shared_models

import (
	"fmt"
	"strings"
	"time"
)

// AIAnalysisResult represents the structured output from an ML-powered link analysis.
// This is used by models.ThreatAnalysis.AIResult
type AIAnalysisResult struct {
	IsSafe           bool                   `json:"is_safe"`
	Confidence       float64                `json:"confidence"`       // Confidence score of the analysis [0.0 - 1.0]
	RiskLevel        string                 `json:"risk_level"`       // E.g., "LOW", "MEDIUM", "HIGH"
	IsShortened      bool                   `json:"is_shortened"`
	HealthScore      float64                `json:"health_score"`     // Overall health score [0.0 - 1.0]
	Threats          []string               `json:"threats"`          // List of detected threats
	Recommendations  []string               `json:"recommendations"`
	Reasoning        string                 `json:"reasoning,omitempty"`       // Explanation of the analysis
	AnalysisID       string                 `json:"analysis_id,omitempty"`     // Unique ID for this analysis
	Timestamp        time.Time              `json:"timestamp,omitempty"`       // When analysis was performed
	ThreatIndicators []string               `json:"threat_indicators,omitempty"` // Specific indicators found
	Metadata         map[string]interface{} `json:"metadata,omitempty"`        // Additional analysis data
	Error            string                 `json:"error,omitempty"`           // Error message if analysis failed
}

// OrchestratorResult represents the outcome of an AI orchestration process.
// This is used by models.ThreatAnalysis.AIOrchestration
type OrchestrationResult struct {
	Success            bool                   `json:"success"`
	Target             string                 `json:"target,omitempty"`           // What was analyzed
	AnalysisTimestamp  time.Time              `json:"analysis_timestamp,omitempty"` // When analysis ran
	TasksExecuted      []string               `json:"tasks_executed"`
	Errors             []string               `json:"errors"`
	PerformanceMetrics map[string]float64     `json:"performance_metrics"`
	Findings           []string               `json:"findings,omitempty"`        // Specific findings from analysis
	Recommendations    []string               `json:"recommendations"`
	NextActions        []string               `json:"next_actions"`
	AIEnhancements     []AIAnalysisResult     `json:"ai_enhancements,omitempty"` // AI-specific results
	RawOutput          string                 `json:"raw_output,omitempty"`
}

// GetDescription provides a summary of the orchestration result for display.
func (or *OrchestrationResult) GetDescription() string {
	if !or.Success {
		if len(or.Errors) > 0 {
			return fmt.Sprintf("Orchestration Failed: %s", strings.Join(or.Errors, "; "))
		}
		return "Orchestration Failed: Unknown error"
	}
	
	if len(or.TasksExecuted) == 0 {
		return "Orchestration completed with no tasks executed"
	}
	
	taskSummary := fmt.Sprintf("%d tasks executed", len(or.TasksExecuted))
	if or.Target != "" {
		taskSummary = fmt.Sprintf("Analysis of %s: %s", or.Target, taskSummary)
	}
	
	if len(or.NextActions) > 0 {
		return fmt.Sprintf("%s. Next: %s", taskSummary, strings.Join(or.NextActions, ", "))
	}
	
	if len(or.Recommendations) > 0 {
		return fmt.Sprintf("%s. Recommendations: %s", taskSummary, strings.Join(or.Recommendations, ", "))
	}
	
	return fmt.Sprintf("%s. Status: Success", taskSummary)
}

// GetHealthScore calculates a simplified health score based on metrics.
func (or *OrchestrationResult) GetHealthScore() float64 {
	// Try to get efficiency score first
	if score, ok := or.PerformanceMetrics["efficiency_score"]; ok {
		return score
	}
	
	// Try to get success rate
	if successRate, ok := or.PerformanceMetrics["success_rate"]; ok {
		return successRate
	}
	
	// Calculate based on errors
	if !or.Success {
		return 0.2 // Default bad score for failure
	}
	
	// If successful but no metrics, use defaults
	if len(or.Errors) > 0 {
		// Success with errors - partial success
		return 0.6
	}
	
	// Success without errors
	return 0.9
}

// AddError adds an error to the orchestrator result
func (or *OrchestrationResult) AddError(err string) {
	or.Errors = append(or.Errors, err)
	or.Success = false
}

// AddTask adds a task to the executed tasks list
func (or *OrchestrationResult) AddTask(task string) {
	or.TasksExecuted = append(or.TasksExecuted, task)
}

// AddFinding adds a finding to the results
func (or *OrchestrationResult) AddFinding(finding string) {
	or.Findings = append(or.Findings, finding)
}

// AddMetric adds or updates a performance metric
func (or *OrchestrationResult) AddMetric(name string, value float64) {
	if or.PerformanceMetrics == nil {
		or.PerformanceMetrics = make(map[string]float64)
	}
	or.PerformanceMetrics[name] = value
}

// AddRecommendation adds a recommendation
func (or *OrchestrationResult) AddRecommendation(recommendation string) {
	or.Recommendations = append(or.Recommendations, recommendation)
}

// AddNextAction adds a next action
func (or *OrchestrationResult) AddNextAction(action string) {
	or.NextActions = append(or.NextActions, action)
}

// AddAIEnhancement adds an AI analysis result
func (or *OrchestrationResult) AddAIEnhancement(result AIAnalysisResult) {
	or.AIEnhancements = append(or.AIEnhancements, result)
}

// IsHighRisk checks if the analysis indicates high risk
func (ar *AIAnalysisResult) IsHighRisk() bool {
	riskLevels := map[string]bool{
		"HIGH":     true,
		"CRITICAL": true,
		"SEVERE":   true,
	}
	return riskLevels[strings.ToUpper(ar.RiskLevel)]
}

// IsMediumRisk checks if the analysis indicates medium risk
func (ar *AIAnalysisResult) IsMediumRisk() bool {
	riskLevels := map[string]bool{
		"MEDIUM": true,
		"MODERATE": true,
	}
	return riskLevels[strings.ToUpper(ar.RiskLevel)]
}

// IsLowRisk checks if the analysis indicates low risk
func (ar *AIAnalysisResult) IsLowRisk() bool {
	riskLevels := map[string]bool{
		"LOW":      true,
		"VERY_LOW": true,
		"SAFE":     true,
	}
	return riskLevels[strings.ToUpper(ar.RiskLevel)]
}

// GetRiskScore converts risk level to numeric score (0-100)
func (ar *AIAnalysisResult) GetRiskScore() int {
	switch strings.ToUpper(ar.RiskLevel) {
	case "CRITICAL", "SEVERE":
		return 90 + int((1.0-ar.Confidence)*10)
	case "HIGH":
		return 70 + int((1.0-ar.Confidence)*20)
	case "MEDIUM", "MODERATE":
		return 40 + int((1.0-ar.Confidence)*30)
	case "LOW":
		return 10 + int((1.0-ar.Confidence)*30)
	case "VERY_LOW", "SAFE":
		return int((1.0 - ar.Confidence) * 10)
	default:
		return int((1.0 - ar.Confidence) * 50) // Default if unknown
	}
}

// AddThreat adds a threat to the analysis result
func (ar *AIAnalysisResult) AddThreat(threat string) {
	ar.Threats = append(ar.Threats, threat)
}

// AddRecommendation adds a recommendation to the analysis result
func (ar *AIAnalysisResult) AddRecommendation(recommendation string) {
	ar.Recommendations = append(ar.Recommendations, recommendation)
}

// AddMetadata adds metadata to the analysis result
func (ar *AIAnalysisResult) AddMetadata(key string, value interface{}) {
	if ar.Metadata == nil {
		ar.Metadata = make(map[string]interface{})
	}
	ar.Metadata[key] = value
}

// GetMetadata retrieves metadata value
func (ar *AIAnalysisResult) GetMetadata(key string) (interface{}, bool) {
	if ar.Metadata == nil {
		return nil, false
	}
	value, exists := ar.Metadata[key]
	return value, exists
}

// HasError checks if there's an error in the analysis
func (ar *AIAnalysisResult) HasError() bool {
	return ar.Error != ""
}

// Summary provides a concise summary of the analysis
func (ar *AIAnalysisResult) Summary() string {
	if ar.HasError() {
		return fmt.Sprintf("AI Analysis Error: %s", ar.Error)
	}
	
	riskEmoji := "✅"
	if ar.IsHighRisk() {
		riskEmoji = "🚨"
	} else if ar.IsMediumRisk() {
		riskEmoji = "⚠️"
	} else if !ar.IsSafe {
		riskEmoji = "⚠️"
	}
	
	safety := "SAFE"
	if !ar.IsSafe {
		safety = "UNSAFE"
	}
	
	return fmt.Sprintf("%s %s | Confidence: %.0f%% | Risk: %s | Threats: %d", 
		riskEmoji, safety, ar.Confidence*100, ar.RiskLevel, len(ar.Threats))
}

// Validate checks if the analysis result is valid
func (ar *AIAnalysisResult) Validate() error {
	if ar.Confidence < 0 || ar.Confidence > 1 {
		return fmt.Errorf("confidence must be between 0 and 1, got %f", ar.Confidence)
	}
	
	if ar.HealthScore < 0 || ar.HealthScore > 1 {
		return fmt.Errorf("health score must be between 0 and 1, got %f", ar.HealthScore)
	}
	
	if ar.RiskLevel == "" {
		return fmt.Errorf("risk level cannot be empty")
	}
	
	return nil
}

// NewAIAnalysisResult creates a new AI analysis result with defaults
func NewAIAnalysisResult() *AIAnalysisResult {
	return &AIAnalysisResult{
		IsSafe:          true,
		Confidence:      0.5,
		RiskLevel:       "UNKNOWN",
		HealthScore:     0.5,
		Threats:         []string{},
		Recommendations: []string{},
		ThreatIndicators: []string{},
		Metadata:        make(map[string]interface{}),
		Timestamp:       time.Now(),
		AnalysisID:      generateAnalysisID(),
	}
}

// NewOrchestrationResult creates a new orchestration result with defaults
func NewOrchestrationResult() *OrchestrationResult {
	return &OrchestrationResult{
		Success:            true,
		TasksExecuted:      []string{},
		Errors:             []string{},
		PerformanceMetrics: make(map[string]float64),
		Findings:           []string{},
		Recommendations:    []string{},
		NextActions:        []string{},
		AIEnhancements:     []AIAnalysisResult{},
		AnalysisTimestamp:  time.Now(),
	}
}

// Helper function to generate analysis ID
func generateAnalysisID() string {
	return fmt.Sprintf("ai-%d", time.Now().UnixNano())
}
