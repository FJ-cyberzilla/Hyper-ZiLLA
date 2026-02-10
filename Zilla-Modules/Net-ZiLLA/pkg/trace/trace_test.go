package trace

import (
	"testing"
)

func TestSpan(t *testing.T) {
	span := StartSpan("test")
	span.SetTag("key", "value")
	span.End()
	
	if span.name != "test" {
		t.Errorf("expected test, got %s", span.name)
	}
}
