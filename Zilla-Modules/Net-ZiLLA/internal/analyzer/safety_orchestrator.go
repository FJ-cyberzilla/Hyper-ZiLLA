package analyzer

import (
	"context"
	"net-zilla/internal/models"
	"net-zilla/internal/network"
	"net-zilla/internal/threat_intel"
	"net-zilla/pkg/logger"
)

type SafetyOrchestrator struct {
	screener *network.SafetyScreener
	intel    *threat_intel.IntelManager
	decoy    *network.DecoyClient
	sandbox  *threat_intel.SandboxManager
	logger   *logger.Logger
}

func NewSafetyOrchestrator(l *logger.Logger) *SafetyOrchestrator {
	return &SafetyOrchestrator{
		screener: network.NewSafetyScreener(),
		intel:    threat_intel.NewIntelManager("", "", ""), // Keys would be injected from config
		decoy:    network.NewDecoyClient(""),
		sandbox:  threat_intel.NewSandboxManager(),
		logger:   l,
	}
}

func (so *SafetyOrchestrator) SecureAnalyze(ctx context.Context, target string) (*models.AdvancedReport, error) {
	so.logger.Info("Starting Secure Analysis Pipeline for: %s", target)

	// 1. Safety Screening (Zero-Touch)
	screening := so.screener.Screen(target)
	if screening.IsSuspicious {
		so.logger.Warn("Initial screening flagged URL: %v", screening.Reasons)
	}

	// 2. Passive Intel Gathering
	reputation := so.intel.MultiCheck(ctx, target)
	so.logger.Info("Intel gathering complete. Positives: %d", reputation.Positives)

	// 3. Escalation Decision
	if screening.RiskScore > 50 || reputation.Malicious {
		so.logger.Warn("High risk detected. Escalating to isolated sandbox...")
		containerID, err := so.sandbox.SpinUpIsolatedBrowser(ctx, target)
		if err == nil {
			defer so.sandbox.DestroySandbox(containerID)
		}
	} else {
		so.logger.Info("Low risk confirmed. Proceeding with standard passive analysis.")
	}

	// 4. Compile Report (Placeholder for data compilation)
	return &models.AdvancedReport{
		Target:   target,
		ReportID: "SO-XYZ",
	}, nil
}
