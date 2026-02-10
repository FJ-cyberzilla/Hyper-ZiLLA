package services

import (
	"context"
	"net-zilla/internal/config"
	"net-zilla/internal/models"
	"net-zilla/internal/storage"
	"net-zilla/pkg/logger"
	"os"
	"testing"
)

func TestAnalysisService_PerformAnalysis(t *testing.T) {
	l := logger.NewLogger()
	cfg := &config.Config{}
	
	// Use a temporary database file
	dbPath := "test_service.db"
	db, err := storage.NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer func() {
		db.Close()
		os.Remove(dbPath)
	}()

	svc := NewAnalysisService(l, db, cfg)

	ctx := context.Background()
	target := "http://example.com"

	report, err := svc.PerformAnalysis(ctx, target)
	if err != nil {
		t.Fatalf("PerformAnalysis failed: %v", err)
	}

	if report == nil {
		t.Fatal("Report is nil")
	}

	if report.Target != target {
		t.Errorf("Expected target %s, got %s", target, report.Target)
	}

	// Verify persistence
	history, err := svc.GetAnalysisHistory(ctx, 10)
	if err != nil {
		t.Fatalf("GetAnalysisHistory failed: %v", err)
	}

	if len(history) == 0 {
		t.Error("No history found in database")
	} else if history[0].URL != target {
		t.Errorf("Expected history URL %s, got %s", target, history[0].URL)
	}
}

func TestAnalysisService_CacheAndLocks(t *testing.T) {
	l := logger.NewLogger()
	cfg := &config.Config{}
	svc := NewAnalysisService(l, nil, cfg)
	
	target := "http://test-cache.com"
	
	// Test Lock
	acquired, lockID := svc.acquireLock(target)
	if !acquired {
		t.Fatal("failed to acquire lock")
	}
	
	if !svc.isAnalysisInProgress(target) {
		t.Error("analysis should be in progress")
	}
	
	svc.releaseLock(target, lockID)
	if svc.isAnalysisInProgress(target) {
		t.Error("analysis should NOT be in progress after release")
	}
	
	// Test Cache
	report := &models.AdvancedReport{Target: target}
	svc.addToCache(target, report)
	
	cached := svc.getFromCache(target)
	if cached == nil || cached.Target != target {
		t.Error("cache retrieval failed")
	}
	
	svc.ClearCache()
	if svc.getFromCache(target) != nil {
		t.Error("cache should be empty after clear")
	}
}

func TestAnalysisService_HealthCheck(t *testing.T) {
	l := logger.NewLogger()
	svc := NewAnalysisService(l, nil, &config.Config{})
	
	health := svc.HealthCheck(context.Background())
	if health["status"] != "healthy" {
		t.Errorf("expected healthy status, got %v", health["status"])
	}
}

func TestAnalysisService_GetAnalysisByID(t *testing.T) {
	l := logger.NewLogger()
	dbPath := "test_id.db"
	db, _ := storage.NewDatabase(dbPath)
	defer os.Remove(dbPath)
	defer db.Close()
	
	svc := NewAnalysisService(l, db, &config.Config{})
	
	analysis := &models.ThreatAnalysis{
		AnalysisID: "test-id-123",
		URL:        "http://test.com",
	}
	db.SaveAnalysis(context.Background(), analysis)
	
	res, err := svc.GetAnalysisByID(context.Background(), "test-id-123")
	if err != nil {
		t.Errorf("GetAnalysisByID failed: %v", err)
	}
	if res == nil || res.URL != "http://test.com" {
		t.Error("analysis retrieval failed")
	}
}

func TestAnalysisService_InstanceID(t *testing.T) {
	svc := NewAnalysisService(logger.NewLogger(), nil, &config.Config{})
	if svc.GetInstanceID() == "" {
		t.Error("expected non-empty instance ID")
	}
}
