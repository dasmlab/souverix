package sbc

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/dasmlab/ims/internal/config"
	"github.com/dasmlab/ims/internal/sip"
	"github.com/dasmlab/ims/internal/stir"
	"github.com/sirupsen/logrus"
)

// SBC is the Session Border Controller / IBCF implementation
type SBC struct {
	config *config.Config
	log    *logrus.Logger

	// SIP listeners
	udpListener *net.UDPConn
	tcpListener net.Listener
	tlsListener net.Listener

	// Rate limiting
	rateLimiter *RateLimiter

	// Topology hiding
	topologyHiding bool

	// STIR/SHAKEN
	stirSigner   *stir.STIRSigner
	stirVerifier *stir.STIRVerifier
	enableSTIR   bool

	// Message handlers
	handlers map[string]MessageHandler

	mu sync.RWMutex
}

// MessageHandler handles SIP messages
type MessageHandler func(*sip.Message) (*sip.Message, error)

// NewSBC creates a new SBC instance
func NewSBC(cfg *config.Config, log *logrus.Logger) (*SBC, error) {
	sbc := &SBC{
		config:         cfg,
		log:            log,
		topologyHiding: cfg.IMS.SBC.TopologyHiding,
		enableSTIR:     cfg.IMS.SBC.EnableSTIR,
		handlers:       make(map[string]MessageHandler),
	}

	// Initialize rate limiter
	if cfg.IMS.SBC.DoSProtection {
		sbc.rateLimiter = NewRateLimiter(
			cfg.IMS.SBC.RateLimitPerIP,
			cfg.IMS.SBC.RateLimitWindow,
			log,
		)
	}

	// Initialize STIR/SHAKEN if enabled
	if cfg.IMS.SBC.EnableSTIR {
		if err := sbc.initSTIR(cfg); err != nil {
			log.WithError(err).Warn("failed to initialize STIR/SHAKEN, continuing without it")
		}
	}

	return sbc, nil
}

// initSTIR initializes STIR/SHAKEN signing and verification
func (s *SBC) initSTIR(cfg *config.Config) error {
	// Initialize ACME certificate manager for STIR/SHAKEN
	acmeMgr, err := stir.NewACMECertificateManager(&cfg.ZeroTrust.ACME, s.log)
	if err != nil {
		return fmt.Errorf("failed to create ACME certificate manager: %w", err)
	}

	// Determine attestation level
	attestation := stir.AttestationFull
	if cfg.IMS.SBC.STIRAttestation == "B" {
		attestation = stir.AttestationPartial
	} else if cfg.IMS.SBC.STIRAttestation == "C" {
		attestation = stir.AttestationGateway
	} else if cfg.IMS.SBC.STIRAttestation == "auto" {
		// Auto-detect based on subscriber info (simplified - always A for now)
		attestation = stir.AttestationFull
	}

	// Create STIR signer
	s.stirSigner = stir.NewSTIRSigner(
		acmeMgr.GetPrivateKey(),
		acmeMgr.GetCertificateURL(),
		attestation,
	)

	// Create STIR verifier
	s.stirVerifier = stir.NewSTIRVerifier(acmeMgr)

	s.log.Info("STIR/SHAKEN initialized with ACME certificate management")
	return nil
}

// Start starts the SBC listeners
func (s *SBC) Start() error {
	// Start UDP listener
	if err := s.startUDP(); err != nil {
		return fmt.Errorf("failed to start UDP listener: %w", err)
	}

	// Start TCP listener
	if err := s.startTCP(); err != nil {
		return fmt.Errorf("failed to start TCP listener: %w", err)
	}

	// Start TLS listener if configured
	if s.config.IMS.SBC.RequireTLS {
		if err := s.startTLS(); err != nil {
			return fmt.Errorf("failed to start TLS listener: %w", err)
		}
	}

	s.log.Info("SBC started")
	return nil
}

// Stop stops the SBC listeners
func (s *SBC) Stop() error {
	if s.udpListener != nil {
		s.udpListener.Close()
	}
	if s.tcpListener != nil {
		s.tcpListener.Close()
	}
	if s.tlsListener != nil {
		s.tlsListener.Close()
	}

	s.log.Info("SBC stopped")
	return nil
}

// RegisterHandler registers a message handler for a SIP method
func (s *SBC) RegisterHandler(method string, handler MessageHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[method] = handler
}

// ProcessMessage processes a SIP message through the SBC
func (s *SBC) ProcessMessage(msg *sip.Message, remoteAddr string) (*sip.Message, error) {
	start := time.Now()
	msg.RemoteAddr = remoteAddr

	// Rate limiting
	if s.rateLimiter != nil {
		if !s.rateLimiter.Allow(remoteAddr) {
			s.log.WithFields(logrus.Fields{
				"remote_addr": remoteAddr,
				"method":      msg.Method,
			}).Warn("rate limit exceeded")

			// Return 503 Service Unavailable
			response := &sip.Message{
				Version:    "SIP/2.0",
				StatusCode: sip.StatusServiceUnavailable,
				StatusText: "Service Unavailable",
				Headers:    make(map[string][]string),
			}
			response.SetHeader("Via", msg.GetHeader("Via"))
			response.SetHeader("From", msg.GetHeader("From"))
			response.SetHeader("To", msg.GetHeader("To"))
			response.SetHeader("Call-ID", msg.GetHeader("Call-ID"))
			response.SetHeader("CSeq", msg.GetHeader("CSeq"))

			return response, nil
		}
	}

	// SIP normalization
	if s.config.IMS.SBC.NormalizeHeaders {
		s.normalizeHeaders(msg)
	}

	// Emergency call handling (highest priority - bypasses everything)
	if msg.IsRequest() && msg.Method == sip.MethodINVITE {
		isEmergency, err := s.processEmergency(msg)
		if err != nil {
			s.log.WithError(err).Error("emergency processing error")
		}
		if isEmergency {
			// Emergency call - skip normal processing and return
			return msg, nil
		}
	}

	// STIR/SHAKEN signing (for outgoing INVITE) - skipped for emergency
	if s.enableSTIR && msg.IsRequest() && msg.Method == sip.MethodINVITE {
		if err := s.signSTIR(msg); err != nil {
			s.log.WithError(err).Warn("failed to sign STIR/SHAKEN")
		}
	}

	// STIR/SHAKEN verification (for incoming INVITE)
	if s.enableSTIR && msg.IsRequest() && msg.Method == sip.MethodINVITE {
		if err := s.verifySTIR(msg); err != nil {
			s.log.WithError(err).Warn("STIR/SHAKEN verification failed")
			// Continue processing but log the failure
		}
	}

	// Topology hiding
	if s.topologyHiding {
		s.hideTopology(msg)
	}

	// Find and call handler
	var handler MessageHandler
	if msg.IsRequest() {
		s.mu.RLock()
		handler = s.handlers[msg.Method]
		s.mu.RUnlock()
	}

	var response *sip.Message
	var err error

	if handler != nil {
		response, err = handler(msg)
	} else {
		// Default handler: forward message
		response, err = s.defaultHandler(msg)
	}

	duration := time.Since(start)
	s.log.WithFields(logrus.Fields{
		"method":   msg.Method,
		"duration": duration,
		"remote":   remoteAddr,
	}).Debug("message processed")

	return response, err
}

// normalizeHeaders normalizes SIP headers
func (s *SBC) normalizeHeaders(msg *sip.Message) {
	// Normalize header names (capitalize first letter of each word)
	// This is a simplified version - real normalization is more complex
	normalized := make(map[string][]string)
	for k, v := range msg.Headers {
		normalized[normalizeHeaderName(k)] = v
	}
	msg.Headers = normalized
}

// normalizeHeaderName normalizes a header name
func normalizeHeaderName(name string) string {
	// Simple capitalization - real implementation would be more sophisticated
	parts := strings.Split(strings.ToLower(name), "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "-")
}

// hideTopology performs topology hiding
func (s *SBC) hideTopology(msg *sip.Message) {
	// Remove or modify headers that reveal topology
	// - Remove Record-Route headers from responses
	// - Modify Via headers
	// - Remove Server/User-Agent headers

	if msg.IsResponse() {
		msg.Headers["Record-Route"] = nil
	}

	// Modify Via header to hide internal topology
	if via := msg.GetHeader("Via"); via != "" {
		// Extract only the transport and sent-by, remove branch and other params
		parts := strings.Split(via, ";")
		if len(parts) > 0 {
			msg.SetHeader("Via", parts[0]+";branch=z9hG4bK"+generateBranch())
		}
	}

	// Remove Server header
	delete(msg.Headers, "Server")
	delete(msg.Headers, "User-Agent")
}

// generateBranch generates a SIP branch parameter
func generateBranch() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// defaultHandler is the default message handler (forward)
func (s *SBC) defaultHandler(msg *sip.Message) (*sip.Message, error) {
	// In a real implementation, this would forward to the next hop
	// For now, return a 200 OK for requests
	if msg.IsRequest() {
		response := &sip.Message{
			Version:    "SIP/2.0",
			StatusCode: sip.StatusOK,
			StatusText: "OK",
			Headers:    make(map[string][]string),
		}

		// Copy required headers
		response.SetHeader("Via", msg.GetHeader("Via"))
		response.SetHeader("From", msg.GetHeader("From"))
		response.SetHeader("To", msg.GetHeader("To"))
		response.SetHeader("Call-ID", msg.GetHeader("Call-ID"))
		response.SetHeader("CSeq", msg.GetHeader("CSeq"))

		return response, nil
	}

	return msg, nil
}

// startUDP starts the UDP listener
func (s *SBC) startUDP() error {
	addr, err := net.ResolveUDPAddr("udp", s.config.Server.SIPAddr)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	s.udpListener = conn

	go s.handleUDP()

	s.log.WithField("addr", addr).Info("UDP listener started")
	return nil
}

// startTCP starts the TCP listener
func (s *SBC) startTCP() error {
	listener, err := net.Listen("tcp", s.config.Server.SIPAddr)
	if err != nil {
		return err
	}

	s.tcpListener = listener

	go s.handleTCP()

	s.log.WithField("addr", s.config.Server.SIPAddr).Info("TCP listener started")
	return nil
}

// startTLS starts the TLS listener
func (s *SBC) startTLS() error {
	// TODO: Implement TLS listener with certificate management
	s.log.Warn("TLS listener not yet implemented")
	return nil
}

// handleUDP handles UDP connections
func (s *SBC) handleUDP() {
	buffer := make([]byte, 65535)
	for {
		n, addr, err := s.udpListener.ReadFromUDP(buffer)
		if err != nil {
			s.log.WithError(err).Error("UDP read error")
			continue
		}

		go s.handleMessage(buffer[:n], addr.String(), "udp")
	}
}

// handleTCP handles TCP connections
func (s *SBC) handleTCP() {
	for {
		conn, err := s.tcpListener.Accept()
		if err != nil {
			s.log.WithError(err).Error("TCP accept error")
			continue
		}

		go s.handleTCPConnection(conn)
	}
}

// handleTCPConnection handles a single TCP connection
func (s *SBC) handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	parser := sip.NewParser()
	for {
		msg, err := parser.ParseMessage(conn)
		if err != nil {
			s.log.WithError(err).Error("failed to parse SIP message")
			break
		}

		response, err := s.ProcessMessage(msg, conn.RemoteAddr().String())
		if err != nil {
			s.log.WithError(err).Error("failed to process message")
			break
		}

		if response != nil {
			if _, err := conn.Write([]byte(response.String())); err != nil {
				s.log.WithError(err).Error("failed to write response")
				break
			}
		}
	}
}

// handleMessage handles a SIP message
func (s *SBC) handleMessage(data []byte, remoteAddr, transport string) {
	parser := sip.NewParser()
	msg, err := parser.ParseMessage(bytes.NewReader(data))
	if err != nil {
		s.log.WithError(err).Error("failed to parse SIP message")
		return
	}

	msg.Transport = transport

	response, err := s.ProcessMessage(msg, remoteAddr)
	if err != nil {
		s.log.WithError(err).Error("failed to process message")
		return
	}

	if response != nil && s.udpListener != nil {
		addr, err := net.ResolveUDPAddr("udp", remoteAddr)
		if err == nil {
			s.udpListener.WriteToUDP([]byte(response.String()), addr)
		}
	}
}
