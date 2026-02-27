package pcscf

import (
	"os"
	"strconv"

	gouverneConfig "github.com/dasmlab/ims/internal/gouverne/config"
)

// ConfigFromGouverne creates a P-CSCF Config from gouverne config
func ConfigFromGouverne(cfg *gouverneConfig.Config) *Config {
	pcscfCfg := &Config{
		SIPAddr:    cfg.Server.SIPAddr,
		SIPTLSAddr: cfg.Server.SIPTLSAddr,
		ICSCFAddr:  getEnv("ICSCF_ADDR", "icscf.ims.local:5060"),
		SCSCFAddr:  getEnv("SCSCF_ADDR", ""),

		RequireTLS:      cfg.IMS.SBC.RequireTLS,
		RequireSRTP:     cfg.IMS.SBC.RequireSRTP,
		DoSProtection:   cfg.IMS.SBC.DoSProtection,
		RateLimitPerIP:  cfg.IMS.SBC.RateLimitPerIP,
		RateLimitWindow: cfg.IMS.SBC.RateLimitWindow,

		EnableNATTraversal: getEnvBool("PCSCF_ENABLE_NAT", true),
		PublicIP:           getEnv("PCSCF_PUBLIC_IP", ""),

		EnableCompression: getEnvBool("PCSCF_ENABLE_COMPRESSION", false),

		EmergencyNumbers: cfg.IMS.Emergency.EmergencyNumbers,
	}

	return pcscfCfg
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
