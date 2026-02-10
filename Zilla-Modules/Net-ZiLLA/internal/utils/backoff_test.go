package utils

import (
	"testing"
	"time"
)

func TestBackoff(t *testing.T) {
	b := NewBackoff(time.Millisecond*10, time.Second)
	
	d1 := b.Next()
	if d1 != time.Millisecond*10 {
		t.Errorf("expected 10ms, got %v", d1)
	}
	
	d2 := b.Next()
	if d2 <= d1 {
		t.Errorf("expected increasing delay, got %v <= %v", d2, d1)
	}
	
	b.Reset()
	if b.Attempts != 0 {
		t.Error("expected 0 attempts after reset")
	}
}
