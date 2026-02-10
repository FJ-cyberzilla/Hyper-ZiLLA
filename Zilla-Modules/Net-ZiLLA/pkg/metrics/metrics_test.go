package metrics

import (
	"testing"
	"time"
)

func TestTracker(t *testing.T) {
	tr := NewTracker()
	tr.IncrementCounter("test_count")
	tr.ObserveDuration("test_dur", time.Second)
	
	if tr.counters["test_count"] != 1 {
		t.Errorf("expected 1, got %d", tr.counters["test_count"])
	}
}
