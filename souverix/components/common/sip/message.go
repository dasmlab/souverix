package sip

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

// Message represents a SIP message (request or response)
type Message struct {
	Method     string            // For requests: INVITE, REGISTER, etc. Empty for responses
	StatusCode int               // For responses: 200, 404, etc. 0 for requests
	StatusText string            // For responses: "OK", "Not Found", etc.
	URI        string            // Request-URI for requests
	Version    string            // SIP version (usually "SIP/2.0")
	Headers    map[string]string // SIP headers
	Body       string            // Message body (SDP, etc.)
	Raw        string            // Raw message string
}

// ParseMessage parses a raw SIP message string into a Message struct
func ParseMessage(raw string) (*Message, error) {
	msg := &Message{
		Headers: make(map[string]string),
		Raw:     raw,
	}

	reader := bufio.NewReader(strings.NewReader(raw))
	
	// Parse start line
	startLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read start line: %w", err)
	}
	startLine = strings.TrimSpace(startLine)

	// Determine if request or response
	parts := strings.Fields(startLine)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid start line: %s", startLine)
	}

	if strings.HasPrefix(parts[0], "SIP/") {
		// Response
		msg.Version = parts[0]
		if len(parts) >= 2 {
			fmt.Sscanf(parts[1], "%d", &msg.StatusCode)
		}
		if len(parts) >= 3 {
			msg.StatusText = strings.Join(parts[2:], " ")
		}
	} else {
		// Request
		msg.Method = parts[0]
		msg.URI = parts[1]
		if len(parts) >= 3 {
			msg.Version = parts[2]
		}
	}

	// Parse headers
	var bodyStart bool
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)

		// Empty line indicates start of body
		if line == "" {
			bodyStart = true
			break
		}

		// Handle header continuation (lines starting with space/tab)
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			// Continuation of previous header
			continue
		}

		// Parse header
		colonIdx := strings.Index(line, ":")
		if colonIdx > 0 {
			key := strings.TrimSpace(line[:colonIdx])
			value := strings.TrimSpace(line[colonIdx+1:])
			msg.Headers[strings.ToLower(key)] = value
		}
	}

	// Parse body if present
	if bodyStart {
		var body strings.Builder
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			body.WriteString(line)
		}
		msg.Body = body.String()
	}

	return msg, nil
}

// String returns the string representation of the SIP message
func (m *Message) String() string {
	var buf strings.Builder

	// Write start line
	if m.Method != "" {
		// Request
		buf.WriteString(fmt.Sprintf("%s %s %s\r\n", m.Method, m.URI, m.Version))
	} else {
		// Response
		buf.WriteString(fmt.Sprintf("%s %d %s\r\n", m.Version, m.StatusCode, m.StatusText))
	}

	// Write headers
	for key, value := range m.Headers {
		// Capitalize header name properly
		headerName := capitalizeHeader(key)
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", headerName, value))
	}

	// Empty line before body
	buf.WriteString("\r\n")

	// Write body if present
	if m.Body != "" {
		buf.WriteString(m.Body)
	}

	return buf.String()
}

// GetHeader returns a header value (case-insensitive)
func (m *Message) GetHeader(name string) string {
	return m.Headers[strings.ToLower(name)]
}

// SetHeader sets a header value
func (m *Message) SetHeader(name, value string) {
	m.Headers[strings.ToLower(name)] = value
}

// AddHeader adds a header (appends if exists)
func (m *Message) AddHeader(name, value string) {
	key := strings.ToLower(name)
	if existing, exists := m.Headers[key]; exists {
		m.Headers[key] = existing + ", " + value
	} else {
		m.Headers[key] = value
	}
}

// GetCallID returns the Call-ID header
func (m *Message) GetCallID() string {
	return m.GetHeader("Call-ID")
}

// GetFrom returns the From header
func (m *Message) GetFrom() string {
	return m.GetHeader("From")
}

// GetTo returns the To header
func (m *Message) GetTo() string {
	return m.GetHeader("To")
}

// GetCSeq returns the CSeq header
func (m *Message) GetCSeq() string {
	return m.GetHeader("CSeq")
}

// GetVia returns the Via header
func (m *Message) GetVia() string {
	return m.GetHeader("Via")
}

// IsRequest returns true if this is a SIP request
func (m *Message) IsRequest() bool {
	return m.Method != ""
}

// IsResponse returns true if this is a SIP response
func (m *Message) IsResponse() bool {
	return m.Method == "" && m.StatusCode > 0
}

// capitalizeHeader capitalizes header name according to SIP conventions
func capitalizeHeader(name string) string {
	parts := strings.Split(name, "-")
	var result strings.Builder
	for i, part := range parts {
		if i > 0 {
			result.WriteString("-")
		}
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(part[:1]))
			if len(part) > 1 {
				result.WriteString(strings.ToLower(part[1:]))
			}
		}
	}
	return result.String()
}

// CreateRequest creates a new SIP request message
func CreateRequest(method, uri string) *Message {
	return &Message{
		Method:  method,
		URI:     uri,
		Version: "SIP/2.0",
		Headers: make(map[string]string),
	}
}

// CreateResponse creates a new SIP response message
func CreateResponse(statusCode int, statusText string, request *Message) *Message {
	resp := &Message{
		StatusCode: statusCode,
		StatusText: statusText,
		Version:    "SIP/2.0",
		Headers:    make(map[string]string),
	}

	// Copy relevant headers from request
	if request != nil {
		if cseq := request.GetCSeq(); cseq != "" {
			resp.SetHeader("CSeq", cseq)
		}
		if callID := request.GetCallID(); callID != "" {
			resp.SetHeader("Call-ID", callID)
		}
		if from := request.GetFrom(); from != "" {
			resp.SetHeader("From", from)
		}
		if to := request.GetTo(); to != "" {
			resp.SetHeader("To", to)
		}
		if via := request.GetVia(); via != "" {
			resp.SetHeader("Via", via)
		}
	}

	return resp
}
