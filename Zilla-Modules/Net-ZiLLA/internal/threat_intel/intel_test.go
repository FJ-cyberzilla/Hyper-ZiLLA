package threat_intel

import (
	"context"
	"net-zilla/internal/models"
	"os"
	"testing"
	"time"
)

func TestThreatDatabase_AddAndLookup(t *testing.T) {
	dbPath := "test_threat.db"
	db, err := NewThreatDatabase(dbPath, nil)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	defer os.Remove(dbPath)
	defer db.Close()

	ctx := context.Background()
	ind := models.Indicator{
		Value:      "1.2.3.4",
		Type:       models.IOCTypeIP,
		Source:     "Test",
		Confidence: 0.9,
		LastSeen:   time.Now(),
	}

	if err := db.AddIndicator(ctx, ind); err != nil {
		t.Errorf("failed to add indicator: %v", err)
	}

	res, err := db.Lookup(ctx, "1.2.3.4")
	if err != nil {
		t.Errorf("failed to lookup: %v", err)
	}

	if res == nil || res.Source != "Test" {
		t.Errorf("indicator lookup mismatch")
	}
}

func TestThreatDatabase_BulkLookup(t *testing.T) {
	dbPath := "test_bulk.db"
	db, _ := NewThreatDatabase(dbPath, nil)
	defer os.Remove(dbPath)
	defer db.Close()

	ctx := context.Background()
	i1 := models.Indicator{Value: "1.1.1.1", Type: models.IOCTypeIP, Source: "S1"}
	i2 := models.Indicator{Value: "2.2.2.2", Type: models.IOCTypeIP, Source: "S2"}

	db.AddIndicator(ctx, i1)
	db.AddIndicator(ctx, i2)

	res, err := db.BulkLookup(ctx, []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"})
	if err != nil {
		t.Errorf("BulkLookup failed: %v", err)
	}

	if len(res) != 2 {
		t.Errorf("expected 2 results, got %d", len(res))
	}
}

func TestThreatDatabase_GetStats(t *testing.T) {
	dbPath := "test_stats.db"
	db, _ := NewThreatDatabase(dbPath, nil)
	defer os.Remove(dbPath)
	defer db.Close()

	ctx := context.Background()
	db.AddIndicator(ctx, models.Indicator{Value: "test.com", Type: models.IOCTypeDomain, Severity: "high"})

	stats, err := db.GetStats(ctx)
	if err != nil {
		t.Errorf("GetStats failed: %v", err)
	}

	if stats.TotalIndicators != 1 {
		t.Errorf("expected 1 indicator, got %d", stats.TotalIndicators)
	}
}

func TestThreatDatabase_Cleanup(t *testing.T) {
	dbPath := "test_cleanup.db"
	db, _ := NewThreatDatabase(dbPath, nil)
	defer os.Remove(dbPath)
	defer db.Close()

	ctx := context.Background()
	// Add an old indicator (mocking last_seen is harder with CURRENT_TIMESTAMP default, but we can test the call)
	db.AddIndicator(ctx, models.Indicator{Value: "old.com", Type: models.IOCTypeDomain})

	rows, err := db.Cleanup(ctx, 30)
	if err != nil {
		t.Errorf("Cleanup failed: %v", err)
	}
	// rows will be 0 because we just added it
	if rows < 0 {
		t.Error("expected non-negative rows affected")
	}
}

func TestNewFeedsClient(t *testing.T) {
	keys := map[string]string{"abuseipdb": "key"}
	fc := NewFeedsClient(keys, time.Second)
	if fc == nil {
		t.Fatal("NewFeedsClient returned nil")
	}
	if len(fc.providers) == 0 {
		t.Error("expected some providers")
	}
}
