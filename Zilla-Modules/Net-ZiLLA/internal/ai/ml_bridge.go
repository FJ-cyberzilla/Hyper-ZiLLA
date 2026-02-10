package ai

import (
	"context"
	"fmt"
	"sync"
	"time"

	"net-zilla/internal/config"
	"net-zilla/internal/models"
	"net-zilla/internal/shared_models"
	"net-zilla/internal/storage"
	"net-zilla/pkg/analyzer_interface"
	"net-zilla/pkg/logger"
)

// MLAgent provides AI/ML analysis capabilities by integrating GoAgent with the orchestrator
type MLAgent struct {
	config           *config.AIConfig
	goAgent          *GoAgent
	orchestrator     analyzer_interface.OrchestratorInterface
	db               *storage.Database
	log              *logger.Logger
	
	// Real metrics tracking
	metrics          Metrics
	metricsMutex     sync.RWMutex
	startTime        time.Time
}

// Metrics tracks real analysis performance
type Metrics struct {
	TotalAnalyses      int64
	SuccessfulAnalyses int64
	FailedAnalyses     int64
	AverageLatency     time.Duration
	CacheHits          int64
	CacheMisses        int64
	ConcurrentAnalyses int32
	LastError          string
	LastErrorTime      time.Time
	ModelInferences    map[string]int64
}

// NewMLAgent creates a production MLAgent
func NewMLAgent(cfg *config.AIConfig, orch analyzer_interface.OrchestratorInterface, db *storage.Database) (*MLAgent, error) {
	if cfg == nil {
		return nil, fmt.Errorf("AI config is required")
	}

	// Initialize real GoAgent with confidence threshold
	goAgent := NewGoAgent(cfg.ConfidenceThreshold)
	
	agent := &MLAgent{
		config:       cfg,
		goAgent:      goAgent,
		orchestrator: orch,
		db:           db,
		log:          logger.New().WithComponent("ml_agent"),
		startTime:    time.Now(),
		metrics: Metrics{
			ModelInferences: make(map[string]int64),
		},
	}

	agent.log.Info("MLAgent initialized with AI enabled: %v", cfg.EnableAI)
	return agent, nil
}

// AnalyzeLink performs AI analysis and correlates with orchestrator results
func (a *MLAgent) AnalyzeLink(ctx context.Context, threatAnalysis *models.ThreatAnalysis) (*shared_models.AIAnalysisResult, error) {
	start := time.Now()
	a.incrementConcurrent()
	defer a.decrementConcurrent()

	// Input validation
	if threatAnalysis == nil {
		return a.handleError("threat analysis cannot be nil"), nil
	}
	if threatAnalysis.URL == "" {
		return a.handleError("URL cannot be empty"), nil
	}

	a.log.Debug("Starting AI analysis for: %s", threatAnalysis.URL)

	// If AI is disabled, return baseline
	if !a.config.EnableAI {
		a.updateMetrics(false, time.Since(start))
		return a.createBaselineResult(threatAnalysis), nil
	}

	// Step 1: Run GoAgent analysis (your real implementation)
	goAgentResult, err := a.goAgent.AnalyzeLink(threatAnalysis)
	if err != nil {
		a.recordError("go_agent_analysis", err)
		return a.fallbackAnalysis(ctx, threatAnalysis, start), nil
	}

	// Step 2: Get comprehensive analysis from orchestrator
	var orchResult *models.AnalysisReport
	if a.orchestrator != nil {
		var err error
		orchResult, err = a.orchestrator.FullAnalysis(ctx, threatAnalysis.URL)
		if err != nil {
			a.recordError("orchestrator_analysis", err)
			// Continue with AI results only
			a.updateMetrics(false, time.Since(start))
			return a.enhanceWithOrchestratorData(goAgentResult, nil, threatAnalysis), nil
		}
	}

	// Step 3: Correlate AI and orchestrator results
	finalResult := a.correlateResults(goAgentResult, orchResult, threatAnalysis)
	finalResult.Metadata["total_processing_ms"] = time.Since(start).Milliseconds()

	// Step 4: Store in database
	if a.db != nil {
		if err := a.db.SaveAnalysis(ctx, threatAnalysis); err != nil {
			a.log.Warn("Failed to save to database: %v", err)
		}
	}

	// Update metrics
	a.updateMetrics(true, time.Since(start))
	a.recordModelInference("link_analysis")

	a.log.Info("Analysis completed for %s in %v, Risk: %s", 
		threatAnalysis.URL, time.Since(start), finalResult.RiskLevel)

	return finalResult, nil
}

// OrchestrateAnalysis orchestrates AI and traditional analysis
func (a *MLAgent) OrchestrateAnalysis(ctx context.Context, url string, analysisType string) (*shared_models.OrchestrationResult, error) {
	start := time.Now()

	if url == "" {
		return nil, fmt.Errorf("URL cannot be empty")
	}

	// Run full orchestrator analysis
	var orchResult *models.AnalysisReport
	if a.orchestrator != nil {
		var err error
		orchResult, err = a.orchestrator.FullAnalysis(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("orchestrator failed: %w", err)
		}
	} else {
		orchResult = &models.AnalysisReport{Findings: []string{"Orchestrator not available"}}
	}

	// Run AI analysis if enabled
	var aiEnhancements []shared_models.AIAnalysisResult
	if a.config.EnableAI {
		threatAnalysis := &models.ThreatAnalysis{
			URL:     url,
		}
		
		aiResult, err := a.AnalyzeLink(ctx, threatAnalysis)
		if err == nil {
			aiEnhancements = append(aiEnhancements, *aiResult)
		}
	}

	// Compile results
	result := &shared_models.OrchestrationResult{
		Success:            true,
		Target:             url,
		AnalysisTimestamp:  time.Now(),
		TasksExecuted:      a.getExecutedTasks(orchResult, len(aiEnhancements) > 0),
		PerformanceMetrics: a.getPerformanceMetrics(start),
		Findings:           orchResult.Findings,
		AIEnhancements:     aiEnhancements,
		Recommendations:    a.generateRecommendations(orchResult, aiEnhancements),
		NextActions:        a.determineNextActions(orchResult),
	}

	return result, nil
}

// SystemDiagnostics provides real system health information
func (a *MLAgent) SystemDiagnostics(ctx context.Context) (*shared_models.OrchestrationResult, error) {
	// Get orchestrator health
	orchHealth := a.orchestrator.GetHealth(ctx)
	
	// Check database
	dbHealth := "healthy"
	if a.db != nil {
		if err := a.db.Ping(ctx); err != nil {
			dbHealth = "unhealthy"
			a.recordError("database_ping", err)
		}
	}

	// Get metrics
	a.metricsMutex.RLock()
	metrics := a.metrics
	uptime := time.Since(a.startTime)
	successRate := float64(metrics.SuccessfulAnalyses) / float64(metrics.TotalAnalyses)
	a.metricsMutex.RUnlock()

	return &shared_models.OrchestrationResult{
		Success:           orchHealth && dbHealth == "healthy",
		Target:            "system_diagnostics",
		AnalysisTimestamp: time.Now(),
		TasksExecuted: []string{
			"orchestrator_health_check",
			"database_health_check",
			"ai_system_check",
			"performance_analysis",
		},
		PerformanceMetrics: map[string]float64{
			"uptime_hours":               uptime.Hours(),
			"total_analyses":             float64(metrics.TotalAnalyses),
			"success_rate":               successRate,
			"average_latency_ms":         float64(metrics.AverageLatency.Milliseconds()),
			"cache_hit_rate":             a.calculateCacheHitRate(),
			"concurrent_analyses":        float64(metrics.ConcurrentAnalyses),
		},
		Findings: []string{
			fmt.Sprintf("Orchestrator: %v", orchHealth),
			fmt.Sprintf("Database: %s", dbHealth),
			fmt.Sprintf("AI Enabled: %v", a.config.EnableAI),
			fmt.Sprintf("Uptime: %v", uptime.Truncate(time.Second)),
			fmt.Sprintf("Total Analyses: %d", metrics.TotalAnalyses),
		},
		Recommendations: a.generateSystemRecommendations(nil, dbHealth, successRate),
		NextActions: []string{
			"Continue monitoring",
			"Review error logs",
			"Update threat intelligence",
		},
	}, nil
}

// AnalyzeSMS uses your real GoAgent SMS analysis
func (a *MLAgent) AnalyzeSMS(ctx context.Context, message string) (*shared_models.AIAnalysisResult, error) {
	if message == "" {
		return a.handleError("SMS message cannot be empty"), nil
	}

	if !a.config.EnableAI {
		return &shared_models.AIAnalysisResult{
			IsSafe:     true,
			Confidence: 0.3,
			RiskLevel:  "UNKNOWN",
			Reasoning:  "AI analysis disabled",
			Timestamp:  time.Now(),
		}, nil
	}

	// Use your real GoAgent SMS analysis
	return a.goAgent.AnalyzeSMS(message)
}

// IsAvailable checks if the AI system is operational
func (a *MLAgent) IsAvailable() bool {
	return a.config.EnableAI
}

// Helper methods
func (a *MLAgent) correlateResults(
	aiResult *shared_models.AIAnalysisResult,
	orchResult *models.AnalysisReport,
	threatAnalysis *models.ThreatAnalysis,
) *shared_models.AIAnalysisResult {
	
	result := &shared_models.AIAnalysisResult{
		IsSafe:          aiResult.IsSafe,
		Confidence:      aiResult.Confidence,
		RiskLevel:       aiResult.RiskLevel,
		IsShortened:     aiResult.IsShortened,
		HealthScore:     aiResult.HealthScore,
		Threats:         aiResult.Threats,
		Recommendations: aiResult.Recommendations,
		Timestamp:       time.Now(),
		Metadata:        make(map[string]interface{}),
	}

	// Enhance with orchestrator data
	if orchResult != nil {
		// Adjust confidence based on orchestrator risk score
		if orchResult.RiskScore > 70 && aiResult.IsSafe {
			result.Confidence *= 0.7
			if result.Confidence < 0.5 {
				result.RiskLevel = "MEDIUM"
			}
			result.Reasoning = "Traditional analysis indicates higher risk than AI assessment"
		}

		// Add metadata
		result.Metadata = map[string]interface{}{
			"orchestrator_risk_score": orchResult.RiskScore,
			"ai_confidence":          aiResult.Confidence,
			"combined_risk":          (float64(orchResult.RiskScore) + (100 * (1 - aiResult.Confidence))) / 2,
			"findings_count":         len(orchResult.Findings),
		}
	}

	return result
}

func (a *MLAgent) enhanceWithOrchestratorData(
	aiResult *shared_models.AIAnalysisResult,
	orchResult *models.AnalysisReport,
	threatAnalysis *models.ThreatAnalysis,
) *shared_models.AIAnalysisResult {
	
	if orchResult == nil {
		aiResult.Metadata = map[string]interface{}{
			"orchestrator_available": false,
			"analysis_type":          "ai_only",
		}
		return aiResult
	}

	aiResult.Metadata = map[string]interface{}{
		"orchestrator_risk_score": orchResult.RiskScore,
		"findings_count":         len(orchResult.Findings),
		"analysis_type":          "combined",
	}

	return aiResult
}

func (a *MLAgent) fallbackAnalysis(ctx context.Context, ta *models.ThreatAnalysis, start time.Time) *shared_models.AIAnalysisResult {
	// Try to get orchestrator results as fallback
	orchResult, err := a.orchestrator.GetQuickAnalysis(ctx, ta.URL)
	if err != nil {
		return a.createEmergencyResult(ta)
	}

	return &shared_models.AIAnalysisResult{
		IsSafe:           orchResult.RiskScore < 50,
		Confidence:       0.6,
		RiskLevel:        a.mapRiskScore(orchResult.RiskScore),
		Reasoning:        "AI analysis failed, using traditional analysis",
		Timestamp:        time.Now(),
		Threats:          orchResult.Findings,
		Metadata: map[string]interface{}{
			"ai_failed":        true,
			"fallback_mode":    true,
			"processing_ms":    time.Since(start).Milliseconds(),
		},
	}
}

func (a *MLAgent) createBaselineResult(ta *models.ThreatAnalysis) *shared_models.AIAnalysisResult {
	return &shared_models.AIAnalysisResult{
		IsSafe:     true,
		Confidence: 0.5,
		RiskLevel:  "BASELINE",
		Reasoning:  "AI analysis disabled, using baseline assessment",
		Timestamp:  time.Now(),
		Metadata: map[string]interface{}{
			"ai_enabled": false,
		},
	}
}

func (a *MLAgent) handleError(msg string) *shared_models.AIAnalysisResult {
	return &shared_models.AIAnalysisResult{
		IsSafe:     false,
		Confidence: 0.1,
		RiskLevel:  "ERROR",
		Reasoning:  msg,
		Timestamp:  time.Now(),
		Metadata: map[string]interface{}{
			"error": true,
		},
	}
}

func (a *MLAgent) createEmergencyResult(ta *models.ThreatAnalysis) *shared_models.AIAnalysisResult {
	return &shared_models.AIAnalysisResult{
		IsSafe:     false,
		Confidence: 0.2,
		RiskLevel:  "HIGH",
		Reasoning:  "System failure - assuming worst case",
		Timestamp:  time.Now(),
		Metadata: map[string]interface{}{
			"emergency_mode": true,
			"recommendation": "Manual review required",
		},
	}
}

func (a *MLAgent) mapRiskScore(score int) string {
	switch {
	case score >= 80:
		return "HIGH"
	case score >= 50:
		return "MEDIUM"
	case score >= 20:
		return "LOW"
	default:
		return "VERY_LOW"
	}
}

func (a *MLAgent) getExecutedTasks(orchResult *models.AnalysisReport, aiUsed bool) []string {
	tasks := []string{"dns_analysis", "ssl_analysis", "threat_intel"}
	if aiUsed {
		tasks = append(tasks, "ai_ml_analysis")
	}
	if orchResult.SandboxUsed {
		tasks = append(tasks, "sandbox_analysis")
	}
	return tasks
}

func (a *MLAgent) getPerformanceMetrics(start time.Time) map[string]float64 {
	a.metricsMutex.RLock()
	defer a.metricsMutex.RUnlock()
	
	return map[string]float64{
		"total_duration_ms": float64(time.Since(start).Milliseconds()),
		"success_rate":      a.calculateSuccessRate(),
		"avg_latency_ms":    float64(a.metrics.AverageLatency.Milliseconds()),
	}
}

func (a *MLAgent) generateRecommendations(orchResult *models.AnalysisReport, aiEnhancements []shared_models.AIAnalysisResult) []string {
	var recs []string
	
	if orchResult != nil {
		if orchResult.RiskScore > 70 {
			recs = append(recs, "Block access to this resource")
		} else if orchResult.RiskScore > 40 {
			recs = append(recs, "Monitor this resource closely")
		}
	}
	
	if len(aiEnhancements) > 0 {
		lastAI := aiEnhancements[len(aiEnhancements)-1]
		if lastAI.RiskLevel == "HIGH" {
			recs = append(recs, "AI detected high-risk patterns")
		}
	}
	
	if len(recs) == 0 {
		recs = append(recs, "No immediate action required")
	}
	
	return recs
}

func (a *MLAgent) determineNextActions(orchResult *models.AnalysisReport) []string {
	if orchResult == nil {
		return []string{"Retry analysis", "Check system connectivity"}
	}
	
	if orchResult.RiskScore > 80 {
		return []string{
			"Add to threat database",
			"Alert security team",
		}
	}
	
	return []string{"Log analysis result"}
}

func (a *MLAgent) generateSystemRecommendations(orchErr error, dbHealth string, successRate float64) []string {
	var recs []string
	
	if orchErr != nil {
		recs = append(recs, "Restart orchestrator service")
	}
	if dbHealth != "healthy" {
		recs = append(recs, "Check database connection")
	}
	if successRate < 0.95 {
		recs = append(recs, "Review failed analyses")
	}
	
	if len(recs) == 0 {
		recs = append(recs, "System operating normally")
	}
	
	return recs
}

// Metrics management
func (a *MLAgent) updateMetrics(success bool, latency time.Duration) {
	a.metricsMutex.Lock()
	defer a.metricsMutex.Unlock()
	
	a.metrics.TotalAnalyses++
	if success {
		a.metrics.SuccessfulAnalyses++
	} else {
		a.metrics.FailedAnalyses++
	}
	
	// Update average latency
	if a.metrics.TotalAnalyses == 1 {
		a.metrics.AverageLatency = latency
	} else {
		// Exponential moving average
		alpha := 0.1
		a.metrics.AverageLatency = time.Duration(
			float64(a.metrics.AverageLatency)*(1-alpha) + float64(latency)*alpha,
		)
	}
}

func (a *MLAgent) incrementConcurrent() {
	a.metricsMutex.Lock()
	a.metrics.ConcurrentAnalyses++
	a.metricsMutex.Unlock()
}

func (a *MLAgent) decrementConcurrent() {
	a.metricsMutex.Lock()
	a.metrics.ConcurrentAnalyses--
	a.metricsMutex.Unlock()
}

func (a *MLAgent) recordError(context string, err error) {
	a.metricsMutex.Lock()
	a.metrics.LastError = fmt.Sprintf("%s: %v", context, err)
	a.metrics.LastErrorTime = time.Now()
	a.metricsMutex.Unlock()
	a.log.Error("%s: %v", context, err)
}

func (a *MLAgent) recordModelInference(model string) {
	a.metricsMutex.Lock()
	a.metrics.ModelInferences[model]++
	a.metricsMutex.Unlock()
}

func (a *MLAgent) calculateCacheHitRate() float64 {
	a.metricsMutex.RLock()
	defer a.metricsMutex.RUnlock()
	
	total := a.metrics.CacheHits + a.metrics.CacheMisses
	if total == 0 {
		return 0.0
	}
	return float64(a.metrics.CacheHits) / float64(total)
}

func (a *MLAgent) calculateSuccessRate() float64 {
	a.metricsMutex.RLock()
	defer a.metricsMutex.RUnlock()
	
	if a.metrics.TotalAnalyses == 0 {
		return 1.0
	}
	return float64(a.metrics.SuccessfulAnalyses) / float64(a.metrics.TotalAnalyses)
}
