package network

import (
	"context"
	"testing"
	"time"

	"net-zilla/pkg/logger"
)

func TestPacketAnalyzer_InspectTraffic(t *testing.T) {
	pa := NewPacketAnalyzer(logger.NewLogger())
	res, err := pa.InspectTraffic(context.Background(), 100*time.Millisecond)
	if err != nil {
		t.Fatalf("InspectTraffic failed: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("expected 0 results, got %d", len(res))
	}
}
