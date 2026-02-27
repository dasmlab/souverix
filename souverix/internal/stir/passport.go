package stir

import (
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AttestationLevel represents the STIR/SHAKEN attestation level
type AttestationLevel string

const (
	AttestationFull    AttestationLevel = "A" // Full attestation - provider knows customer & number
	AttestationPartial AttestationLevel = "B" // Partial - knows customer but not full number control
	AttestationGateway AttestationLevel = "C" // Gateway - call originated externally
)

// PASSporT represents a PASSporT token (RFC 8225)
type PASSporT struct {
	// Standard JWT claims
	jwt.RegisteredClaims

	// PASSporT specific claims
	Orig OrigClaim `json:"orig"`
	Dest DestClaim `json:"dest"`
	Attest AttestationLevel `json:"attest"`
	OrigID string `json:"origid,omitempty"`
}

// OrigClaim represents the originating telephone number claim
type OrigClaim struct {
	TN string `json:"tn"` // Telephone number
}

// DestClaim represents the destination telephone number claim
type DestClaim struct {
	TN []string `json:"tn"` // Array of destination numbers
}

// STIRSigner signs SIP INVITE messages with STIR/SHAKEN
type STIRSigner struct {
	privateKey *ecdsa.PrivateKey
	certURL    string // URL to fetch public certificate
	attestation AttestationLevel
}

// NewSTIRSigner creates a new STIR signer
func NewSTIRSigner(privateKey *ecdsa.PrivateKey, certURL string, attestation AttestationLevel) *STIRSigner {
	return &STIRSigner{
		privateKey: privateKey,
		certURL:    certURL,
		attestation: attestation,
	}
}

// SignINVITE signs an INVITE message with a PASSporT token
func (s *STIRSigner) SignINVITE(origTN, destTN string, callID string) (string, error) {
	// Create PASSporT token
	passport := &PASSporT{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)), // Short-lived
			ID:        callID,
		},
		Orig: OrigClaim{
			TN: origTN,
		},
		Dest: DestClaim{
			TN: []string{destTN},
		},
		Attest: s.attestation,
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, passport)

	// Set header
	token.Header["x5u"] = s.certURL // Certificate URL for verification

	// Sign token
	tokenString, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign PASSporT token: %w", err)
	}

	return tokenString, nil
}

// STIRVerifier verifies STIR/SHAKEN signatures
type STIRVerifier struct {
	certFetcher CertificateFetcher
}

// CertificateFetcher fetches certificates for verification
type CertificateFetcher interface {
	FetchCertificate(certURL string) (*ecdsa.PublicKey, error)
}

// NewSTIRVerifier creates a new STIR verifier
func NewSTIRVerifier(fetcher CertificateFetcher) *STIRVerifier {
	return &STIRVerifier{
		certFetcher: fetcher,
	}
}

// VerifyINVITE verifies the Identity header in an INVITE message
func (v *STIRVerifier) VerifyINVITE(identityHeader string) (*PASSporT, error) {
	// Parse token
	token, err := jwt.Parse(identityHeader, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get certificate URL from header
		certURL, ok := token.Header["x5u"].(string)
		if !ok {
			return nil, fmt.Errorf("missing x5u header")
		}

		// Fetch public key
		publicKey, err := v.certFetcher.FetchCertificate(certURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch certificate: %w", err)
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims format")
	}

	// Build PASSporT from claims
	passport := &PASSporT{}
	
	if orig, ok := claims["orig"].(map[string]interface{}); ok {
		if tn, ok := orig["tn"].(string); ok {
			passport.Orig.TN = tn
		}
	}

	if dest, ok := claims["dest"].(map[string]interface{}); ok {
		if tnArray, ok := dest["tn"].([]interface{}); ok {
			passport.Dest.TN = make([]string, 0, len(tnArray))
			for _, tn := range tnArray {
				if tnStr, ok := tn.(string); ok {
					passport.Dest.TN = append(passport.Dest.TN, tnStr)
				}
			}
		}
	}

	if attest, ok := claims["attest"].(string); ok {
		passport.Attest = AttestationLevel(attest)
	}

	return passport, nil
}

// DetermineAttestationLevel determines the attestation level based on subscriber info
func DetermineAttestationLevel(subscriberKnown bool, numberControl bool, externalOrigin bool) AttestationLevel {
	if externalOrigin {
		return AttestationGateway
	}
	if subscriberKnown && numberControl {
		return AttestationFull
	}
	if subscriberKnown {
		return AttestationPartial
	}
	return AttestationGateway
}

// FormatIdentityHeader formats the Identity header for SIP
func FormatIdentityHeader(token string) string {
	// Base64 encode the token (some implementations require this)
	encoded := base64.URLEncoding.EncodeToString([]byte(token))
	return encoded
}

// ParseIdentityHeader parses the Identity header from SIP
func ParseIdentityHeader(header string) (string, error) {
	// Try to decode if it's base64 encoded
	decoded, err := base64.URLEncoding.DecodeString(header)
	if err == nil {
		return string(decoded), nil
	}
	// If not base64, return as-is
	return header, nil
}
