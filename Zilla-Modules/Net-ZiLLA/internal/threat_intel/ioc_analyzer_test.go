package threat_intel

import (
	"context"
	"net-zilla/internal/models"
	"net-zilla/pkg/logger"
	"testing"
)

func TestIOCAnalyzer(t *testing.T) {
	db, _ := NewThreatDatabase("test_ioc.db", nil)
	defer db.Close()
	
	l := logger.NewLogger()
	analyzer := NewIOCAnalyzer(db, nil, l)
	
	ctx := context.Background()
	
	// Test IP
	db.AddIndicator(ctx, models.Indicator{Value: "8.8.8.8", Type: models.IOCTypeIP, Severity: "low"})
	res1, err := analyzer.Analyze(ctx, "8.8.8.8")
	if err != nil {
		t.Errorf("Analyze failed: %v", err)
	}
	if res1 == nil {
		t.Error("expected result for known IP")
	}
	
	// Test unknown
	res2, err := analyzer.Analyze(ctx, "unknown")
	if err != nil {
		t.Errorf("Analyze failed: %v", err)
	}
	if res2 != nil {
		t.Error("expected nil result for unknown")
	}
}
