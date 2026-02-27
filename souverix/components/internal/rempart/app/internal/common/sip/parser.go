package sip

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Parser parses SIP messages
type Parser struct{}

// NewParser creates a new SIP parser
func NewParser() *Parser {
	return &Parser{}
}

// ParseMessage parses a SIP message from a reader
func (p *Parser) ParseMessage(reader io.Reader) (*Message, error) {
	msg := &Message{
		Headers: make(map[string][]string),
		Version: "SIP/2.0",
	}

	scanner := bufio.NewScanner(reader)
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty message")
	}

	// Parse start line
	startLine := scanner.Text()
	if err := p.parseStartLine(msg, startLine); err != nil {
		return nil, err
	}

	// Parse headers
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break // Empty line indicates end of headers
		}

		// Handle continuation lines (lines starting with space/tab)
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			// Append to last header
			if len(msg.Headers) > 0 {
				// Find last header key
				var lastKey string
				for k := range msg.Headers {
					lastKey = k
					break
				}
				if lastKey != "" && len(msg.Headers[lastKey]) > 0 {
					msg.Headers[lastKey][len(msg.Headers[lastKey])-1] += " " + strings.TrimSpace(line)
				}
			}
			continue
		}

		// Parse header
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue // Skip malformed headers
		}

		name := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Add header (support multiple values)
		if existing, ok := msg.Headers[name]; ok {
			msg.Headers[name] = append(existing, value)
		} else {
			msg.Headers[name] = []string{value}
		}
	}

	// Parse body if Content-Length is present
	if cl := msg.GetHeader("Content-Length"); cl != "" {
		if length, err := strconv.Atoi(strings.TrimSpace(cl)); err == nil && length > 0 {
			body := make([]byte, length)
			if n, err := io.ReadFull(reader, body); err == nil {
				msg.Body = string(body[:n])
			}
		}
	}

	return msg, scanner.Err()
}

// parseStartLine parses the start line (request or response)
func (p *Parser) parseStartLine(msg *Message, line string) error {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return fmt.Errorf("invalid start line: %s", line)
	}

	// Check if it's a request or response
	if strings.HasPrefix(parts[0], "SIP/") {
		// Response
		msg.Version = parts[0]
		statusCode, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("invalid status code: %s", parts[1])
		}
		msg.StatusCode = statusCode
		msg.StatusText = strings.Join(parts[2:], " ")
	} else {
		// Request
		msg.Method = parts[0]
		msg.URI = parts[1]
		msg.Version = parts[2]
	}

	return nil
}

// ParseRequestLine parses a SIP request line
func ParseRequestLine(line string) (method, uri, version string, err error) {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return "", "", "", fmt.Errorf("invalid request line: %s", line)
	}
	return parts[0], parts[1], parts[2], nil
}

// ParseStatusLine parses a SIP status line
func ParseStatusLine(line string) (version string, statusCode int, statusText string, err error) {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return "", 0, "", fmt.Errorf("invalid status line: %s", line)
	}

	version = parts[0]
	statusCode, err = strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, "", fmt.Errorf("invalid status code: %s", parts[1])
	}

	statusText = strings.Join(parts[2:], " ")
	return version, statusCode, statusText, nil
}
