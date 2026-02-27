package sip

import (
	"fmt"
	"strings"
)

// Message represents a SIP message (request or response)
type Message struct {
	// Request line (for requests)
	Method  string
	URI     string
	Version string

	// Status line (for responses)
	StatusCode int
	StatusText string

	// Headers
	Headers map[string][]string

	// Body
	Body string

	// Transport info
	Transport string // "udp", "tcp", "tls"
	RemoteAddr string
}

// IsRequest returns true if this is a SIP request
func (m *Message) IsRequest() bool {
	return m.Method != ""
}

// IsResponse returns true if this is a SIP response
func (m *Message) IsResponse() bool {
	return m.StatusCode > 0
}

// GetHeader returns the first value for a header (case-insensitive)
func (m *Message) GetHeader(name string) string {
	for k, v := range m.Headers {
		if strings.EqualFold(k, name) && len(v) > 0 {
			return v[0]
		}
	}
	return ""
}

// GetHeaderAll returns all values for a header (case-insensitive)
func (m *Message) GetHeaderAll(name string) []string {
	for k, v := range m.Headers {
		if strings.EqualFold(k, name) {
			return v
		}
	}
	return nil
}

// SetHeader sets a header value
func (m *Message) SetHeader(name, value string) {
	if m.Headers == nil {
		m.Headers = make(map[string][]string)
	}
	m.Headers[name] = []string{value}
}

// AddHeader adds a header value
func (m *Message) AddHeader(name, value string) {
	if m.Headers == nil {
		m.Headers = make(map[string][]string)
	}
	m.Headers[name] = append(m.Headers[name], value)
}

// String returns a string representation of the SIP message
func (m *Message) String() string {
	var sb strings.Builder

	if m.IsRequest() {
		sb.WriteString(fmt.Sprintf("%s %s %s\r\n", m.Method, m.URI, m.Version))
	} else {
		sb.WriteString(fmt.Sprintf("%s %d %s\r\n", m.Version, m.StatusCode, m.StatusText))
	}

	for name, values := range m.Headers {
		for _, value := range values {
			sb.WriteString(fmt.Sprintf("%s: %s\r\n", name, value))
		}
	}

	sb.WriteString("\r\n")
	if m.Body != "" {
		sb.WriteString(m.Body)
	}

	return sb.String()
}

// Common SIP methods
const (
	MethodINVITE   = "INVITE"
	MethodACK      = "ACK"
	MethodBYE      = "BYE"
	MethodCANCEL   = "CANCEL"
	MethodOPTIONS  = "OPTIONS"
	MethodREGISTER = "REGISTER"
	MethodINFO     = "INFO"
	MethodUPDATE   = "UPDATE"
	MethodPRACK    = "PRACK"
	MethodREFER    = "REFER"
	MethodNOTIFY   = "NOTIFY"
	MethodSUBSCRIBE = "SUBSCRIBE"
)

// Common SIP status codes
const (
	StatusTrying                = 100
	StatusRinging               = 180
	StatusCallIsBeingForwarded  = 181
	StatusQueued                = 182
	StatusSessionProgress       = 183
	StatusOK                    = 200
	StatusAccepted              = 202
	StatusMultipleChoices       = 300
	StatusMovedPermanently      = 301
	StatusMovedTemporarily      = 302
	StatusUseProxy              = 305
	StatusAlternativeService    = 380
	StatusBadRequest            = 400
	StatusUnauthorized          = 401
	StatusPaymentRequired       = 402
	StatusForbidden             = 403
	StatusNotFound              = 404
	StatusMethodNotAllowed      = 405
	StatusNotAcceptable         = 406
	StatusProxyAuthRequired     = 407
	StatusRequestTimeout        = 408
	StatusGone                  = 410
	StatusRequestEntityTooLarge = 413
	StatusRequestURITooLong     = 414
	StatusUnsupportedMediaType  = 415
	StatusUnsupportedURIScheme  = 416
	StatusBadExtension          = 420
	StatusExtensionRequired     = 421
	StatusIntervalTooBrief      = 423
	StatusTemporarilyUnavailable = 480
	StatusCallLegTransactionDoesNotExist = 481
	StatusLoopDetected          = 482
	StatusTooManyHops           = 483
	StatusAddressIncomplete     = 484
	StatusAmbiguous             = 485
	StatusBusyHere              = 486
	StatusRequestTerminated     = 487
	StatusNotAcceptableHere     = 488
	StatusBadEvent              = 489
	StatusRequestPending        = 491
	StatusUndecipherable        = 493
	StatusInternalServerError   = 500
	StatusNotImplemented        = 501
	StatusBadGateway            = 502
	StatusServiceUnavailable    = 503
	StatusServerTimeout         = 504
	StatusVersionNotSupported   = 505
	StatusMessageTooLarge       = 513
	StatusBusyEverywhere        = 600
	StatusDecline               = 603
	StatusDoesNotExistAnywhere   = 604
	StatusNotAcceptableGlobally  = 606
)
