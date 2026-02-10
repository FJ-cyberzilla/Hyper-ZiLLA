package patterns

import (
	"net-zilla/internal/models"
)

type PatternMatcher struct {
	signatureEngine *SignatureEngine
	yara            *YaraManager
}

func NewPatternMatcher() *PatternMatcher {
	return &PatternMatcher{
		signatureEngine: NewSignatureEngine(),
		yara:            NewYaraManager(),
	}
}

func (pm *PatternMatcher) AnalyzeContent(content string) *models.BehaviorAnalysis {
	analysis := &models.BehaviorAnalysis{}

	// Scan with built-in regex signatures
	analysis.Patterns = pm.signatureEngine.Scan(content)

	if len(analysis.Patterns) > 0 {
		analysis.RiskSignature = "PATTERN_MATCH_FOUND"
	}

	return analysis
}
