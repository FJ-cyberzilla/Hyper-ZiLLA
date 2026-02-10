package correlation

import (
	"net-zilla/internal/models"
	"testing"
)

func TestAnomalyDetector_CalculateEntropy(t *testing.T) {
	ad := NewAnomalyDetector()

	// Low entropy
	low := ad.CalculateEntropy("aaaaa")
	// High entropy
	high := ad.CalculateEntropy("a1b2c3d4e5")

	if high <= low {
		t.Errorf("entropy calculation logic failed: high(%f) <= low(%f)", high, low)
	}
}

func TestAnomalyDetector_FlagAnomalies(t *testing.T) {
	ad := NewAnomalyDetector()

	analysis := &models.ThreatAnalysis{
		URL:           "http://very-long-random-string-that-should-have-high-entropy-1234567890.xyz",
		RedirectCount: 10,
	}

	flags := ad.FlagAnomalies(analysis)
	if len(flags) == 0 {
		t.Error("expected anomalies to be flagged, got none")
	}
}

func TestAnomalyDetector_DetectLookalikes(t *testing.T) {
	ad := NewAnomalyDetector()
	
	tests := []struct {
		domain string
		want   bool
	}{
		{"google.com", false},
		{"paypa111-000.com", true},
	}

	for _, tt := range tests {
		if got := ad.DetectLookalikes(tt.domain); got != tt.want {
			t.Errorf("DetectLookalikes(%s) = %v, want %v", tt.domain, got, tt.want)
		}
	}
}