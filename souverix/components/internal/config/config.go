package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds the application configuration
type Config struct {
	// Server configuration
	Server ServerConfig

	// IMS component configuration
	IMS IMSConfig

	// Zero Trust configuration
	ZeroTrust ZeroTrustConfig

	// Logging
	LogLevel string

	// Version info
	Version string
}

// ServerConfig holds HTTP/SIP server configuration
type ServerConfig struct {
	APIAddr     string
	MetricsAddr string
	SIPAddr     string
	SIPTLSAddr  string
}

// IMSConfig holds IMS-specific configuration
type IMSConfig struct {
	// Domain
	Domain string

	// Component roles
	EnablePCSCF bool
	EnableICSCF bool
	EnableSCSCF bool
	EnableHSS   bool
	EnableSBC   bool
	EnableIBCF  bool
	EnableMGCF  bool
	EnableBGCF  bool

	// HSS configuration
	HSS HSSConfig

	// SBC/IBCF configuration
	SBC SBCConfig

	// Lawful Intercept configuration
	LI LIConfig

	// Emergency Services configuration
	Emergency EmergencyConfig
}

// HSSConfig holds HSS configuration
type HSSConfig struct {
	Backend string // "memory", "postgres", "redis"
	DSN     string // Connection string for persistent backends
}

// SBCConfig holds Session Border Controller configuration
type SBCConfig struct {
	// Topology hiding
	TopologyHiding bool

	// SIP normalization
	NormalizeHeaders bool

	// Security
	RequireTLS      bool
	RequireSRTP     bool
	DoSProtection   bool
	RateLimitPerIP  int
	RateLimitWindow time.Duration

	// STIR/SHAKEN
	EnableSTIR      bool
	STIRAttestation string // "A", "B", "C" or "auto"
}

// ZeroTrustConfig holds Zero Trust Architecture configuration
type ZeroTrustConfig struct {
	Enabled bool

	// CA configuration
	CAProvider string // "internal", "vault", "acme"

	// Internal CA
	InternalCA InternalCAConfig

	// Vault CA (if using HashiCorp Vault)
	VaultCA VaultCAConfig

	// ACME (Let's Encrypt, etc.)
	ACME ACMEConfig
}

// InternalCAConfig holds internal CA configuration
type InternalCAConfig struct {
	CertPath string
	KeyPath  string
}

// VaultCAConfig holds Vault CA configuration
type VaultCAConfig struct {
	Address   string
	RoleID    string
	SecretID  string
	PKIPath   string
	CertPath  string
}

// ACMEConfig holds ACME configuration
type ACMEConfig struct {
	Provider   string // "letsencrypt", "zerossl", "custom"
	Email      string
	Staging    bool
	Domain     string
	CertPath   string
	KeyPath    string
}

// LIConfig holds Lawful Intercept configuration
type LIConfig struct {
	Enabled          bool
	MediationDevice  string // MD endpoint
	HandoverInterface string // HI endpoint
	AuditLogging     bool
}

// EmergencyConfig holds Emergency Services configuration
type EmergencyConfig struct {
	Enabled        bool
	EmergencyNumbers []string // e.g., ["911", "112", "999"]
	PSAPRoutes     map[string]string // number -> route
	BypassSTIR     bool
	BypassFraud    bool
	BypassRateLimit bool
}

// Load loads configuration from environment variables
func Load() *Config {
	cfg := &Config{
		Server: ServerConfig{
			APIAddr:     getEnv("API_ADDR", ":8080"),
			MetricsAddr: getEnv("METRICS_ADDR", ":9443"),
			SIPAddr:     getEnv("SIP_ADDR", ":5060"),
			SIPTLSAddr:  getEnv("SIP_TLS_ADDR", ":5061"),
		},
		IMS: IMSConfig{
			Domain:       getEnv("IMS_DOMAIN", "ims.local"),
			EnablePCSCF:  getEnvBool("ENABLE_PCSCF", true),
			EnableICSCF:  getEnvBool("ENABLE_ICSCF", true),
			EnableSCSCF:  getEnvBool("ENABLE_SCSCF", true),
			EnableHSS:    getEnvBool("ENABLE_HSS", true),
			EnableSBC:    getEnvBool("ENABLE_SBC", true),
			EnableIBCF:   getEnvBool("ENABLE_IBCF", false),
			EnableMGCF:   getEnvBool("ENABLE_MGCF", false),
			EnableBGCF:   getEnvBool("ENABLE_BGCF", false),
			HSS: HSSConfig{
				Backend: getEnv("HSS_BACKEND", "memory"),
				DSN:     getEnv("HSS_DSN", ""),
			},
			SBC: SBCConfig{
				TopologyHiding:  getEnvBool("SBC_TOPOLOGY_HIDING", true),
				NormalizeHeaders: getEnvBool("SBC_NORMALIZE_HEADERS", true),
				RequireTLS:       getEnvBool("SBC_REQUIRE_TLS", false),
				RequireSRTP:      getEnvBool("SBC_REQUIRE_SRTP", false),
				DoSProtection:    getEnvBool("SBC_DOS_PROTECTION", true),
				RateLimitPerIP:   getEnvInt("SBC_RATE_LIMIT_PER_IP", 100),
				RateLimitWindow:  getEnvDuration("SBC_RATE_LIMIT_WINDOW", 60*time.Second),
				EnableSTIR:       getEnvBool("SBC_ENABLE_STIR", false),
				STIRAttestation:  getEnv("SBC_STIR_ATTESTATION", "auto"),
			},
			LI: LIConfig{
				Enabled:          getEnvBool("LI_ENABLED", false),
				MediationDevice:  getEnv("LI_MEDIATION_DEVICE", ""),
				HandoverInterface: getEnv("LI_HANDOVER_INTERFACE", ""),
				AuditLogging:     getEnvBool("LI_AUDIT_LOGGING", true),
			},
			Emergency: EmergencyConfig{
				Enabled:          getEnvBool("EMERGENCY_ENABLED", true),
				EmergencyNumbers: []string{"911", "112", "999", "000"},
				PSAPRoutes:       make(map[string]string),
				BypassSTIR:       getEnvBool("EMERGENCY_BYPASS_STIR", true),
				BypassFraud:      getEnvBool("EMERGENCY_BYPASS_FRAUD", true),
				BypassRateLimit: getEnvBool("EMERGENCY_BYPASS_RATE_LIMIT", true),
			},
		},
		ZeroTrust: ZeroTrustConfig{
			Enabled:    getEnvBool("ZERO_TRUST_MODE", false),
			CAProvider: getEnv("ZTA_CA_PROVIDER", "internal"),
			InternalCA: InternalCAConfig{
				CertPath: getEnv("ZTA_INTERNAL_CA_CERT", "/etc/ims/certs/ca.crt"),
				KeyPath:  getEnv("ZTA_INTERNAL_CA_KEY", "/etc/ims/certs/ca.key"),
			},
			VaultCA: VaultCAConfig{
				Address:  getEnv("ZTA_VAULT_ADDR", ""),
				RoleID:   getEnv("ZTA_VAULT_ROLE_ID", ""),
				SecretID: getEnv("ZTA_VAULT_SECRET_ID", ""),
				PKIPath:  getEnv("ZTA_VAULT_PKI_PATH", "pki/issue/ims"),
				CertPath: getEnv("ZTA_VAULT_CERT_PATH", "/etc/ims/certs"),
			},
			ACME: ACMEConfig{
				Provider: getEnv("ZTA_ACME_PROVIDER", "letsencrypt"),
				Email:    getEnv("ZTA_ACME_EMAIL", ""),
				Staging:  getEnvBool("ZTA_ACME_STAGING", false),
				Domain:   getEnv("ZTA_ACME_DOMAIN", ""),
				CertPath: getEnv("ZTA_ACME_CERT_PATH", "/etc/ims/certs/tls.crt"),
				KeyPath:  getEnv("ZTA_ACME_KEY_PATH", "/etc/ims/certs/tls.key"),
			},
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Version:  getEnv("VERSION", "dev"),
	}

	return cfg
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
