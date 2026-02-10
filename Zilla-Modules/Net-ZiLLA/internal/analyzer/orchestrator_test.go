package analyzer

import (
	"context"
	"testing"

	"net-zilla/internal/config"
	"net-zilla/internal/models"
	"net-zilla/internal/network"
	"net-zilla/pkg/logger"
)

func TestAnalysisOrchestrator_Orchestrate(t *testing.T) {
	l := logger.NewLogger()
	cfg := &config.Config{}
	ao := NewAnalysisOrchestrator(l, cfg)

	ctx := context.Background()
	target := "http://example.com"

	report, err := ao.Orchestrate(ctx, target)
	if err != nil {
		t.Fatalf("Orchestrate failed: %v", err)
	}

	if report.Target != target {
		t.Errorf("expected target %s, got %s", target, report.Target)
	}
}

func TestAnalysisOrchestrator_CalculateFinalScore(t *testing.T) {
	ao := &AnalysisOrchestrator{}
	
	tests := []struct {
		name      string
		screening network.ScreeningResult
		intel     int
		patterns  int
		want      float64
	}{
		{
			name:      "Clean",
			screening: network.ScreeningResult{RiskScore: 0},
			intel:     0,
			patterns:  0,
			want:      0.0,
		},
		{
			name:      "High Risk Screening",
			screening: network.ScreeningResult{RiskScore: 80},
			intel:     0,
			patterns:  0,
			want:      0.8,
		},
		{
			name:      "Combined Risks",
			screening: network.ScreeningResult{RiskScore: 40},
			intel:     1,
			patterns:  1,
			want:      0.9, // 40 + 30 + 20 = 90
		},
		{
			name:      "Capped at 1.0",
			screening: network.ScreeningResult{RiskScore: 60},
			intel:     2,
			patterns:  1,
			want:      1.0, // 60 + 30 + 20 = 110 -> 100
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := &models.AdvancedReport{}
			if tt.intel > 0 {
				report.ThreatIntelligence = &models.IOCRegistry{TotalFound: tt.intel}
			}
			if tt.patterns > 0 {
				report.BehavioralAnalysis = &models.BehaviorAnalysis{
					Patterns: make([]models.BehavioralPattern, tt.patterns),
				}
			}
			if got := ao.calculateFinalScore(tt.screening, report); got != tt.want {
				t.Errorf("calculateFinalScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnalysisOrchestrator_CalculateRiskLevelFromScore(t *testing.T) {
	ao := &AnalysisOrchestrator{}
	tests := []struct {
		score float64
		want  string
	}{
		{0.9, "CRITICAL"},
		{0.6, "HIGH"},
		{0.3, "MEDIUM"},
		{0.1, "LOW"},
	}

	for _, tt := range tests {
		if got := ao.calculateRiskLevelFromScore(tt.score); got != tt.want {
			t.Errorf("calculateRiskLevelFromScore(%v) = %v, want %v", tt.score, got, tt.want)
		}
	}
}

func TestAnalysisOrchestrator_FullAnalysis(t *testing.T) {
	ao := NewAnalysisOrchestrator(logger.NewLogger(), &config.Config{})
	res, err := ao.FullAnalysis(context.Background(), "https://example.com")
	if err != nil {
		t.Errorf("FullAnalysis failed: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil result")
	}
}

func TestAnalysisOrchestrator_GetQuickAnalysis(t *testing.T) {
	ao := NewAnalysisOrchestrator(logger.NewLogger(), &config.Config{})
	res, err := ao.GetQuickAnalysis(context.Background(), "https://example.com")
	if err != nil {
		t.Errorf("GetQuickAnalysis failed: %v", err)
	}
	if res == nil {
		t.Error("expected non-nil result")
	}
}

func TestAnalysisOrchestrator_GetHealth(t *testing.T) {
	ao := NewAnalysisOrchestrator(logger.NewLogger(), &config.Config{})
	if !ao.GetHealth(context.Background()) {
		t.Error("expected healthy status")
	}
}
