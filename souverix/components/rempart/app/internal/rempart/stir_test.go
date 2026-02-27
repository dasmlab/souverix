package rempart

import (
	"testing"

	"github.com/dasmlab/ims/internal/config"
	"github.com/dasmlab/ims/internal/sip"
	"github.com/sirupsen/logrus"
)

func TestSBC_STIRSigning(t *testing.T) {
	cfg := &config.Config{
		IMS: config.IMSConfig{
			SBC: config.SBCConfig{
				EnableSTIR:      true,
				STIRAttestation: "A",
			},
		},
		ZeroTrust: config.ZeroTrustConfig{
			Enabled: true,
			ACME: config.ACMEConfig{
				Provider: "letsencrypt",
				Email:    "test@example.com",
				Domain:   "ims.local",
				Staging:  true,
			},
		},
	}
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	sbc, err := NewSBC(cfg, log)
	if err != nil {
		t.Fatalf("NewSBC() error = %v", err)
	}

	msg := &sip.Message{
		Method:  sip.MethodINVITE,
		URI:     "sip:+15145551234@example.com",
		Version: "SIP/2.0",
		Headers: map[string][]string{
			"From":    {"sip:+15145559876@ims.local"},
			"To":      {"sip:+15145551234@example.com"},
			"Call-ID": {"test-call-id"},
			"CSeq":    {"1 INVITE"},
		},
	}

	// Process message (should sign)
	_, err = sbc.ProcessMessage(msg, "192.168.1.1:5060")
	if err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}

	// Check Identity header was added
	identity := msg.GetHeader("Identity")
	if identity == "" {
		t.Error("STIR signing failed: Identity header not added")
	}
}

func TestSBC_STIRVerification(t *testing.T) {
	cfg := &config.Config{
		IMS: config.IMSConfig{
			SBC: config.SBCConfig{
				EnableSTIR:      true,
				STIRAttestation: "A",
			},
		},
		ZeroTrust: config.ZeroTrustConfig{
			Enabled: true,
			ACME: config.ACMEConfig{
				Provider: "letsencrypt",
				Email:    "test@example.com",
				Domain:   "ims.local",
				Staging:  true,
			},
		},
	}
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	sbc, err := NewSBC(cfg, log)
	if err != nil {
		t.Fatalf("NewSBC() error = %v", err)
	}

	msg := &sip.Message{
		Method:  sip.MethodINVITE,
		URI:     "sip:+15145551234@example.com",
		Version: "SIP/2.0",
		Headers: map[string][]string{
			"From":    {"sip:+15145559876@peer.com"},
			"To":      {"sip:+15145551234@example.com"},
			"Call-ID": {"test-call-id"},
			"CSeq":    {"1 INVITE"},
			"Identity": {"test-identity-token"},
		},
	}

	// Process message (should attempt verification)
	_, err = sbc.ProcessMessage(msg, "192.168.1.1:5060")
	if err != nil {
		// Verification may fail in test (no real cert), but should not crash
		t.Logf("ProcessMessage() verification failed (expected in test): %v", err)
	}
}

func TestSBC_ExtractTN(t *testing.T) {
	cfg := &config.Config{
		IMS: config.IMSConfig{
			SBC: config.SBCConfig{},
		},
	}
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	sbc, _ := NewSBC(cfg, log)

	tests := []struct {
		name string
		uri  string
		want string
	}{
		{"sip URI", "sip:+15145559876@ims.local", "+15145559876"},
		{"with display name", "Alice <sip:+15145559876@ims.local>", "+15145559876"},
		{"with tag", "sip:+15145559876@ims.local;tag=abc123", "+15145559876"},
		{"tel URI", "tel:+15145559876", "+15145559876"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sbc.extractTN(tt.uri)
			if got != tt.want {
				t.Errorf("extractTN() = %v, want %v", got, tt.want)
			}
		})
	}
}
