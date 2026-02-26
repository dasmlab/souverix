package stir

import (
	"testing"
	"time"

	"github.com/dasmlab/ims/internal/config"
	"github.com/sirupsen/logrus"
)

func TestNewACMECertificateManager(t *testing.T) {
	cfg := &config.ACMEConfig{
		Provider: "letsencrypt",
		Email:    "test@example.com",
		Domain:   "ims.local",
		Staging:  true,
	}
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	mgr, err := NewACMECertificateManager(cfg, log)
	if err != nil {
		t.Fatalf("NewACMECertificateManager() error = %v", err)
	}

	if mgr == nil {
		t.Error("NewACMECertificateManager() returned nil")
	}
}

func TestACMECertificateManager_GetPrivateKey(t *testing.T) {
	cfg := &config.ACMEConfig{
		Provider: "letsencrypt",
		Email:    "test@example.com",
		Domain:   "ims.local",
		Staging:  true,
	}
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	mgr, _ := NewACMECertificateManager(cfg, log)

	key := mgr.GetPrivateKey()
	if key == nil {
		t.Error("GetPrivateKey() returned nil")
	}
}

func TestACMECertificateManager_GetCertificateURL(t *testing.T) {
	cfg := &config.ACMEConfig{
		Provider: "letsencrypt",
		Email:    "test@example.com",
		Domain:   "ims.local",
		Staging:  true,
	}
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	mgr, _ := NewACMECertificateManager(cfg, log)

	url := mgr.GetCertificateURL()
	if url == "" {
		t.Error("GetCertificateURL() returned empty string")
	}

	if url != "https://ims.local/.well-known/stir/cert.pem" {
		t.Errorf("GetCertificateURL() = %v, want https://ims.local/.well-known/stir/cert.pem", url)
	}
}

func TestACMECertificateManager_CertificateExpiry(t *testing.T) {
	cfg := &config.ACMEConfig{
		Provider: "letsencrypt",
		Email:    "test@example.com",
		Domain:   "ims.local",
		Staging:  true,
	}
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	mgr, _ := NewACMECertificateManager(cfg, log)

	expiry := mgr.CertificateExpiry()
	if expiry.IsZero() {
		t.Error("CertificateExpiry() returned zero time")
	}

	// Should be in the future
	if expiry.Before(time.Now()) {
		t.Error("CertificateExpiry() should be in the future")
	}
}

func TestACMECertificateManager_RenewCertificate(t *testing.T) {
	cfg := &config.ACMEConfig{
		Provider: "letsencrypt",
		Email:    "test@example.com",
		Domain:   "ims.local",
		Staging:  true,
	}
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	mgr, _ := NewACMECertificateManager(cfg, log)

	// Get original expiry
	originalExpiry := mgr.CertificateExpiry()

	// Renew
	err := mgr.RenewCertificate()
	if err != nil {
		t.Fatalf("RenewCertificate() error = %v", err)
	}

	// New expiry should be in the future
	newExpiry := mgr.CertificateExpiry()
	if newExpiry.Before(time.Now()) {
		t.Error("RenewCertificate() should set expiry in the future")
	}

	// Should be different (or at least not before original)
	if newExpiry.Before(originalExpiry) {
		t.Error("RenewCertificate() should extend expiry")
	}
}
