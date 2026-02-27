// Package pcscf implements the Proxy CSCF (P-CSCF) node.
// P-CSCF is the first contact point for User Equipment (UE) in the IMS network.
// Per 3GPP TS 23.228
//
// CI/CD Pipeline: Build → Unit Test (Diagnostic API) → Publish → System Test → Stable
// See docs/DESIGN_PHILOSOPHY.md for SDLC methodology
package pcscf

// 🚀 Triggering P-CSCF container build
import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dasmlab/ims/internal/common/node"
	"github.com/dasmlab/ims/internal/common/sip"
	"github.com/sirupsen/logrus"
)

// Config holds P-CSCF configuration
type Config struct {
	// SIP listener addresses
	SIPAddr    string // UDP/TCP address (e.g., ":5060")
	SIPTLSAddr string // TLS address (e.g., ":5061")

	// Next hop addresses
	ICSCFAddr string // I-CSCF address for registration
	SCSCFAddr string // S-CSCF address (if known)

	// Security settings
	RequireTLS      bool
	RequireSRTP     bool
	DoSProtection   bool
	RateLimitPerIP  int
	RateLimitWindow time.Duration

	// NAT traversal
	EnableNATTraversal bool
	PublicIP           string // Public IP for NAT traversal

	// Compression
	EnableCompression bool

	// Emergency services
	EmergencyNumbers []string
}

// PCSCF implements the Proxy CSCF node.
type PCSCF struct {
	*node.BaseNode
	config *Config
	log    *logrus.Logger

	// SIP listeners
	udpListener *net.UDPConn
	tcpListener net.Listener
	tlsListener net.Listener

	// Rate limiting
	rateLimiter *RateLimiter

	// NAT traversal mapping (Contact -> public address)
	natMapping map[string]string
	natMu      sync.RWMutex

	// Active sessions
	sessions map[string]*Session
	sessMu   sync.RWMutex

	// Message handlers
	handlers map[string]MessageHandler

	// Parser
	parser *sip.Parser

	// Shutdown
	shutdown chan struct{}
	wg       sync.WaitGroup

	// Diagnostic server
	diagnostic *DiagnosticServer

	mu sync.RWMutex
}

// Session represents an active SIP session
type Session struct {
	CallID       string
	From         string
	To           string
	Contact      string
	RemoteAddr   string
	Transport    string
	CreatedAt    time.Time
	LastActivity time.Time
}

// MessageHandler handles SIP messages
type MessageHandler func(*sip.Message, *Session) (*sip.Message, error)

// New creates a new P-CSCF instance.
func New(cfg *Config, log *logrus.Logger) *PCSCF {
	pcscf := &PCSCF{
		BaseNode:   node.NewBaseNode("pcscf"),
		config:     cfg,
		log:        log,
		natMapping: make(map[string]string),
		sessions:   make(map[string]*Session),
		handlers:   make(map[string]MessageHandler),
		parser:     sip.NewParser(),
		shutdown:   make(chan struct{}),
	}

	// Initialize rate limiter if DoS protection is enabled
	if cfg.DoSProtection {
		pcscf.rateLimiter = NewRateLimiter(
			cfg.RateLimitPerIP,
			cfg.RateLimitWindow,
			log,
		)
	}

	// Register default handlers
	pcscf.registerDefaultHandlers()

	return pcscf
}

// Start initializes and starts the P-CSCF node.
func (p *PCSCF) Start(ctx context.Context) error {
	p.log.Info("starting P-CSCF")

	// Start UDP listener
	if err := p.startUDP(); err != nil {
		return fmt.Errorf("failed to start UDP listener: %w", err)
	}

	// Start TCP listener
	if err := p.startTCP(); err != nil {
		return fmt.Errorf("failed to start TCP listener: %w", err)
	}

	// Start TLS listener if configured
	if p.config.RequireTLS {
		if err := p.startTLS(); err != nil {
			return fmt.Errorf("failed to start TLS listener: %w", err)
		}
	}

	// Start session cleanup goroutine
	p.wg.Add(1)
	go p.sessionCleanup()

	// Start diagnostic server if enabled
	diagEnabled := os.Getenv("DIAGNOSTICS_ENABLED") == "true"
	if diagEnabled {
		p.diagnostic = NewDiagnosticServer(p, true)
		diagAddr := os.Getenv("DIAGNOSTICS_ADDR")
		if diagAddr == "" {
			diagAddr = ":8081" // Default diagnostic port
		}
		if err := p.diagnostic.Start(diagAddr); err != nil {
			p.log.WithError(err).Warn("failed to start diagnostic server")
		}
	}

	p.BaseNode.SetHealth("healthy", map[string]interface{}{
		"started":     true,
		"sip_udp":     p.config.SIPAddr,
		"sip_tcp":     p.config.SIPAddr,
		"sip_tls":     p.config.SIPTLSAddr,
		"nat_enabled": p.config.EnableNATTraversal,
		"compression": p.config.EnableCompression,
		"diagnostics": diagEnabled,
	})

	p.log.Info("P-CSCF started")
	return nil
}

// Stop gracefully stops the P-CSCF node.
func (p *PCSCF) Stop(ctx context.Context) error {
	p.log.Info("stopping P-CSCF")

	// Signal shutdown
	close(p.shutdown)

	// Stop diagnostic server
	if p.diagnostic != nil {
		if err := p.diagnostic.Stop(ctx); err != nil {
			p.log.WithError(err).Warn("error stopping diagnostic server")
		}
	}

	// Close listeners
	if p.udpListener != nil {
		p.udpListener.Close()
	}
	if p.tcpListener != nil {
		p.tcpListener.Close()
	}
	if p.tlsListener != nil {
		p.tlsListener.Close()
	}

	// Wait for goroutines to finish
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		return ctx.Err()
	}

	p.BaseNode.SetHealth("unhealthy", map[string]interface{}{
		"stopped": true,
	})

	p.log.Info("P-CSCF stopped")
	return nil
}

// ProcessMessage processes an incoming SIP message.
func (p *PCSCF) ProcessMessage(msg *sip.Message, remoteAddr, transport string) (*sip.Message, error) {
	start := time.Now()
	msg.RemoteAddr = remoteAddr
	msg.Transport = transport

	// Rate limiting
	if p.rateLimiter != nil {
		if !p.rateLimiter.Allow(remoteAddr) {
			p.log.WithFields(logrus.Fields{
				"remote_addr": remoteAddr,
				"method":      msg.Method,
			}).Warn("rate limit exceeded")

			return p.createErrorResponse(msg, sip.StatusServiceUnavailable, "Service Unavailable"), nil
		}
	}

	// Security enforcement
	if err := p.enforceSecurity(msg, remoteAddr); err != nil {
		p.log.WithError(err).Warn("security check failed")
		return p.createErrorResponse(msg, sip.StatusForbidden, "Forbidden"), nil
	}

	// Get or create session
	session := p.getOrCreateSession(msg, remoteAddr, transport)

	// NAT traversal - update Contact header if needed
	if p.config.EnableNATTraversal && msg.IsRequest() {
		p.handleNATTraversal(msg, remoteAddr)
	}

	// Emergency call detection
	if msg.IsRequest() && msg.Method == sip.MethodINVITE {
		if p.isEmergencyCall(msg) {
			p.log.WithFields(logrus.Fields{
				"call_id": msg.GetHeader("Call-ID"),
				"from":    msg.GetHeader("From"),
			}).Info("emergency call detected")
			// Emergency calls bypass normal processing
		}
	}

	// Find and call handler
	var response *sip.Message
	var err error

	if msg.IsRequest() {
		handler := p.handlers[msg.Method]
		if handler != nil {
			response, err = handler(msg, session)
		} else {
			// Default: forward to I-CSCF or S-CSCF
			response, err = p.forwardMessage(msg, session)
		}
	} else {
		// Response - forward back to originator
		response, err = p.forwardResponse(msg, session)
	}

	duration := time.Since(start)
	p.log.WithFields(logrus.Fields{
		"method":   msg.Method,
		"duration": duration,
		"remote":   remoteAddr,
	}).Debug("message processed")

	p.BaseNode.IncrementMessages()
	if err != nil {
		p.BaseNode.IncrementErrors()
	}

	return response, err
}

// startUDP starts the UDP listener
func (p *PCSCF) startUDP() error {
	addr, err := net.ResolveUDPAddr("udp", p.config.SIPAddr)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	p.udpListener = conn

	p.wg.Add(1)
	go p.handleUDP()

	p.log.WithField("addr", addr).Info("UDP listener started")
	return nil
}

// startTCP starts the TCP listener
func (p *PCSCF) startTCP() error {
	listener, err := net.Listen("tcp", p.config.SIPAddr)
	if err != nil {
		return err
	}

	p.tcpListener = listener

	p.wg.Add(1)
	go p.handleTCP()

	p.log.WithField("addr", p.config.SIPAddr).Info("TCP listener started")
	return nil
}

// startTLS starts the TLS listener
func (p *PCSCF) startTLS() error {
	// TODO: Implement TLS listener with certificate management
	p.log.Warn("TLS listener not yet implemented")
	return nil
}

// handleUDP handles UDP connections
func (p *PCSCF) handleUDP() {
	defer p.wg.Done()
	buffer := make([]byte, 65535)
	for {
		select {
		case <-p.shutdown:
			return
		default:
		}

		p.udpListener.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, addr, err := p.udpListener.ReadFromUDP(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			p.log.WithError(err).Error("UDP read error")
			continue
		}

		go p.handleMessage(buffer[:n], addr.String(), "udp")
	}
}

// handleTCP handles TCP connections
func (p *PCSCF) handleTCP() {
	defer p.wg.Done()
	for {
		select {
		case <-p.shutdown:
			return
		default:
		}

		conn, err := p.tcpListener.Accept()
		if err != nil {
			select {
			case <-p.shutdown:
				return
			default:
				p.log.WithError(err).Error("TCP accept error")
				continue
			}
		}

		p.wg.Add(1)
		go p.handleTCPConnection(conn)
	}
}

// handleTCPConnection handles a single TCP connection
func (p *PCSCF) handleTCPConnection(conn net.Conn) {
	defer p.wg.Done()
	defer conn.Close()

	for {
		select {
		case <-p.shutdown:
			return
		default:
		}

		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		msg, err := p.parser.ParseMessage(conn)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			p.log.WithError(err).Error("failed to parse SIP message")
			break
		}

		response, err := p.ProcessMessage(msg, conn.RemoteAddr().String(), "tcp")
		if err != nil {
			p.log.WithError(err).Error("failed to process message")
			break
		}

		if response != nil {
			if _, err := conn.Write([]byte(response.String())); err != nil {
				p.log.WithError(err).Error("failed to write response")
				break
			}
		}
	}
}

// handleMessage handles a SIP message
func (p *PCSCF) handleMessage(data []byte, remoteAddr, transport string) {
	msg, err := p.parser.ParseMessage(bytes.NewReader(data))
	if err != nil {
		p.log.WithError(err).Error("failed to parse SIP message")
		return
	}

	response, err := p.ProcessMessage(msg, remoteAddr, transport)
	if err != nil {
		p.log.WithError(err).Error("failed to process message")
		return
	}

	if response != nil && p.udpListener != nil && transport == "udp" {
		addr, err := net.ResolveUDPAddr("udp", remoteAddr)
		if err == nil {
			p.udpListener.WriteToUDP([]byte(response.String()), addr)
		}
	}
}

// registerDefaultHandlers registers default message handlers
func (p *PCSCF) registerDefaultHandlers() {
	// REGISTER handler
	p.handlers[sip.MethodREGISTER] = p.handleREGISTER

	// INVITE handler
	p.handlers[sip.MethodINVITE] = p.handleINVITE

	// BYE handler
	p.handlers[sip.MethodBYE] = p.handleBYE

	// CANCEL handler
	p.handlers[sip.MethodCANCEL] = p.handleCANCEL

	// ACK handler (stateless)
	p.handlers[sip.MethodACK] = p.handleACK

	// OPTIONS handler
	p.handlers[sip.MethodOPTIONS] = p.handleOPTIONS
}

// handleREGISTER handles REGISTER requests
func (p *PCSCF) handleREGISTER(msg *sip.Message, session *Session) (*sip.Message, error) {
	p.log.WithFields(logrus.Fields{
		"call_id": msg.GetHeader("Call-ID"),
		"from":    msg.GetHeader("From"),
	}).Info("handling REGISTER request")

	// Add P-CSCF address to Path header
	pcscfAddr := p.getPCSCFAddress()
	pathHeader := fmt.Sprintf("<sip:%s;lr>", pcscfAddr)
	msg.AddHeader("Path", pathHeader)

	// Forward to I-CSCF
	return p.forwardToICSCF(msg)
}

// handleINVITE handles INVITE requests
func (p *PCSCF) handleINVITE(msg *sip.Message, session *Session) (*sip.Message, error) {
	p.log.WithFields(logrus.Fields{
		"call_id": msg.GetHeader("Call-ID"),
		"from":    msg.GetHeader("From"),
		"to":      msg.GetHeader("To"),
	}).Info("handling INVITE request")

	// Add Record-Route header
	pcscfAddr := p.getPCSCFAddress()
	recordRoute := fmt.Sprintf("<sip:%s;lr>", pcscfAddr)
	msg.AddHeader("Record-Route", recordRoute)

	// Forward to S-CSCF (if known) or I-CSCF
	if p.config.SCSCFAddr != "" {
		return p.forwardToSCSCF(msg)
	}
	return p.forwardToICSCF(msg)
}

// handleBYE handles BYE requests
func (p *PCSCF) handleBYE(msg *sip.Message, session *Session) (*sip.Message, error) {
	p.log.WithFields(logrus.Fields{
		"call_id": msg.GetHeader("Call-ID"),
	}).Info("handling BYE request")

	// Forward to S-CSCF
	return p.forwardToSCSCF(msg)
}

// handleCANCEL handles CANCEL requests
func (p *PCSCF) handleCANCEL(msg *sip.Message, session *Session) (*sip.Message, error) {
	p.log.WithFields(logrus.Fields{
		"call_id": msg.GetHeader("Call-ID"),
	}).Info("handling CANCEL request")

	// Forward to S-CSCF
	return p.forwardToSCSCF(msg)
}

// handleACK handles ACK requests (stateless)
func (p *PCSCF) handleACK(msg *sip.Message, session *Session) (*sip.Message, error) {
	// ACK is stateless - just forward
	return p.forwardToSCSCF(msg)
}

// handleOPTIONS handles OPTIONS requests
func (p *PCSCF) handleOPTIONS(msg *sip.Message, session *Session) (*sip.Message, error) {
	// Respond to OPTIONS locally
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

	// Add Allow header
	response.SetHeader("Allow", "INVITE, ACK, CANCEL, BYE, REGISTER, OPTIONS, INFO, UPDATE, PRACK, REFER, NOTIFY, SUBSCRIBE")

	return response, nil
}

// forwardMessage forwards a message to the next hop
func (p *PCSCF) forwardMessage(msg *sip.Message, session *Session) (*sip.Message, error) {
	// Add Via header
	via := p.createViaHeader(msg.GetHeader("Via"))
	msg.SetHeader("Via", via)

	// Forward to I-CSCF or S-CSCF
	if p.config.SCSCFAddr != "" {
		return p.forwardToSCSCF(msg)
	}
	return p.forwardToICSCF(msg)
}

// forwardResponse forwards a response back to the originator
func (p *PCSCF) forwardResponse(msg *sip.Message, session *Session) (*sip.Message, error) {
	// Remove top Via header
	via := msg.GetHeader("Via")
	if via != "" {
		// Parse and remove top Via
		parts := strings.Split(via, ",")
		if len(parts) > 1 {
			msg.SetHeader("Via", strings.Join(parts[1:], ","))
		} else {
			delete(msg.Headers, "Via")
		}
	}

	// Forward back to originator (session contains remote address)
	return msg, nil
}

// forwardToICSCF forwards a message to I-CSCF
func (p *PCSCF) forwardToICSCF(msg *sip.Message) (*sip.Message, error) {
	// TODO: Implement actual forwarding via network
	// For now, return a 200 OK response
	p.log.WithField("next_hop", p.config.ICSCFAddr).Debug("forwarding to I-CSCF")
	return p.createOKResponse(msg), nil
}

// forwardToSCSCF forwards a message to S-CSCF
func (p *PCSCF) forwardToSCSCF(msg *sip.Message) (*sip.Message, error) {
	// TODO: Implement actual forwarding via network
	// For now, return a 200 OK response
	p.log.WithField("next_hop", p.config.SCSCFAddr).Debug("forwarding to S-CSCF")
	return p.createOKResponse(msg), nil
}

// enforceSecurity enforces security policies
func (p *PCSCF) enforceSecurity(msg *sip.Message, remoteAddr string) error {
	// Basic security checks
	// TODO: Implement full security enforcement
	// - TLS requirement checking
	// - SRTP requirement checking
	// - Authentication validation
	// - Message integrity checking

	return nil
}

// handleNATTraversal handles NAT traversal by updating Contact header
func (p *PCSCF) handleNATTraversal(msg *sip.Message, remoteAddr string) {
	contact := msg.GetHeader("Contact")
	if contact == "" {
		return
	}

	// Extract IP from remote address
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return
	}

	// Check if this is a private IP (behind NAT)
	if isPrivateIP(host) && p.config.PublicIP != "" {
		// Update Contact header with public IP
		// This is simplified - real NAT traversal is more complex
		p.natMu.Lock()
		p.natMapping[contact] = p.config.PublicIP
		p.natMu.Unlock()
	}
}

// isEmergencyCall checks if a call is an emergency call
func (p *PCSCF) isEmergencyCall(msg *sip.Message) bool {
	// Check Request-URI for emergency numbers
	uri := msg.URI
	for _, emergencyNum := range p.config.EmergencyNumbers {
		if strings.Contains(uri, emergencyNum) {
			return true
		}
	}

	// Check for emergency indication in headers
	if msg.GetHeader("Emergency") == "true" {
		return true
	}

	return false
}

// getOrCreateSession gets or creates a session for a message
func (p *PCSCF) getOrCreateSession(msg *sip.Message, remoteAddr, transport string) *Session {
	callID := msg.GetHeader("Call-ID")
	if callID == "" {
		return nil
	}

	p.sessMu.RLock()
	session, exists := p.sessions[callID]
	p.sessMu.RUnlock()

	if exists {
		session.LastActivity = time.Now()
		return session
	}

	// Create new session
	session = &Session{
		CallID:       callID,
		From:         msg.GetHeader("From"),
		To:           msg.GetHeader("To"),
		Contact:      msg.GetHeader("Contact"),
		RemoteAddr:   remoteAddr,
		Transport:    transport,
		CreatedAt:    time.Now(),
		LastActivity: time.Now(),
	}

	p.sessMu.Lock()
	p.sessions[callID] = session
	p.sessMu.Unlock()

	return session
}

// sessionCleanup periodically cleans up old sessions
func (p *PCSCF) sessionCleanup() {
	defer p.wg.Done()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-p.shutdown:
			return
		case <-ticker.C:
			p.sessMu.Lock()
			now := time.Now()
			for callID, session := range p.sessions {
				if now.Sub(session.LastActivity) > 30*time.Minute {
					delete(p.sessions, callID)
				}
			}
			p.sessMu.Unlock()
		}
	}
}

// createViaHeader creates a Via header
func (p *PCSCF) createViaHeader(existingVia string) string {
	branch := generateBranch()
	pcscfAddr := p.getPCSCFAddress()
	via := fmt.Sprintf("SIP/2.0/UDP %s;branch=z9hG4bK%s", pcscfAddr, branch)
	if existingVia != "" {
		via = via + "," + existingVia
	}
	return via
}

// getPCSCFAddress returns the P-CSCF address
func (p *PCSCF) getPCSCFAddress() string {
	if p.config.PublicIP != "" {
		return p.config.PublicIP
	}
	// Extract host from SIP address
	host, _, err := net.SplitHostPort(p.config.SIPAddr)
	if err != nil || host == "" || host == "0.0.0.0" {
		return "pcscf.ims.local"
	}
	return host
}

// createOKResponse creates a 200 OK response
func (p *PCSCF) createOKResponse(msg *sip.Message) *sip.Message {
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

	return response
}

// createErrorResponse creates an error response
func (p *PCSCF) createErrorResponse(msg *sip.Message, code int, text string) *sip.Message {
	response := &sip.Message{
		Version:    "SIP/2.0",
		StatusCode: code,
		StatusText: text,
		Headers:    make(map[string][]string),
	}

	// Copy required headers
	response.SetHeader("Via", msg.GetHeader("Via"))
	response.SetHeader("From", msg.GetHeader("From"))
	response.SetHeader("To", msg.GetHeader("To"))
	response.SetHeader("Call-ID", msg.GetHeader("Call-ID"))
	response.SetHeader("CSeq", msg.GetHeader("CSeq"))

	return response
}

// generateBranch generates a SIP branch parameter
func generateBranch() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		// Fallback to timestamp-based branch if crypto/rand fails
		return hex.EncodeToString([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	}
	return hex.EncodeToString(b)
}

// isPrivateIP checks if an IP address is private
func isPrivateIP(ip string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}

	return parsed.IsLoopback() ||
		parsed.IsPrivate() ||
		parsed.IsLinkLocalUnicast() ||
		parsed.IsLinkLocalMulticast()
}
