package visualization

import (
	"strings"
	"testing"
)

func TestGraphRenderer_RenderASCIIChain(t *testing.T) {
	gr := NewGraphRenderer()
	
	// Case 1: No hops
	got1 := gr.RenderASCIIChain("example.com", nil)
	if got1 != "[example.com]" {
		t.Errorf("expected [example.com], got %s", got1)
	}

	// Case 2: With hops
	hops := []string{"hop1", "final"}
	got2 := gr.RenderASCIIChain("start", hops)
	if !strings.Contains(got2, "(START) [start]") || !strings.Contains(got2, "(FINAL) [final]") {
		t.Error("rendered chain missing key elements")
	}
}
