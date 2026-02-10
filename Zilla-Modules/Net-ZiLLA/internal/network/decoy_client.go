package network

import (
	"context"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type DecoyClient struct {
	client     *http.Client
	userAgents []string
	proxyURL   string
	killSwitch bool
}

func NewDecoyClient(proxy string) *DecoyClient {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Avoid TLS fingerprinting blocking
	}

	if proxy != "" {
		p, _ := url.Parse(proxy)
		transport.Proxy = http.ProxyURL(p)
	}

	return &DecoyClient{
		killSwitch: false,
		userAgents: []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
			"Googlebot/2.1 (+http://www.google.com/bot.html)",
			"Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)",
		},
		client: &http.Client{
			Transport: transport,
			Timeout:   15 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 5 {
					return fmt.Errorf("stopped after 5 redirects for safety")
				}
				return nil
			},
		},
	}
}

func (dc *DecoyClient) SafeGet(ctx context.Context, target string) (*http.Response, error) {
	if dc.killSwitch {
		return nil, fmt.Errorf("kill switch active: requests blocked")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", target, nil)
	if err != nil {
		return nil, err
	}

	// Rotate User-Agent
	ua := dc.userAgents[rand.Intn(len(dc.userAgents))]
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	return dc.client.Do(req)
}

func (dc *DecoyClient) ActivateKillSwitch() {
	dc.killSwitch = true
}
