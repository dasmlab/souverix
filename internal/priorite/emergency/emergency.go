package priorite

import (
	"fmt"
	"strings"
	"sync"

	"github.com/dasmlab/ims/internal/sip"
	"github.com/sirupsen/logrus"
)

// EmergencyNumber represents an emergency number pattern
type EmergencyNumber struct {
	Number    string // e.g., "911", "112", "999"
	Country   string
	PSAPRoute string // Route to PSAP
	Priority  int    // Higher = more priority
}

// EmergencyDetector detects emergency calls
type EmergencyDetector struct {
	numbers map[string]*EmergencyNumber
	log     *logrus.Logger
	mu      sync.RWMutex
}

// NewEmergencyDetector creates a new emergency detector
func NewEmergencyDetector(log *logrus.Logger) *EmergencyDetector {
	ed := &EmergencyDetector{
		numbers: make(map[string]*EmergencyNumber),
		log:     log,
	}

	// Seed common emergency numbers
	ed.seedEmergencyNumbers()

	return ed
}

// seedEmergencyNumbers seeds common emergency numbers
func (ed *EmergencyDetector) seedEmergencyNumbers() {
	numbers := []*EmergencyNumber{
		{Number: "911", Country: "US", PSAPRoute: "psap-us", Priority: 100},
		{Number: "112", Country: "EU", PSAPRoute: "psap-eu", Priority: 100},
		{Number: "999", Country: "UK", PSAPRoute: "psap-uk", Priority: 100},
		{Number: "000", Country: "AU", PSAPRoute: "psap-au", Priority: 100},
	}

	for _, num := range numbers {
		ed.numbers[num.Number] = num
	}
}

// IsEmergency checks if a dialed number is an emergency number
func (ed *EmergencyDetector) IsEmergency(dialedNumber string) (*EmergencyNumber, bool) {
	ed.mu.RLock()
	defer ed.mu.RUnlock()

	// Normalize number (remove formatting)
	normalized := normalizeNumber(dialedNumber)

	// Check exact match
	if num, ok := ed.numbers[normalized]; ok {
		return num, true
	}

	// Check if starts with emergency number (for extensions like 9111)
	// But only if normalized is longer than the emergency number
	for numStr, num := range ed.numbers {
		if len(normalized) >= len(numStr) && strings.HasPrefix(normalized, numStr) {
			return num, true
		}
	}

	return nil, false
}

// normalizeNumber normalizes a phone number for comparison
func normalizeNumber(number string) string {
	// Remove common formatting
	number = strings.ReplaceAll(number, "-", "")
	number = strings.ReplaceAll(number, " ", "")
	number = strings.ReplaceAll(number, "(", "")
	number = strings.ReplaceAll(number, ")", "")
	number = strings.ReplaceAll(number, ".", "")

	// Remove + prefix if present
	if strings.HasPrefix(number, "+") {
		number = number[1:]
		// For +1911, we want to check "911" not "1911"
		// So if it's a country code + emergency, extract just the emergency part
		if strings.HasSuffix(number, "911") && len(number) > 3 {
			number = number[len(number)-3:]
		} else if strings.HasSuffix(number, "112") && len(number) > 3 {
			number = number[len(number)-3:]
		} else if strings.HasSuffix(number, "999") && len(number) > 3 {
			number = number[len(number)-3:]
		} else if strings.HasSuffix(number, "000") && len(number) > 3 {
			number = number[len(number)-3:]
		}
	}

	return number
}

// EmergencyRouter routes emergency calls to PSAP
type EmergencyRouter struct {
	detector *EmergencyDetector
	log      *logrus.Logger
}

// NewEmergencyRouter creates a new emergency router
func NewEmergencyRouter(detector *EmergencyDetector, log *logrus.Logger) *EmergencyRouter {
	return &EmergencyRouter{
		detector: detector,
		log:      log,
	}
}

// RouteEmergency routes an emergency call
func (er *EmergencyRouter) RouteEmergency(msg *sip.Message) (string, error) {
	// Extract dialed number from Request-URI
	uri := msg.URI
	dialedNumber := extractNumberFromURI(uri)

	// Detect emergency
	emergencyNum, isEmergency := er.detector.IsEmergency(dialedNumber)
	if !isEmergency {
		return "", fmt.Errorf("not an emergency number: %s", dialedNumber)
	}

	er.log.WithFields(logrus.Fields{
		"number":      dialedNumber,
		"emergency":   emergencyNum.Number,
		"country":     emergencyNum.Country,
		"psap_route":  emergencyNum.PSAPRoute,
		"call_id":     msg.GetHeader("Call-ID"),
	}).Info("emergency call detected")

	return emergencyNum.PSAPRoute, nil
}

// extractNumberFromURI extracts the number from a SIP URI
func extractNumberFromURI(uri string) string {
	// Remove sip: prefix
	uri = strings.TrimPrefix(uri, "sip:")
	uri = strings.TrimPrefix(uri, "tel:")

	// Extract user part (before @)
	parts := strings.Split(uri, "@")
	if len(parts) > 0 {
		return parts[0]
	}

	return uri
}

// EmergencyPolicy enforces emergency call policies
type EmergencyPolicy struct {
	log *logrus.Logger
}

// NewEmergencyPolicy creates a new emergency policy
func NewEmergencyPolicy(log *logrus.Logger) *EmergencyPolicy {
	return &EmergencyPolicy{
		log: log,
	}
}

// ShouldBypassRestrictions checks if emergency call should bypass restrictions
func (ep *EmergencyPolicy) ShouldBypassRestrictions(isEmergency bool) bool {
	// Emergency calls ALWAYS bypass restrictions
	return isEmergency
}

// ShouldBypassSTIR checks if emergency call should bypass STIR verification
func (ep *EmergencyPolicy) ShouldBypassSTIR(isEmergency bool) bool {
	// Emergency calls bypass STIR if it would block the call
	return isEmergency
}

// ShouldBypassFraudDetection checks if emergency call should bypass fraud detection
func (ep *EmergencyPolicy) ShouldBypassFraudDetection(isEmergency bool) bool {
	// Emergency calls ALWAYS bypass fraud detection
	return isEmergency
}

// ShouldBypassRateLimit checks if emergency call should bypass rate limiting
func (ep *EmergencyPolicy) ShouldBypassRateLimit(isEmergency bool) bool {
	// Emergency calls ALWAYS bypass rate limiting
	return isEmergency
}

// GetPriority returns the priority for an emergency call
func (ep *EmergencyPolicy) GetPriority(isEmergency bool) int {
	if isEmergency {
		return 1000 // Highest priority
	}
	return 100 // Normal priority
}

// EmergencyLocationHandler handles location information for emergency calls
type EmergencyLocationHandler struct {
	log *logrus.Logger
}

// NewEmergencyLocationHandler creates a new location handler
func NewEmergencyLocationHandler(log *logrus.Logger) *EmergencyLocationHandler {
	return &EmergencyLocationHandler{
		log: log,
	}
}

// ExtractLocation extracts location from SIP headers
func (elh *EmergencyLocationHandler) ExtractLocation(msg *sip.Message) (string, error) {
	// Check for P-Access-Network-Info header
	location := msg.GetHeader("P-Access-Network-Info")
	if location != "" {
		return location, nil
	}

	// Check for Geolocation header
	location = msg.GetHeader("Geolocation")
	if location != "" {
		return location, nil
	}

	// Check for Geolocation-Routing header
	location = msg.GetHeader("Geolocation-Routing")
	if location != "" {
		return location, nil
	}

	return "", fmt.Errorf("no location information found")
}

// ValidateLocation validates location information
func (elh *EmergencyLocationHandler) ValidateLocation(location string) bool {
	// Basic validation - in production would validate format
	return location != ""
}

// PreserveLocation preserves location in routing
func (elh *EmergencyLocationHandler) PreserveLocation(msg *sip.Message, location string) {
	// Ensure location headers are preserved
	if msg.GetHeader("P-Access-Network-Info") == "" {
		msg.SetHeader("P-Access-Network-Info", location)
	}
}
