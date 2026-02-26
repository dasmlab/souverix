package sbc

import (
	"fmt"
	"strings"

	"github.com/dasmlab/ims/internal/sip"
	"github.com/dasmlab/ims/internal/stir"
	"github.com/sirupsen/logrus"
)

// signSTIR signs an INVITE message with STIR/SHAKEN
func (s *SBC) signSTIR(msg *sip.Message) error {
	if s.stirSigner == nil {
		return fmt.Errorf("STIR signer not initialized")
	}

	// Extract calling and called numbers from SIP headers
	from := msg.GetHeader("From")
	to := msg.GetHeader("To")
	callID := msg.GetHeader("Call-ID")

	// Parse telephone numbers from SIP URIs
	origTN := s.extractTN(from)
	destTN := s.extractTN(to)

	if origTN == "" || destTN == "" {
		s.log.Debug("skipping STIR signing - missing telephone numbers")
		return nil
	}

	// Sign the INVITE
	token, err := s.stirSigner.SignINVITE(origTN, destTN, callID)
	if err != nil {
		return fmt.Errorf("failed to sign INVITE: %w", err)
	}

	// Add Identity header
	msg.SetHeader("Identity", token)

	s.log.WithFields(logrus.Fields{
		"orig_tn": origTN,
		"dest_tn": destTN,
		"call_id": callID,
	}).Debug("STIR/SHAKEN signature added")

	return nil
}

// verifySTIR verifies the STIR/SHAKEN signature in an INVITE message
func (s *SBC) verifySTIR(msg *sip.Message) error {
	if s.stirVerifier == nil {
		return fmt.Errorf("STIR verifier not initialized")
	}

	// Get Identity header
	identityHeader := msg.GetHeader("Identity")
	if identityHeader == "" {
		s.log.Debug("no Identity header found, skipping STIR verification")
		return nil
	}

	// Parse Identity header (may be base64 encoded)
	identityToken, err := stir.ParseIdentityHeader(identityHeader)
	if err != nil {
		return fmt.Errorf("failed to parse Identity header: %w", err)
	}

	// Verify the token
	passport, err := s.stirVerifier.VerifyINVITE(identityToken)
	if err != nil {
		return fmt.Errorf("STIR verification failed: %w", err)
	}

	// Log verification result
	s.log.WithFields(logrus.Fields{
		"orig_tn":  passport.Orig.TN,
		"dest_tn":  passport.Dest.TN,
		"attest":   passport.Attest,
		"verified": true,
	}).Info("STIR/SHAKEN verification successful")

	// Add verification result to message headers for downstream processing
	msg.SetHeader("X-STIR-Attestation", string(passport.Attest))
	msg.SetHeader("X-STIR-Verified", "true")

	return nil
}

// extractTN extracts a telephone number from a SIP URI
func (s *SBC) extractTN(sipURI string) string {
	// Simple extraction - in production, use proper SIP URI parsing
	// Format: "sip:+1234567890@domain.com" or "<sip:+1234567890@domain.com>"
	
	// Remove angle brackets if present
	sipURI = strings.Trim(sipURI, "<>")
	
	// Extract the user part (before @)
	parts := strings.Split(sipURI, "@")
	if len(parts) == 0 {
		return ""
	}
	
	userPart := parts[0]
	
	// Remove "sip:" prefix if present
	userPart = strings.TrimPrefix(userPart, "sip:")
	
	// Extract display name if present (format: "Display Name" <sip:number@domain>)
	if strings.Contains(userPart, "<") {
		parts = strings.Split(userPart, "<")
		if len(parts) > 1 {
			userPart = strings.Trim(parts[1], ">")
			userPart = strings.TrimPrefix(userPart, "sip:")
		}
	}
	
	// Remove any parameters (e.g., ;tag=...)
	userPart = strings.Split(userPart, ";")[0]
	
	// Clean up the number (remove + if needed, or keep it)
	return strings.TrimSpace(userPart)
}
