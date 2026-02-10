package analyzer

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"net-zilla/internal/network"
	"net-zilla/pkg/logger"
)

type ContentAnalyzer struct {
	logger     *logger.Logger
	httpClient *network.HTTPClient
}

func NewContentAnalyzer(logger *logger.Logger) *ContentAnalyzer {
	return &ContentAnalyzer{
		logger:     logger,
		httpClient: network.NewHTTPClient(logger),
	}
}

func (ca *ContentAnalyzer) Analyze(ctx context.Context, url string) ([]string, int) {
	score := 0
	var warnings []string

	// 1. Analyze URL path for suspicious keywords
	lowerURL := strings.ToLower(url)
	
	suspiciousKeywords := []string{
		"login", "signin", "sign-in", "log-in",
		"password", "passwd", "pwd",
		"verify", "verification", "confirm",
		"update", "account", "profile",
		"secure", "security", "auth", "authentication",
		"bank", "paypal", "payment", "pay",
		"wallet", "credit", "card", "cvv",
		"social", "security", "ssn",
		"admin", "administrator", "root",
		"phpmyadmin", "wp-admin", "administrator",
		"cgi-bin", "bin", "cmd", "exec",
	}
	
	for _, kw := range suspiciousKeywords {
		if strings.Contains(lowerURL, kw) {
			score += 10
			warnings = append(warnings, "Suspicious keyword in URL: "+kw)
		}
	}

	// 2. Check for URL obfuscation techniques
	if ca.detectObfuscation(url) {
		score += 20
		warnings = append(warnings, "URL obfuscation detected")
	}

	// 3. Check for encoded characters
	if ca.hasEncodedCharacters(url) {
		score += 15
		warnings = append(warnings, "Encoded/obfuscated characters in URL")
	}

	// 4. Check for excessive special characters
	specialChars := ca.countSpecialCharacters(url)
	if specialChars > 5 {
		score += specialChars * 2
		warnings = append(warnings, 
			fmt.Sprintf("Excessive special characters in URL: %d", specialChars))
	}

	// 5. Check for IP addresses in URL (potential direct IP access)
	if ca.hasIPAddress(url) {
		score += 25
		warnings = append(warnings, "Direct IP address access detected")
	}

	// 6. Check for file extensions
	maliciousExtensions := []string{
		".exe", ".msi", ".bat", ".cmd", ".ps1", ".vbs", ".js", ".jar",
		".scr", ".pif", ".com", ".hta", ".wsf", ".sh", ".bash",
	}
	
	for _, ext := range maliciousExtensions {
		if strings.HasSuffix(lowerURL, ext) || 
		   strings.Contains(lowerURL, ext+"?") ||
		   strings.Contains(lowerURL, ext+"&") ||
		   strings.Contains(lowerURL, ext+"#") {
			score += 40
			warnings = append(warnings, "Suspicious file extension: "+ext)
			break
		}
	}

	// 7. Check URL length (potential for obfuscation)
	if len(url) > 200 {
		score += 15
		warnings = append(warnings, "Unusually long URL (potential obfuscation)")
	}

	// 8. Check for @ symbol (userinfo obfuscation)
	if strings.Contains(url, "@") {
		score += 20
		warnings = append(warnings, "@ symbol in URL (possible obfuscation)")
	}

	// 9. Check for double slashes or encoded slashes
	if strings.Contains(url, "//") && !strings.HasPrefix(url, "http://") && 
	   !strings.HasPrefix(url, "https://") {
		score += 10
		warnings = append(warnings, "Double slashes in path")
	}
	
	if strings.Contains(url, "%2f") || strings.Contains(url, "%2F") {
		score += 10
		warnings = append(warnings, "Encoded slashes in URL")
	}

	// 10. Check for NULL bytes or control characters
	if ca.hasControlCharacters(url) {
		score += 30
		warnings = append(warnings, "Control characters in URL")
	}

	// 11. Check for punycode/homograph attacks
	if ca.isPunycode(url) {
		score += 25
		warnings = append(warnings, "Punycode/IDN domain detected")
	}

	// 12. Check for port numbers (non-standard)
	if ca.hasNonStandardPort(url) {
		score += 10
		warnings = append(warnings, "Non-standard port in URL")
	}

	// 13. Check for redirect parameters
	redirectParams := []string{
		"redirect=", "redirect_to=", "redirect_uri=", "redirect_url=",
		"return=", "return_to=", "return_url=",
		"goto=", "go=", "next=", "forward=",
		"url=", "link=", "target=", "dest=", "destination=",
	}
	
	for _, param := range redirectParams {
		if strings.Contains(lowerURL, param) {
			score += 15
			warnings = append(warnings, "Redirect parameter detected: "+param)
			break
		}
	}

	// 14. Check for data/javascript/file URIs
	badSchemes := []string{
		"data:", "javascript:", "file:", "ftp:", "telnet:", "gopher:",
		"mailto:", "news:", "nntp:", "irc:", "ircs:", "ssh:",
	}
	
	for _, scheme := range badSchemes {
		if strings.HasPrefix(strings.ToLower(url), scheme) {
			score += 50
			warnings = append(warnings, "Dangerous URI scheme: "+scheme)
			break
		}
	}

	// 15. Check for local/internal addresses
	if ca.isLocalAddress(url) {
		score += 100
		warnings = append(warnings, "Local/internal network address")
	}

	// Cap score at 100
	if score > 100 {
		score = 100
	}

	// Add severity note based on score
	if score >= 80 {
		warnings = append(warnings, "⚠️ HIGH RISK: Multiple suspicious indicators detected")
	} else if score >= 50 {
		warnings = append(warnings, "⚠️ MEDIUM RISK: Several suspicious indicators")
	} else if score >= 20 {
		warnings = append(warnings, "⚠️ LOW RISK: Some suspicious indicators")
	}

	return warnings, score
}

// Helper methods using only standard Go libraries
func (ca *ContentAnalyzer) detectObfuscation(url string) bool {
	// Check for hex encoding
	hexRegex := regexp.MustCompile(`%[0-9a-fA-F]{2}`)
	hexMatches := hexRegex.FindAllString(url, -1)
	if len(hexMatches) > 3 {
		return true
	}

	// Check for Unicode encoding
	unicodeRegex := regexp.MustCompile(`\\u[0-9a-fA-F]{4}|\\x[0-9a-fA-F]{2}`)
	if unicodeRegex.MatchString(url) {
		return true
	}

	// Check for double encoding
	if strings.Contains(url, "%25") { // %25 = encoded %
		return true
	}

	return false
}

func (ca *ContentAnalyzer) hasEncodedCharacters(url string) bool {
	// Check for URL encoded characters
	encodedPatterns := []string{
		"%20", "%21", "%22", "%23", "%24", "%25", "%26", "%27", "%28", "%29",
		"%2A", "%2B", "%2C", "%2D", "%2E", "%2F",
		"%3A", "%3B", "%3C", "%3D", "%3E", "%3F",
		"%40", "%5B", "%5C", "%5D", "%5E", "%5F",
		"%60", "%7B", "%7C", "%7D", "%7E",
	}
	
	for _, pattern := range encodedPatterns {
		if strings.Contains(strings.ToUpper(url), pattern) {
			return true
		}
	}
	
	return false
}

func (ca *ContentAnalyzer) countSpecialCharacters(url string) int {
	count := 0
	specialChars := "!@#$%^&*()_+-=[]{}|;:'\",.<>?/~`"
	
	for _, char := range url {
		if strings.ContainsRune(specialChars, char) {
			count++
		}
	}
	
	return count
}

func (ca *ContentAnalyzer) hasIPAddress(url string) bool {
	// Simple IPv4 pattern
	ipv4Regex := regexp.MustCompile(`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`)
	return ipv4Regex.MatchString(url)
}

func (ca *ContentAnalyzer) hasControlCharacters(url string) bool {
	for _, char := range url {
		if char < 32 || char == 127 { // Control characters
			return true
		}
	}
	return false
}

func (ca *ContentAnalyzer) isPunycode(url string) bool {
	// Check for xn-- prefix (punycode)
	punycodeRegex := regexp.MustCompile(`(?i)xn--`)
	return punycodeRegex.MatchString(url)
}

func (ca *ContentAnalyzer) hasNonStandardPort(url string) bool {
	// Check for port numbers other than 80, 443
	portRegex := regexp.MustCompile(`:\d+/`)
	matches := portRegex.FindStringSubmatch(url)
	if len(matches) > 0 {
		portStr := strings.TrimSuffix(strings.TrimPrefix(matches[0], ":"), "/")
		// Check if it's a non-standard port
		if portStr != "80" && portStr != "443" && portStr != "" {
			return true
		}
	}
	return false
}

func (ca *ContentAnalyzer) isLocalAddress(url string) bool {
	lowerURL := strings.ToLower(url)
	
	localPatterns := []string{
		"localhost",
		"127.0.0.1",
		"0.0.0.0",
		"::1",
		"fe80::",
		"169.254.",
		"192.168.",
		"10.",
		"172.16.", "172.17.", "172.18.", "172.19.",
		"172.20.", "172.21.", "172.22.", "172.23.",
		"172.24.", "172.25.", "172.26.", "172.27.",
		"172.28.", "172.29.", "172.30.", "172.31.",
	}
	
	for _, pattern := range localPatterns {
		if strings.Contains(lowerURL, pattern) {
			return true
		}
	}
	
	return false
}

// Additional method for deeper analysis if needed
func (ca *ContentAnalyzer) AnalyzeContentString(content string) ([]string, int) {
	score := 0
	var warnings []string
	
	if content == "" {
		return warnings, score
	}
	
	lowerContent := strings.ToLower(content)
	
	// Check for suspicious patterns in content
	patterns := []struct {
		pattern string
		score   int
		message string
	}{
		// Phishing patterns
		{`password\s*:\s*`, 20, "Password field in content"},
		{`username\s*:\s*`, 15, "Username field in content"},
		{`credit\s*card`, 25, "Credit card reference"},
		{`social\s*security`, 25, "Social security number reference"},
		{`account\s*number`, 20, "Account number reference"},
		
		// Malware patterns
		{`eval\s*\(`, 30, "JavaScript eval() function"},
		{`document\.write`, 15, "Dynamic content writing"},
		{`innerHTML\s*=`, 15, "InnerHTML manipulation"},
		{`fromCharCode`, 20, "Character code obfuscation"},
		{`unescape\s*\(`, 20, "JavaScript unescape()"},
		
		// Obfuscation patterns
		{`\\x[0-9a-f]{2}`, 25, "Hex encoded characters"},
		{`\\u[0-9a-f]{4}`, 25, "Unicode encoded characters"},
		{`String\.fromCharCode`, 20, "Character code construction"},
		
		// Iframe injection
		{`<iframe[^>]*>`, 15, "Iframe element"},
		{`style\s*=\s*["']display:\s*none["']`, 20, "Hidden element"},
		{`visibility\s*:\s*hidden`, 20, "Invisible element"},
		
		// Form phishing
		{`<form[^>]*>`, 10, "Form element"},
		{`action\s*=\s*["']http`, 15, "Form submits to external URL"},
	}
	
	for _, p := range patterns {
		regex := regexp.MustCompile(p.pattern)
		if regex.MatchString(lowerContent) {
			score += p.score
			warnings = append(warnings, p.message)
		}
	}
	
	// Check content length (very short content could be malicious)
	if len(content) < 100 {
		score += 10
		warnings = append(warnings, "Very short content")
	}
	
	// Check for excessive whitespace (obfuscation)
	spaceRatio := float64(strings.Count(content, " ")+strings.Count(content, "\t")+strings.Count(content, "\n")) / float64(len(content))
	if spaceRatio > 0.3 {
		score += 15
		warnings = append(warnings, "Excessive whitespace (potential obfuscation)")
	}
	
	if score > 100 {
		score = 100
	}
	
	return warnings, score
}
