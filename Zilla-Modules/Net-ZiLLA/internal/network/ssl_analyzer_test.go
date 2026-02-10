package network

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"testing"
	"time"

	"net-zilla/pkg/logger"
)

func TestNewSSLAnalyzer(t *testing.T) {
	l := logger.NewLogger()
	sa := NewSSLAnalyzer(l)
	if sa == nil {
		t.Fatal("NewSSLAnalyzer returned nil")
	}
	if sa.logger != l {
		t.Error("SSLAnalyzer logger mismatch")
	}
}

func TestSSLAnalyzer_GradeCertificate(t *testing.T) {
	sa := &SSLAnalyzer{}

	// Helper to create dummy certificates
	createCert := func(key interface{}, notAfter time.Time) *x509.Certificate {
		return &x509.Certificate{
			PublicKey: key,
			NotAfter:  notAfter,
		}
	}

	rsa2048, _ := rsa.GenerateKey(rand.Reader, 2048)
	rsa4096, _ := rsa.GenerateKey(rand.Reader, 4096)
	rsa1024, _ := rsa.GenerateKey(rand.Reader, 1024)
	ecdsaP256, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	future := time.Now().Add(365 * 24 * time.Hour)
	soon := time.Now().Add(10 * 24 * time.Hour)

	tests := []struct {
		name string
		cert *x509.Certificate
		want string
	}{
		{"Expiring Soon", createCert(&rsa2048.PublicKey, soon), "F"},
		{"Strong RSA", createCert(&rsa4096.PublicKey, future), "A+"},
		{"Standard RSA", createCert(&rsa2048.PublicKey, future), "A"},
		{"Weak RSA", createCert(&rsa1024.PublicKey, future), "C"},
		{"Standard ECDSA", createCert(&ecdsaP256.PublicKey, future), "C"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sa.gradeCertificate(tt.cert); got != tt.want {
				t.Errorf("gradeCertificate() = %v, want %v", got, tt.want)
			}
		})
	}
}
