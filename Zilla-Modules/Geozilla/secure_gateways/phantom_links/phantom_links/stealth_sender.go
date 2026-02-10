// secure_gateways/phantom_links/stealth_sender.go
package phantom_links

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"time"
)

type StealthSender struct {
	OriginObfuscation bool
	RotatingUserAgent bool
	RequestRandomization bool
	TorIntegration    bool
}

// SendStealthRequest avoids spam detection
func (s *StealthSender) SendStealthRequest(payload []byte, targetURL string) error {
	// Rotate through multiple techniques
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	// Create stealth request
	req, err := http.NewRequest("POST", targetURL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	// Apply stealth headers
	s.applyStealthHeaders(req)

	// Random delays to avoid rate limiting
	s.randomizedDelay()

	// Execute with retry logic
	return s.executeWithRetry(client, req)
}

// applyStealthHeaders uses legitimate-looking headers
func (s *StealthSender) applyStealthHeaders(req *http.Request) {
	// Rotating user agents
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
	}

	req.Header.Set("User-Agent", userAgents[time.Now().Unix()%int64(len(userAgents))])
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Referer", "https://example.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
}

// executeWithRetry handles failures intelligently
func (s *StealthSender) executeWithRetry(client *http.Client, req *http.Request) error {
	maxRetries := 3
	baseDelay := time.Second * 2

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err := client.Do(req)
		if err != nil {
			if attempt == maxRetries-1 {
				return fmt.Errorf("final attempt failed: %v", err)
			}
			
			// Exponential backoff with jitter
			delay := baseDelay * time.Duration(1<<attempt)
			jitter := time.Duration(rand.Int63n(1000)) * time.Millisecond
			time.Sleep(delay + jitter)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil // Success
		}

		// Handle specific HTTP status codes
		switch resp.StatusCode {
		case 429: // Rate limited
			retryAfter := resp.Header.Get("Retry-After")
			if retryAfter != "" {
				if seconds, err := strconv.Atoi(retryAfter); err == nil {
					time.Sleep(time.Duration(seconds) * time.Second)
				}
			}
		case 403: // Forbidden - change approach
			s.rotateStealthTechnique()
		}

		time.Sleep(time.Duration(1<<attempt) * time.Second)
	}

	return fmt.Errorf("max retries exceeded")
}
