package pcscf

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestNewPCSCF(t *testing.T) {
	cfg := &Config{
		SIPAddr:          ":5060",
		SIPTLSAddr:       ":5061",
		ICSCFAddr:        "icscf.ims.local:5060",
		SCSCFAddr:        "scscf.ims.local:5060",
		DoSProtection:    true,
		RateLimitPerIP:   100,
		RateLimitWindow:  60 * time.Second,
		EmergencyNumbers: []string{"911", "112", "999"},
	}

	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel) // Reduce noise in tests

	pcscf := New(cfg, log)
	if pcscf == nil {
		t.Fatal("New() returned nil")
	}

	if pcscf.Name() != "pcscf" {
		t.Errorf("expected name 'pcscf', got '%s'", pcscf.Name())
	}
}

func TestPCSCFStartStop(t *testing.T) {
	cfg := &Config{
		SIPAddr:          ":0", // Use port 0 for testing (OS assigns available port)
		SIPTLSAddr:       ":0",
		ICSCFAddr:        "icscf.ims.local:5060",
		DoSProtection:    false, // Disable for faster tests
		RateLimitPerIP:   100,
		RateLimitWindow:  60 * time.Second,
		EmergencyNumbers: []string{"911"},
	}

	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	pcscf := New(cfg, log)
	ctx := context.Background()

	// Start should succeed
	if err := pcscf.Start(ctx); err != nil {
		t.Fatalf("Start() failed: %v", err)
	}

	// Give it a moment to initialize
	time.Sleep(100 * time.Millisecond)

	// Check health
	health := pcscf.Health()
	if health.Status != "healthy" {
		t.Errorf("expected health status 'healthy', got '%s'", health.Status)
	}

	// Stop should succeed
	stopCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pcscf.Stop(stopCtx); err != nil {
		t.Errorf("Stop() failed: %v", err)
	}
}

func TestRateLimiter(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	rl := NewRateLimiter(5, 1*time.Second, log)

	// First 5 requests should be allowed
	for i := 0; i < 5; i++ {
		if !rl.Allow("192.168.1.1") {
			t.Errorf("request %d should be allowed", i+1)
		}
	}

	// 6th request should be rate limited
	if rl.Allow("192.168.1.1") {
		t.Error("6th request should be rate limited")
	}

	// Different IP should still be allowed
	if !rl.Allow("192.168.1.2") {
		t.Error("different IP should be allowed")
	}
}

func TestConfigFromGouverne(t *testing.T) {
	// This test requires the gouverne config package
	// For now, just verify the function exists and doesn't panic
	// Full test would require importing gouverne config
	t.Skip("Requires gouverne config package - integration test")
}
