package testrig

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// PIXITConfig holds Protocol Implementation eXtra Information for Testing
type PIXITConfig struct {
	Timers    TimerConfig    `yaml:"timers"`
	Codec     CodecConfig    `yaml:"codec"`
	TLS       TLSConfig      `yaml:"tls"`
	Peer      PeerConfig     `yaml:"peer"`
	STIR      STIRConfig      `yaml:"stir"`
	LI        LIConfig        `yaml:"li"`
	Emergency EmergencyConfig `yaml:"emergency"`
	Chaos     ChaosConfig     `yaml:"chaos"`
}

// TimerConfig holds SIP timer settings
type TimerConfig struct {
	T1          time.Duration `yaml:"t1"`
	T2          time.Duration `yaml:"t2"`
	TimerB      time.Duration `yaml:"timer_b"`
	TimerF      time.Duration `yaml:"timer_f"`
	SessionTimer time.Duration `yaml:"session_timer"`
	MaxDialogs  int           `yaml:"max_dialogs"`
	MaxCPS      int           `yaml:"max_cps"`
}

// CodecConfig holds codec policy settings
type CodecConfig struct {
	AudioCodecs []string `yaml:"audio"`
	VideoCodecs []string `yaml:"video"`
	Transcoding  bool     `yaml:"transcoding"`
	SRTPRequired bool     `yaml:"srtp_required"`
	DTMFMode     string   `yaml:"dtmf_mode"` // "RFC2833" or "SIP-INFO"
}

// TLSConfig holds TLS/security settings
type TLSConfig struct {
	Version      []string `yaml:"version"`      // ["1.2", "1.3"]
	CipherSuites []string `yaml:"cipher_suites"`
	MTLSRequired bool     `yaml:"mtls_required"`
	CertReloadMode string  `yaml:"cert_reload_mode"` // "hot" or "restart"
	OCSPMode      string   `yaml:"ocsp_mode"`       // "hard-fail" or "soft-fail"
	STIREnforcement string `yaml:"stir_enforcement"`  // "soft" or "hard"
}

// PeerConfig holds peer profile settings
type PeerConfig struct {
	PeerID              string `yaml:"peer_id"`
	Transport           string `yaml:"transport"` // "SIP-TLS", "SIP-TCP", "SIP-UDP"
	IPVersion           string `yaml:"ip_version"` // "IPv4" or "IPv6"
	AuthMode            string `yaml:"auth_mode"`  // "mTLS", "IP", "None"
	STIRTrustLevel      string `yaml:"stir_trust_level"` // "Trusted" or "External"
	EmergencyRouting    bool   `yaml:"emergency_routing"`
	LICooperation       bool   `yaml:"li_cooperation"`
	MaxCPS              int    `yaml:"max_cps"`
	TopologyHidingMode  string `yaml:"topology_hiding_mode"` // "Full", "Partial", "None"
}

// STIRConfig holds STIR/SHAKEN settings
type STIRConfig struct {
	AttestationPolicy   string        `yaml:"attestation_policy"`   // "auto", "A", "B", "C"
	IATSkewTolerance    time.Duration `yaml:"iat_skew"`
	CertCacheTTL        time.Duration `yaml:"cert_cache_ttl"`
	SigningKeySource    string        `yaml:"signing_key_source"`   // "HSM", "File", "Vault"
	ReSignTransitCalls  bool          `yaml:"re_sign_transit_calls"`
	IdentityHeaderMaxSize int          `yaml:"identity_header_max_size"` // bytes
}

// LIConfig holds Lawful Intercept settings
type LIConfig struct {
	Mode              string `yaml:"mode"`                // "disabled", "signaling", "full"
	MediationIP       string `yaml:"mediation_ip"`
	InterceptTargetList string `yaml:"intercept_target_list"` // "dynamic" or "static"
	OverloadPolicy    string `yaml:"overload_policy"`     // "preserve" or "drop"
	AuditLogRetention int    `yaml:"audit_log_retention"` // days
	TLSDecryptForLI   bool   `yaml:"tls_decrypt_for_li"`
}

// EmergencyConfig holds Emergency Services settings
type EmergencyConfig struct {
	Numbers              []string `yaml:"numbers"`
	PriorityQueue        bool     `yaml:"priority_queue"`
	STIROverride         bool     `yaml:"stir_override"`
	FraudOverride        bool     `yaml:"fraud_override"`
	EmergencyRoute       string   `yaml:"emergency_route"`
	LocationHeaderMandatory bool   `yaml:"location_header_mandatory"`
}

// ChaosConfig holds chaos injection settings
type ChaosConfig struct {
	PacketLoss          float64 `yaml:"packet_loss"`          // percentage
	Jitter              int     `yaml:"jitter"`               // milliseconds
	DNSFailure          bool    `yaml:"dns_failure"`
	TLSCertExpirySim    bool    `yaml:"tls_cert_expiry_sim"`
	SigningServiceCrash bool    `yaml:"signing_service_crash"`
	StateStorePartition bool    `yaml:"state_store_partition"`
	MediaRelayKill      bool    `yaml:"media_relay_kill"`
}

// DefaultPIXIT returns default PIXIT configuration
func DefaultPIXIT() *PIXITConfig {
	return &PIXITConfig{
		Timers: TimerConfig{
			T1:          500 * time.Millisecond,
			T2:          4 * time.Second,
			TimerB:      32 * time.Second,
			TimerF:      32 * time.Second,
			SessionTimer: 1800 * time.Second,
			MaxDialogs:  50000,
			MaxCPS:      2000,
		},
		Codec: CodecConfig{
			AudioCodecs:  []string{"AMR-WB", "G.711"},
			VideoCodecs:  []string{"H.264"},
			Transcoding:  false,
			SRTPRequired: true,
			DTMFMode:     "RFC2833",
		},
		TLS: TLSConfig{
			Version:        []string{"1.2", "1.3"},
			CipherSuites:   []string{"ECDHE-ECDSA-AES256-GCM"},
			MTLSRequired:   true,
			CertReloadMode: "hot",
			OCSPMode:       "soft-fail",
			STIREnforcement: "soft",
		},
		Peer: PeerConfig{
			PeerID:             "PEER-A",
			Transport:         "SIP-TLS",
			IPVersion:         "IPv4",
			AuthMode:          "mTLS",
			STIRTrustLevel:    "Trusted",
			EmergencyRouting:  true,
			LICooperation:     true,
			MaxCPS:            1000,
			TopologyHidingMode: "Full",
		},
		STIR: STIRConfig{
			AttestationPolicy:   "auto",
			IATSkewTolerance:   60 * time.Second,
			CertCacheTTL:       24 * time.Hour,
			SigningKeySource:    "HSM",
			ReSignTransitCalls:  false,
			IdentityHeaderMaxSize: 8 * 1024, // 8 KB
		},
		LI: LIConfig{
			Mode:              "disabled",
			MediationIP:       "10.1.1.10",
			InterceptTargetList: "dynamic",
			OverloadPolicy:    "preserve",
			AuditLogRetention: 180,
			TLSDecryptForLI:   true,
		},
		Emergency: EmergencyConfig{
			Numbers:                []string{"911", "112"},
			PriorityQueue:         true,
			STIROverride:          true,
			FraudOverride:         true,
			EmergencyRoute:        "PSAP-A",
			LocationHeaderMandatory: true,
		},
		Chaos: ChaosConfig{
			PacketLoss:          0.0,
			Jitter:              0,
			DNSFailure:          false,
			TLSCertExpirySim:    false,
			SigningServiceCrash: false,
			StateStorePartition: false,
			MediaRelayKill:      false,
		},
	}
}

// LoadPIXIT loads PIXIT configuration from YAML
func LoadPIXIT(filename string) (*PIXITConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read PIXIT file: %w", err)
	}

	var config PIXITConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse PIXIT YAML: %w", err)
	}

	return &config, nil
}

// SavePIXIT saves PIXIT configuration to YAML
func (p *PIXITConfig) SavePIXIT(filename string) error {
	data, err := yaml.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal PIXIT: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write PIXIT file: %w", err)
	}

	return nil
}

// MergePIXIT merges another PIXIT config into this one
func (p *PIXITConfig) MergePIXIT(other *PIXITConfig) {
	// Merge timers
	if other.Timers.T1 > 0 {
		p.Timers.T1 = other.Timers.T1
	}
	if other.Timers.T2 > 0 {
		p.Timers.T2 = other.Timers.T2
	}
	// ... (merge other fields as needed)
}

// ValidatePIXIT validates PIXIT configuration
func (p *PIXITConfig) ValidatePIXIT() error {
	if p.Timers.T1 < 100*time.Millisecond || p.Timers.T1 > 1000*time.Millisecond {
		return fmt.Errorf("T1 must be between 100ms and 1000ms")
	}

	if p.Timers.MaxCPS < 100 || p.Timers.MaxCPS > 20000 {
		return fmt.Errorf("MaxCPS must be between 100 and 20000")
	}

	if p.Chaos.PacketLoss < 0 || p.Chaos.PacketLoss > 20 {
		return fmt.Errorf("packet loss must be between 0% and 20%")
	}

	return nil
}
