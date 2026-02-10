package models

import "time"

// ReputationSource details results from a specific threat intel provider
type ReputationSource struct {
	Provider      string    `json:"provider"`
	IP            string    `json:"ip,omitempty"`
	Score         int       `json:"score"`  // 0-100
	Status        string    `json:"status"` // Clean, Malicious, Suspicious
	Country       string    `json:"country,omitempty"`
	ISP           string    `json:"isp,omitempty"`
	Domain        string    `json:"domain,omitempty"`
	Hostnames     []string  `json:"hostnames,omitempty"`
	Reports       int       `json:"reports,omitempty"`
	Users         int       `json:"users,omitempty"`
	LastSeen      time.Time `json:"last_seen,omitempty"`
	IsPublic      bool      `json:"is_public"`
	IsWhitelisted bool      `json:"is_whitelisted"`
	CheckedAt     time.Time `json:"checked_at"`
	Details       string    `json:"details,omitempty"`
}

// ReputationSummary aggregates scores from multiple providers
type ReputationSummary struct {
	AggregateScore int                `json:"aggregate_score"`
	Verdict        string             `json:"verdict"`
	Sources        []ReputationSource `json:"sources"`
	Blacklisted    bool               `json:"blacklisted"`
}
