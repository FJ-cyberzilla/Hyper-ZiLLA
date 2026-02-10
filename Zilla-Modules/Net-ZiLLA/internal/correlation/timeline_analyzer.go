package correlation

import (
	"net-zilla/internal/models"
	"time"
)

// TimelineAnalyzer evaluates the timing and sequence of network events.
type TimelineAnalyzer struct{}

func NewTimelineAnalyzer() *TimelineAnalyzer {
	return &TimelineAnalyzer{}
}

// AnalyzeTemporalPatterns checks for automated vs human-like behavior.
func (ta *TimelineAnalyzer) AnalyzeTemporalPatterns(analysis *models.ThreatAnalysis) map[string]string {
	insights := make(map[string]string)

	if len(analysis.RedirectChain) > 1 {
		avgDuration := ta.calculateAvgHopTime(analysis.RedirectChain)

		// Flash redirects (< 50ms) usually indicate automated bot infrastructure
		if avgDuration < 50*time.Millisecond {
			insights["timing_verdict"] = "Automated (Flash Redirects detected)"
		} else {
			insights["timing_verdict"] = "Human-compatible sequence"
		}
	}

	return insights
}

func (ta *TimelineAnalyzer) calculateAvgHopTime(chain []models.RedirectDetail) time.Duration {
	var total time.Duration
	for _, hop := range chain {
		total += hop.Duration
	}
	return total / time.Duration(len(chain))
}
