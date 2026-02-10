package network

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecoyClient_SafeGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	dc := NewDecoyClient("")
	resp, err := dc.SafeGet(context.Background(), ts.URL)
	if err != nil {
		t.Fatalf("SafeGet failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestDecoyClient_KillSwitch(t *testing.T) {
	dc := NewDecoyClient("")
	dc.ActivateKillSwitch()
	_, err := dc.SafeGet(context.Background(), "http://example.com")
	if err == nil {
		t.Error("expected error when kill switch is active")
	}
}

func TestDecoyClient_RedirectLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	}))
	defer ts.Close()

	dc := NewDecoyClient("")
	_, err := dc.SafeGet(context.Background(), ts.URL)
	if err == nil {
		t.Error("expected error for too many redirects")
	}
}
