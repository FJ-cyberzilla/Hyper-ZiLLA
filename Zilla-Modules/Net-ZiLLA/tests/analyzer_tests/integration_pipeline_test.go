package analyzer_tests

import (
	"context"
	"net-zilla/internal/config"
	"net-zilla/internal/services"
	"net-zilla/internal/storage"
	"net-zilla/pkg/logger"
	"os"
	"testing"
)

func TestFullAnalysisPipeline(t *testing.T) {
	// 1. Setup Environment
	l := logger.NewLogger()
	cfg := &config.Config{}
	dbPath := "test_full_pipeline.db"
	db, err := storage.NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to init database: %v", err)
	}
	defer os.Remove(dbPath)
	defer db.Close()

	// 2. Initialize Service (Brain)
	service := services.NewAnalysisService(l, db, cfg)

	// 3. Run Pipeline for a known target
	ctx := context.Background()
	target := "http://xn--pypal-4ve.com/secure-login" // Suspicious punycode target

	report, err := service.PerformAnalysis(ctx, target)
	if err != nil {
		t.Fatalf("Pipeline execution failed: %v", err)
	}

	// 4. Assertions
	if report.Target != target {
		t.Errorf("Target mismatch: expected %s, got %s", target, report.Target)
	}

	// Verify Static Screener flagged it (it should have Punycode + HTTP)
	if report.RiskAssessment.RiskScore < 0.2 {
		t.Errorf("Risk score should be elevated for suspicious punycode target, got %.2f", report.RiskAssessment.RiskScore)
	}

	// Verify persistence in database
	history, err := service.GetAnalysisHistory(ctx, 1)
	if err != nil || len(history) == 0 {
		t.Fatalf("History not saved to database")
	}

	if history[0].URL != target {
		t.Errorf("Database record mismatch: expected %s, got %s", target, history[0].URL)
	}
}
