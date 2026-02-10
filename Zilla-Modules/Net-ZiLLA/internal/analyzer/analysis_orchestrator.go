package analyzer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"net-zilla/internal/config"
	"net-zilla/internal/correlation"
	"net-zilla/internal/models"
	"net-zilla/internal/network"
	"net-zilla/internal/patterns"
	"net-zilla/internal/threat_intel"
	"net-zilla/pkg/logger"
)

// AnalysisOrchestrator manages the execution flow of a security analysis.
type AnalysisOrchestrator struct {
	logger     *logger.Logger
	screener   *network.SafetyScreener
	intel      *threat_intel.IntelManager
	matcher    *patterns.PatternMatcher
	correlator *correlation.EventCorrelator
	sandbox    *threat_intel.SandboxManager
}

func NewAnalysisOrchestrator(l *logger.Logger, cfg *config.Config) *AnalysisOrchestrator {
	return &AnalysisOrchestrator{
		logger:     l,
		screener:   network.NewSafetyScreener(),
		intel:      threat_intel.NewIntelManager("", "", ""), // Keys would come from cfg in real prod
		matcher:    patterns.NewPatternMatcher(),
		correlator: correlation.NewEventCorrelator(),
		sandbox:    threat_intel.NewSandboxManager(),
	}
}

// Orchestrate runs the multi-stage analysis pipeline concurrently.
func (ao *AnalysisOrchestrator) Orchestrate(ctx context.Context, target string) (*models.AdvancedReport, error) {
	// Global timeout for the entire orchestration to prevent hanging
	ctx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	ao.logger.Info("Initializing orchestration for: %s", target)

	report := &models.AdvancedReport{
		ReportID:  fmt.Sprintf("NZ-%d", time.Now().Unix()),
		Target:    target,
		Timestamp: time.Now(),
		Metadata:  make(map[string]interface{}),
	}

	// STAGE 1: Safety Screening (Synchronous as it is fast and foundational)
	screening := ao.screener.Screen(target)
	report.Metadata["screening_risk_score"] = fmt.Sprintf("%d", screening.RiskScore)

	var wg sync.WaitGroup
	wg.Add(2)

	// STAGE 2: Passive Intelligence (Concurrent)
	go func() {
		defer wg.Done()
		intelResult := ao.intel.MultiCheck(ctx, target)
		report.ThreatIntelligence = &models.IOCRegistry{
			TotalFound: intelResult.Positives,
		}
		if intelResult.Malicious {
			report.Metadata["intel_malicious"] = "true"
		}
	}()

	// STAGE 3: Behavioral Pattern Matching (Concurrent)
	go func() {
		defer wg.Done()
		behavior := ao.matcher.AnalyzeContent(target)
		report.BehavioralAnalysis = behavior
	}()

	// Wait for concurrent stages with timeout protection
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All concurrent stages finished
	case <-ctx.Done():
		return nil, fmt.Errorf("orchestration timed out: %w", ctx.Err())
	}

	// STAGE 4: Risk-Based Escalation (Sandbox)
	// Only escalate if initial findings are highly suspicious
	if screening.RiskScore > 60 || (report.ThreatIntelligence != nil && report.ThreatIntelligence.TotalFound > 0) {
		ao.logger.Info("Risk threshold exceeded. Escalating to isolated sandbox...")
		containerID, err := ao.sandbox.SpinUpIsolatedBrowser(ctx, target)
		if err == nil {
			report.Metadata["sandbox_escalated"] = "true"
			report.Metadata["sandbox_container_id"] = containerID
			// In production, we'd have a cleanup worker, but for now we defer
			defer ao.sandbox.DestroySandbox(containerID)
		}
	}

	// STAGE 5: Correlation
	ao.correlator.Correlate(report)

	// STAGE 6: Final Risk Assessment
	report.RiskAssessment = &models.RiskAssessment{
		RiskScore:        ao.calculateFinalScore(screening, report),
		OverallRiskLevel: ao.calculateRiskLevelFromScore(ao.calculateFinalScore(screening, report)),
		Summary:          "Analysis completed through concurrent production pipeline.",
	}

	return report, nil
}

func (ao *AnalysisOrchestrator) calculateFinalScore(s network.ScreeningResult, r *models.AdvancedReport) float64 {
	base := float64(s.RiskScore)
	if r.ThreatIntelligence != nil && r.ThreatIntelligence.TotalFound > 0 {
		base += 30
	}
	if r.BehavioralAnalysis != nil && len(r.BehavioralAnalysis.Patterns) > 0 {
		base += 20
	}
	if base > 100 {
		base = 100
	}
	return base / 100.0
}

func (ao *AnalysisOrchestrator) calculateRiskLevelFromScore(score float64) string {
	s := score * 100
	switch {
	case s >= 80:
		return "CRITICAL"
	case s >= 50:
		return "HIGH"
	case s >= 20:
		return "MEDIUM"
	default:
		return "LOW"
	}
}

// FullAnalysis provides a comprehensive analysis for the ML bridge
func (ao *AnalysisOrchestrator) FullAnalysis(ctx context.Context, url string) (*models.AnalysisReport, error) {
	report, err := ao.Orchestrate(ctx, url)
	if err != nil {
		return nil, err
	}

	return &models.AnalysisReport{
		RiskScore:   int(report.RiskAssessment.RiskScore * 100),
		Findings:    report.Findings,
		SandboxUsed: report.Metadata["sandbox_escalated"] == "true",
	}, nil
}

// GetQuickAnalysis provides a fast analysis for fallback or preview
func (ao *AnalysisOrchestrator) GetQuickAnalysis(ctx context.Context, url string) (*models.AnalysisReport, error) {
	screening := ao.screener.Screen(url)
	return &models.AnalysisReport{
		RiskScore: screening.RiskScore,
		Findings:  screening.Reasons,
	}, nil
}

// GetHealth returns the health status of the orchestrator
func (ao *AnalysisOrchestrator) GetHealth(ctx context.Context) bool {
	return true
}

