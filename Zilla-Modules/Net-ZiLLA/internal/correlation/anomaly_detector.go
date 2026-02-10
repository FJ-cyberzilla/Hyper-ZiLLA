package correlation

import (
	"math"
	"net-zilla/internal/models"
	"strings"
)

// AnomalyDetector identifies behavioral and structural outliers.
type AnomalyDetector struct{}

func NewAnomalyDetector() *AnomalyDetector {
	return &AnomalyDetector{}
}

// DetectLookalikes checks for domain squatting/bit-squatting.
func (ad *AnomalyDetector) DetectLookalikes(domain string) bool {
	// Professional heuristic: Check for excessive hyphens or number substitutions
	// E.g., paypa1.com vs paypal.com
	confusingChars := []string{"1", "0", "l", "o", "-"}
	count := 0
	for _, char := range confusingChars {
		count += strings.Count(domain, char)
	}

	return count > 5
}

// CalculateEntropy measures the randomness of a string (URL/Domain).
// High entropy often correlates with DGA (Domain Generation Algorithms).
func (ad *AnomalyDetector) CalculateEntropy(s string) float64 {
	m := make(map[rune]float64)
	for _, r := range s {
		m[r]++
	}
	var entropy float64
	lenStr := float64(len(s))
	for _, count := range m {
		p := count / lenStr
		entropy -= p * math.Log2(p)
	}
	return entropy
}

// FlagAnomalies runs all detection algorithms on the analysis.
func (ad *AnomalyDetector) FlagAnomalies(ba *models.ThreatAnalysis) []string {
	var flags []string

	if ad.CalculateEntropy(ba.URL) > 4.5 {
		flags = append(flags, "High Entropy URL (Potential DGA)")
	}

	if ba.RedirectCount > 5 {
		flags = append(flags, "Unusually long redirect chain")
	}

	return flags
}
