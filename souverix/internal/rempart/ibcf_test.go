package rempart

import (
	"testing"

	"github.com/dasmlab/ims/internal/config"
	"github.com/dasmlab/ims/internal/sip"
	"github.com/dasmlab/ims/internal/stir"
	"github.com/sirupsen/logrus"
)

func TestNewIBCF(t *testing.T) {
	cfg := &config.Config{
		IMS: config.IMSConfig{
			Domain: "ims.local",
			SBC: config.SBCConfig{
				TopologyHiding: true,
				RequireTLS:     false,
				EnableSTIR:     false,
			},
		},
	}
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	ibcf, err := NewIBCF(cfg, log)
	if err != nil {
		t.Fatalf("NewIBCF() error = %v", err)
	}

	if ibcf == nil {
		t.Error("NewIBCF() returned nil")
	}
}

func TestIBCF_ProcessMessage_Validation(t *testing.T) {
	cfg := &config.Config{
		IMS: config.IMSConfig{
			Domain: "ims.local",
			SBC: config.SBCConfig{
				TopologyHiding: true,
			},
		},
	}
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	ibcf, _ := NewIBCF(cfg, log)

	// Invalid message (missing headers)
	msg := &sip.Message{
		Method:  sip.MethodINVITE,
		URI:     "sip:bob@example.com",
		Version: "SIP/2.0",
		Headers: make(map[string][]string),
	}

	response, err := ibcf.ProcessMessage(msg, "192.168.1.1:5060")
	if err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}

	// Should return error response
	if response == nil || response.StatusCode != sip.StatusBadRequest {
		t.Error("ProcessMessage() should return 400 for invalid message")
	}
}

func TestIBCF_ProcessMessage_Valid(t *testing.T) {
	cfg := &config.Config{
		IMS: config.IMSConfig{
			Domain: "ims.local",
			SBC: config.SBCConfig{
				TopologyHiding: true,
			},
		},
	}
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	ibcf, _ := NewIBCF(cfg, log)

	msg := &sip.Message{
		Method:  sip.MethodINVITE,
		URI:     "sip:bob@example.com",
		Version: "SIP/2.0",
		Headers: map[string][]string{
			"Via":     {"SIP/2.0/UDP 192.168.1.1:5060"},
			"From":    {"sip:alice@example.com"},
			"To":      {"sip:bob@example.com"},
			"Call-ID": {"test-call-id"},
			"CSeq":    {"1 INVITE"},
		},
	}

	response, err := ibcf.ProcessMessage(msg, "192.168.1.1:5060")
	if err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}

	if response == nil {
		t.Error("ProcessMessage() returned nil")
	}
}

func TestIBCF_TopologyHiding(t *testing.T) {
	cfg := &config.Config{
		IMS: config.IMSConfig{
			Domain: "internal.ims.local",
			SBC: config.SBCConfig{
				TopologyHiding: true,
			},
		},
	}
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	ibcf, _ := NewIBCF(cfg, log)

	msg := &sip.Message{
		Method:  sip.MethodINVITE,
		URI:     "sip:bob@example.com",
		Version: "SIP/2.0",
		Headers: map[string][]string{
			"Via":     {"SIP/2.0/UDP internal.ims.local:5060"},
			"From":    {"sip:alice@internal.ims.local"},
			"To":      {"sip:bob@example.com"},
			"Call-ID": {"test-call-id"},
			"CSeq":    {"1 INVITE"},
			"Server":  {"Internal-Server/1.0"},
		},
	}

	_, err := ibcf.ProcessMessage(msg, "192.168.1.1:5060")
	if err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}

	// Check topology hiding
	if msg.GetHeader("Server") != "" {
		t.Error("Topology hiding failed: Server header not removed")
	}
}

func TestSimplePolicyEngine_IsPeerAllowed(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	peers := []string{"peer1.com", "peer2.com"}
	policy := NewSimplePolicyEngine(peers, false, stir.AttestationFull, log)

	// Allowed peers
	if !policy.IsPeerAllowed("peer1.com") {
		t.Error("IsPeerAllowed() should allow peer1.com")
	}

	// Disallowed peer
	if policy.IsPeerAllowed("unknown.com") {
		t.Error("IsPeerAllowed() should reject unknown.com")
	}

	// Empty list allows all (for development)
	policyEmpty := NewSimplePolicyEngine([]string{}, false, stir.AttestationFull, log)
	if !policyEmpty.IsPeerAllowed("any.com") {
		t.Error("IsPeerAllowed() with empty list should allow all")
	}
}
