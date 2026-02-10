package analyzer

import (
	"context"
	"testing"
	"time"

	"net-zilla/internal/models"
	"net-zilla/pkg/logger"
)

func TestThreatAnalyzer_NormalizeURL(t *testing.T) {
	ta := &ThreatAnalyzer{}
	
	tests := []struct {
		raw  string
		want string
	}{
		{"google.com", "https://google.com"},
		{"http://example.com", "http://example.com"},
		{"https://test.org/path", "https://test.org/path"},
	}

	for _, tt := range tests {
		got, err := ta.normalizeURL(tt.raw)
		if err != nil {
			t.Errorf("normalizeURL(%s) error: %v", tt.raw, err)
		}
		if got != tt.want {
			t.Errorf("normalizeURL(%s) = %v, want %v", tt.raw, got, tt.want)
		}
	}
}

func TestThreatAnalyzer_DetermineThreatLevel(t *testing.T) {
	ta := &ThreatAnalyzer{}
	
	tests := []struct {
		score int
		want  models.ThreatLevel
	}{
		{90, models.ThreatLevelCritical},
		{70, models.ThreatLevelHigh},
		{50, models.ThreatLevelMedium},
		{30, models.ThreatLevelLow},
		{10, models.ThreatLevelSafe},
	}

	for _, tt := range tests {
		if got := ta.determineThreatLevel(tt.score); got != tt.want {
			t.Errorf("determineThreatLevel(%d) = %v, want %v", tt.score, got, tt.want)
		}
	}
}

func TestThreatAnalyzer_GenerateSafetyRecommendations(t *testing.T) {
	ta := &ThreatAnalyzer{}
	
	analysis := &models.ThreatAnalysis{ThreatScore: 80}
	ta.generateSafetyRecommendations(analysis)
	
	if len(analysis.SafetyTips) < 3 {
		t.Errorf("expected more safety tips for high score, got %d", len(analysis.SafetyTips))
	}
}

func TestThreatAnalyzer_HelperMethods(t *testing.T) {
	l := logger.NewLogger()
	ta := NewThreatAnalyzer(nil, l, nil)
	
	if id := generateAnalysisID(); id == "" {
		t.Error("expected non-empty analysis ID")
	}

	// Just verify these don't panic and return something
	ctx := context.Background()
	target := "example.com"
	
	_, _ = ta.PerformDNSLookup(ctx, target)
	_, _ = ta.PerformWhoisLookup(ctx, target)
	_, _ = ta.PerformIPGeolocation(ctx, target)
	_, _ = ta.PerformTLSAnalysis(ctx, target)
	_, _ = ta.PerformTraceroute(ctx, target)
}

func TestThreatAnalyzer_ComprehensiveAnalysis(t *testing.T) {
	l := logger.NewLogger()
	ta := NewThreatAnalyzer(nil, l, nil)
	
	ctx := context.Background()
	target := "https://example.com"
	
	analysis, err := ta.ComprehensiveAnalysis(ctx, target)
	if err != nil {
		t.Errorf("ComprehensiveAnalysis failed: %v", err)
	}
	
	if analysis == nil {
		t.Fatal("expected non-nil analysis")
	}
	
	if analysis.URL == "" {
		t.Error("expected non-empty URL")
	}
}

func TestThreatAnalyzer_CircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker(2, time.Second)
	
	if !cb.Allow() {
		t.Error("should allow when closed")
	}
	
	cb.RecordFailure()
	if !cb.Allow() {
		t.Error("should allow after 1 failure")
	}
	
	cb.RecordFailure()
	if cb.Allow() {
		t.Error("should NOT allow after 2 failures")
	}
	
	time.Sleep(1100 * time.Millisecond)
	if !cb.Allow() {
		t.Error("should allow after reset timeout")
	}
	
	cb.RecordSuccess()
	if cb.state != "CLOSED" {
		t.Errorf("expected CLOSED state, got %s", cb.state)
	}
}

func TestThreatAnalyzer_AnalysisCache(t *testing.T) {
	cache := NewAnalysisCache(time.Second)
	analysis := &models.ThreatAnalysis{URL: "test.com"}
	
	cache.SetWithTTL("test.com", analysis, time.Second)
	
	got := cache.Get("test.com")
	if got == nil || got.URL != "test.com" {
		t.Error("cache retrieval failed")
	}
	
	time.Sleep(1500 * time.Millisecond)
	if cache.Get("test.com") != nil {
		t.Error("cache should have expired")
	}
}

func TestThreatAnalyzer_StoreComponentScore(t *testing.T) {
	ta := &ThreatAnalyzer{}
	analysis := &models.ThreatAnalysis{}
	
	ta.storeComponentScore(analysis, 50)
	if analysis.ComponentScores["score_0"] != 50 {
		t.Errorf("expected score 50, got %d", analysis.ComponentScores["score_0"])
	}
}

func TestThreatAnalyzer_CalculateWeightedScore(t *testing.T) {
	ta := NewThreatAnalyzer(nil, logger.NewLogger(), nil)
	analysis := &models.ThreatAnalysis{
		ComponentScores: map[string]int{"score_0": 80},
		RedirectCount:   3,
	}
	
	score := ta.calculateWeightedScore(analysis)
	if score == 0 {
		t.Error("expected non-zero weighted score")
	}
}

func TestThreatAnalyzer_PerformThreatAnalysis(t *testing.T) {
	l := logger.NewLogger()
	ta := NewThreatAnalyzer(nil, l, nil)
	analysis := &models.ThreatAnalysis{}
	
	// Should not panic even with nil components
	score, _ := ta.performThreatAnalysis(context.Background(), "https://example.com", analysis)
	_ = score
}

func TestThreatAnalyzer_PerformMLAnalysis(t *testing.T) {
	ta := NewThreatAnalyzer(nil, logger.NewLogger(), nil)
	analysis := &models.ThreatAnalysis{}
	
	score, _ := ta.performMLAnalysis(context.Background(), "https://example.com", analysis)
	if score != 0 {
		t.Error("expected 0 score for nil ML agent")
	}
}

func TestThreatAnalyzer_ExecuteStandardAnalysis(t *testing.T) {
	ta := NewThreatAnalyzer(nil, logger.NewLogger(), nil)
	analysis := &models.ThreatAnalysis{}
	
	// This will trigger concurrent tasks
	ta.executeStandardAnalysis(context.Background(), "https://example.com", analysis)
	
	if len(analysis.ComponentScores) == 0 {
		t.Error("expected some component scores")
	}
}
