package models

import "time"

// IOCType defines the category of an Indicator of Compromise
type IOCType string

const (
	IOCTypeIP     IOCType = "IP"
	IOCTypeDomain IOCType = "DOMAIN"
	IOCTypeURL    IOCType = "URL"
	IOCTypeHash   IOCType = "HASH"
)

// Indicator represents a single threat indicator
type Indicator struct {
	Type                   IOCType   `json:"type"`
	Value                  string    `json:"value"`
	Source                 string    `json:"source"`
	Confidence             float64   `json:"confidence"`
	Severity               string    `json:"severity"`
	LastSeen               time.Time `json:"last_seen"`
	FirstSeen              time.Time `json:"first_seen"`
	Description            string    `json:"description"`
	Country                string    `json:"country,omitempty"`
	ISP                    string    `json:"isp,omitempty"`
	Domain                 string    `json:"domain,omitempty"`
	Reports                int       `json:"reports,omitempty"`
	Tags                   []string  `json:"tags,omitempty"`
	References             []string  `json:"references,omitempty"`
	AbuseConfidenceScore   float64   `json:"abuse_confidence_score,omitempty"`
}

// ThreatDBStats provides metrics for the threat database
type ThreatDBStats struct {
	TotalIndicators  int64            `json:"total_indicators"`
	CountByType      map[string]int64 `json:"count_by_type"`
	CountBySeverity  map[string]int64 `json:"count_by_severity"`
	RecentActivity7d int64            `json:"recent_activity_7d"`
	CacheHits        int64            `json:"cache_hits"`
	CacheMisses      int64            `json:"cache_misses"`
	TotalQueries     int64            `json:"total_queries"`
	Timestamp        time.Time        `json:"timestamp"`
}

// IOCRegistry holds a collection of detected indicators
type IOCRegistry struct {
	Indicators []Indicator `json:"indicators"`
	TotalFound int         `json:"total_found"`
}
