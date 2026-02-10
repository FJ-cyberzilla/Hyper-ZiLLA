package middleware

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"net-zilla/pkg/logger"
)

// Middleware is a function that takes an http.Handler and returns an http.Handler.
type Middleware func(http.Handler) http.Handler

// MiddlewareStack holds configured middleware services.
type MiddlewareStack struct {
	logger *logger.Logger
}

// NewMiddleware creates a new MiddlewareStack.
func NewMiddleware(logger *logger.Logger) *MiddlewareStack {
	return &MiddlewareStack{
		logger: logger,
	}
}

// Chain applies a list of middleware to a http.Handler.
func (ms *MiddlewareStack) Chain(h http.Handler, middleware ...Middleware) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

// LoggerMiddleware logs incoming HTTP requests with enhanced details.
func LoggerMiddleware(logger *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Use response writer wrapper to capture status code
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			
			// Process request
			next.ServeHTTP(rw, r)
			
			// Calculate duration
			duration := time.Since(start)
			
			// Log with different levels based on status code and duration
			clientIP := getClientIP(r)
			userAgent := r.UserAgent()
			
			// Truncate user agent if too long
			if len(userAgent) > 100 {
				userAgent = userAgent[:97] + "..."
			}
			
			// Enhanced logging
			logEntry := map[string]interface{}{
				"method":     r.Method,
				"path":       r.URL.Path,
				"query":      r.URL.RawQuery,
				"ip":         clientIP,
				"user_agent": userAgent,
				"status":     rw.statusCode,
				"duration":   duration.String(),
				"duration_ms": duration.Milliseconds(),
				"referer":    r.Referer(),
				"time":       time.Now().Format(time.RFC3339),
			}
			
			// Log at appropriate level
			if rw.statusCode >= 500 {
				logger.Error("HTTP %d %s %s from %s (%s) - %v", 
					rw.statusCode, r.Method, r.URL.Path, clientIP, duration, userAgent)
			} else if rw.statusCode >= 400 {
				logger.Warn("HTTP %d %s %s from %s (%s)", 
					rw.statusCode, r.Method, r.URL.Path, clientIP, duration)
			} else {
				logger.Info("HTTP %d %s %s from %s (%s)", 
					rw.statusCode, r.Method, r.URL.Path, clientIP, duration)
			}
			
			// Debug log with full details
			if logger.IsDebug() {
				logger.Debug("HTTP Request Details: %v", logEntry)
			}
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// AuthMiddleware handles API key or token authentication with improvements.
func AuthMiddleware(authConfig ...AuthConfig) Middleware {
	config := AuthConfig{
		AuthHeader: "Authorization",
		AuthPrefix: "Bearer ",
		ApiKeys:    make(map[string]bool),
	}
	
	if len(authConfig) > 0 {
		config = authConfig[0]
	}
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check for API key in header
			token := r.Header.Get(config.AuthHeader)
			
			// If no header, check query parameter
			if token == "" && config.AllowQueryParam {
				token = r.URL.Query().Get("api_key")
			}
			
			if token == "" {
				RespondWithError(w, http.StatusUnauthorized, "Authorization token required")
				return
			}
			
			// Remove prefix if present
			if config.AuthPrefix != "" && strings.HasPrefix(token, config.AuthPrefix) {
				token = strings.TrimPrefix(token, config.AuthPrefix)
			}
			
			// Validate token
			var userID string
			var valid bool
			
			if len(config.ApiKeys) > 0 {
				// Check against predefined API keys
				valid = config.ApiKeys[token]
				if valid {
					userID = "api_user"
				}
			} else {
				// Use default validation
				valid = isValidToken(token)
				if valid {
					userID = "authenticated_user"
				}
			}
			
			if !valid {
				RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}
			
			// Add user to context
			ctx := context.WithValue(r.Context(), "user", userID)
			ctx = context.WithValue(ctx, "token", token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AuthConfig provides configuration for authentication middleware
type AuthConfig struct {
	AuthHeader      string
	AuthPrefix      string
	ApiKeys         map[string]bool
	AllowQueryParam bool
}

// RateLimiter applies rate limiting based on client IP with improvements.
type RateLimiter struct {
	clients        map[string]*Client
	mu             sync.RWMutex
	rate           int
	interval       time.Duration
	maxBurst       int
	blockDuration  time.Duration
	cleanupTicker  *time.Ticker
	cleanupStop    chan bool
}

type Client struct {
	lastRequest    time.Time
	requests       int
	firstRequest   time.Time
	blockedUntil   time.Time
	isBlocked      bool
	blockCount     int
}

// RateLimitConfig provides configuration for rate limiting
type RateLimitConfig struct {
	Rate           int           // Requests per interval
	Interval       time.Duration // Time interval
	MaxBurst       int           // Maximum burst capacity
	BlockDuration  time.Duration // How long to block after limit
	CleanupPeriod  time.Duration // How often to clean old entries
}

func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	if config.Rate <= 0 {
		config.Rate = 60
	}
	if config.Interval <= 0 {
		config.Interval = time.Minute
	}
	if config.MaxBurst <= 0 {
		config.MaxBurst = config.Rate
	}
	if config.BlockDuration <= 0 {
		config.BlockDuration = 5 * time.Minute
	}
	if config.CleanupPeriod <= 0 {
		config.CleanupPeriod = 10 * time.Minute
	}
	
	rl := &RateLimiter{
		clients:       make(map[string]*Client),
		rate:          config.Rate,
		interval:      config.Interval,
		maxBurst:      config.MaxBurst,
		blockDuration: config.BlockDuration,
		cleanupStop:   make(chan bool),
	}
	
	// Start cleanup goroutine
	rl.cleanupTicker = time.NewTicker(config.CleanupPeriod)
	go rl.cleanupOldEntries()
	
	return rl
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	client, found := rl.clients[ip]
	
	if !found {
		rl.clients[ip] = &Client{
			lastRequest:  now,
			firstRequest: now,
			requests:     1,
		}
		return true
	}
	
	// Check if client is blocked
	if client.isBlocked {
		if now.After(client.blockedUntil) {
			// Unblock client
			client.isBlocked = false
			client.requests = 0
			client.lastRequest = now
			client.firstRequest = now
			return true
		}
		return false
	}
	
	// Check if interval has passed
	if now.Sub(client.firstRequest) > rl.interval {
		// Reset for new interval
		client.requests = 1
		client.firstRequest = now
		client.lastRequest = now
		return true
	}
	
	// Check if within rate limit
	if client.requests < rl.rate {
		client.requests++
		client.lastRequest = now
		return true
	}
	
	// Check burst capacity
	if client.requests < rl.maxBurst {
		client.requests++
		client.lastRequest = now
		return true
	}
	
	// Exceeded limit, block client
	client.isBlocked = true
	client.blockCount++
	client.blockedUntil = now.Add(time.Duration(client.blockCount) * rl.blockDuration)
	
	return false
}

func (rl *RateLimiter) cleanupOldEntries() {
	for {
		select {
		case <-rl.cleanupTicker.C:
			rl.mu.Lock()
			now := time.Now()
			for ip, client := range rl.clients {
				// Remove entries older than 24 hours
				if now.Sub(client.lastRequest) > 24*time.Hour {
					delete(rl.clients, ip)
				}
			}
			rl.mu.Unlock()
		case <-rl.cleanupStop:
			rl.cleanupTicker.Stop()
			return
		}
	}
}

func (rl *RateLimiter) Stop() {
	rl.cleanupStop <- true
}

func RateLimitMiddleware(rateLimit int, config ...RateLimitConfig) Middleware {
	var limiter *RateLimiter
	
	if len(config) > 0 {
		limiter = NewRateLimiter(config[0])
	} else {
		limiter = NewRateLimiter(RateLimitConfig{
			Rate:     rateLimit,
			Interval: time.Minute,
		})
	}
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)
			
			// Apply rate limiting
			if !limiter.Allow(clientIP) {
				// Add rate limit headers
				w.Header().Set("X-RateLimit-Limit", "60")
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("X-RateLimit-Reset", time.Now().Add(time.Minute).Format(time.RFC1123))
				
				RespondWithError(w, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later.")
				return
			}
			
			// Add rate limit info headers
			w.Header().Set("X-RateLimit-Limit", "60")
			
			next.ServeHTTP(w, r)
		})
	}
}

// CORSHeaderMiddleware adds CORS headers with configurable options.
func CORSHeaderMiddleware(config ...CORSConfig) Middleware {
	corsConfig := CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "X-RateLimit-Limit", "X-RateLimit-Remaining"},
		AllowCredentials: false,
		MaxAge:           300, // 5 minutes
	}
	
	if len(config) > 0 {
		corsConfig = config[0]
	}
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			if len(corsConfig.AllowOrigins) > 0 {
				if corsConfig.AllowOrigins[0] == "*" {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else {
					// Check if request origin is allowed
					origin := r.Header.Get("Origin")
					for _, allowed := range corsConfig.AllowOrigins {
						if allowed == origin {
							w.Header().Set("Access-Control-Allow-Origin", origin)
							break
						}
					}
				}
			}
			
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowHeaders, ", "))
			
			if len(corsConfig.ExposeHeaders) > 0 {
				w.Header().Set("Access-Control-Expose-Headers", strings.Join(corsConfig.ExposeHeaders, ", "))
			}
			
			if corsConfig.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			
			if corsConfig.MaxAge > 0 {
				w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", corsConfig.MaxAge))
			}
			
			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

// CORSConfig provides configuration for CORS middleware
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// SecurityHeadersMiddleware adds security headers.
func SecurityHeadersMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Security headers
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("Content-Security-Policy", "default-src 'self'")
			
			// Remove server header
			w.Header().Del("Server")
			w.Header().Del("X-Powered-By")
			
			next.ServeHTTP(w, r)
		})
	}
}

// RequestIDMiddleware adds a unique request ID to each request.
func RequestIDMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := generateRequestID()
			w.Header().Set("X-Request-ID", requestID)
			
			ctx := context.WithValue(r.Context(), "request_id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// TimeoutMiddleware adds a timeout to requests.
func TimeoutMiddleware(timeout time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			
			// Create a channel to detect handler completion
			done := make(chan bool, 1)
			
			go func() {
				next.ServeHTTP(w, r.WithContext(ctx))
				done <- true
			}()
			
			select {
			case <-done:
				// Handler completed normally
				return
			case <-ctx.Done():
				// Timeout occurred
				if ctx.Err() == context.DeadlineExceeded {
					RespondWithError(w, http.StatusGatewayTimeout, "Request timeout")
				}
			}
		})
	}
}

// RecoveryMiddleware recovers from panics.
func RecoveryMiddleware(logger *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("Panic recovered: %v", err)
					
					// Log stack trace
					// debug.PrintStack() // Uncomment for detailed stack traces
					
					RespondWithError(w, http.StatusInternalServerError, "Internal server error")
				}
			}()
			
			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extracts client IP from request with improved logic.
func getClientIP(r *http.Request) string {
	// Check for Cloudflare
	if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
		return ip
	}
	
	// Check for X-Forwarded-For
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	
	// Check for X-Real-IP
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	
	// Check for Forwarded header (RFC 7239)
	if forwarded := r.Header.Get("Forwarded"); forwarded != "" {
		// Parse Forwarded header
		parts := strings.Split(forwarded, ";")
		for _, part := range parts {
			if strings.HasPrefix(strings.TrimSpace(part), "for=") {
				ip := strings.TrimPrefix(strings.TrimSpace(part), "for=")
				// Remove quotes and brackets
				ip = strings.Trim(ip, `"[]`)
				return ip
			}
		}
	}
	
	// Fall back to remote address
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// isValidToken validates authentication tokens with constant-time comparison.
func isValidToken(token string) bool {
	// Use constant-time comparison to prevent timing attacks
	expected := "your-expected-token-here" // Should come from config
	return subtle.ConstantTimeCompare([]byte(token), []byte(expected)) == 1
}

// generateRequestID generates a unique request ID.
func generateRequestID() string {
	return strings.ReplaceAll(time.Now().Format("20060102150405.000000"), ".", "") + 
		"-" + randomString(8)
}

// randomString generates a random string of specified length.
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		// Simple random, not cryptographically secure
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// RespondWithError sends an error response in JSON format.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]interface{}{
		"error":   message,
		"code":    code,
		"success": false,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// RespondWithJSON sends a JSON response.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Failed to marshal JSON response","code":500,"success":false}`))
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
