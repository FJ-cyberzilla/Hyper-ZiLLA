package threat_intel

import (
	"context"
	"testing"
	"time"
)

func TestIntelManager_MultiCheck_NoKeys(t *testing.T) {
	im := NewIntelManager("", "", "")
	res := im.MultiCheck(context.Background(), "example.com")
	
	if res.TotalEngineCount != 2 {
		t.Errorf("expected 2 engines checked (even if skipped), got %d", res.TotalEngineCount)
	}
	if res.Malicious {
		t.Error("expected not malicious with no keys")
	}
}

func TestIntelManager_MultiCheck_WithKeys(t *testing.T) {
	// Our "mock" implementation in manager.go checks if keys are not empty
	im := NewIntelManager("fake_vt", "fake_abuse", "")
	res := im.MultiCheck(context.Background(), "example.com")
	
	if res.TotalEngineCount != 2 {
		t.Errorf("expected 2 engines checked, got %d", res.TotalEngineCount)
	}
	// Our mock VT returns malicious if key is present
	if !res.Malicious {
		t.Error("expected malicious result from mock VT")
	}
}

func TestIntelManager_MultiCheck_Timeout(t *testing.T) {
	im := NewIntelManager("fake_vt", "fake_abuse", "")
	
	// Create a context that times out immediately
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-1*time.Second))
	defer cancel()

	res := im.MultiCheck(ctx, "example.com")
	
	if res.Details["error"] != "threat intel timeout" {
		t.Errorf("expected timeout error, got %v", res.Details["error"])
	}
}
