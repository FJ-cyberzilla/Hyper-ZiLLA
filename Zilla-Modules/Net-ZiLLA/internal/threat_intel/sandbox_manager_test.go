package threat_intel

import (
	"context"
	"testing"
)

func TestSandboxManager(t *testing.T) {
	sm := NewSandboxManager()
	
	ctx := context.Background()
	url := "https://example.com"
	
	// Default mode
	id1, err := sm.SpinUpIsolatedBrowser(ctx, url)
	if err != nil {
		t.Errorf("default mode failed: %v", err)
	}
	
	// Hybrid mode
	sm.SetMode("hybrid")
	id2, err := sm.SpinUpIsolatedBrowser(ctx, url)
	if err != nil {
		t.Errorf("hybrid mode failed: %v", err)
	}
	
	if id1 == "" || id2 == "" {
		t.Error("expected non-empty container IDs")
	}
	
	sm.DestroySandbox(id1)
}
