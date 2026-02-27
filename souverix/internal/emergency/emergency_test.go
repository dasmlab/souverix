package emergency

import (
	"testing"

	"github.com/dasmlab/ims/internal/sip"
	"github.com/sirupsen/logrus"
)

func TestNewEmergencyDetector(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	detector := NewEmergencyDetector(log)
	if detector == nil {
		t.Error("NewEmergencyDetector() returned nil")
	}
}

func TestEmergencyDetector_IsEmergency(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	detector := NewEmergencyDetector(log)

	tests := []struct {
		name   string
		number string
		want   bool
	}{
		{"US 911", "911", true},
		{"EU 112", "112", true},
		{"UK 999", "999", true},
		{"AU 000", "000", true},
		{"Normal number", "5551234", false},
		{"Formatted 911", "9-1-1", true},
		{"Formatted 911", "(911)", true},
		{"With country code", "+1911", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := detector.IsEmergency(tt.number)
			if got != tt.want {
				t.Errorf("IsEmergency(%v) = %v, want %v", tt.number, got, tt.want)
			}
		})
	}
}

func TestEmergencyRouter_RouteEmergency(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	detector := NewEmergencyDetector(log)
	router := NewEmergencyRouter(detector, log)

	msg := &sip.Message{
		Method:  sip.MethodINVITE,
		URI:     "sip:911@ims.local",
		Version: "SIP/2.0",
		Headers: map[string][]string{
			"Call-ID": {"test-call-id"},
		},
	}

	route, err := router.RouteEmergency(msg)
	if err != nil {
		t.Fatalf("RouteEmergency() error = %v", err)
	}

	if route == "" {
		t.Error("RouteEmergency() should return PSAP route")
	}
}

func TestEmergencyPolicy_ShouldBypassRestrictions(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	policy := NewEmergencyPolicy(log)

	// Emergency should always bypass
	if !policy.ShouldBypassRestrictions(true) {
		t.Error("ShouldBypassRestrictions(true) should return true")
	}

	// Non-emergency should not bypass
	if policy.ShouldBypassRestrictions(false) {
		t.Error("ShouldBypassRestrictions(false) should return false")
	}
}

func TestEmergencyPolicy_ShouldBypassSTIR(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	policy := NewEmergencyPolicy(log)

	// Emergency should bypass STIR
	if !policy.ShouldBypassSTIR(true) {
		t.Error("ShouldBypassSTIR(true) should return true")
	}
}

func TestEmergencyPolicy_ShouldBypassFraudDetection(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	policy := NewEmergencyPolicy(log)

	// Emergency should bypass fraud detection
	if !policy.ShouldBypassFraudDetection(true) {
		t.Error("ShouldBypassFraudDetection(true) should return true")
	}
}

func TestEmergencyPolicy_ShouldBypassRateLimit(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	policy := NewEmergencyPolicy(log)

	// Emergency should bypass rate limiting
	if !policy.ShouldBypassRateLimit(true) {
		t.Error("ShouldBypassRateLimit(true) should return true")
	}
}

func TestEmergencyPolicy_GetPriority(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	policy := NewEmergencyPolicy(log)

	// Emergency should have highest priority
	priority := policy.GetPriority(true)
	if priority != 1000 {
		t.Errorf("GetPriority(true) = %v, want 1000", priority)
	}

	// Normal should have lower priority
	priority = policy.GetPriority(false)
	if priority != 100 {
		t.Errorf("GetPriority(false) = %v, want 100", priority)
	}
}

func TestEmergencyLocationHandler_ExtractLocation(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	handler := NewEmergencyLocationHandler(log)

	msg := &sip.Message{
		Headers: map[string][]string{
			"P-Access-Network-Info": {"3GPP-UTRAN-FDD; utran-cell-id-3gpp=234151234567890"},
		},
	}

	location, err := handler.ExtractLocation(msg)
	if err != nil {
		t.Fatalf("ExtractLocation() error = %v", err)
	}

	if location == "" {
		t.Error("ExtractLocation() should return location")
	}
}

func TestEmergencyLocationHandler_PreserveLocation(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	handler := NewEmergencyLocationHandler(log)

	msg := &sip.Message{
		Headers: make(map[string][]string),
	}

	location := "3GPP-UTRAN-FDD; utran-cell-id-3gpp=234151234567890"
	handler.PreserveLocation(msg, location)

	if msg.GetHeader("P-Access-Network-Info") != location {
		t.Error("PreserveLocation() should preserve location header")
	}
}
