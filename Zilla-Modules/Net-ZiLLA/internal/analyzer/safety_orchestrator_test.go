package analyzer

import (
	"context"
	"testing"

	"net-zilla/pkg/logger"
)

func TestSafetyOrchestrator_SecureAnalyze(t *testing.T) {
	so := NewSafetyOrchestrator(logger.NewLogger())
	
	report, err := so.SecureAnalyze(context.Background(), "http://example.com")
	if err != nil {
		t.Fatalf("SecureAnalyze failed: %v", err)
	}

	if report.Target != "http://example.com" {
		t.Errorf("expected target example.com, got %s", report.Target)
	}
}

func TestSafetyOrchestrator_SecureAnalyze_HighRisk(t *testing.T) {
	so := NewSafetyOrchestrator(logger.NewLogger())
	
	// Target that triggers high risk screening
	target := "http://xn--pypal-4ve.com" 
	report, err := so.SecureAnalyze(context.Background(), target)
	if err != nil {
		t.Fatalf("SecureAnalyze failed: %v", err)
	}

	if report.Target != target {
		t.Errorf("expected target %s, got %s", target, report.Target)
	}
}

func TestNewSafetyOrchestrator(t *testing.T) {
	l := logger.NewLogger()
	so := NewSafetyOrchestrator(l)
	if so == nil {
		t.Fatal("NewSafetyOrchestrator returned nil")
	}
}
