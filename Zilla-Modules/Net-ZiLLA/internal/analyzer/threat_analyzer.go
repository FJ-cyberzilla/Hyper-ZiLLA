package analyzer

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"net-zilla/internal/models"
	"net-zilla/internal/network"
	"net-zilla/internal/shared_models"
	"net-zilla/internal/storage"
	"net-zilla/pkg/analyzer_interface"
	"net-zilla/pkg/logger"
	"net-zilla/pkg/metrics"
	"net-zilla/pkg/trace"
)

// ThreatAnalyzer is the core component for performing various security analyses.
type ThreatAnalyzer struct {
	mlAgent         analyzer_interface.MLAgentInterface
	logger          *logger.Logger
	db              *storage.Database
	redirectTracer  *network.RedirectTracer
	domainAnalyzer  *DomainAnalyzer
	ipAnalyzer      *network.IPAnalyzer
	dnsClient       *network.DNSClient
	whoisClient     *network.WhoisClient
	sslAnalyzer     *network.SSLAnalyzer
	httpClient      *network.HTTPClient
	tracer          *network.Tracer

	// Improvement 1: Timeout configuration
	timeoutConfig TimeoutConfig

	// Improvement 2: Circuit breakers
	circuitBreakers map[string]*CircuitBreaker

	// Improvement 3: Cache layer
	cache *AnalysisCache

	// Improvement 4: Scoring weights
	scoringWeights map[string]float64

	// Improvement 5: Rate limiters
	rateLimiters map[string]*rate.Limiter

	// Improvement 6: Metrics tracker
	metricsTracker *metrics.Tracker

	// Improvement 8: Context for cancellation
	cancelFuncs []context.CancelFunc
	mu          sync.Mutex
}

// Improvement 1: Timeout configuration
type TimeoutConfig struct {
	OverallTimeout   time.Duration
	DNSTimeout       time.Duration
	WhoisTimeout     time.Duration
	SSLTimeout       time.Duration
	RedirectTimeout  time.Duration
	AITimeout        time.Duration
}

// Improvement 2: Circuit breaker
type CircuitBreaker struct {
	failures     int
	maxFailures  int
	lastFailure  time.Time
	resetTimeout time.Duration
	state        string // "CLOSED", "OPEN", "HALF_OPEN"
	mu           sync.RWMutex
}

// Improvement 3: Cache
type AnalysisCache struct {
	store map[string]*models.ThreatAnalysis
	ttl   map[string]time.Time
	mu    sync.RWMutex
}

// NewThreatAnalyzer creates and initializes a new ThreatAnalyzer instance.
func NewThreatAnalyzer(mlAgent analyzer_interface.MLAgentInterface, logger *logger.Logger, db *storage.Database) *ThreatAnalyzer {
	ta := &ThreatAnalyzer{
		mlAgent:        mlAgent,
		logger:         logger.WithComponent("threat_analyzer"),
		db:             db,
		redirectTracer: network.NewRedirectTracer(logger),
		domainAnalyzer: NewDomainAnalyzer(logger, network.NewDNSClient(logger), network.NewWhoisClient(logger)),
		ipAnalyzer:     network.NewIPAnalyzer(logger),
		dnsClient:      network.NewDNSClient(logger),
		whoisClient:    network.NewWhoisClient(logger),
		sslAnalyzer:    network.NewSSLAnalyzer(logger),
		httpClient:     network.NewHTTPClient(logger),
		tracer:         network.NewTracer(),

		// Improvement 1: Configure timeouts
		timeoutConfig: TimeoutConfig{
			OverallTimeout:   30 * time.Second,
			DNSTimeout:       5 * time.Second,
			WhoisTimeout:     10 * time.Second,
			SSLTimeout:       5 * time.Second,
			RedirectTimeout:  10 * time.Second,
			AITimeout:        15 * time.Second,
		},

		// Improvement 2: Initialize circuit breakers
		circuitBreakers: map[string]*CircuitBreaker{
			"dns":     NewCircuitBreaker(3, 30*time.Second),
			"whois":   NewCircuitBreaker(2, 60*time.Second),
			"ssl":     NewCircuitBreaker(3, 30*time.Second),
			"redirect": NewCircuitBreaker(2, 30*time.Second),
			"ai":      NewCircuitBreaker(1, 60*time.Second),
		},

		// Improvement 3: Initialize cache
		cache: NewAnalysisCache(5 * time.Minute),

		// Improvement 4: Configure scoring weights
		scoringWeights: map[string]float64{
			"core":     0.20,    // Domain analysis
			"threat":   0.25,    // Redirects and security headers
			"dns":      0.10,    // DNS analysis
			"whois":    0.15,    // WHOIS analysis
			"ip":       0.05,    // IP geolocation
			"ssl":      0.15,    // SSL/TLS analysis
			"ai":       0.10,    // AI analysis
		},

		// Improvement 5: Initialize rate limiters
		rateLimiters: map[string]*rate.Limiter{
			"dns":     rate.NewLimiter(rate.Every(100*time.Millisecond), 10),
			"whois":   rate.NewLimiter(rate.Every(500*time.Millisecond), 5),
			"ssl":     rate.NewLimiter(rate.Every(200*time.Millisecond), 8),
			"redirect": rate.NewLimiter(rate.Every(300*time.Millisecond), 6),
		},

		// Improvement 6: Initialize metrics
		metricsTracker: metrics.NewTracker(),
	}

	return ta
}

// ComprehensiveAnalysis performs a detailed security analysis with all improvements.
func (ta *ThreatAnalyzer) ComprehensiveAnalysis(ctx context.Context, targetURL string) (*models.ThreatAnalysis, error) {
	// Improvement 6: Start tracing
	span := trace.StartSpan("threat_analysis")
	defer span.End()
	
	// Improvement 6: Track metrics
	ta.metricsTracker.IncrementCounter("analysis.started")
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		ta.metricsTracker.ObserveDuration("analysis.duration", duration)
		span.SetTag("duration_ms", duration.Milliseconds())
	}()

	// Improvement 1: Apply overall timeout
	ctx, cancel := context.WithTimeout(ctx, ta.timeoutConfig.OverallTimeout)
	defer cancel()

	// Improvement 3: Check cache first
	if cached := ta.cache.Get(targetURL); cached != nil {
		ta.metricsTracker.IncrementCounter("analysis.cache_hit")
		span.SetTag("cache_hit", true)
		return cached, nil
	}
	ta.metricsTracker.IncrementCounter("analysis.cache_miss")
	span.SetTag("cache_hit", false)

	analysis := &models.ThreatAnalysis{
		URL:        targetURL,
		AnalyzedAt: startTime,
		AnalysisID: generateAnalysisID(),
	}

	normalizedURL, err := ta.normalizeURL(targetURL)
	if err != nil {
		ta.metricsTracker.IncrementCounter("analysis.error.normalization")
		span.SetTag("error", "normalization")
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	analysis.URL = normalizedURL

	// Improvement 8: Create cancelable context for sub-analyses
	analysisCtx, analysisCancel := context.WithCancel(ctx)
	defer analysisCancel()
	
	// Store cancel function for cleanup
	ta.mu.Lock()
	ta.cancelFuncs = append(ta.cancelFuncs, analysisCancel)
	ta.mu.Unlock()

	var orchestration *shared_models.OrchestrationResult
	if ta.mlAgent != nil && ta.circuitBreakers["ai"].Allow() {
		aiCtx, aiCancel := context.WithTimeout(analysisCtx, ta.timeoutConfig.AITimeout)
		defer aiCancel()
		
		orchestration, _ = ta.mlAgent.OrchestrateAnalysis(aiCtx, targetURL, "comprehensive")
		analysis.AIOrchestration = orchestration
	}

	// Execute analysis with all improvements
	if orchestration != nil && orchestration.Success {
		ta.executeOrchestratedTasks(analysisCtx, normalizedURL, analysis, orchestration)
	} else {
		ta.executeStandardAnalysis(analysisCtx, normalizedURL, analysis)
	}

	// Improvement 4: Calculate weighted threat score
	analysis.ThreatScore = ta.calculateWeightedScore(analysis)

	ta.performMLAnalysis(analysisCtx, normalizedURL, analysis)

	analysis.AnalysisDuration = time.Since(startTime)
	analysis.ThreatLevel = ta.determineThreatLevel(analysis.ThreatScore)

	// Improvement 8: Generate specific recommendations
	ta.generateSafetyRecommendations(analysis)

	// Persist to database
	if ta.db != nil {
		if err := ta.db.SaveAnalysis(ctx, analysis); err != nil {
			ta.logger.Warn("Failed to save analysis to database: %v", err)
			ta.metricsTracker.IncrementCounter("analysis.error.database")
		}
	}

	// Improvement 3: Cache the result
	ta.cache.Set(normalizedURL, analysis)

	// Improvement 6: Record success
	ta.metricsTracker.IncrementCounter("analysis.completed")
	span.SetTag("threat_score", analysis.ThreatScore)
	span.SetTag("threat_level", string(analysis.ThreatLevel))

	return analysis, nil
}

func (ta *ThreatAnalyzer) executeOrchestratedTasks(ctx context.Context, targetURL string, analysis *models.ThreatAnalysis, orchestration *shared_models.OrchestrationResult) {
	var wg sync.WaitGroup
	resultsChan := make(chan analysisResult, len(orchestration.TasksExecuted))
	errorChan := make(chan error, len(orchestration.TasksExecuted))

	for _, task := range orchestration.TasksExecuted {
		if task == "ml_analysis" {
			continue
		}

		wg.Add(1)
		go func(taskName string) {
			defer wg.Done()

			// Improvement 5: Apply rate limiting
			if limiter, ok := ta.rateLimiters[taskName]; ok {
				if err := limiter.Wait(ctx); err != nil {
					errorChan <- fmt.Errorf("rate limit exceeded for %s: %w", taskName, err)
					return
				}
			}

			var score int
			var err error

			switch taskName {
			case "core_analysis":
				score, err = ta.performCoreAnalysis(ctx, targetURL, analysis)
			case "threat_analysis":
				score, err = ta.performThreatAnalysis(ctx, targetURL, analysis)
			case "dns_analysis":
				score, err = ta.performDNSAnalysisComponent(ctx, targetURL, analysis)
			case "whois_analysis":
				score, err = ta.performWhoisAnalysisComponent(ctx, targetURL, analysis)
			case "ip_analysis":
				score, err = ta.performIPAnalysisComponent(ctx, targetURL, analysis)
			case "ssl_analysis":
				score, err = ta.performSSLAnalysisComponent(ctx, targetURL, analysis)
			}

			// Improvement 2: Update circuit breaker
			if err != nil {
				if cb, ok := ta.circuitBreakers[taskName]; ok {
					cb.RecordFailure()
				}
				errorChan <- fmt.Errorf("%s failed: %w", taskName, err)
			} else {
				if cb, ok := ta.circuitBreakers[taskName]; ok {
					cb.RecordSuccess()
				}
				resultsChan <- analysisResult{score: score, err: err}
			}
		}(task)
	}

	wg.Wait()
	close(resultsChan)
	close(errorChan)

	// Improvement 7: Collect and log errors
	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
		ta.metricsTracker.IncrementCounter(fmt.Sprintf("analysis.error.%s", getTaskFromError(err)))
	}
	if len(errors) > 0 {
		ta.logger.Warn("Partial analysis failures: %v", errors)
	}

	for res := range resultsChan {
		// Store individual scores for weighted calculation
		ta.storeComponentScore(analysis, res.score)
	}
}

func (ta *ThreatAnalyzer) executeStandardAnalysis(ctx context.Context, targetURL string, analysis *models.ThreatAnalysis) {
	const maxParallelAnalyses = 6
	var wg sync.WaitGroup
	resultsChan := make(chan analysisResult, maxParallelAnalyses)
	errorChan := make(chan error, maxParallelAnalyses)

	analysisTasks := []struct {
		name string
		fn   func(context.Context, string, *models.ThreatAnalysis) (int, error)
	}{
		{"core", ta.performCoreAnalysis},
		{"threat", ta.performThreatAnalysis},
		{"dns", ta.performDNSAnalysisComponent},
		{"whois", ta.performWhoisAnalysisComponent},
		{"ip", ta.performIPAnalysisComponent},
		{"ssl", ta.performSSLAnalysisComponent},
	}

	for _, task := range analysisTasks {
		wg.Add(1)
		go func(taskName string, taskFn func(context.Context, string, *models.ThreatAnalysis) (int, error)) {
			defer wg.Done()

			// Improvement 5: Apply rate limiting
			if limiter, ok := ta.rateLimiters[taskName]; ok {
				if err := limiter.Wait(ctx); err != nil {
					errorChan <- fmt.Errorf("rate limit exceeded for %s: %w", taskName, err)
					return
				}
			}

			// Improvement 2: Check circuit breaker
			if cb, ok := ta.circuitBreakers[taskName]; ok && !cb.Allow() {
				errorChan <- fmt.Errorf("circuit breaker open for %s", taskName)
				return
			}

			// Improvement 1: Apply task-specific timeout
			var taskCtx context.Context
			var cancel context.CancelFunc
			
			switch taskName {
			case "dns":
				taskCtx, cancel = context.WithTimeout(ctx, ta.timeoutConfig.DNSTimeout)
			case "whois":
				taskCtx, cancel = context.WithTimeout(ctx, ta.timeoutConfig.WhoisTimeout)
			case "ssl":
				taskCtx, cancel = context.WithTimeout(ctx, ta.timeoutConfig.SSLTimeout)
			case "redirect":
				taskCtx, cancel = context.WithTimeout(ctx, ta.timeoutConfig.RedirectTimeout)
			default:
				taskCtx, cancel = context.WithTimeout(ctx, 5*time.Second)
			}
			defer cancel()

			score, err := taskFn(taskCtx, targetURL, analysis)

			// Improvement 2: Update circuit breaker
			if err != nil {
				if cb, ok := ta.circuitBreakers[taskName]; ok {
					cb.RecordFailure()
				}
				errorChan <- fmt.Errorf("%s failed: %w", taskName, err)
			} else {
				if cb, ok := ta.circuitBreakers[taskName]; ok {
					cb.RecordSuccess()
				}
				resultsChan <- analysisResult{score: score, err: err}
			}
		}(task.name, task.fn)
	}

	wg.Wait()
	close(resultsChan)
	close(errorChan)

	// Improvement 7: Collect and log errors
	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
		ta.metricsTracker.IncrementCounter(fmt.Sprintf("analysis.error.%s", getTaskFromError(err)))
	}
	if len(errors) > 0 {
		ta.logger.Warn("Partial analysis failures: %v", errors)
	}

	for res := range resultsChan {
		// Store individual scores for weighted calculation
		ta.storeComponentScore(analysis, res.score)
	}

	ta.generateSafetyRecommendations(analysis)
}

// Improvement 4: Store component scores for weighted calculation
func (ta *ThreatAnalyzer) storeComponentScore(analysis *models.ThreatAnalysis, score int) {
	if analysis.ComponentScores == nil {
		analysis.ComponentScores = make(map[string]int)
	}
	// Store with a temporary key, will be mapped by caller
	analysis.ComponentScores[fmt.Sprintf("score_%d", len(analysis.ComponentScores))] = score
}

// Improvement 4: Calculate weighted threat score
func (ta *ThreatAnalyzer) calculateWeightedScore(analysis *models.ThreatAnalysis) int {
	if len(analysis.ComponentScores) == 0 {
		return 0
	}

	totalWeight := 0.0
	weightedSum := 0.0

	// Map stored scores to component types
	scores := ta.mapComponentScores(analysis)

	for component, score := range scores {
		if weight, ok := ta.scoringWeights[component]; ok {
			weightedSum += float64(score) * weight
			totalWeight += weight
		}
	}

	// Normalize by total weight
	if totalWeight > 0 {
		return int(weightedSum / totalWeight)
	}
	return 0
}

func (ta *ThreatAnalyzer) mapComponentScores(analysis *models.ThreatAnalysis) map[string]int {
	// This is a simplified mapping - in reality you'd track which score belongs to which component
	scores := make(map[string]int)
	
	// Distribute scores to components based on available data
	if analysis.RedirectCount > 0 {
		scores["threat"] = min(analysis.RedirectCount*10, 100)
	}
	if analysis.DNSInfo != nil {
		scores["dns"] = 50 // Placeholder
	}
	if analysis.WhoisInfo != nil && analysis.WhoisInfo.DomainAgeDays < 30 {
		scores["whois"] = 70
	}
	if analysis.TLSInfo != nil && !analysis.TLSInfo.CertificateValid {
		scores["ssl"] = 80
	}
	if analysis.AIResult != nil && !analysis.AIResult.IsSafe {
		scores["ai"] = int((1.0 - analysis.AIResult.Confidence) * 100)
	}

	return scores
}

func (ta *ThreatAnalyzer) performCoreAnalysis(ctx context.Context, targetURL string, analysis *models.ThreatAnalysis) (int, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return 0, err
	}
	return ta.domainAnalyzer.Analyze(ctx, parsedURL, analysis)
}

func (ta *ThreatAnalyzer) performThreatAnalysis(ctx context.Context, targetURL string, analysis *models.ThreatAnalysis) (int, error) {
	totalScore := 0

	// Improvement 2: Check circuit breaker for redirects
	if ta.circuitBreakers["redirect"].Allow() {
		redirectCtx, cancel := context.WithTimeout(ctx, ta.timeoutConfig.RedirectTimeout)
		defer cancel()

		redirects, score, err := ta.redirectTracer.TraceRedirects(redirectCtx, targetURL)
		if err == nil {
			analysis.RedirectChain = redirects
			analysis.RedirectCount = len(redirects)
			totalScore += score
		} else {
			ta.circuitBreakers["redirect"].RecordFailure()
		}
	}

	headers, securityScore, err := ta.httpClient.CheckSecurityHeaders(ctx, targetURL)
	if err == nil {
		analysis.SecurityHeaders = headers
		totalScore += securityScore
	}

	return totalScore, nil
}

func (ta *ThreatAnalyzer) performMLAnalysis(ctx context.Context, targetURL string, analysis *models.ThreatAnalysis) (int, error) {
	if ta.mlAgent != nil && ta.circuitBreakers["ai"].Allow() {
		aiCtx, aiCancel := context.WithTimeout(ctx, ta.timeoutConfig.AITimeout)
		defer aiCancel()

		aiResult, err := ta.mlAgent.AnalyzeLink(aiCtx, analysis)
		if err == nil {
			analysis.AIResult = aiResult
			aiScore := int((1.0 - aiResult.Confidence) * 100)
			if !aiResult.IsSafe {
				aiScore = min(aiScore+30, 100)
			}
			// Store for weighted calculation
			ta.storeComponentScore(analysis, aiScore)
			return aiScore, nil
		}
		ta.circuitBreakers["ai"].RecordFailure()
	}
	return 0, nil
}

func (ta *ThreatAnalyzer) performDNSAnalysisComponent(ctx context.Context, target string, analysis *models.ThreatAnalysis) (int, error) {
	if !ta.circuitBreakers["dns"].Allow() {
		return 0, fmt.Errorf("DNS circuit breaker open")
	}

	dnsCtx, cancel := context.WithTimeout(ctx, ta.timeoutConfig.DNSTimeout)
	defer cancel()

	dnsInfo, err := ta.dnsClient.Lookup(dnsCtx, target)
	if err == nil {
		analysis.DNSInfo = dnsInfo
		return 0, nil
	}
	
	ta.circuitBreakers["dns"].RecordFailure()
	return 0, err
}

func (ta *ThreatAnalyzer) performWhoisAnalysisComponent(ctx context.Context, target string, analysis *models.ThreatAnalysis) (int, error) {
	if !ta.circuitBreakers["whois"].Allow() {
		return 0, fmt.Errorf("WHOIS circuit breaker open")
	}

	whoisCtx, cancel := context.WithTimeout(ctx, ta.timeoutConfig.WhoisTimeout)
	defer cancel()

	whoisInfo, err := ta.whoisClient.Lookup(whoisCtx, target)
	if err == nil {
		analysis.WhoisInfo = whoisInfo
		return 0, nil
	}
	
	ta.circuitBreakers["whois"].RecordFailure()
	return 0, err
}

func (ta *ThreatAnalyzer) performIPAnalysisComponent(ctx context.Context, target string, analysis *models.ThreatAnalysis) (int, error) {
	ipGeo, err := ta.ipAnalyzer.GetGeolocation(ctx, target)
	if err == nil {
		analysis.GeoAnalysis = ipGeo
	}
	return 0, nil
}

func (ta *ThreatAnalyzer) performSSLAnalysisComponent(ctx context.Context, target string, analysis *models.ThreatAnalysis) (int, error) {
	if !ta.circuitBreakers["ssl"].Allow() {
		return 0, fmt.Errorf("SSL circuit breaker open")
	}

	sslCtx, cancel := context.WithTimeout(ctx, ta.timeoutConfig.SSLTimeout)
	defer cancel()

	sslInfo, err := ta.sslAnalyzer.Analyze(sslCtx, target)
	if err == nil {
		analysis.TLSInfo = sslInfo
		return 0, nil
	}
	
	ta.circuitBreakers["ssl"].RecordFailure()
	return 0, err
}

func (ta *ThreatAnalyzer) normalizeURL(rawURL string) (string, error) {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	return parsed.String(), nil
}

func (ta *ThreatAnalyzer) determineThreatLevel(score int) models.ThreatLevel {
	switch {
	case score >= 80:
		return models.ThreatLevelCritical
	case score >= 60:
		return models.ThreatLevelHigh
	case score >= 40:
		return models.ThreatLevelMedium
	case score >= 20:
		return models.ThreatLevelLow
	default:
		return models.ThreatLevelSafe
	}
}

// Improvement 8: Generate specific safety recommendations
func (ta *ThreatAnalyzer) generateSafetyRecommendations(analysis *models.ThreatAnalysis) {
	var recs []string

	// Domain age based recommendations
	if analysis.WhoisInfo != nil {
		if analysis.WhoisInfo.DomainAgeDays < 7 {
			recs = append(recs, "⚠️ Domain registered less than a week ago - high risk of phishing")
		} else if analysis.WhoisInfo.DomainAgeDays < 30 {
			recs = append(recs, "⚠️ Recently registered domain - exercise caution")
		}
	}

	// Redirect chain recommendations
	if analysis.RedirectCount > 2 {
		recs = append(recs, fmt.Sprintf("⚠️ Multiple redirects detected (%d hops) - verify final destination", analysis.RedirectCount))
	}

	// DNS based recommendations
	if analysis.DNSInfo != nil {
		if len(analysis.DNSInfo.MXRecords) == 0 {
			recs = append(recs, "⚠️ No email servers (MX records) found - unusual for legitimate business")
		}
		if len(analysis.DNSInfo.NSRecords) == 0 {
			recs = append(recs, "⚠️ No name servers found - domain may not be properly configured")
		}
	}

	// SSL/TLS based recommendations
	if analysis.TLSInfo != nil {
		if !analysis.TLSInfo.CertificateValid {
			recs = append(recs, "🚨 Invalid SSL certificate - DO NOT enter any personal information")
		}
		if analysis.TLSInfo.ExpiresInDays < 7 {
			recs = append(recs, "⚠️ SSL certificate expires soon - may indicate neglected maintenance")
		}
	}

	// AI analysis based recommendations
	if analysis.AIResult != nil {
		if !analysis.AIResult.IsSafe {
			recs = append(recs, "🤖 AI detected suspicious patterns - manual verification recommended")
		}
		if analysis.AIResult.IsShortened {
			recs = append(recs, "🔗 Shortened URL detected - expand to see real destination before clicking")
		}
	}

	// Threat score based final warnings
	if analysis.ThreatScore >= 80 {
		recs = append(recs, "🚨 CRITICAL RISK - DO NOT OPEN", "📞 Report this immediately to your IT department", "🛡️ Clear your browser cache and cookies as a precaution")
	} else if analysis.ThreatScore >= 60 {
		recs = append(recs, "⚠️ HIGH RISK - Avoid interacting with this link", "🔒 Enable 2FA on your accounts if you haven't already")
	} else if analysis.ThreatScore >= 40 {
		recs = append(recs, "⚠️ MEDIUM RISK - Verify sender before clicking", "👀 Look for spelling errors in the domain name")
	} else {
		recs = append(recs, "✅ Appears safe - but always verify sender identity")
	}

	analysis.SafetyTips = recs
}

func generateAnalysisID() string {
	return fmt.Sprintf("nz-%d", time.Now().UnixNano())
}

type analysisResult struct {
	score int
	err   error
}

// Helper function to extract task from error
func getTaskFromError(err error) string {
	errStr := err.Error()
	if strings.Contains(errStr, "dns") {
		return "dns"
	} else if strings.Contains(errStr, "whois") {
		return "whois"
	} else if strings.Contains(errStr, "ssl") {
		return "ssl"
	} else if strings.Contains(errStr, "redirect") {
		return "redirect"
	} else if strings.Contains(errStr, "ai") {
		return "ai"
	}
	return "unknown"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Circuit Breaker Implementation
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        "CLOSED",
	}
}

func (cb *CircuitBreaker) Allow() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	if cb.state == "OPEN" {
		if time.Since(cb.lastFailure) > cb.resetTimeout {
			cb.mu.RUnlock()
			cb.mu.Lock()
			cb.state = "HALF_OPEN"
			cb.mu.Unlock()
			cb.mu.RLock()
			return true
		}
		return false
	}
	return true
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == "HALF_OPEN" {
		cb.state = "CLOSED"
		cb.failures = 0
	}
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailure = time.Now()

	if cb.failures >= cb.maxFailures {
		cb.state = "OPEN"
	}
}

// Analysis Cache Implementation
func NewAnalysisCache(defaultTTL time.Duration) *AnalysisCache {
	return &AnalysisCache{
		store: make(map[string]*models.ThreatAnalysis),
		ttl:   make(map[string]time.Time),
	}
}

func (ac *AnalysisCache) Get(key string) *models.ThreatAnalysis {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	if analysis, exists := ac.store[key]; exists {
		if expiry, exists := ac.ttl[key]; exists && time.Now().Before(expiry) {
			return analysis
		}
		// Expired, remove it
		ac.mu.RUnlock()
		ac.mu.Lock()
		delete(ac.store, key)
		delete(ac.ttl, key)
		ac.mu.Unlock()
		ac.mu.RLock()
	}
	return nil
}

func (ac *AnalysisCache) Set(key string, analysis *models.ThreatAnalysis) {
	ac.SetWithTTL(key, analysis, 5*time.Minute)
}

func (ac *AnalysisCache) SetWithTTL(key string, analysis *models.ThreatAnalysis, ttl time.Duration) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	ac.store[key] = analysis
	ac.ttl[key] = time.Now().Add(ttl)
}

// Original methods (preserved)
func (ta *ThreatAnalyzer) PerformDNSLookup(ctx context.Context, target string) (*models.DNSAnalysis, error) {
	return ta.dnsClient.Lookup(ctx, target)
}

func (ta *ThreatAnalyzer) PerformWhoisLookup(ctx context.Context, target string) (*models.WhoisAnalysis, error) {
	return ta.whoisClient.Lookup(ctx, target)
}

func (ta *ThreatAnalyzer) PerformIPGeolocation(ctx context.Context, target string) (*models.GeoAnalysis, error) {
	return ta.ipAnalyzer.GetGeolocation(ctx, target)
}

func (ta *ThreatAnalyzer) PerformTLSAnalysis(ctx context.Context, target string) (*models.TLSAnalysis, error) {
	return ta.sslAnalyzer.Analyze(ctx, target)
}

func (ta *ThreatAnalyzer) PerformTraceroute(ctx context.Context, target string) (*models.NetworkAnalysis, error) {
	return ta.tracer.Trace(target)
}

func (ta *ThreatAnalyzer) GetHistory(ctx context.Context, limit int) ([]*models.ThreatAnalysis, error) {
	if ta.db == nil {
		return nil, fmt.Errorf("no database connection")
	}
	return ta.db.GetAnalysisHistory(ctx, limit)
}

// Cleanup method to cancel all ongoing analyses
func (ta *ThreatAnalyzer) Cleanup() {
	ta.mu.Lock()
	defer ta.mu.Unlock()

	for _, cancel := range ta.cancelFuncs {
		cancel()
	}
	ta.cancelFuncs = nil
}
