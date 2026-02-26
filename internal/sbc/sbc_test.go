package sbc

import (
	"testing"

	"github.com/dasmlab/ims/internal/config"
	"github.com/dasmlab/ims/internal/sip"
	"github.com/sirupsen/logrus"
)

func TestNewSBC(t *testing.T) {
	cfg := &config.Config{
		IMS: config.IMSConfig{
			SBC: config.SBCConfig{
				TopologyHiding:  true,
				NormalizeHeaders: true,
				DoSProtection:    true,
				RateLimitPerIP:   100,
			},
		},
	}
	log := logrus.New()

	sbc, err := NewSBC(cfg, log)
	if err != nil {
		t.Fatalf("NewSBC() error = %v", err)
	}

	if sbc == nil {
		t.Error("NewSBC() returned nil")
	}
}

func TestSBC_ProcessMessage(t *testing.T) {
	cfg := &config.Config{
		IMS: config.IMSConfig{
			SBC: config.SBCConfig{
				TopologyHiding:  true,
				NormalizeHeaders: true,
				DoSProtection:    false, // Disable for unit test
			},
		},
	}
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel) // Reduce noise

	sbc, _ := NewSBC(cfg, log)

	msg := &sip.Message{
		Method:  sip.MethodINVITE,
		URI:     "sip:bob@example.com",
		Version: "SIP/2.0",
		Headers: map[string][]string{
			"From": {"sip:alice@example.com"},
			"To":   {"sip:bob@example.com"},
			"Call-ID": {"test-call-id"},
			"CSeq": {"1 INVITE"},
		},
	}

	response, err := sbc.ProcessMessage(msg, "192.168.1.1:5060")
	if err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}

	if response == nil {
		t.Error("ProcessMessage() returned nil response")
	}
}

func TestSBC_TopologyHiding(t *testing.T) {
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

	sbc, _ := NewSBC(cfg, log)

	msg := &sip.Message{
		Method:  sip.MethodINVITE,
		URI:     "sip:bob@example.com",
		Version: "SIP/2.0",
		Headers: map[string][]string{
			"Via":     {"SIP/2.0/UDP internal.ims.local:5060"},
			"Contact": {"sip:alice@internal.ims.local"},
			"Server":  {"Internal-Server/1.0"},
		},
	}

	_, err := sbc.ProcessMessage(msg, "192.168.1.1:5060")
	if err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}

	// Check that Server header was removed
	if msg.GetHeader("Server") != "" {
		t.Error("Topology hiding failed: Server header not removed")
	}
}
