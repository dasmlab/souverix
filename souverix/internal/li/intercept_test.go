package li

import (
	"testing"
	"time"

	"github.com/dasmlab/ims/internal/sip"
	"github.com/sirupsen/logrus"
)

// mockMediationDevice is a mock MD for testing
type mockMediationDevice struct {
	available bool
	signaling []*sip.Message
	media     [][]byte
}

func (m *mockMediationDevice) SendSignaling(msg *sip.Message, targetID string) error {
	m.signaling = append(m.signaling, msg)
	return nil
}

func (m *mockMediationDevice) SendMedia(rtpData []byte, targetID string) error {
	m.media = append(m.media, rtpData)
	return nil
}

func (m *mockMediationDevice) IsAvailable() bool {
	return m.available
}

func TestNewInterceptController(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	md := &mockMediationDevice{available: true}
	controller := NewInterceptController(md, log)

	if controller == nil {
		t.Error("NewInterceptController() returned nil")
	}
}

func TestInterceptController_ActivateWarrant(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	md := &mockMediationDevice{available: true}
	controller := NewInterceptController(md, log)

	target := &InterceptTarget{
		IMPI:        "alice@ims.local",
		TN:          "+15145559876",
		WarrantID:   "WARRANT-001",
		WarrantType: "signaling+media",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}

	err := controller.ActivateWarrant(target)
	if err != nil {
		t.Fatalf("ActivateWarrant() error = %v", err)
	}

	// Check if target is intercepted
	intercepted, isTarget := controller.IsTarget("alice@ims.local")
	if !isTarget {
		t.Error("IsTarget() should return true for activated warrant")
	}

	if intercepted.WarrantID != target.WarrantID {
		t.Errorf("IsTarget() WarrantID = %v, want %v", intercepted.WarrantID, target.WarrantID)
	}
}

func TestInterceptController_InterceptMessage(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	md := &mockMediationDevice{available: true}
	controller := NewInterceptController(md, log)

	// Activate warrant
	target := &InterceptTarget{
		IMPI:        "alice@ims.local",
		WarrantID:   "WARRANT-001",
		WarrantType: "signaling",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}
	controller.ActivateWarrant(target)

	// Create SIP message
	msg := &sip.Message{
		Method:  sip.MethodINVITE,
		URI:     "sip:bob@example.com",
		Version: "SIP/2.0",
		Headers: map[string][]string{
			"From":    {"sip:alice@ims.local"},
			"To":      {"sip:bob@example.com"},
			"Call-ID": {"test-call-id"},
			"CSeq":    {"1 INVITE"},
		},
	}

	// Intercept
	err := controller.InterceptMessage(msg, "alice@ims.local", "bob@example.com")
	if err != nil {
		t.Fatalf("InterceptMessage() error = %v", err)
	}

	// Check if message was sent to MD
	if len(md.signaling) == 0 {
		t.Error("InterceptMessage() should send message to MD")
	}
}

func TestInterceptController_NonTarget(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	md := &mockMediationDevice{available: true}
	controller := NewInterceptController(md, log)

	msg := &sip.Message{
		Method:  sip.MethodINVITE,
		URI:     "sip:bob@example.com",
		Version: "SIP/2.0",
		Headers: map[string][]string{
			"From": {"sip:alice@ims.local"},
			"To":   {"sip:bob@example.com"},
		},
	}

	// Should not intercept (no warrant)
	err := controller.InterceptMessage(msg, "alice@ims.local", "bob@example.com")
	if err != nil {
		t.Errorf("InterceptMessage() should not error for non-target: %v", err)
	}

	// Should not send to MD
	if len(md.signaling) > 0 {
		t.Error("InterceptMessage() should not send non-target messages to MD")
	}
}

func TestInterceptController_DeactivateWarrant(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	md := &mockMediationDevice{available: true}
	controller := NewInterceptController(md, log)

	// Activate
	target := &InterceptTarget{
		IMPI:      "alice@ims.local",
		WarrantID: "WARRANT-001",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	controller.ActivateWarrant(target)

	// Deactivate
	err := controller.DeactivateWarrant("WARRANT-001")
	if err != nil {
		t.Fatalf("DeactivateWarrant() error = %v", err)
	}

	// Should not be intercepted
	_, isTarget := controller.IsTarget("alice@ims.local")
	if isTarget {
		t.Error("IsTarget() should return false after deactivation")
	}
}

func TestAuditLogger_Log(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	auditLogger := NewAuditLogger(log)

	entry := AuditLog{
		WarrantID: "WARRANT-001",
		Action:    "activate",
		Target:    "alice@ims.local",
		Operator:  "operator1",
	}

	auditLogger.Log(entry)

	logs := auditLogger.GetLogs("WARRANT-001", time.Now().Add(-1*time.Hour), time.Now().Add(1*time.Hour))
	if len(logs) == 0 {
		t.Error("GetLogs() should return logged entries")
	}
}
