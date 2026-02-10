package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	rl := NewRateLimiter(RateLimitConfig{
		Rate:          2,
		Interval:      time.Second,
		BlockDuration: time.Millisecond,
	})
	ip := "127.0.0.1"

	if !rl.Allow(ip) {
		t.Error("first request should be allowed")
	}
	if !rl.Allow(ip) {
		t.Error("second request should be allowed")
	}
	if rl.Allow(ip) {
		t.Error("third request should be blocked")
	}

	time.Sleep(1100 * time.Millisecond)
	if !rl.Allow(ip) {
		t.Error("request after interval should be allowed")
	}
}

func TestAuthMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mw := AuthMiddleware()
	handler := mw(nextHandler)

	// Case 1: No token
	req1 := httptest.NewRequest("GET", "/", nil)
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)
	if rr1.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr1.Code)
	}

	// Case 2: Valid token (must match constant-time comparison in isValidToken)
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.Header.Set("Authorization", "Bearer your-expected-token-here")
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr2.Code)
	}
}
