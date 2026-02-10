package utils

import "testing"

func TestColorConstants(t *testing.T) {
	if ColorOrange == "" {
		t.Error("ColorOrange should not be empty")
	}
	if ColorReset == "" {
		t.Error("ColorReset should not be empty")
	}
}

func TestDisplayBanner(t *testing.T) {
	// Simple smoke test to ensure no panic
	DisplayBanner()
}
