package storage

import (
	"context"
	"net-zilla/internal/models"
	"os"
	"testing"
)

func TestDatabase_SaveAndRetrieve(t *testing.T) {
	dbPath := "test_db.db"
	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	defer os.Remove(dbPath)
	defer db.Close()

	analysis := &models.ThreatAnalysis{
		AnalysisID:  "test-id",
		URL:         "https://test.com",
		ThreatScore: 45,
		ThreatLevel: "MEDIUM",
	}

	ctx := context.Background()
	if err := db.SaveAnalysis(ctx, analysis); err != nil {
		t.Errorf("failed to save analysis: %v", err)
	}

	history, err := db.GetAnalysisHistory(ctx, 1)
	if err != nil {
		t.Errorf("failed to get history: %v", err)
	}

	if len(history) == 0 {
		t.Fatal("history is empty")
	}

	if history[0].URL != analysis.URL {
		t.Errorf("expected URL %s, got %s", analysis.URL, history[0].URL)
	}
}
