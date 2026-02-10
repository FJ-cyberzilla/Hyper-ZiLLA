package processor

import (
	"testing"
)

func TestURLParser_ParseAndAnalyze(t *testing.T) {
	p := NewURLParser()

	tests := []struct {
		name         string
		url          string
		isShortened  bool
		isObfuscated bool
		hasIP        bool
	}{
		{
			name:         "Standard URL",
			url:          "https://www.google.com/search?q=test",
			isShortened:  false,
			isObfuscated: false,
			hasIP:        false,
		},
		{
			name:         "Shortened URL",
			url:          "http://bit.ly/xxxx",
			isShortened:  true,
			isObfuscated: false,
			hasIP:        false,
		},
		{
			name:         "Obfuscated with @",
			url:          "https://google.com@malicious.com",
			isShortened:  false,
			isObfuscated: true,
			hasIP:        false,
		},
		{
			name:         "IP Address Host",
			url:          "http://192.168.1.1/admin",
			isShortened:  false,
			isObfuscated: true, // regex matches IP in obfuscationPatterns too
			hasIP:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := p.ParseAndAnalyze(tt.url)
			if err != nil {
				t.Fatalf("ParseAndAnalyze() error = %v", err)
			}
			if parsed.IsShortened != tt.isShortened {
				t.Errorf("IsShortened = %v, want %v", parsed.IsShortened, tt.isShortened)
			}
			if parsed.IsObfuscated != tt.isObfuscated {
				t.Errorf("IsObfuscated = %v, want %v", parsed.IsObfuscated, tt.isObfuscated)
			}
			if parsed.HasIP != tt.hasIP {
				t.Errorf("HasIP = %v, want %v", parsed.HasIP, tt.hasIP)
			}
		})
	}
}

func TestURLParser_DetectHomographAttack(t *testing.T) {
	p := NewURLParser()

	tests := []struct {
		domain string
		want   bool
	}{
		{"google.com", false},
		{"pypal.com", false}, // just misspelling, not homograph
		{"xn--pypal-4ve.com", false}, // punycode itself is just latin/digits/hyphen
		{"google.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.domain, func(t *testing.T) {
			if got := p.DetectHomographAttack(tt.domain); got != tt.want {
				t.Errorf("DetectHomographAttack() = %v, want %v", got, tt.want)
			}
		})
	}

	// Manual check for a real homograph (mixed script)
	// 'а' is Cyrillic small letter a (U+0430)
	homograph := "google.а" // mixed latin 'google' + cyrillic 'а'
	if !p.DetectHomographAttack(homograph) {
		t.Errorf("DetectHomographAttack(%s) should be true", homograph)
	}
}
