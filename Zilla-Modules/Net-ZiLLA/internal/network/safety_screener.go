package network

import (
	"net/url"
	"regexp"
	"strings"
)

type SafetyScreener struct {
	dangerousTLDs []string
}

func NewSafetyScreener() *SafetyScreener {
	return &SafetyScreener{
		dangerousTLDs: []string{".tk", ".ml", ".ga", ".cf", ".gq", ".xyz", ".top", ".buzz", ".work"},
	}
}

type ScreeningResult struct {
	IsSuspicious bool
	Reasons      []string
	RiskScore    int
}

func (ss *SafetyScreener) Screen(rawURL string) ScreeningResult {
	result := ScreeningResult{Reasons: []string{}}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		result.IsSuspicious = true
		result.Reasons = append(result.Reasons, "Invalid URL structure")
		result.RiskScore = 100
		return result
	}

	host := strings.ToLower(parsed.Hostname())

	// 1. Check for Hex/Octal/IP encoding
	ipPattern := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	if ipPattern.MatchString(host) {
		result.Reasons = append(result.Reasons, "Uses direct IP address instead of domain")
		result.RiskScore += 40
	}

	// 2. Check for suspicious TLDs
	for _, tld := range ss.dangerousTLDs {
		if strings.HasSuffix(host, tld) {
			result.Reasons = append(result.Reasons, "High-risk TLD detected: "+tld)
			result.RiskScore += 35 // Elevated score to trigger suspicious flag
		}
	}

	// 3. Check for excessive subdomains or length
	if strings.Count(host, ".") > 4 {
		result.Reasons = append(result.Reasons, "Excessive subdomain nesting")
		result.RiskScore += 20
	}

	// 4. Check for Punycode (homograph attacks)
	if strings.HasPrefix(host, "xn--") || strings.Contains(host, ".xn--") {
		result.Reasons = append(result.Reasons, "Punycode detected (Potential homograph attack)")
		result.RiskScore += 50
	}

	// 5. Check for obfuscation patterns in path
	if strings.Contains(rawURL, "%") {
		result.Reasons = append(result.Reasons, "URL encoding detected")
		result.RiskScore += 10
	}
	if strings.Contains(rawURL, `\`) {
		result.Reasons = append(result.Reasons, "Backslash obfuscation detected")
		result.RiskScore += 10
	}

	result.IsSuspicious = result.RiskScore > 30
	return result
}
