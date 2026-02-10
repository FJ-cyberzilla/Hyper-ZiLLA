package analyzer

import (
	"context"
	"net-zilla/pkg/logger"
	"testing"
)

func TestContentAnalyzer_Analyze(t *testing.T) {
	l := logger.NewLogger()
	ca := NewContentAnalyzer(l)

	testCases := []struct {
		url           string
		expectedScore int
	}{
		{"https://legit-site.com", 0},
		{"https://malicious-site.com/login-verify", 10}, // Score in code is 10 for "login"
	}

	for _, tc := range testCases {
		_, score := ca.Analyze(context.Background(), tc.url)
		if score < tc.expectedScore {
			t.Errorf("URL %s: Expected score >= %d, got %d", tc.url, tc.expectedScore, score)
		}
	}
}
