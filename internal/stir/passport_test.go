package stir

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

func TestNewSTIRSigner(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	signer := NewSTIRSigner(privateKey, "https://example.com/cert.pem", AttestationFull)

	if signer == nil {
		t.Error("NewSTIRSigner() returned nil")
	}
}

func TestSTIRSigner_SignINVITE(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	signer := NewSTIRSigner(privateKey, "https://example.com/cert.pem", AttestationFull)

	token, err := signer.SignINVITE("+15145559876", "+15145551234", "test-call-id")
	if err != nil {
		t.Fatalf("SignINVITE() error = %v", err)
	}

	if token == "" {
		t.Error("SignINVITE() returned empty token")
	}
}

func TestSTIRSigner_AttestationLevels(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	tests := []struct {
		name       string
		attestation AttestationLevel
	}{
		{"Full", AttestationFull},
		{"Partial", AttestationPartial},
		{"Gateway", AttestationGateway},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signer := NewSTIRSigner(privateKey, "https://example.com/cert.pem", tt.attestation)
			token, err := signer.SignINVITE("+15145559876", "+15145551234", "test-call-id")
			if err != nil {
				t.Fatalf("SignINVITE() error = %v", err)
			}
			if token == "" {
				t.Error("SignINVITE() returned empty token")
			}
		})
	}
}

func TestSTIRVerifier_VerifyINVITE(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	signer := NewSTIRSigner(privateKey, "https://example.com/cert.pem", AttestationFull)

	// Sign a token
	token, err := signer.SignINVITE("+15145559876", "+15145551234", "test-call-id")
	if err != nil {
		t.Fatalf("SignINVITE() error = %v", err)
	}

	// Create a mock certificate fetcher
	fetcher := &mockCertFetcher{publicKey: &privateKey.PublicKey}
	verifier := NewSTIRVerifier(fetcher)

	// Verify the token
	passport, err := verifier.VerifyINVITE(token)
	if err != nil {
		t.Fatalf("VerifyINVITE() error = %v", err)
	}

	if passport.Orig.TN != "+15145559876" {
		t.Errorf("VerifyINVITE() orig TN = %v, want +15145559876", passport.Orig.TN)
	}

	if passport.Attest != AttestationFull {
		t.Errorf("VerifyINVITE() attestation = %v, want %v", passport.Attest, AttestationFull)
	}
}

func TestSTIRVerifier_InvalidSignature(t *testing.T) {
	privateKey1, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	privateKey2, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	signer := NewSTIRSigner(privateKey1, "https://example.com/cert.pem", AttestationFull)
	token, _ := signer.SignINVITE("+15145559876", "+15145551234", "test-call-id")

	// Use different key for verification
	fetcher := &mockCertFetcher{publicKey: &privateKey2.PublicKey}
	verifier := NewSTIRVerifier(fetcher)

	_, err := verifier.VerifyINVITE(token)
	if err == nil {
		t.Error("VerifyINVITE() should fail with invalid signature")
	}
}

func TestDetermineAttestationLevel(t *testing.T) {
	tests := []struct {
		name           string
		subscriberKnown bool
		numberControl   bool
		externalOrigin  bool
		want           AttestationLevel
	}{
		{"Full", true, true, false, AttestationFull},
		{"Partial", true, false, false, AttestationPartial},
		{"Gateway", false, false, true, AttestationGateway},
		{"Gateway External", true, true, true, AttestationGateway},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetermineAttestationLevel(tt.subscriberKnown, tt.numberControl, tt.externalOrigin)
			if got != tt.want {
				t.Errorf("DetermineAttestationLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatIdentityHeader(t *testing.T) {
	token := "test-token"
	formatted := FormatIdentityHeader(token)

	if formatted == "" {
		t.Error("FormatIdentityHeader() returned empty string")
	}
}

func TestParseIdentityHeader(t *testing.T) {
	token := "test-token"
	formatted := FormatIdentityHeader(token)

	parsed, err := ParseIdentityHeader(formatted)
	if err != nil {
		t.Fatalf("ParseIdentityHeader() error = %v", err)
	}

	if parsed != token {
		t.Errorf("ParseIdentityHeader() = %v, want %v", parsed, token)
	}
}

// mockCertFetcher is a mock implementation for testing
type mockCertFetcher struct {
	publicKey *ecdsa.PublicKey
}

func (m *mockCertFetcher) FetchCertificate(certURL string) (*ecdsa.PublicKey, error) {
	return m.publicKey, nil
}
