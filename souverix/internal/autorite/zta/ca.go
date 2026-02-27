package autorite

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/dasmlab/ims/internal/config"
	"github.com/sirupsen/logrus"
)

// CA manages certificate authority operations for Zero Trust Mode
type CA struct {
	config *config.ZeroTrustConfig
	log    *logrus.Logger
}

// NewCA creates a new CA instance
func NewCA(cfg *config.ZeroTrustConfig, log *logrus.Logger) (*CA, error) {
	ca := &CA{
		config: cfg,
		log:    log,
	}

	if !cfg.Enabled {
		return ca, nil
	}

	// Initialize CA based on provider
	switch cfg.CAProvider {
	case "internal":
		if err := ca.initInternalCA(); err != nil {
			return nil, fmt.Errorf("failed to initialize internal CA: %w", err)
		}
	case "vault":
		if err := ca.initVaultCA(); err != nil {
			return nil, fmt.Errorf("failed to initialize Vault CA: %w", err)
		}
	case "acme":
		if err := ca.initACMECA(); err != nil {
			return nil, fmt.Errorf("failed to initialize ACME CA: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown CA provider: %s", cfg.CAProvider)
	}

	return ca, nil
}

// initInternalCA initializes an internal CA
func (ca *CA) initInternalCA() error {
	certPath := ca.config.InternalCA.CertPath
	keyPath := ca.config.InternalCA.KeyPath

	// Check if CA certificate already exists
	if _, err := os.Stat(certPath); err == nil {
		ca.log.WithField("path", certPath).Info("using existing internal CA certificate")
		return nil
	}

	ca.log.Info("generating new internal CA certificate")

	// Generate CA private key
	caKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("failed to generate CA key: %w", err)
	}

	// Create CA certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"IMS Core"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{""},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
			CommonName:    "IMS Core CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // 10 years
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            2,
	}

	// Create CA certificate
	caCertDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &caKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("failed to create CA certificate: %w", err)
	}

	// Save CA certificate
	if err := ca.saveCertificate(certPath, caCertDER); err != nil {
		return err
	}

	// Save CA private key
	if err := ca.savePrivateKey(keyPath, caKey); err != nil {
		return err
	}

	ca.log.WithFields(logrus.Fields{
		"cert": certPath,
		"key":  keyPath,
	}).Info("internal CA certificate generated")

	return nil
}

// initVaultCA initializes Vault CA integration
func (ca *CA) initVaultCA() error {
	// TODO: Implement Vault CA integration
	ca.log.Warn("Vault CA integration not yet implemented")
	return nil
}

// initACMECA initializes ACME CA integration
func (ca *CA) initACMECA() error {
	// TODO: Implement ACME CA integration
	ca.log.Warn("ACME CA integration not yet implemented")
	return nil
}

// IssueCertificate issues a certificate for a domain
func (ca *CA) IssueCertificate(domain string) (*x509.Certificate, *rsa.PrivateKey, error) {
	if !ca.config.Enabled {
		return nil, nil, fmt.Errorf("Zero Trust Mode is not enabled")
	}

	switch ca.config.CAProvider {
	case "internal":
		return ca.issueInternalCertificate(domain)
	case "vault":
		return ca.issueVaultCertificate(domain)
	case "acme":
		return ca.issueACMECertificate(domain)
	default:
		return nil, nil, fmt.Errorf("unknown CA provider: %s", ca.config.CAProvider)
	}
}

// issueInternalCertificate issues a certificate using the internal CA
func (ca *CA) issueInternalCertificate(domain string) (*x509.Certificate, *rsa.PrivateKey, error) {
	// Load CA certificate and key
	caCert, caKey, err := ca.loadCA()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load CA: %w", err)
	}

	// Generate server key
	serverKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate server key: %w", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			CommonName: domain,
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0), // 1 year
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		DNSNames:     []string{domain},
		IPAddresses:  []net.IP{}, // Can be extended
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, caCert, &serverKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	// Parse certificate
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return cert, serverKey, nil
}

// Helper functions
func (ca *CA) saveCertificate(path string, certDER []byte) error {
	certFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer certFile.Close()

	return pem.Encode(certFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})
}

func (ca *CA) savePrivateKey(path string, key *rsa.PrivateKey) error {
	keyFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer keyFile.Close()

	keyDER := x509.MarshalPKCS1PrivateKey(key)
	return pem.Encode(keyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyDER,
	})
}

func (ca *CA) loadCA() (*x509.Certificate, *rsa.PrivateKey, error) {
	// Load certificate
	certPEM, err := os.ReadFile(ca.config.InternalCA.CertPath)
	if err != nil {
		return nil, nil, err
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, nil, fmt.Errorf("failed to decode CA certificate PEM")
	}

	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, nil, err
	}

	// Load private key
	keyPEM, err := os.ReadFile(ca.config.InternalCA.KeyPath)
	if err != nil {
		return nil, nil, err
	}

	block, _ = pem.Decode(keyPEM)
	if block == nil {
		return nil, nil, fmt.Errorf("failed to decode CA key PEM")
	}

	caKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, err
	}

	return caCert, caKey, nil
}

func (ca *CA) issueVaultCertificate(domain string) (*x509.Certificate, *rsa.PrivateKey, error) {
	// TODO: Implement
	return nil, nil, fmt.Errorf("Vault CA not yet implemented")
}

func (ca *CA) issueACMECertificate(domain string) (*x509.Certificate, *rsa.PrivateKey, error) {
	// TODO: Implement
	return nil, nil, fmt.Errorf("ACME CA not yet implemented")
}
