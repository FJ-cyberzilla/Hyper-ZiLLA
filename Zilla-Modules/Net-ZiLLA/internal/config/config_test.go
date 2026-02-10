package config

import "testing"

func TestConfig_Load(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.AI.ConfidenceThreshold <= 0 {
		t.Errorf("expected valid confidence threshold, got %f", cfg.AI.ConfidenceThreshold)
	}
}
