package analyzer_interface

import (
	"context"
	"net-zilla/internal/models"
	"net-zilla/internal/shared_models"
)

type MLAgentInterface interface {
	AnalyzeLink(ctx context.Context, threatAnalysis *models.ThreatAnalysis) (*shared_models.AIAnalysisResult, error)
	OrchestrateAnalysis(ctx context.Context, url string, analysisType string) (*shared_models.OrchestrationResult, error)
}

type OrchestratorInterface interface {
	FullAnalysis(ctx context.Context, url string) (*models.AnalysisReport, error)
	GetQuickAnalysis(ctx context.Context, url string) (*models.AnalysisReport, error)
	GetHealth(ctx context.Context) bool
}
