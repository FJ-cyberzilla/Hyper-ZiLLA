package utils

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"net-zilla/internal/services"
	"net-zilla/pkg/logger"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
	ColorOrange = "\033[38;5;208m"
	ColorBold   = "\033[1m"
)

func DisplayBanner() {
	fmt.Printf("%s\n  ░█▀█░█▀▀░▀█▀░▀▀█░▀█▀░█░░░█░░░█▀█░░\n  ░█░█░█▀▀░░█░░▄▀░░░█░░█░░░█░░░█▀█░░\n  ░▀░▀░▀▀▀░░▀░░▀▀▀░▀▀▀░▀▀▀░▀▀▀░▀░▀░░%s\n", ColorOrange, ColorReset)
}

type Menu struct {
	service *services.AnalysisService
	logger  *logger.Logger
	reader  *bufio.Reader
	running bool
}

func NewMenu(service *services.AnalysisService, l *logger.Logger) *Menu {
	return &Menu{
		service: service,
		logger:  l,
		reader:  bufio.NewReader(os.Stdin),
		running: true,
	}
}

func (m *Menu) Run() error {
	for m.running {
		DisplayBanner()
		fmt.Printf("\n1. 🔍 Secure Analysis\n2. 📜 History\n0. Exit\n\nOption: ")
		choice, _ := m.readInput()
		switch choice {
		case "1":
			m.performAnalysis()
		case "2":
			m.showHistory()
		case "0":
			m.running = false
		}
	}
	return nil
}

func (m *Menu) readInput() (string, error) {
	input, _ := m.reader.ReadString('\n')
	return strings.TrimSpace(input), nil
}

func (m *Menu) performAnalysis() {
	fmt.Print("🔗 Target URL/IP: ")
	target, _ := m.readInput()
	if target == "" {
		return
	}

	fmt.Printf("\n%s[*] Initiating Secure Pipeline...%s\n", ColorCyan, ColorReset)
	report, err := m.service.PerformAnalysis(context.Background(), target)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("\n%sVERDICT: %s (Score: %.2f)%s\n", ColorBold, report.RiskAssessment.OverallRiskLevel, report.RiskAssessment.RiskScore, ColorReset)

	if report.BasicAnalysis != nil && report.BasicAnalysis.URLEnrichment != nil {
		en := report.BasicAnalysis.URLEnrichment
		fmt.Printf("\n%s🔍 URL Enrichment:%s\n", ColorCyan, ColorReset)
		fmt.Printf(" - Entropy: %.2f\n", en.Entropy)
		fmt.Printf(" - TLD Risk: %.2f\n", en.TLDRisk)
		if en.HomographAttack {
			fmt.Printf(" - %s⚠️  Homograph Attack Detected%s\n", ColorRed, ColorReset)
		}
		if len(en.KeywordsFound) > 0 {
			fmt.Printf(" - Keywords: %s\n", strings.Join(en.KeywordsFound, ", "))
		}
	}
}

func (m *Menu) showHistory() {
	history, _ := m.service.GetAnalysisHistory(context.Background(), 10)
	for _, h := range history {
		fmt.Printf("[%s] %-30s | %s\n", h.AnalyzedAt.Format("15:04"), h.URL, h.ThreatLevel)
	}
}
