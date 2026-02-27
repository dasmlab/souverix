package sip

import (
	"strings"
	"testing"
)

func TestParser_ParseMessage_Request(t *testing.T) {
	parser := NewParser()
	
	request := "INVITE sip:bob@example.com SIP/2.0\r\n" +
		"Via: SIP/2.0/UDP 192.168.1.1:5060\r\n" +
		"From: <sip:alice@example.com>;tag=abc123\r\n" +
		"To: <sip:bob@example.com>\r\n" +
		"Call-ID: test-call-id@example.com\r\n" +
		"CSeq: 1 INVITE\r\n" +
		"Content-Length: 0\r\n" +
		"\r\n"

	msg, err := parser.ParseMessage(strings.NewReader(request))
	if err != nil {
		t.Fatalf("ParseMessage() error = %v", err)
	}

	if !msg.IsRequest() {
		t.Error("ParseMessage() should parse as request")
	}
	if msg.Method != MethodINVITE {
		t.Errorf("ParseMessage() method = %v, want %v", msg.Method, MethodINVITE)
	}
	if msg.URI != "sip:bob@example.com" {
		t.Errorf("ParseMessage() URI = %v, want sip:bob@example.com", msg.URI)
	}
	if msg.GetHeader("From") == "" {
		t.Error("ParseMessage() missing From header")
	}
}

func TestParser_ParseMessage_Response(t *testing.T) {
	parser := NewParser()
	
	response := "SIP/2.0 200 OK\r\n" +
		"Via: SIP/2.0/UDP 192.168.1.1:5060\r\n" +
		"From: <sip:alice@example.com>;tag=abc123\r\n" +
		"To: <sip:bob@example.com>;tag=def456\r\n" +
		"Call-ID: test-call-id@example.com\r\n" +
		"CSeq: 1 INVITE\r\n" +
		"Content-Length: 0\r\n" +
		"\r\n"

	msg, err := parser.ParseMessage(strings.NewReader(response))
	if err != nil {
		t.Fatalf("ParseMessage() error = %v", err)
	}

	if !msg.IsResponse() {
		t.Error("ParseMessage() should parse as response")
	}
	if msg.StatusCode != StatusOK {
		t.Errorf("ParseMessage() status code = %v, want %v", msg.StatusCode, StatusOK)
	}
	if msg.StatusText != "OK" {
		t.Errorf("ParseMessage() status text = %v, want OK", msg.StatusText)
	}
}

func TestParser_ParseMessage_WithBody(t *testing.T) {
	parser := NewParser()
	
	request := "INVITE sip:bob@example.com SIP/2.0\r\n" +
		"Content-Type: application/sdp\r\n" +
		"Content-Length: 10\r\n" +
		"\r\n" +
		"v=0\r\no=test"

	msg, err := parser.ParseMessage(strings.NewReader(request))
	if err != nil {
		t.Fatalf("ParseMessage() error = %v", err)
	}

	if msg.Body != "v=0\r\no=test" {
		t.Errorf("ParseMessage() body = %v, want v=0\\r\\no=test", msg.Body)
	}
}

func TestParseRequestLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantErr bool
	}{
		{"valid", "INVITE sip:bob@example.com SIP/2.0", false},
		{"invalid", "INVITE sip:bob@example.com", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, err := ParseRequestLine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRequestLine() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseStatusLine(t *testing.T) {
	tests := []struct {
		name       string
		line        string
		wantCode    int
		wantErr     bool
	}{
		{"200 OK", "SIP/2.0 200 OK", 200, false},
		{"404 Not Found", "SIP/2.0 404 Not Found", 404, false},
		{"invalid", "SIP/2.0 abc OK", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, code, _, err := ParseStatusLine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStatusLine() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && code != tt.wantCode {
				t.Errorf("ParseStatusLine() code = %v, want %v", code, tt.wantCode)
			}
		})
	}
}
