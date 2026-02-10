package analyzer_tests

import (
	"context"
	"net-zilla/internal/ai"
	"net-zilla/internal/analyzer"
	"net-zilla/internal/config"
	"net-zilla/internal/storage"
	"net-zilla/pkg/logger"
	"os"
	"testing"
)

func TestSystem_Integration(t *testing.T) {
	// Setup
	l := logger.NewLogger()
	dbPath := "test_integration.db"
	db, err := storage.NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to init db: %v", err)
	}
	defer os.Remove(dbPath)
	defer db.Close()

	cfg := &config.AIConfig{EnableAI: true, ConfidenceThreshold: 0.7}
	orch := analyzer.NewAnalysisOrchestrator(l, &config.Config{})
	mlAgent, _ := ai.NewMLAgent(cfg, orch, db)
	ta := analyzer.NewThreatAnalyzer(mlAgent, l, db)

	// Execute
	ctx := context.Background()
	res, err := ta.ComprehensiveAnalysis(ctx, "https://example.com")
	if err != nil {
		t.Fatalf("ComprehensiveAnalysis failed: %v", err)
	}

	// Verify persistence
	history, err := ta.GetHistory(ctx, 1)
	if err != nil || len(history) == 0 {
		t.Errorf("History not found in database")
	}

	if history[0].AnalysisID != res.AnalysisID {
		t.Errorf("Stored ID mismatch. Expected %s, got %s", res.AnalysisID, history[0].AnalysisID)
	}
}
