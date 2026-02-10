package visualization

import (
	"fmt"
	"strings"
)

// GraphRenderer handles visualization of network relationships.
type GraphRenderer struct{}

func NewGraphRenderer() *GraphRenderer {
	return &GraphRenderer{}
}

// RenderASCIIChain visualizes a redirect or network chain.
func (gr *GraphRenderer) RenderASCIIChain(target string, hops []string) string {
	if len(hops) == 0 {
		return fmt.Sprintf("[%s]", target)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\n[*] NETWORK PATH VISUALIZATION:\n\n   (START) [%s]\n", target))

	for i, hop := range hops {
		sb.WriteString("             │\n")
		sb.WriteString("             ▼\n")
		if i == len(hops)-1 {
			sb.WriteString(fmt.Sprintf("   (FINAL) [%s]\n", hop))
		} else {
			sb.WriteString(fmt.Sprintf("   (HOP %d) [%s]\n", i+1, hop))
		}
	}

	return sb.String()
}
