package ai

import (
	"math"
	"net/url"
	"regexp"
	"strings"

	"net-zilla/internal/models"
	"net-zilla/internal/shared_models"
)

// GoAgent implements advanced analysis logic in native Go.
type GoAgent struct {
	confidenceThreshold float64
}

// NewGoAgent creates a new GoAgent.
func NewGoAgent(confidenceThreshold float64) *GoAgent {
	return &GoAgent{
		confidenceThreshold: confidenceThreshold,
	}
}

// AnalyzeLink performs native Go-based analysis using extracted features and component results.
func (g *GoAgent) AnalyzeLink(analysis *models.ThreatAnalysis) (*shared_models.AIAnalysisResult, error) {
	features := g.extractFeatures(analysis)

	// Calculate a heuristic-based "health score" from actual component data
	healthScore := g.calculateHealthScore(features)

	ipRisk := 0.1
	if features.IsProxy {
		ipRisk = 0.8
	}

	overallConfidence := (healthScore + (1.0 - ipRisk)) / 2.0
	isSafe := overallConfidence > 0.7

	riskLevel := "LOW"
	if overallConfidence < 0.3 {
		riskLevel = "CRITICAL"
	} else if overallConfidence < 0.5 {
		riskLevel = "HIGH"
	} else if overallConfidence < 0.7 {
		riskLevel = "MEDIUM"
	}

	threats := g.generateThreats(features, healthScore, ipRisk)
	recommendations := g.generateRecommendations(isSafe, riskLevel, features.IsShortened)

	result := &shared_models.AIAnalysisResult{
		IsSafe:          isSafe,
		Confidence:      overallConfidence,
		RiskLevel:       riskLevel,
		IsShortened:     features.IsShortened,
		HealthScore:     healthScore,
		Threats:         threats,
		Recommendations: recommendations,
	}

	if result.Confidence < g.confidenceThreshold && !isSafe {
		result.RiskLevel = "UNKNOWN (Low Confidence)"
		result.Threats = append(result.Threats, "AI analysis confidence below threshold")
	}

	return result, nil
}

type AdvancedFeatures struct {
	URLLength       int
	NumSpecialChars int
	HasIP           bool
	HasRedirect     bool
	TLDRisk         float64
	Entropy         float64
	KeywordMatches  int
	DomainAgeDays   int
	SSLValid        bool
	IsProxy         bool
	IsShortened     bool
	IsPunycode      bool
}

func (g *GoAgent) extractFeatures(a *models.ThreatAnalysis) AdvancedFeatures {
	parsed, _ := url.Parse(a.URL)
	lowerURL := strings.ToLower(a.URL)
	host := ""
	if parsed != nil {
		host = strings.ToLower(parsed.Hostname())
	}

	f := AdvancedFeatures{
		URLLength:       len(a.URL),
		NumSpecialChars: strings.Count(lowerURL, "@") + strings.Count(lowerURL, "%") + strings.Count(lowerURL, "&") + strings.Count(lowerURL, "=") + strings.Count(lowerURL, "?") + strings.Count(lowerURL, "#"),
		SSLValid:        true,
		IsPunycode:      strings.HasPrefix(host, "xn--") || strings.Contains(host, ".xn--"),
	}

	ipRegex := regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)
	f.HasIP = ipRegex.MatchString(a.URL)

	redirectKeywords := []string{"redirect=", "url=", "goto=", "next="}
	for _, kw := range redirectKeywords {
		if strings.Contains(lowerURL, kw) {
			f.HasRedirect = true
			break
		}
	}

	f.TLDRisk = g.assessTLDRisk(parsed)
	f.Entropy = g.calculateEntropy(a.URL)

	suspiciousKeywords := []string{"login", "verify", "account", "secure", "bank", "paypal", "update", "confirm"}
	for _, kw := range suspiciousKeywords {
		if strings.Contains(lowerURL, kw) {
			f.KeywordMatches++
		}
	}

	if f.IsPunycode {
		f.KeywordMatches += 5 // Force high suspicion
	}

	if a.WhoisInfo != nil {
		if strings.Contains(a.WhoisInfo.DomainAge, "years") {
			f.DomainAgeDays = 365 * 2
		} else if strings.Contains(a.WhoisInfo.DomainAge, "days") {
			f.DomainAgeDays = 10
		} else {
			f.DomainAgeDays = 365
		}
	} else {
		f.DomainAgeDays = 365
	}

	if a.TLSInfo != nil {
		f.SSLValid = a.TLSInfo.CertificateValid
	}

	if a.GeoAnalysis != nil {
		f.IsProxy = a.GeoAnalysis.IsProxy
	}

	shorteners := []string{"bit.ly", "tinyurl.com", "goo.gl", "ow.ly", "t.co"}
	for _, s := range shorteners {
		if strings.Contains(lowerURL, s) {
			f.IsShortened = true
			break
		}
	}

	return f
}

func (g *GoAgent) assessTLDRisk(u *url.URL) float64 {
	if u == nil {
		return 0.1
	}
	highRisk := []string{".tk", ".ml", ".ga", ".cf", ".gq", ".xyz", ".top"}
	host := strings.ToLower(u.Hostname())
	for _, tld := range highRisk {
		if host == tld[1:] || strings.HasSuffix(host, "."+tld[1:]) {
			return 0.9
		}
	}
	return 0.1
}

func (g *GoAgent) calculateEntropy(s string) float64 {
	counts := make(map[rune]int)
	for _, r := range s {
		counts[r]++
	}
	entropy := 0.0
	length := float64(len(s))
	for _, count := range counts {
		p := float64(count) / length
		entropy -= p * math.Log2(p)
	}
	return entropy
}

func (g *GoAgent) calculateHealthScore(f AdvancedFeatures) float64 {
	score := 1.0
	if f.HasIP {
		score -= 0.3
	}
	if f.IsShortened {
		score -= 0.2
	}
	if f.TLDRisk > 0.5 {
		score -= 0.2
	}
	if f.DomainAgeDays < 30 {
		score -= 0.3
	}
	if !f.SSLValid {
		score -= 0.4
	}
	if f.IsProxy {
		score -= 0.2
	}
	if f.IsPunycode {
		score -= 0.6
	} // High penalty
	if f.KeywordMatches > 3 {
		score -= 0.2
	}

	if score < 0 {
		score = 0
	}
	return score
}

func (g *GoAgent) generateThreats(f AdvancedFeatures, healthScore, ipRisk float64) []string {
	var threats []string
	if f.HasIP {
		threats = append(threats, "Direct IP access detected")
	}
	if f.IsShortened {
		threats = append(threats, "URL shortener obscures destination")
	}
	if f.TLDRisk > 0.7 {
		threats = append(threats, "High-risk top-level domain")
	}
	if f.DomainAgeDays < 30 && f.DomainAgeDays > 0 {
		threats = append(threats, "Very recently registered domain")
	}
	if !f.SSLValid {
		threats = append(threats, "Invalid or missing SSL/TLS certificate")
	}
	if f.IsProxy {
		threats = append(threats, "Accessing through a known proxy/VPN")
	}
	if f.IsPunycode {
		threats = append(threats, "Punycode (homograph) domain detected")
	}
	return threats
}

func (g *GoAgent) generateRecommendations(isSafe bool, riskLevel string, isShortened bool) []string {
	var recs []string
	if !isSafe {
		recs = append(recs, "DO NOT visit this link", "Delete the message")
	}
	if isShortened {
		recs = append(recs, "Verify original destination before clicking")
	}
	recs = append(recs, "Always look for the official domain name")
	return recs
}

func (g *GoAgent) AnalyzeSMS(message string) (*shared_models.AIAnalysisResult, error) {
	lower := strings.ToLower(message)
	confidence := 0.6
	isScam := false
	var threats []string

	urgency := []string{"urgent", "immediate", "action required", "locked", "suspended"}
	for _, w := range urgency {
		if strings.Contains(lower, w) {
			isScam = true
			confidence += 0.1
			threats = append(threats, "Urgent/Threatening language")
			break
		}
	}

	lures := []string{"prize", "winner", "money", "refund", "gift"}
	for _, w := range lures {
		if strings.Contains(lower, w) {
			isScam = true
			confidence += 0.1
			threats = append(threats, "Financial lure")
			break
		}
	}

	risk := "LOW"
	if isScam {
		if confidence > 0.8 {
			risk = "HIGH"
		} else {
			risk = "MEDIUM"
		}
	}

	return &shared_models.AIAnalysisResult{
		IsSafe:          !isScam,
		Confidence:      confidence,
		RiskLevel:       risk,
		HealthScore:     1.0 - (confidence * 0.5),
		Threats:         threats,
		Recommendations: []string{"Do not reply to suspicious messages", "Contact the company through official channels"},
	}, nil
}
