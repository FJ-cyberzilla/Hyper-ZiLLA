package network

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"net-zilla/internal/models"
	"net-zilla/pkg/logger"
)

func TestExtractDomain(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"https://google.com/search", "google.com"},
		{"http://sub.domain.co.uk:8080/path", "sub.domain.co.uk"},
		{"invalid url", ""},
	}
	for _, tt := range tests {
		if got := extractDomain(tt.url); got != tt.want {
			t.Errorf("extractDomain(%s) = %v, want %v", tt.url, got, tt.want)
		}
	}
}

func TestIsObfuscatedURL(t *testing.T) {
	tests := []struct {
		url  string
		want bool
	}{
		{"https://google.com", false},
		{"http://192.168.1.1/test", true},
		{"http://example.com/%2525encoded", true},
		{"http://verylongsubdomainnamethatismorethantwentycharacters.com", true},
	}
	for _, tt := range tests {
		if got := isObfuscatedURL(tt.url); got != tt.want {
			t.Errorf("isObfuscatedURL(%s) = %v, want %v", tt.url, got, tt.want)
		}
	}
}

func TestRedirectTracer_ResolveNextURL(t *testing.T) {
	rt := &RedirectTracer{}
	base := "https://example.com/path/file.html"
	
	tests := []struct {
		location string
		want     string
	}{
		{"/newpath", "https://example.com/newpath"},
		{"other.html", "https://example.com/path/other.html"},
		{"https://otherdomain.com", "https://otherdomain.com/"},
	}

	for _, tt := range tests {
		got, err := rt.resolveNextURL(base, tt.location)
		if err != nil {
			t.Errorf("resolveNextURL failed: %v", err)
		}
		if strings.TrimSuffix(got, "/") != strings.TrimSuffix(tt.want, "/") {
			t.Errorf("resolveNextURL(%s) = %s, want %s", tt.location, got, tt.want)
		}
	}
}

func TestTraceRedirects(t *testing.T) {
	l := logger.NewLogger()
	rt := NewRedirectTracer(l)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/start" {
			http.Redirect(w, r, "/end", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	redirects, score, err := rt.TraceRedirects(context.Background(), ts.URL+"/start")
	if err != nil {
		t.Fatalf("TraceRedirects failed: %v", err)
	}

	if len(redirects) != 2 {
		t.Errorf("expected 2 hops, got %d", len(redirects))
	}
	if score == 0 {
		t.Error("expected non-zero threat score (at least for redirect itself)")
	}
}

func TestAnalyzeRedirectChain(t *testing.T) {
	chain := []models.RedirectDetail{
		{URL: "http://example.com/start"},
		{URL: "http://other.com/login.php"},
	}
	analysis := AnalyzeRedirectChain(chain)
	if analysis.Hops != 2 {
		t.Errorf("expected 2 hops, got %d", analysis.Hops)
	}
	if analysis.UniqueDomains != 2 {
		t.Errorf("expected 2 unique domains, got %d", analysis.UniqueDomains)
	}
	if analysis.SuspiciousURLs != 1 {
		t.Errorf("expected 1 suspicious URL, got %d", analysis.SuspiciousURLs)
	}
}
