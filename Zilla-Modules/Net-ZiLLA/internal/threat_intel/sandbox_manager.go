package threat_intel

import (
	"context"
	"fmt"
	"time"
)

type SandboxManager struct {
	dockerImage string
	timeout     time.Duration
	mode        string // docker, anyrun, hybrid
	apiKey      string
}

func NewSandboxManager() *SandboxManager {
	return &SandboxManager{
		dockerImage: "netzilla/isolated-browser:latest",
		timeout:     5 * time.Minute,
		mode:        "docker",
	}
}

func (sm *SandboxManager) SetMode(mode string) {
	sm.mode = mode
}

// SpinUpIsolatedBrowser prepares an ephemeral container or remote sandbox for analysis
func (sm *SandboxManager) SpinUpIsolatedBrowser(ctx context.Context, url string) (string, error) {
	switch sm.mode {
	case "hybrid":
		fmt.Printf("[*] Hybrid Analysis: Starting local container and remote task for: %s\n", url)
		// Start local and remote in parallel (simplified here)
		id, _ := sm.startDockerSandbox(ctx, url)
		sm.startRemoteSandbox(ctx, url)
		return id, nil
	case "anyrun":
		return sm.startRemoteSandbox(ctx, url)
	default:
		return sm.startDockerSandbox(ctx, url)
	}
}

func (sm *SandboxManager) startDockerSandbox(ctx context.Context, url string) (string, error) {
	fmt.Printf("[*] Escalating to Docker sandbox: Spinning up %s for URL: %s\n", sm.dockerImage, url)
	return "DOCKER-ID-" + fmt.Sprintf("%d", time.Now().Unix()), nil
}

func (sm *SandboxManager) startRemoteSandbox(ctx context.Context, url string) (string, error) {
	fmt.Printf("[*] Escalating to Remote sandbox (ANY.RUN) for URL: %s\n", url)
	return "REMOTE-ID-" + fmt.Sprintf("%d", time.Now().Unix()), nil
}

func (sm *SandboxManager) DestroySandbox(id string) error {
	fmt.Printf("[*] Sanitizing environment: Destroying sandbox %s\n", id)
	return nil
}
