package patterns

import "regexp"

type ThreatRegex struct {
	Name        string
	Regex       *regexp.Regexp
	Severity    string
	Description string
}

func GetDefaultPatterns() []ThreatRegex {
	return []ThreatRegex{
		{
			Name:        "JavaScript Obfuscation",
			Regex:       regexp.MustCompile(`(?i)eval\s*\(\s*atob|String\.fromCharCode`),
			Severity:    "High",
			Description: "Detects common JS obfuscation techniques used in phishing.",
		},
		{
			Name:        "Suspicious Redirect",
			Regex:       regexp.MustCompile(`(?i)window\.location\s*=|document\.location\.replace`),
			Severity:    "Medium",
			Description: "Detects client-side redirects.",
		},
		{
			Name:        "Credential Phishing",
			Regex:       regexp.MustCompile(`(?i)password|login|signin|verify_account`),
			Severity:    "Low",
			Description: "Common keywords associated with credential harvesting.",
		},
		{
			Name:        "Data Exfiltration",
			Regex:       regexp.MustCompile(`(?i)XMLHttpRequest|fetch\(|navigator\.sendBeacon`),
			Severity:    "Medium",
			Description: "Potential background data transfer detected.",
		},
	}
}
