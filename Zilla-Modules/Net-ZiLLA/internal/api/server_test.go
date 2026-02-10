package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"net-zilla/internal/config"
	"net-zilla/internal/services"
	"net-zilla/internal/storage"
	"net-zilla/pkg/logger"
)

func TestHealthHandler(t *testing.T) {
	l := logger.NewLogger()
	cfg := &config.Config{}
	svc := services.NewAnalysisService(l, nil, cfg)
	server := NewServer(svc, l, cfg)

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.healthHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if response["status"] != "UP" {
		t.Errorf("expected status UP, got %v", response["status"])
	}
}

func TestAnalyzeHandler_MethodNotAllowed(t *testing.T) {
	l := logger.NewLogger()
	cfg := &config.Config{}
	svc := services.NewAnalysisService(l, nil, cfg)
	server := NewServer(svc, l, cfg)

	req, err := http.NewRequest("GET", "/api/v1/analyze", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.analyzeHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestAnalyzeHandler_InvalidJSON(t *testing.T) {
	l := logger.NewLogger()
	cfg := &config.Config{}
	svc := services.NewAnalysisService(l, nil, cfg)
	server := NewServer(svc, l, cfg)

	req, err := http.NewRequest("POST", "/api/v1/analyze", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.analyzeHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestAnalyzeHandler_Success(t *testing.T) {
	l := logger.NewLogger()
	cfg := &config.Config{}
	// Need a real-ish DB for service
	dbPath := "test_api.db"
	db, _ := storage.NewDatabase(dbPath)
	defer db.Close()
	
	svc := services.NewAnalysisService(l, db, cfg)
	server := NewServer(svc, l, cfg)

	body, _ := json.Marshal(map[string]string{"target": "http://example.com"})
	req, err := http.NewRequest("POST", "/api/v1/analyze", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.analyzeHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
