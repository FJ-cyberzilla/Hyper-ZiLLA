package threat_intel

import (
	"context"
	"net/http"
	"time"
)

type IntelManager struct {
	vtKey      string
	abuseKey   string
	avKey      string
	httpClient *http.Client
}

func NewIntelManager(vt, abuse, av string) *IntelManager {
	return &IntelManager{
		vtKey:    vt,
		abuseKey: abuse,
		avKey:    av,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type GlobalReputation struct {
	TotalEngineCount int
	Positives        int
	Malicious        bool
	Details          map[string]string
}

func (im *IntelManager) MultiCheck(ctx context.Context, indicator string) GlobalReputation {
	res := GlobalReputation{Details: make(map[string]string)}

	// In production, these would be real HTTP calls to APIs
	// For this logic, we implement a structured mock that respects context
	
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	results := make(chan string, 2)

	if im.vtKey != "" {
		go func() {
			// Simulate VT Check
			time.Sleep(100 * time.Millisecond)
			results <- "VT: Malicious"
		}()
	} else {
		results <- "VT: Skipped (No Key)"
	}

	if im.abuseKey != "" {
		go func() {
			// Simulate AbuseIPDB Check
			time.Sleep(150 * time.Millisecond)
			results <- "AbuseIPDB: Clean"
		}()
	} else {
		results <- "AbuseIPDB: Skipped (No Key)"
	}

	for i := 0; i < 2; i++ {
		select {
		case r := <-results:
			if r == "VT: Malicious" {
				res.Positives++
				res.Malicious = true
			}
			res.TotalEngineCount++
		case <-ctx.Done():
			res.Details["error"] = "threat intel timeout"
			return res
		}
	}

	return res
}

