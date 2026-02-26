package ibcf

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/dasmlab/ims/internal/config"
	"github.com/dasmlab/ims/internal/sip"
	"github.com/dasmlab/ims/internal/stir"
	"github.com/sirupsen/logrus"
)

// IBCF is the Interconnection Border Control Function (3GPP TS 23.228)
// It provides standardized border control between IMS networks
type IBCF struct {
	config *config.Config
	log    *logrus.Logger

	// SIP listeners
	udpListener *net.UDPConn
	tcpListener net.Listener
	tlsListener net.Listener

	// Security
	allowedPeers map[string]bool // Whitelist of allowed peer domains
	requireTLS   bool
	requireSTIR  bool
	minAttestation stir.AttestationLevel

	// Topology hiding
	topologyHiding bool
	internalDomain string

	// Policy enforcement
	policy PolicyEngine

	// STIR/SHAKEN
	stirSigner   *stir.STIRSigner
	stirVerifier *stir.STIRVerifier

	mu sync.RWMutex
}

// PolicyEngine enforces inter-operator peering policies
type PolicyEngine interface {
	// IsPeerAllowed checks if a peer is allowed
	IsPeerAllowed(peerDomain string) bool

	// IsCallAllowed checks if a call is allowed based on policy
	IsCallAllowed(msg *sip.Message) (bool, string)

	// GetAttestationRequirement returns minimum required attestation
	GetAttestationRequirement() stir.AttestationLevel
}

// SimplePolicyEngine is a basic policy engine implementation
type SimplePolicyEngine struct {
	allowedPeers   map[string]bool
	requireSTIR    bool
	minAttestation stir.AttestationLevel
	log            *logrus.Logger
}

// NewSimplePolicyEngine creates a simple policy engine
func NewSimplePolicyEngine(allowedPeers []string, requireSTIR bool, minAttestation stir.AttestationLevel, log *logrus.Logger) *SimplePolicyEngine {
	peers := make(map[string]bool)
	for _, peer := range allowedPeers {
		peers[peer] = true
	}

	return &SimplePolicyEngine{
		allowedPeers:   peers,
		requireSTIR:    requireSTIR,
		minAttestation: minAttestation,
		log:            log,
	}
}

// IsPeerAllowed checks if a peer domain is allowed
func (p *SimplePolicyEngine) IsPeerAllowed(peerDomain string) bool {
	// If no peers configured, allow all (for development)
	if len(p.allowedPeers) == 0 {
		return true
	}

	// Check exact match
	if p.allowedPeers[peerDomain] {
		return true
	}

	// Check subdomain matches
	for allowed := range p.allowedPeers {
		if strings.HasSuffix(peerDomain, "."+allowed) {
			return true
		}
	}

	return false
}

// IsCallAllowed checks if a call is allowed
func (p *SimplePolicyEngine) IsCallAllowed(msg *sip.Message) (bool, string) {
	// Extract peer domain from From or Contact header
	from := msg.GetHeader("From")
	peerDomain := extractDomain(from)

	if !p.IsPeerAllowed(peerDomain) {
		return false, fmt.Sprintf("peer domain not allowed: %s", peerDomain)
	}

	// Check STIR/SHAKEN requirement
	if p.requireSTIR {
		identity := msg.GetHeader("Identity")
		if identity == "" {
			return false, "STIR/SHAKEN Identity header required"
		}

		// Verify attestation level (would need to parse token)
		// For now, just check presence
	}

	return true, ""
}

// GetAttestationRequirement returns minimum required attestation
func (p *SimplePolicyEngine) GetAttestationRequirement() stir.AttestationLevel {
	return p.minAttestation
}

// NewIBCF creates a new IBCF instance
func NewIBCF(cfg *config.Config, log *logrus.Logger) (*IBCF, error) {
	ibcf := &IBCF{
		config:         cfg,
		log:            log,
		topologyHiding: cfg.IMS.SBC.TopologyHiding,
		internalDomain: cfg.IMS.Domain,
		allowedPeers:   make(map[string]bool),
		requireTLS:     cfg.IMS.SBC.RequireTLS,
		requireSTIR:    cfg.IMS.SBC.EnableSTIR,
	}

	// Initialize policy engine
	allowedPeersEnv := os.Getenv("IBCF_ALLOWED_PEERS")
	if allowedPeersEnv == "" {
		allowedPeersEnv = ""
	}
	allowedPeersList := strings.Split(allowedPeersEnv, ",")
	peers := make([]string, 0)
	for _, peer := range allowedPeersList {
		if peer = strings.TrimSpace(peer); peer != "" {
			peers = append(peers, peer)
			ibcf.allowedPeers[peer] = true
		}
	}

	minAttestation := stir.AttestationFull
	if cfg.IMS.SBC.STIRAttestation == "B" {
		minAttestation = stir.AttestationPartial
	} else if cfg.IMS.SBC.STIRAttestation == "C" {
		minAttestation = stir.AttestationGateway
	}

	ibcf.policy = NewSimplePolicyEngine(peers, ibcf.requireSTIR, minAttestation, log)

	// Initialize STIR/SHAKEN if enabled
	if cfg.IMS.SBC.EnableSTIR {
		if err := ibcf.initSTIR(cfg); err != nil {
			log.WithError(err).Warn("failed to initialize STIR/SHAKEN in IBCF")
		}
	}

	ibcf.log.Info("IBCF initialized")
	return ibcf, nil
}

// initSTIR initializes STIR/SHAKEN for IBCF
func (i *IBCF) initSTIR(cfg *config.Config) error {
	// Similar to SBC STIR initialization
	acmeMgr, err := stir.NewACMECertificateManager(&cfg.ZeroTrust.ACME, i.log)
	if err != nil {
		return fmt.Errorf("failed to create ACME certificate manager: %w", err)
	}

	attestation := stir.AttestationFull
	if cfg.IMS.SBC.STIRAttestation == "B" {
		attestation = stir.AttestationPartial
	} else if cfg.IMS.SBC.STIRAttestation == "C" {
		attestation = stir.AttestationGateway
	}

	i.stirSigner = stir.NewSTIRSigner(
		acmeMgr.GetPrivateKey(),
		acmeMgr.GetCertificateURL(),
		attestation,
	)

	i.stirVerifier = stir.NewSTIRVerifier(acmeMgr)

	i.log.Info("IBCF STIR/SHAKEN initialized")
	return nil
}

// ProcessMessage processes a SIP message through the IBCF
// This implements the core IBCF functions per 3GPP TS 23.228
func (i *IBCF) ProcessMessage(msg *sip.Message, remoteAddr string) (*sip.Message, error) {
	// 1. Message Inspection & Enforcement
	if err := i.validateMessage(msg); err != nil {
		i.log.WithError(err).Warn("message validation failed")
		return i.createErrorResponse(msg, sip.StatusBadRequest, "Invalid message"), nil
	}

	// 2. Policy Enforcement (Inter-Operator Peering Control)
	if allowed, reason := i.policy.IsCallAllowed(msg); !allowed {
		i.log.WithFields(logrus.Fields{
			"reason": reason,
			"from":   msg.GetHeader("From"),
		}).Warn("call rejected by policy")
		return i.createErrorResponse(msg, sip.StatusForbidden, reason), nil
	}

	// 3. STIR/SHAKEN Verification (for inbound)
	if msg.IsRequest() && msg.Method == sip.MethodINVITE {
		if i.requireSTIR && i.stirVerifier != nil {
			if err := i.verifySTIR(msg); err != nil {
				i.log.WithError(err).Warn("STIR verification failed")
				// Continue but log the failure
			}
		}
	}

	// 4. Topology Hiding
	if i.topologyHiding {
		i.hideTopology(msg)
	}

	// 5. SIP Header Normalization
	i.normalizeHeaders(msg)

	// 6. STIR/SHAKEN Signing (for outbound)
	if msg.IsRequest() && msg.Method == sip.MethodINVITE {
		if i.requireSTIR && i.stirSigner != nil {
			if err := i.signSTIR(msg); err != nil {
				i.log.WithError(err).Warn("STIR signing failed")
			}
		}
	}

	return msg, nil
}

// validateMessage validates SIP message structure
func (i *IBCF) validateMessage(msg *sip.Message) error {
	// Check required headers
	if msg.IsRequest() {
		if msg.Method == "" {
			return fmt.Errorf("missing method")
		}
		if msg.URI == "" {
			return fmt.Errorf("missing request URI")
		}
	} else {
		if msg.StatusCode == 0 {
			return fmt.Errorf("missing status code")
		}
	}

	// Check required headers
	requiredHeaders := []string{"Via", "From", "To", "Call-ID", "CSeq"}
	for _, header := range requiredHeaders {
		if msg.GetHeader(header) == "" {
			return fmt.Errorf("missing required header: %s", header)
		}
	}

	// Validate SIP method (IBCF may restrict certain methods)
	allowedMethods := map[string]bool{
		sip.MethodINVITE:  true,
		sip.MethodACK:     true,
		sip.MethodBYE:     true,
		sip.MethodCANCEL:  true,
		sip.MethodOPTIONS: true,
		sip.MethodUPDATE:  true,
		sip.MethodPRACK:   true,
	}

	if msg.IsRequest() && !allowedMethods[msg.Method] {
		return fmt.Errorf("method not allowed: %s", msg.Method)
	}

	return nil
}

// hideTopology performs topology hiding per 3GPP requirements
func (i *IBCF) hideTopology(msg *sip.Message) {
	// Remove Record-Route headers (reveal internal routing)
	if msg.IsResponse() {
		delete(msg.Headers, "Record-Route")
	}

	// Rewrite Via headers to hide internal hops
	if via := msg.GetHeader("Via"); via != "" {
		// Extract only transport and sent-by, remove branch and other params
		parts := strings.Split(via, ";")
		if len(parts) > 0 {
			// Replace internal domain with external domain
			viaBase := parts[0]
			viaBase = strings.ReplaceAll(viaBase, i.internalDomain, "border.ims.local")
			msg.SetHeader("Via", viaBase+";branch=z9hG4bK"+generateBranch())
		}
	}

	// Replace Contact URIs to hide internal addresses
	if contact := msg.GetHeader("Contact"); contact != "" {
		// Replace internal domain/IP with border gateway
		contact = strings.ReplaceAll(contact, i.internalDomain, "border.ims.local")
		// Remove internal IP addresses
		contact = removeInternalIPs(contact)
		msg.SetHeader("Contact", contact)
	}

	// Remove Server/User-Agent headers
	delete(msg.Headers, "Server")
	delete(msg.Headers, "User-Agent")

	// Remove internal routing information from other headers
	// (implementation specific)
}

// normalizeHeaders normalizes SIP headers
func (i *IBCF) normalizeHeaders(msg *sip.Message) {
	// Normalize From/To URIs
	if from := msg.GetHeader("From"); from != "" {
		normalized := i.normalizeURI(from)
		msg.SetHeader("From", normalized)
	}

	if to := msg.GetHeader("To"); to != "" {
		normalized := i.normalizeURI(to)
		msg.SetHeader("To", normalized)
	}

	// Normalize header names (capitalize properly)
	normalized := make(map[string][]string)
	for k, v := range msg.Headers {
		normalized[normalizeHeaderName(k)] = v
	}
	msg.Headers = normalized
}

// normalizeURI normalizes a SIP URI
func (i *IBCF) normalizeURI(uri string) string {
	// Replace internal domains with external domain
	uri = strings.ReplaceAll(uri, "@"+i.internalDomain, "@ims.local")
	return uri
}

// signSTIR signs an INVITE with STIR/SHAKEN
func (i *IBCF) signSTIR(msg *sip.Message) error {
	if i.stirSigner == nil {
		return fmt.Errorf("STIR signer not initialized")
	}

	from := msg.GetHeader("From")
	to := msg.GetHeader("To")
	callID := msg.GetHeader("Call-ID")

	origTN := extractTN(from)
	destTN := extractTN(to)

	if origTN == "" || destTN == "" {
		return nil // Skip if no TNs
	}

	token, err := i.stirSigner.SignINVITE(origTN, destTN, callID)
	if err != nil {
		return err
	}

	msg.SetHeader("Identity", token)
	return nil
}

// verifySTIR verifies STIR/SHAKEN signature
func (i *IBCF) verifySTIR(msg *sip.Message) error {
	if i.stirVerifier == nil {
		return fmt.Errorf("STIR verifier not initialized")
	}

	identity := msg.GetHeader("Identity")
	if identity == "" {
		return nil // No identity header, skip
	}

	identityToken, err := stir.ParseIdentityHeader(identity)
	if err != nil {
		return err
	}

	passport, err := i.stirVerifier.VerifyINVITE(identityToken)
	if err != nil {
		return err
	}

	// Check attestation level requirement
	minAttestation := i.policy.GetAttestationRequirement()
	if !isAttestationSufficient(passport.Attest, minAttestation) {
		return fmt.Errorf("attestation level insufficient: got %s, required %s", passport.Attest, minAttestation)
	}

	msg.SetHeader("X-STIR-Attestation", string(passport.Attest))
	msg.SetHeader("X-STIR-Verified", "true")

	return nil
}

// Helper functions
func extractDomain(uri string) string {
	parts := strings.Split(uri, "@")
	if len(parts) > 1 {
		domain := strings.Split(parts[1], ";")[0]
		domain = strings.Trim(domain, ">")
		return domain
	}
	return ""
}

func extractTN(uri string) string {
	uri = strings.Trim(uri, "<>")
	parts := strings.Split(uri, "@")
	if len(parts) > 0 {
		userPart := strings.TrimPrefix(parts[0], "sip:")
		userPart = strings.Split(userPart, ";")[0]
		return strings.TrimSpace(userPart)
	}
	return ""
}

func removeInternalIPs(contact string) string {
	// Remove private IP addresses (simplified)
	privateIPs := []string{
		"192.168.",
		"10.",
		"172.16.",
		"172.17.",
		"172.18.",
		"172.19.",
		"172.20.",
		"172.21.",
		"172.22.",
		"172.23.",
		"172.24.",
		"172.25.",
		"172.26.",
		"172.27.",
		"172.28.",
		"172.29.",
		"172.30.",
		"172.31.",
	}

	for _, prefix := range privateIPs {
		if strings.Contains(contact, prefix) {
			// Replace with border gateway
			contact = strings.ReplaceAll(contact, prefix, "border.ims.local")
		}
	}

	return contact
}

func normalizeHeaderName(name string) string {
	parts := strings.Split(strings.ToLower(name), "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "-")
}

func generateBranch() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func isAttestationSufficient(got, required stir.AttestationLevel) bool {
	levels := map[stir.AttestationLevel]int{
		stir.AttestationFull:    3,
		stir.AttestationPartial: 2,
		stir.AttestationGateway: 1,
	}

	return levels[got] >= levels[required]
}

func (i *IBCF) createErrorResponse(original *sip.Message, statusCode int, reason string) *sip.Message {
	response := &sip.Message{
		Version:    "SIP/2.0",
		StatusCode: statusCode,
		StatusText: reason,
		Headers:    make(map[string][]string),
	}

	// Copy required headers
	response.SetHeader("Via", original.GetHeader("Via"))
	response.SetHeader("From", original.GetHeader("From"))
	response.SetHeader("To", original.GetHeader("To"))
	response.SetHeader("Call-ID", original.GetHeader("Call-ID"))
	response.SetHeader("CSeq", original.GetHeader("CSeq"))

	return response
}
