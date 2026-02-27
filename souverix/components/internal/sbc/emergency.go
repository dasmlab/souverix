package sbc

import (
	"strings"

	"github.com/dasmlab/ims/internal/emergency"
	"github.com/dasmlab/ims/internal/sip"
	"github.com/sirupsen/logrus"
)

// processEmergency handles emergency calls with highest priority
// This must be called BEFORE any rate limiting or restrictions
func (s *SBC) processEmergency(msg *sip.Message) (bool, error) {
	if !s.config.IMS.Emergency.Enabled {
		return false, nil
	}

	// Extract dialed number
	dialedNumber := extractNumberFromRequestURI(msg.URI)

	// Check if emergency
	detector := emergency.NewEmergencyDetector(s.log)
	emergencyNum, isEmergency := detector.IsEmergency(dialedNumber)

	if !isEmergency {
		return false, nil
	}

	s.log.WithFields(logrus.Fields{
		"number":    dialedNumber,
		"emergency": emergencyNum.Number,
		"call_id":   msg.GetHeader("Call-ID"),
	}).Warn("EMERGENCY CALL DETECTED")

	// Apply emergency policy
	policy := emergency.NewEmergencyPolicy(s.log)

	// Bypass all restrictions
	if policy.ShouldBypassRateLimit(true) && s.rateLimiter != nil {
		// Emergency calls bypass rate limiting
		// (rate limiter should check emergency flag)
	}

	// Bypass STIR if configured
	if s.config.IMS.Emergency.BypassSTIR && s.enableSTIR {
		// Skip STIR verification for emergency
		s.log.Debug("emergency call bypassing STIR verification")
	}

	// Route to PSAP
	router := emergency.NewEmergencyRouter(detector, s.log)
	psapRoute, err := router.RouteEmergency(msg)
	if err != nil {
		s.log.WithError(err).Error("failed to route emergency call")
		return false, err
	}

	// Preserve location
	locationHandler := emergency.NewEmergencyLocationHandler(s.log)
	if location, err := locationHandler.ExtractLocation(msg); err == nil {
		locationHandler.PreserveLocation(msg, location)
	} else {
		s.log.WithError(err).Warn("no location information for emergency call")
	}

	// Set priority header
	msg.SetHeader("Priority", "emergency")
	msg.SetHeader("X-Emergency-Route", psapRoute)

	return true, nil
}

// extractNumberFromRequestURI extracts number from SIP Request-URI
func extractNumberFromRequestURI(uri string) string {
	// Remove sip: or tel: prefix
	uri = strings.TrimPrefix(uri, "sip:")
	uri = strings.TrimPrefix(uri, "tel:")

	// Extract user part (before @)
	parts := strings.Split(uri, "@")
	if len(parts) > 0 {
		return parts[0]
	}

	return uri
}
