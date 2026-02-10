package network

import "testing"

func TestSafetyScreener_Screen(t *testing.T) {
	screener := NewSafetyScreener()

	tests := []struct {
		url          string
		isSuspicious bool
	}{
		{"https://microsoft.com", false},
		{"http://192.168.1.1", true},
		{"https://something.tk", true},
		{"http://xn--pypal-4ve.com", true},
	}

	for _, tt := range tests {
		result := screener.Screen(tt.url)
		if result.IsSuspicious != tt.isSuspicious {
			t.Errorf("URL %s: expected suspicious=%v, got %v", tt.url, tt.isSuspicious, result.IsSuspicious)
		}
	}
}
