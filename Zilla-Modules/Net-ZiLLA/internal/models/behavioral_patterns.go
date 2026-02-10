package models

import "time"

// PatternType defines the nature of the detected behavior
type PatternType string

const (
	PatternPhishing PatternType = "PHISHING"
	PatternMalware  PatternType = "MALWARE"
	PatternTracking PatternType = "TRACKING"
	PatternEvasion  PatternType = "EVASION"
)

// BehavioralPattern describes a detected security pattern
type BehavioralPattern struct {
	Type        PatternType `json:"type"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Weight      int         `json:"weight"` // Impact on final score
	Evidence    []string    `json:"evidence"`
	Details     string      `json:"details,omitempty"`
	Match       string      `json:"match,omitempty"`
	Category    string      `json:"category,omitempty"`
}

// BehaviorAnalysis stores all patterns found during a scan
type BehaviorAnalysis struct {
	Patterns      []BehavioralPattern `json:"patterns"`
	Anomalies     []string            `json:"anomalies"`
	RiskSignature string              `json:"risk_signature"`
	RiskScore     int                 `json:"risk_score"`
	Severity      string              `json:"severity"`
	Confidence    float64             `json:"confidence"`
	ScriptType    string              `json:"script_type,omitempty"`
	Timestamp     time.Time           `json:"timestamp,omitempty"`
	ContentLength int64               `json:"content_length,omitempty"`
	FileSize      int64               `json:"file_size,omitempty"`
	FileName      string              `json:"file_name,omitempty"`
	FileHashes    map[string]string   `json:"file_hashes,omitempty"`
}