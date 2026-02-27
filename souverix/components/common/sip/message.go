package sip

import (
	"fmt"
	"strings"
)

// Message represents a SIP message
type Message struct {
	Method   string
	URI      string
	Version  string
	Headers  map[string]string
	Body     string
	From     string
	To       string
	CallID   string
	CSeq     string
	Contact  string
	Route    []string
	RecordRoute []string
}

// NewINVITE creates a new SIP INVITE message
func NewINVITE(from, to, callID string) *Message {
	return &Message{
		Method:  "INVITE",
		URI:     to,
		Version: "SIP/2.0",
		Headers: make(map[string]string),
		From:    from,
		To:      to,
		CallID:  callID,
		CSeq:    "1 INVITE",
		Route:   make([]string, 0),
		RecordRoute: make([]string, 0),
	}
}

// New200OK creates a 200 OK response
func New200OK(invite *Message, contact string) *Message {
	return &Message{
		Method:  "200",
		Version: "SIP/2.0",
		Headers: make(map[string]string),
		From:    invite.From,
		To:      invite.To,
		CallID:  invite.CallID,
		CSeq:    invite.CSeq,
		Contact: contact,
		Route:   invite.Route,
		RecordRoute: invite.RecordRoute,
	}
}

// New180Ringing creates a 180 Ringing response
func New180Ringing(invite *Message) *Message {
	return &Message{
		Method:  "180",
		Version: "SIP/2.0",
		Headers: make(map[string]string),
		From:    invite.From,
		To:      invite.To,
		CallID:  invite.CallID,
		CSeq:    invite.CSeq,
		Route:   invite.RecordRoute,
	}
}

// String returns the SIP message as a string
func (m *Message) String() string {
	var sb strings.Builder
	
	// Request/Response line
	if m.Method == "200" || m.Method == "180" {
		sb.WriteString(fmt.Sprintf("SIP/2.0 %s\r\n", m.Method))
	} else {
		sb.WriteString(fmt.Sprintf("%s %s %s\r\n", m.Method, m.URI, m.Version))
	}
	
	// Headers
	if m.From != "" {
		sb.WriteString(fmt.Sprintf("From: %s\r\n", m.From))
	}
	if m.To != "" {
		sb.WriteString(fmt.Sprintf("To: %s\r\n", m.To))
	}
	if m.CallID != "" {
		sb.WriteString(fmt.Sprintf("Call-ID: %s\r\n", m.CallID))
	}
	if m.CSeq != "" {
		sb.WriteString(fmt.Sprintf("CSeq: %s\r\n", m.CSeq))
	}
	if m.Contact != "" {
		sb.WriteString(fmt.Sprintf("Contact: %s\r\n", m.Contact))
	}
	for _, route := range m.Route {
		sb.WriteString(fmt.Sprintf("Route: %s\r\n", route))
	}
	for _, rr := range m.RecordRoute {
		sb.WriteString(fmt.Sprintf("Record-Route: %s\r\n", rr))
	}
	
	sb.WriteString("\r\n")
	if m.Body != "" {
		sb.WriteString(m.Body)
	}
	
	return sb.String()
}

// IsTelURI checks if the URI is a tel: URI (PSTN destination)
func (m *Message) IsTelURI() bool {
	return strings.HasPrefix(m.URI, "tel:") || strings.HasPrefix(m.To, "tel:")
}

// AddRecordRoute adds a Record-Route header
func (m *Message) AddRecordRoute(route string) {
	m.RecordRoute = append(m.RecordRoute, route)
}

// AddRoute adds a Route header
func (m *Message) AddRoute(route string) {
	m.Route = append(m.Route, route)
}
