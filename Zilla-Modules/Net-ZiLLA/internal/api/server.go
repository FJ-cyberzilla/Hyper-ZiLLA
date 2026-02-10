package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"net-zilla/internal/config"
	"net-zilla/internal/middleware"
	"net-zilla/internal/services"
	"net-zilla/pkg/logger"
)

type APIServer struct {
	server          *http.Server
	analysisService *services.AnalysisService
	logger          *logger.Logger
	config          *config.Config
	middleware      *middleware.MiddlewareStack
}

func NewServer(analysisService *services.AnalysisService, logger *logger.Logger, cfg *config.Config) *APIServer {
	mux := http.NewServeMux()
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	s := &APIServer{
		analysisService: analysisService,
		logger:          logger,
		config:          cfg,
		middleware:      middleware.NewMiddleware(logger),
	}

	s.server = &http.Server{
		Addr:         serverAddr,
		Handler:      mux,
		ReadTimeout:  cfg.Security.RequestTimeout,
		WriteTimeout: cfg.Security.RequestTimeout,
	}

	s.setupRoutes(mux)
	return s
}

func (s *APIServer) setupRoutes(mux *http.ServeMux) {
	mux.Handle("/api/v1/analyze", s.middleware.Chain(http.HandlerFunc(s.analyzeHandler), middleware.LoggerMiddleware(s.logger)))
	mux.HandleFunc("/health", s.healthHandler)
}

func (s *APIServer) healthHandler(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status    string `json:"status"`
		Timestamp string `json:"timestamp"`
		Version   string `json:"version"`
	}{
		Status:    "UP",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (s *APIServer) analyzeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	var req struct {
		Target string `json:"target"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON body"})
		return
	}

	if req.Target == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Target is required"})
		return
	}

	report, err := s.analysisService.PerformAnalysis(r.Context(), req.Target)
	if err != nil {
		s.logger.Error("Analysis failed for %s: %v", req.Target, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal analysis error"})
		return
	}

	json.NewEncoder(w).Encode(report)
}


func (s *APIServer) Run(ctx context.Context) error {
	s.logger.Info("🚀 Net-Zilla API server starting on %s", s.server.Addr)
	go s.server.ListenAndServe()
	<-ctx.Done()
	return s.server.Shutdown(context.Background())
}
