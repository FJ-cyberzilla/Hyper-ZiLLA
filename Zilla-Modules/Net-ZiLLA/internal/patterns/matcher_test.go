package patterns

import "testing"

func TestPatternMatcher_AnalyzeContent(t *testing.T) {
	pm := NewPatternMatcher()

	content := "This is a page with eval(atob('cGFzc3dvcmQ=')) and a password field."
	analysis := pm.AnalyzeContent(content)

	if len(analysis.Patterns) == 0 {
		t.Error("expected patterns to be detected")
	}

	foundObfuscation := false
	for _, p := range analysis.Patterns {
		if p.Name == "JavaScript Obfuscation" {
			foundObfuscation = true
		}
	}

	if !foundObfuscation {
		t.Error("expected JavaScript Obfuscation pattern to be found")
	}
}
