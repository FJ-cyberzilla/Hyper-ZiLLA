package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"net-zilla/pkg/logger"
)

func TestMiddlewareStack_Chain(t *testing.T) {
	ms := NewMiddleware(logger.NewLogger())
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Test", "base")
		w.WriteHeader(http.StatusOK)
	})

	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Test", "mw1")
			next.ServeHTTP(w, r)
		})
	}

	chained := ms.Chain(handler, mw1)
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	chained.ServeHTTP(rr, req)

	vals := rr.Header().Values("X-Test")
	if len(vals) != 2 || vals[0] != "mw1" || vals[1] != "base" {
		t.Errorf("chaining order or execution failed: %v", vals)
	}
}

func TestLoggerMiddleware(t *testing.T) {
	l := logger.NewLogger()
	mw := LoggerMiddleware(l)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/log-test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Error("expected 200 OK")
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	mw := RateLimitMiddleware(1)
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "1.2.3.4:1234"

	// First allowed
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req)
	if rr1.Code != http.StatusOK {
		t.Errorf("first request should be allowed, got %d", rr1.Code)
	}

	// Second blocked
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req)
	if rr2.Code != http.StatusTooManyRequests {
		t.Errorf("second request should be blocked, got %d", rr2.Code)
	}
}

func TestCORSHeaderMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mw := CORSHeaderMiddleware()
	handler := mw(nextHandler)

	// Case 1: OPTIONS
	req1 := httptest.NewRequest("OPTIONS", "/", nil)
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)
	if rr1.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("expected CORS header on OPTIONS")
	}

	// Case 2: GET
	req2 := httptest.NewRequest("GET", "/", nil)
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	if rr2.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("expected CORS header on GET")
	}
}
