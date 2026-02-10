package threat_intel

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"net-zilla/internal/models"
)

// Provider represents a threat intelligence provider
type Provider struct {
	Name    string
	BaseURL string
	APIKey  string
	Enabled bool
}

// FeedCache stores fetched indicators with TTL
type FeedCache struct {
	mu          sync.RWMutex
	indicators  map[string]models.Indicator
	lastUpdated time.Time
	ttl         time.Duration
}

// FeedsClient manages connections to external threat intelligence providers
type FeedsClient struct {
	client     *http.Client
	providers  []Provider
	cache      *FeedCache
	maxWorkers int
	timeout    time.Duration
}

// NewFeedsClient creates a new threat intelligence client
func NewFeedsClient(apiKeys map[string]string, timeout time.Duration) *FeedsClient {
	// Configure providers
	providers := []Provider{
		{
			Name:    "AbuseIPDB",
			BaseURL: "https://api.abuseipdb.com/api/v2",
			APIKey:  apiKeys["abuseipdb"],
			Enabled: apiKeys["abuseipdb"] != "",
		},
		{
			Name:    "VirusTotal",
			BaseURL: "https://www.virustotal.com/api/v3",
			APIKey:  apiKeys["virustotal"],
			Enabled: apiKeys["virustotal"] != "",
		},
		{
			Name:    "AlienVaultOTX",
			BaseURL: "https://otx.alienvault.com/api/v1",
			APIKey:  apiKeys["alienvault"],
			Enabled: apiKeys["alienvault"] != "",
		},
	}

	return &FeedsClient{
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxConnsPerHost:     10,
				IdleConnTimeout:     90 * time.Second,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		},
		providers:  providers,
		cache:      newFeedCache(15 * time.Minute),
		maxWorkers: 5,
		timeout:    timeout,
	}
}

// newFeedCache creates a new feed cache with TTL
func newFeedCache(ttl time.Duration) *FeedCache {
	return &FeedCache{
		indicators: make(map[string]models.Indicator),
		ttl:        ttl,
	}
}

// FetchUpdates fetches threat indicators from all enabled providers
func (fc *FeedsClient) FetchUpdates(ctx context.Context) ([]models.Indicator, error) {
	// Check cache first
	if indicators := fc.getCachedIndicators(); len(indicators) > 0 {
		return indicators, nil
	}

	var allIndicators []models.Indicator
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(fc.providers))
	semaphore := make(chan struct{}, fc.maxWorkers)

	for _, provider := range fc.providers {
		if !provider.Enabled {
			continue
		}

		wg.Add(1)
		go func(p Provider) {
			defer wg.Done()
			
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			default:
				indicators, err := fc.fetchFromProvider(ctx, p)
				if err != nil {
					errChan <- fmt.Errorf("%s: %w", p.Name, err)
					return
				}

				mu.Lock()
				allIndicators = append(allIndicators, indicators...)
				mu.Unlock()
			}
		}(provider)
	}

	wg.Wait()
	close(errChan)

	// Collect errors (non-blocking, continue with successful fetches)
	var errs []string
	for err := range errChan {
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	// Update cache
	fc.updateCache(allIndicators)

	if len(errs) > 0 {
		return allIndicators, fmt.Errorf("partial fetch completed with errors: %s", strings.Join(errs, "; "))
	}

	return allIndicators, nil
}

// CheckIP checks an IP address against all enabled providers
func (fc *FeedsClient) CheckIP(ctx context.Context, ip string) ([]models.ReputationSource, error) {
	// Validate IP format
	if !isValidIP(ip) {
		return nil, fmt.Errorf("invalid IP address: %s", ip)
	}

	var results []models.ReputationSource
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(fc.providers))

	for _, provider := range fc.providers {
		if !provider.Enabled {
			continue
		}

		wg.Add(1)
		go func(p Provider) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			default:
				result, err := fc.checkIPWithProvider(ctx, p, ip)
				if err != nil {
					errChan <- fmt.Errorf("%s: %w", p.Name, err)
					return
				}

				mu.Lock()
				results = append(results, *result)
				mu.Unlock()
			}
		}(provider)
	}

	wg.Wait()
	close(errChan)

	// Collect errors
	var errs []string
	for err := range errChan {
		if err != nil && err != context.Canceled {
			errs = append(errs, err.Error())
		}
	}

	if len(results) == 0 && len(errs) > 0 {
		return nil, fmt.Errorf("all providers failed: %s", strings.Join(errs, "; "))
	}

	return results, nil
}

// fetchFromProvider fetches indicators from a specific provider
func (fc *FeedsClient) fetchFromProvider(ctx context.Context, provider Provider) ([]models.Indicator, error) {
	switch provider.Name {
	case "AbuseIPDB":
		return fc.fetchAbuseIPDB(ctx, provider)
	case "VirusTotal":
		return fc.fetchVirusTotal(ctx, provider)
	case "AlienVaultOTX":
		return fc.fetchAlienVaultOTX(ctx, provider)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider.Name)
	}
}

// checkIPWithProvider checks IP with a specific provider
func (fc *FeedsClient) checkIPWithProvider(ctx context.Context, provider Provider, ip string) (*models.ReputationSource, error) {
	switch provider.Name {
	case "AbuseIPDB":
		return fc.checkIPWithAbuseIPDB(ctx, provider, ip)
	case "VirusTotal":
		return fc.checkIPWithVirusTotal(ctx, provider, ip)
	case "AlienVaultOTX":
		return fc.checkIPWithAlienVaultOTX(ctx, provider, ip)
	default:
		return &models.ReputationSource{
			Provider:  provider.Name,
			Status:    "Unknown",
			Score:     -1,
			CheckedAt: time.Now(),
		}, nil
	}
}

// AbuseIPDB Implementation
func (fc *FeedsClient) fetchAbuseIPDB(ctx context.Context, provider Provider) ([]models.Indicator, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", 
		fmt.Sprintf("%s/blacklist", provider.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Key", provider.APIKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Net-ZiLLA/1.0")

	resp, err := fc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []struct {
			IPAddress           string `json:"ipAddress"`
			CountryCode         string `json:"countryCode"`
			AbuseConfidenceScore int    `json:"abuseConfidenceScore"`
			LastReportedAt     string `json:"lastReportedAt"`
			ISP                string `json:"isp"`
			Domain             string `json:"domain"`
			TotalReports       int    `json:"totalReports"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}

	var indicators []models.Indicator
	for _, item := range result.Data {
		lastReported, _ := time.Parse(time.RFC3339, item.LastReportedAt)
		
		indicators = append(indicators, models.Indicator{
			Value:       item.IPAddress,
			Type:        "ipv4",
			Source:      "AbuseIPDB",
			Confidence:  float64(item.AbuseConfidenceScore) / 100.0,
			Severity:    calculateSeverity(item.AbuseConfidenceScore),
			LastSeen:    lastReported,
			Country:     item.CountryCode,
			ISP:         item.ISP,
			Domain:      item.Domain,
			Reports:     item.TotalReports,
			Tags:        []string{"malicious", "abuse"},
		})
	}

	return indicators, nil
}

func (fc *FeedsClient) checkIPWithAbuseIPDB(ctx context.Context, provider Provider, ip string) (*models.ReputationSource, error) {
	reqURL := fmt.Sprintf("%s/check?ipAddress=%s", provider.BaseURL, url.QueryEscape(ip))
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Key", provider.APIKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Net-ZiLLA/1.0")

	resp, err := fc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data struct {
			IPAddress           string `json:"ipAddress"`
			IsPublic            bool   `json:"isPublic"`
			IPVersion           int    `json:"ipVersion"`
			IsWhitelisted       bool   `json:"isWhitelisted"`
			AbuseConfidenceScore int    `json:"abuseConfidenceScore"`
			CountryCode         string `json:"countryCode"`
			UsageType           string `json:"usageType"`
			ISP                 string `json:"isp"`
			Domain              string `json:"domain"`
			Hostnames           []string `json:"hostnames"`
			TotalReports        int    `json:"totalReports"`
			NumDistinctUsers    int    `json:"numDistinctUsers"`
			LastReportedAt      string `json:"lastReportedAt"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}

	lastReported, _ := time.Parse(time.RFC3339, result.Data.LastReportedAt)
	
	return &models.ReputationSource{
		Provider:      "AbuseIPDB",
		IP:            result.Data.IPAddress,
		Status:        getAbuseIPDBStatus(result.Data.AbuseConfidenceScore),
		Score:         result.Data.AbuseConfidenceScore,
		Country:       result.Data.CountryCode,
		ISP:           result.Data.ISP,
		Domain:        result.Data.Domain,
		Hostnames:     result.Data.Hostnames,
		Reports:       result.Data.TotalReports,
		Users:         result.Data.NumDistinctUsers,
		LastSeen:      lastReported,
		IsPublic:      result.Data.IsPublic,
		IsWhitelisted: result.Data.IsWhitelisted,
		CheckedAt:     time.Now(),
	}, nil
}

// Helper functions
func (fc *FeedsClient) getCachedIndicators() []models.Indicator {
	fc.cache.mu.RLock()
	defer fc.cache.mu.RUnlock()

	// Check if cache is valid
	if time.Since(fc.cache.lastUpdated) > fc.cache.ttl {
		return nil
	}

	indicators := make([]models.Indicator, 0, len(fc.cache.indicators))
	for _, indicator := range fc.cache.indicators {
		indicators = append(indicators, indicator)
	}

	return indicators
}

func (fc *FeedsClient) updateCache(indicators []models.Indicator) {
	fc.cache.mu.Lock()
	defer fc.cache.mu.Unlock()

	fc.cache.indicators = make(map[string]models.Indicator)
	for _, indicator := range indicators {
		fc.cache.indicators[indicator.Value] = indicator
	}
	fc.cache.lastUpdated = time.Now()
}

func isValidIP(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}
	
	for _, part := range parts {
		if len(part) == 0 || len(part) > 3 {
			return false
		}
		for _, ch := range part {
			if ch < '0' || ch > '9' {
				return false
			}
		}
		if num := atoi(part); num < 0 || num > 255 {
			return false
		}
	}
	return true
}

func atoi(s string) int {
	n := 0
	for _, ch := range s {
		n = n*10 + int(ch-'0')
	}
	return n
}

func calculateSeverity(confidence int) string {
	switch {
	case confidence >= 80:
		return "Critical"
	case confidence >= 60:
		return "High"
	case confidence >= 40:
		return "Medium"
	case confidence >= 20:
		return "Low"
	default:
		return "Info"
	}
}

func getAbuseIPDBStatus(confidence int) string {
	switch {
	case confidence >= 80:
		return "Malicious"
	case confidence >= 60:
		return "Suspicious"
	case confidence >= 40:
		return "Suspicious"
	case confidence >= 20:
		return "Monitor"
	default:
		return "Clean"
	}
}

// Placeholder implementations for other providers
func (fc *FeedsClient) fetchVirusTotal(ctx context.Context, provider Provider) ([]models.Indicator, error) {
	// Implement VirusTotal API integration
	return []models.Indicator{}, nil
}

func (fc *FeedsClient) fetchAlienVaultOTX(ctx context.Context, provider Provider) ([]models.Indicator, error) {
	// Implement AlienVault OTX API integration
	return []models.Indicator{}, nil
}

func (fc *FeedsClient) checkIPWithVirusTotal(ctx context.Context, provider Provider, ip string) (*models.ReputationSource, error) {
	// Implement VirusTotal IP check
	return &models.ReputationSource{
		Provider: "VirusTotal",
		Status:   "Not Implemented",
		Score:    -1,
		CheckedAt: time.Now(),
	}, nil
}

func (fc *FeedsClient) checkIPWithAlienVaultOTX(ctx context.Context, provider Provider, ip string) (*models.ReputationSource, error) {
	// Implement AlienVault OTX IP check
	return &models.ReputationSource{
		Provider: "AlienVault OTX",
		Status:   "Not Implemented",
		Score:    -1,
		CheckedAt: time.Now(),
	}, nil
}
