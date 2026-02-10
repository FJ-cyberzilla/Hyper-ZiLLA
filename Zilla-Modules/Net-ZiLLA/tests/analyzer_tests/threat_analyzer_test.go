package analyzer_tests

import (
	"context"
	"net-zilla/internal/ai"
	"net-zilla/internal/analyzer"
	"net-zilla/internal/config"
	"net-zilla/pkg/logger"
	"testing"
)

func TestThreatAnalyzer_ComprehensiveAnalysis(t *testing.T) {
	l := logger.NewLogger()
	cfg := &config.AIConfig{EnableAI: true, ConfidenceThreshold: 0.7}
	mlAgent, _ := ai.NewMLAgent(cfg, nil, nil)

	// db is nil for unit test
	ta := analyzer.NewThreatAnalyzer(mlAgent, l, nil)

	ctx := context.Background()
	target := "https://google.com"

	analysis, err := ta.ComprehensiveAnalysis(ctx, target)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if analysis.URL != target {
		t.Errorf("Expected URL %s, got %s", target, analysis.URL)
	}

	if analysis.ThreatScore < 0 || analysis.ThreatScore > 100 {
		t.Errorf("Invalid threat score: %d", analysis.ThreatScore)
	}
}
