package network

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"net-zilla/pkg/logger"
)

func TestHTTPClient_SafeGetRequest(t *testing.T) {
	l := logger.NewLogger()
	hc := NewHTTPClient(l)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test-Header", "test-value")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	}))
	defer ts.Close()

	resp, err := hc.SafeGetRequest(context.Background(), ts.URL)
	if err != nil {
		t.Fatalf("SafeGetRequest failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
	if resp.Headers["X-Test-Header"] != "test-value" {
		t.Errorf("expected header value test-value, got %s", resp.Headers["X-Test-Header"])
	}
}

func TestHTTPClient_CheckSecurityHeaders(t *testing.T) {
	l := logger.NewLogger()
	hc := NewHTTPClient(l)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	headers, score, err := hc.CheckSecurityHeaders(context.Background(), ts.URL)
	if err != nil {
		t.Fatalf("CheckSecurityHeaders failed: %v", err)
	}

	if len(headers) < 2 {
		t.Errorf("expected at least 2 security headers, got %d", len(headers))
	}
	// Note: score will be low because httptest server doesn't have TLS by default
	_ = score
}

func TestHTTPClient_SafeHeadRequest(t *testing.T) {
	l := logger.NewLogger()
	hc := NewHTTPClient(l)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	resp, err := hc.SafeHeadRequest(context.Background(), ts.URL)
	if err != nil {
		t.Fatalf("SafeHeadRequest failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestHTTPClient_RequestError(t *testing.T) {
	l := logger.NewLogger()
	hc := NewHTTPClient(l)

	_, err := hc.SafeGetRequest(context.Background(), "http://invalid-url-that-does-not-exist.nz")
	if err == nil {
		t.Error("expected error for invalid URL")
	}
}
