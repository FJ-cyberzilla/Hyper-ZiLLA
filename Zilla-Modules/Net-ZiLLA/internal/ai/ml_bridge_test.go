package ai

import (
	"context"
	"net-zilla/internal/config"
	"net-zilla/internal/models"
	"testing"
)

type mockOrchestrator struct{}

func (m *mockOrchestrator) FullAnalysis(ctx context.Context, url string) (*models.AnalysisReport, error) {
	return &models.AnalysisReport{
		RiskScore: 20,
		Findings:  []string{"clean"},
	}, nil
}

func (m *mockOrchestrator) GetQuickAnalysis(ctx context.Context, url string) (*models.AnalysisReport, error) {
	return &models.AnalysisReport{RiskScore: 10}, nil
}

func (m *mockOrchestrator) GetHealth(ctx context.Context) bool {
	return true
}

func TestMLAgent_AnalyzeLink(t *testing.T) {
	cfg := &config.AIConfig{
		EnableAI:            true,
		ConfidenceThreshold: 0.5,
	}
	
	agent, err := NewMLAgent(cfg, &mockOrchestrator{}, nil)
	if err != nil {
		t.Fatalf("failed to create agent: %v", err)
	}

	ta := &models.ThreatAnalysis{
		URL: "https://example.com",
	}

	result, err := agent.AnalyzeLink(context.Background(), ta)
	if err != nil {
		t.Errorf("AnalyzeLink failed: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.RiskLevel == "" {
		t.Error("expected non-empty risk level")
	}
}

func TestMLAgent_OrchestrateAnalysis(t *testing.T) {
	cfg := &config.AIConfig{EnableAI: true}
	agent, _ := NewMLAgent(cfg, &mockOrchestrator{}, nil)

	res, err := agent.OrchestrateAnalysis(context.Background(), "https://example.com", "comprehensive")
	if err != nil {
		t.Errorf("OrchestrateAnalysis failed: %v", err)
	}

	if !res.Success {
		t.Error("expected success")
	}
}

func TestMLAgent_SystemDiagnostics(t *testing.T) {
	cfg := &config.AIConfig{EnableAI: true}
	agent, _ := NewMLAgent(cfg, &mockOrchestrator{}, nil)

	res, err := agent.SystemDiagnostics(context.Background())
	if err != nil {
		t.Errorf("SystemDiagnostics failed: %v", err)
	}

	if !res.Success {
		t.Error("expected success")
	}
}
