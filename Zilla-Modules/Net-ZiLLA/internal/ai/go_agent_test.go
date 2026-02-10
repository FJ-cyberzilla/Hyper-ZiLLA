package ai

import (
	"net-zilla/internal/models"
	"testing"
)

func TestGoAgent_AnalyzeLink(t *testing.T) {
	agent := NewGoAgent(0.7)

	tests := []struct {
		name       string
		analysis   *models.ThreatAnalysis
		expectSafe bool
	}{
		{
			name: "Safe URL",
			analysis: &models.ThreatAnalysis{
				URL: "https://google.com",
			},
			expectSafe: true,
		},
		{
			name: "Suspicious Punycode",
			analysis: &models.ThreatAnalysis{
				URL: "http://xn--pypal-4ve.com",
			},
			expectSafe: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.AnalyzeLink(tt.analysis)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.IsSafe != tt.expectSafe {
				t.Errorf("expected IsSafe=%v, got %v", tt.expectSafe, result.IsSafe)
			}
		})
	}
}

func TestGoAgent_AnalyzeSMS(t *testing.T) {
	agent := NewGoAgent(0.7)

	tests := []struct {
		name       string
		message    string
		expectSafe bool
	}{
		{
			name:       "Normal message",
			message:    "Hey, are we still meeting for lunch?",
			expectSafe: true,
		},
		{
			name:       "Scam message",
			message:    "URGENT: Your account is suspended. Click here to verify now!",
			expectSafe: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.AnalyzeSMS(tt.message)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.IsSafe != tt.expectSafe {
				t.Errorf("expected IsSafe=%v, got %v", tt.expectSafe, result.IsSafe)
			}
		})
	}
}
