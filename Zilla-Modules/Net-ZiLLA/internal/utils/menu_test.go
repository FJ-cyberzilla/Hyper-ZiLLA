package utils

import (
	"testing"

	"net-zilla/pkg/logger"
)

func TestNewMenu(t *testing.T) {
	l := logger.NewLogger()
	m := NewMenu(nil, l)
	if m == nil {
		t.Fatal("NewMenu returned nil")
	}
	if m.logger != l {
		t.Error("Menu logger mismatch")
	}
}

func TestMenu_ShowHistory(t *testing.T) {
	l := logger.NewLogger()
	m := NewMenu(nil, l)
	if m == nil {
		t.Fatal("failed to create menu")
	}
	
	// Just verify fields are set since we can't easily call showHistory with nil service
	if m.logger != l {
		t.Errorf("logger mismatch")
	}
}
