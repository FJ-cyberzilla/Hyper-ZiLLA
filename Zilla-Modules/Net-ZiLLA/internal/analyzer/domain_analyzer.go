package analyzer

import (
	"context"
	"net/url"
	"strings"
	"sync"
	"time"

	"net-zilla/internal/models"
	"net-zilla/internal/network"
	"net-zilla/internal/processor"
	"net-zilla/pkg/logger"
)

// DomainAnalyzer is responsible for performing domain-specific analyses.
type DomainAnalyzer struct {
	logger      *logger.Logger
	dnsClient   *network.DNSClient
	whoisClient *network.WhoisClient
	urlParser   *processor.URLParser
	
	// Simple cache to avoid redundant lookups
	dnsCache   map[string]*models.DNSAnalysis
	whoisCache map[string]*models.WhoisAnalysis
	cacheMutex sync.RWMutex
}

// NewDomainAnalyzer creates a new DomainAnalyzer.
func NewDomainAnalyzer(logger *logger.Logger, dnsClient *network.DNSClient, whoisClient *network.WhoisClient) *DomainAnalyzer {
	return &DomainAnalyzer{
		logger:      logger,
		dnsClient:   dnsClient,
		whoisClient: whoisClient,
		urlParser:   processor.NewURLParser(),
		dnsCache:    make(map[string]*models.DNSAnalysis),
		whoisCache:  make(map[string]*models.WhoisAnalysis),
	}
}

// Analyze performs basic domain analysis and enriches the ThreatAnalysis object.
func (da *DomainAnalyzer) Analyze(ctx context.Context, parsedURL *url.URL, analysis *models.ThreatAnalysis) (int, error) {
	score := 0
	
	// Extract hostname
	hostname := parsedURL.Hostname()
	if hostname == "" {
		return 0, nil
	}
	
	// 1. Analyze URL structure
	score += da.analyzeURLStructure(parsedURL, analysis)
	
	// 2. Check WHOIS for domain age
	if analysis.WhoisInfo == nil {
		// Check cache first
		da.cacheMutex.RLock()
		cachedWhois, found := da.whoisCache[hostname]
		da.cacheMutex.RUnlock()
		
		if found && time.Since(cachedWhois.LastUpdated) < time.Hour {
			analysis.WhoisInfo = cachedWhois
		} else {
			whoisInfo, err := da.whoisClient.Lookup(ctx, hostname)
			if err == nil {
				analysis.WhoisInfo = whoisInfo
				// Cache it
				da.cacheMutex.Lock()
				da.whoisCache[hostname] = whoisInfo
				da.cacheMutex.Unlock()
				
				// Analyze WHOIS data
				if whoisInfo.DomainAge == "Unknown" || whoisInfo.DomainAge == "Less than 30 days" {
					score += 10
					analysis.Warnings = append(analysis.Warnings, "Domain is very new, potential risk")
				}
				
				// Check for suspicious registrar
				if da.isSuspiciousRegistrar(whoisInfo.Registrar) {
					score += 5
					analysis.Warnings = append(analysis.Warnings, "Suspicious registrar detected")
				}
			} else {
				da.logger.Warn("Failed to perform WHOIS lookup for domain analyzer: %v", err)
				score += 5
			}
		}
	}
	
	// 3. Check DNS for suspicious records
	if analysis.DNSInfo == nil {
		// Check cache first
		da.cacheMutex.RLock()
		cachedDNS, found := da.dnsCache[hostname]
		da.cacheMutex.RUnlock()
		
		if found && time.Since(cachedDNS.LastUpdated) < 5*time.Minute {
			analysis.DNSInfo = cachedDNS
		} else {
			dnsInfo, err := da.dnsClient.Lookup(ctx, hostname)
			if err == nil {
				analysis.DNSInfo = dnsInfo
				// Cache it
				da.cacheMutex.Lock()
				da.dnsCache[hostname] = dnsInfo
				da.cacheMutex.Unlock()
				
				// Analyze DNS records
				if len(dnsInfo.TXTRecords) == 0 {
					score += 2
					analysis.Warnings = append(analysis.Warnings, "No TXT records found")
				}
				
				if len(dnsInfo.MXRecords) == 0 {
					score += 3
					analysis.Warnings = append(analysis.Warnings, "No MX records found (no email servers)")
				}
				
				// Check for suspicious CNAME records
				for _, cname := range dnsInfo.CNAMERecords {
					if da.isSuspiciousCNAME(cname) {
						score += 10
						analysis.Warnings = append(analysis.Warnings, 
							"CNAME points to suspicious domain: "+cname)
					}
				}
			} else {
				da.logger.Warn("Failed to perform DNS lookup for domain analyzer: %v", err)
				score += 3
			}
		}
	}
	
	// 4. Analyze domain name for suspicious patterns
	score += da.analyzeDomainName(hostname, analysis)
	
	// 5. Check TLD risk
	score += da.analyzeTLD(hostname, analysis)
	
	// 6. Check for homograph attacks
	if da.isHomographAttack(hostname) {
		score += 30
		analysis.Warnings = append(analysis.Warnings, 
			"Possible homograph attack detected (Punycode/IDN)")
	}

	// 7. Perform URL enrichment
	analysis.URLEnrichment = da.EnrichURL(parsedURL)
	
	// Cap score at 100
	if score > 100 {
		score = 100
	}
	
	return score, nil
}

func (da *DomainAnalyzer) EnrichURL(u *url.URL) *models.URLEnrichment {
	enrichment := &models.URLEnrichment{
		Entropy:         da.urlParser.CalculateEntropy(u.String()),
		HomographAttack: da.isHomographAttack(u.Hostname()),
	}

	// Calculate TLD risk
	tld := da.extractTLD(u.Hostname())
	enrichment.TLDRisk = da.getTLDRiskScore(tld)

	// Find keywords
	suspiciousKeywords := []string{"login", "bank", "paypal", "verify", "secure", "account"}
	lowerURL := strings.ToLower(u.String())
	for _, kw := range suspiciousKeywords {
		if strings.Contains(lowerURL, kw) {
			enrichment.KeywordsFound = append(enrichment.KeywordsFound, kw)
		}
	}

	// Find suspicious params
	suspiciousParams := []string{"redirect", "url", "next", "return", "dest"}
	query := u.Query()
	for p := range query {
		lowerParam := strings.ToLower(p)
		for _, sp := range suspiciousParams {
			if strings.Contains(lowerParam, sp) {
				enrichment.SuspiciousParams = append(enrichment.SuspiciousParams, p)
				break
			}
		}
	}

	return enrichment
}

func (da *DomainAnalyzer) getTLDRiskScore(tld string) float64 {
	highRisk := map[string]float64{
		".tk": 0.9, ".ml": 0.9, ".ga": 0.9, ".cf": 0.8, ".gq": 0.8,
		".xyz": 0.7, ".top": 0.7, ".club": 0.6,
	}
	if score, ok := highRisk[tld]; ok {
		return score
	}
	return 0.1
}

// analyzeURLStructure analyzes URL for suspicious patterns
func (da *DomainAnalyzer) analyzeURLStructure(parsedURL *url.URL, analysis *models.ThreatAnalysis) int {
	score := 0
	hostname := parsedURL.Hostname()
	
	// Check for IP address instead of domain
	if da.isIPAddress(hostname) {
		score += 15
		analysis.Warnings = append(analysis.Warnings, 
			"Direct IP address used instead of domain name")
	}
	
	// Check for port number
	if parsedURL.Port() != "" {
		if parsedURL.Port() != "80" && parsedURL.Port() != "443" {
			score += 5
			analysis.Warnings = append(analysis.Warnings, 
				"Non-standard port used: "+parsedURL.Port())
		}
	}
	
	// Check URL length
	fullURL := parsedURL.String()
	if len(fullURL) > 200 {
		score += 10
		analysis.Warnings = append(analysis.Warnings, 
			"Unusually long URL (potential obfuscation)")
	}
	
	// Check for @ symbol
	if strings.Contains(fullURL, "@") {
		score += 20
		analysis.Warnings = append(analysis.Warnings, 
			"@ symbol in URL (possible credential injection)")
	}
	
	// Check path for suspicious patterns
	path := strings.ToLower(parsedURL.Path)
	suspiciousPaths := []struct {
		pattern string
		score   int
		message string
	}{
		{"/login", 5, "Login page detected"},
		{"/admin", 10, "Admin page detected"},
		{"/wp-admin", 8, "WordPress admin page"},
		{"/cgi-bin", 15, "CGI-BIN directory"},
		{".php", 5, "PHP file extension"},
		{".exe", 25, "Executable file in URL"},
		{".js", 15, "JavaScript file in URL"},
	}
	
	for _, sp := range suspiciousPaths {
		if strings.Contains(path, sp.pattern) {
			score += sp.score
			analysis.Warnings = append(analysis.Warnings, sp.message)
		}
	}
	
	return score
}

// analyzeDomainName checks domain name for suspicious patterns
func (da *DomainAnalyzer) analyzeDomainName(hostname string, analysis *models.ThreatAnalysis) int {
	score := 0
	domainWithoutTLD := da.removeTLD(hostname)
	
	// Check for numbers in domain
	if da.containsNumbers(domainWithoutTLD) {
		score += 5
		analysis.Warnings = append(analysis.Warnings, 
			"Numbers in domain name")
	}
	
	// Check hyphen count
	hyphenCount := strings.Count(domainWithoutTLD, "-")
	if hyphenCount > 1 {
		score += hyphenCount * 3
		analysis.Warnings = append(analysis.Warnings, 
			"Multiple hyphens in domain")
	}
	
	// Check domain length
	if len(domainWithoutTLD) < 3 {
		score += 15
		analysis.Warnings = append(analysis.Warnings, 
			"Extremely short domain name")
	} else if len(domainWithoutTLD) > 30 {
		score += 10
		analysis.Warnings = append(analysis.Warnings, 
			"Very long domain name")
	}
	
	// Check for brand impersonation
	if da.isBrandImpersonation(domainWithoutTLD) {
		score += 25
		analysis.Warnings = append(analysis.Warnings, 
			"Potential brand impersonation")
	}
	
	// Check subdomain count
	parts := strings.Split(hostname, ".")
	if len(parts) > 3 {
		score += (len(parts) - 3) * 2
		analysis.Warnings = append(analysis.Warnings, 
			"Multiple subdomains")
	}
	
	return score
}

// analyzeTLD checks TLD for risk
func (da *DomainAnalyzer) analyzeTLD(hostname string, analysis *models.ThreatAnalysis) int {
	score := 0
	tld := da.extractTLD(hostname)
	
	if tld == "" {
		return score
	}
	
	// High-risk TLDs
	highRiskTLDs := []string{
		".tk", ".ml", ".ga", ".cf", ".gq",
		".xyz", ".top", ".club", ".win", ".bid",
	}
	
	// Medium-risk TLDs
	mediumRiskTLDs := []string{
		".info", ".biz", ".work", ".online", ".site",
		".website", ".space", ".tech", ".store", ".shop",
	}
	
	for _, riskyTLD := range highRiskTLDs {
		if strings.HasSuffix(hostname, riskyTLD) {
			score += 20
			analysis.Warnings = append(analysis.Warnings, 
				"High-risk TLD: "+riskyTLD)
			return score
		}
	}
	
	for _, mediumTLD := range mediumRiskTLDs {
		if strings.HasSuffix(hostname, mediumTLD) {
			score += 10
			analysis.Warnings = append(analysis.Warnings, 
				"Medium-risk TLD: "+mediumTLD)
			return score
		}
	}
	
	return score
}

// Helper methods
func (da *DomainAnalyzer) isIPAddress(s string) bool {
	// Simple IPv4 check
	parts := strings.Split(s, ".")
	if len(parts) != 4 {
		return false
	}
	for _, part := range parts {
		if len(part) == 0 || len(part) > 3 {
			return false
		}
		for _, ch := range part {
			if ch < '0' || ch > '9' {
				return false
			}
		}
	}
	return true
}

func (da *DomainAnalyzer) isSuspiciousRegistrar(registrar string) bool {
	if registrar == "" {
		return false
	}
	
	suspiciousRegistrars := []string{
		"porkbun", "namecheap", "godaddy", "enom",
	}
	
	registrarLower := strings.ToLower(registrar)
	for _, suspicious := range suspiciousRegistrars {
		if strings.Contains(registrarLower, suspicious) {
			return true
		}
	}
	return false
}

func (da *DomainAnalyzer) isSuspiciousCNAME(cname string) bool {
	suspiciousPatterns := []string{
		"amazonaws.com", "cloudfront.net", "azurewebsites.net",
		"herokuapp.com", "github.io", "netlify.app",
	}
	
	cnameLower := strings.ToLower(cname)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(cnameLower, pattern) {
			return true
		}
	}
	return false
}

func (da *DomainAnalyzer) removeTLD(hostname string) string {
	lastDot := strings.LastIndex(hostname, ".")
	if lastDot == -1 {
		return hostname
	}
	return hostname[:lastDot]
}

func (da *DomainAnalyzer) containsNumbers(s string) bool {
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			return true
		}
	}
	return false
}

func (da *DomainAnalyzer) isBrandImpersonation(domain string) bool {
	brands := []string{
		"google", "facebook", "amazon", "apple", "microsoft",
		"paypal", "netflix", "twitter", "instagram", "whatsapp",
		"bank", "chase", "wellsfargo", "citi", "boa",
	}
	
	domainLower := strings.ToLower(domain)
	for _, brand := range brands {
		if strings.Contains(domainLower, brand) {
			return true
		}
	}
	return false
}

func (da *DomainAnalyzer) isHomographAttack(hostname string) bool {
	// Check for punycode
	if strings.HasPrefix(strings.ToLower(hostname), "xn--") {
		return true
	}
	
	// Check for mixed script characters
	hasLatin := false
	hasNonLatin := false
	
	for _, ch := range hostname {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			hasLatin = true
		} else if ch > 127 {
			hasNonLatin = true
		}
	}
	
	return hasLatin && hasNonLatin
}

func (da *DomainAnalyzer) extractTLD(hostname string) string {
	lastDot := strings.LastIndex(hostname, ".")
	if lastDot == -1 {
		return ""
	}
	return hostname[lastDot:]
}

func (da *DomainAnalyzer) extractDomain(rawURL string) string {
	parsed, err := da.urlParser.ParseAndAnalyze(rawURL)
	if err != nil {
		da.logger.Warn("Failed to parse URL %s for domain extraction: %v", rawURL, err)
		return ""
	}
	return parsed.Domain
}
