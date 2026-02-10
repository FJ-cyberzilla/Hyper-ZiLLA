package patterns

import (
	"net-zilla/internal/models"
)

type SignatureEngine struct {
	patterns []ThreatRegex
}

func NewSignatureEngine() *SignatureEngine {
	return &SignatureEngine{
		patterns: GetDefaultPatterns(),
	}
}

func (se *SignatureEngine) Scan(content string) []models.BehavioralPattern {
	var findings []models.BehavioralPattern

	for _, p := range se.patterns {
		if p.Regex.MatchString(content) {
			findings = append(findings, models.BehavioralPattern{
				Name:        p.Name,
				Type:        models.PatternMalware,
				Description: p.Description,
				Weight:      se.getWeight(p.Severity),
				Evidence:    []string{p.Regex.FindString(content)},
			})
		}
	}

	return findings
}

func (se *SignatureEngine) getWeight(severity string) int {
	switch severity {
	case "High":
		return 40
	case "Medium":
		return 20
	case "Low":
		return 10
	default:
		return 0
	}
}
